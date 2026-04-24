// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/marco-m/flightplan/internal/goof"
	"github.com/marco-m/flightplan/resources"
	"github.com/mitchellh/copystructure"
)

var (
	ErrMissingNewPipeline    = errors.New("must use NewPipeline to create a pipeline")
	ErrCreatePipelineDir     = errors.New("cannot create pipeline directory")
	ErrDuplicateJobName      = errors.New("pipeline cannot have duplicate job name")
	ErrDuplicateResourceName = errors.New("pipeline cannot have duplicate resource name")
	ErrNoJobs                = errors.New("pipeline must have at least one job")
	ErrEmptyPipelineName     = errors.New("pipeline name cannot be empty")
	ErrEmptyJobName          = errors.New("job name cannot be empty")
	ErrTaskBothConfigAndFile = errors.New("task cannot have both Config and File")
	ErrTaskNoConfigNoFile    = errors.New("task must have Config or File")
	ErrTaskNoName            = errors.New("task field 'Task' cannot be empty")
	ErrImageResource         = errors.New("Config.ImageResource cannot be empty")
	ErrSetImageResourceType  = errors.New("Config.ImageResource.Type cannot be set (will be set by Source.Type)")
	ErrImageResourceSource   = errors.New("Config.ImageResource.Source cannot be empty")
	ErrEmptyResourceName     = errors.New("Resource Name cannot be empty")
	ErrSetResourceType       = errors.New("Resource Type cannot be set (will be set by Source.Type)")
	ErrSourceValidation      = errors.New("validating Source")
	ErrMissingSource         = errors.New("field Source cannot be empty")
	ErrGetValidation         = errors.New("Get validation")
	ErrPutValidation         = errors.New("Put validation")
	ErrDuplicateExtTaskFile  = errors.New("duplicate task file")
	ErrNotASentinelDir       = errors.New("handle.Source does not implement SentinelDir")
	ErrRepoRootNotFound      = errors.New("repository root not found")
	ErrSystem                = errors.New("system error")
	ErrInternal              = errors.New("flightplan internal error, please report")
)

// Pipeline is used to construct a Concourse pipeline. Use [NewPipeline] to instantiate.
type Pipeline struct {
	// The current working directory.
	cwd string
	// The directory set on the command-line
	dir string
	// The name of the pipeline.
	name string
	// The object that will be serialized by [Pipeline.Render].
	po       pipelineObject
	extTasks []*externalTask
	// Errors collected during construction and returned by [Pipeline.Render].
	errs        []error
	newPipeline bool // True if instantiated by NewPipeline.
}

// NewPipeline parses the command-line and returns a *[Pipeline] ready to use.
//
// Param 'name' is the the name of the pipeline, used to construct its file name as
// <name.json>. Can be overridden on the command-line.
//
// Note: NewPipeline can call [os.Exit] during command-line parsing.
//
// Usage:
//
// pipeline := plan.NewPipeline("default-name", "default-dir", os.Args[1:])
func NewPipeline(name string, args []string) *Pipeline {
	pl := &Pipeline{name: name, newPipeline: true}

	cwd, err := os.Getwd()
	if err != nil {
		pl.errs = append(pl.errs, goof.Wrap("NewPipeline: %w", ErrSystem))
	}
	if pl.name == "" {
		pl.errs = append(pl.errs, goof.Wrap("NewPipeline: %w", ErrEmptyPipelineName))
	}
	// TODO ADD OTHER CHECKS: should be only 0-9 A-Z a-z _-
	// For sure not / anywhere
	// no . as first character
	// no . as last character

	parseCommandLine(pl, args)
	// Make pl.dir absolute
	if !filepath.IsAbs(pl.dir) {
		pl.dir = filepath.Join(cwd, pl.dir)
	}

	if err := os.MkdirAll(pl.dir, 0o775); err != nil {
		pl.errs = append(pl.errs, goof.Wrap("NewPipeline: %w: %s", ErrCreatePipelineDir, err))
	}

	return pl
}

func (pl *Pipeline) AddJob(job Job) JobHandle {
	// Since parameter 'job' contains a slice and we modify it, we need to make a deep
	// copy of it, to avoid surprising the caller.
	dup, err := copystructure.Copy(job)
	if err != nil {
		panic(fmt.Errorf("%w: copying Job: %w", ErrInternal, err))
	}
	job2 := dup.(Job)

	if job2.Name == "" {
		pl.errs = append(pl.errs, goof.Wrap("AddJob: %w", ErrEmptyJobName))
		return ""
	}
	for _, step := range job2.Plan {
		if err := step.Validate(pl); err != nil {
			pl.errs = append(pl.errs, goof.Wrap("AddJob: %w", err))
			// TODO Is the return needed ???
			return ""
		}
		if task, ok := step.(Task); ok {
			extTask, err := task.Process(pl.extTasks)
			if err != nil {
				pl.errs = append(pl.errs, goof.Wrap("AddJob: %w", err))
				continue
			}
			if extTask != nil {
				pl.extTasks = append(pl.extTasks, extTask)
			}
		}
	}
	for _, j := range pl.po.Jobs {
		if j.Name == job2.Name {
			pl.errs = append(pl.errs, goof.Wrap("AddJob: %w: %q", ErrDuplicateJobName, j.Name))
		}
	}
	pl.po.Jobs = append(pl.po.Jobs, job2)
	return JobHandle(job2.Name)
}

// Job returns a copy of the [Job] associated with 'handle'.
// Client code doesn't need to call this function.
func (pl *Pipeline) Job(handle JobHandle) (res Job, found bool) {
	for _, r := range pl.po.Jobs {
		if JobHandle(r.Name) == handle {
			return r, true
		}
	}
	return Job{}, false
}

// AddResource adds resource 'res' to the pipeline. Any error will be returned by
// [Pipeline.Render].
func (pl *Pipeline) AddResource(res resources.Resource) *resources.Handle {
	if res.Name == "" {
		pl.errs = append(pl.errs,
			goof.Wrap("AddResource: %w", ErrEmptyResourceName))
		return nil
	}
	if res.Type != "" {
		pl.errs = append(pl.errs,
			goof.Wrap("AddResource: %s: %w: %q", res.Name, ErrSetResourceType, res.Type))
		return nil
	}
	if res.Source == nil {
		pl.errs = append(pl.errs,
			goof.Wrap("AddResource: %s: %w", res.Name, ErrMissingSource))
		return nil
	}
	if err := res.Source.Validate(); err != nil {
		pl.errs = append(pl.errs,
			goof.Wrap("AddResource: %s: %w: %w", res.Name, ErrSourceValidation, err))
		return nil
	}
	res.Type = res.Source.Type()
	for _, r := range pl.po.Resources {
		if r.Name == res.Name {
			pl.errs = append(pl.errs,
				goof.Wrap("AddResource: %s: %w", res.Name, ErrDuplicateResourceName))
			return nil
		}
	}
	pl.po.Resources = append(pl.po.Resources, res)
	return &resources.Handle{Resource: res}
}

// Resource returns a copy of the [Resource] associated with 'handle'.
// Client code doesn't need to call this function.
func (pl *Pipeline) Resource(handle *resources.Handle) (res resources.Resource, found bool) {
	for _, r := range pl.po.Resources {
		if r.Name == handle.Name {
			return r, true
		}
	}
	return resources.Resource{}, false
}

// Render generates the pipeline and (eventually) the external taskfiles and writes them
// by default to the same directory where the program is running. The directory can be
// changed to a relative or absolute path with the command-line option --directory.
// Render retuns all the errors collected during the usage of the flightplan API.
func (pl *Pipeline) Render() error {
	if !pl.newPipeline {
		// prepend
		pl.errs = slices.Insert(pl.errs, 0, goof.Wrap("Render: %w", ErrMissingNewPipeline))
	}
	if len(pl.po.Jobs) == 0 {
		pl.errs = append(pl.errs, goof.Wrap("Render: %w", ErrNoJobs))
	}

	if err := pl.Errors(); err != nil {
		return err
	}

	writeOne := func(dstPath string, data any) error {
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o775); err != nil {
			return err
		}
		wr, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer wr.Close()

		enc := json.NewEncoder(wr)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			return err
		}
		return nil
	}

	dst := pl.Path()
	if err := writeOne(dst, pl.po); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote pipeline to %s\n", dst)

	for _, extTask := range pl.extTasks {
		root, relpath := reconstructRepoRoot(pl.dir, extTask.File)
		dst := filepath.Join(root, relpath)
		if err := writeOne(dst, extTask.FileConfig); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "wrote taskfile to %s\n", dst)
	}
	return nil
}

// Path returns the path of the rendered pipeline.
func (pl *Pipeline) Path() string {
	return filepath.Join(pl.dir, pl.name) + ".json"
}

// RelDir returns the relative path from the root of the SCM repository to the directory
// containing the pipeline, with prepended the name of the SCM repository. For this to
// work, 'handle.Source' must be a SCM repository, such as [resources.Git]. Client code
// can use it to construct correct paths for fields [Task.File] and [TaskCommand.Path].
func (pl *Pipeline) RelDir(handle *resources.Handle) string {
	if sentinel, ok := handle.Source.(resources.SentinelDir); ok {
		// Assumption: pl.dir is absolute
		repoRoot, err := findRepoRoot(pl.dir, sentinel.SentinelDir())
		if err != nil {
			pl.errs = append(pl.errs, goof.Wrap("RelDir: %w", err))
			return ""
		}
		rel, err := filepath.Rel(repoRoot, pl.dir)
		if err != nil {
			pl.errs = append(pl.errs, goof.Wrap("RelDir: %w", err))
			return ""
		}
		return path.Join(handle.Name, rel)
	}
	pl.errs = append(pl.errs, goof.Wrap("RelDir: %w (type: %s)", ErrNotASentinelDir, handle.Type))
	return ""
}

// Errors returns all the errors so far, joined into a single error.
// Does not drain the errors: [Pipeline.Render] will still return all of them.
// Client code doesn't need to call this function.
func (pl *Pipeline) Errors() error {
	return errors.Join(pl.errs...)
}

// Parses the command-line and overrides settings in 'pl' based on the parse.
// Can call directly or indirectly [os.Exit].
func parseCommandLine(pl *Pipeline, args []string) {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s -- generate a Concourse pipeline with flightplan\n\n",
			fs.Name())
		fmt.Fprintf(os.Stderr, "Usage: %s [ARGS]\n\n", fs.Name())
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}
	fs.StringVar(&pl.name, "name", pl.name, "Name of the pipeline file")
	fs.StringVar(&pl.dir, "directory", ".", "Directory in which to write the pipeline file")
	fs.Parse(args)
	if len(fs.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "unrecognized arguments: %v\n", strings.Join(fs.Args(), " "))
		os.Exit(1)
	}
}

func findRepoRoot(path, sentinelDir string) (string, error) {
	candidateRoot := path
loop:
	for {
		entries, err := os.ReadDir(candidateRoot)
		if err != nil {
			return "", goof.Wrap("findRepoRoot: %w", err)
		}
		for _, entry := range entries {
			if entry.IsDir() && entry.Name() == sentinelDir {
				// found
				break loop
			}
		}
		// iterate
		previous := candidateRoot
		candidateRoot = filepath.Dir(candidateRoot)
		if candidateRoot == previous {
			// Arrived at the root of the filesystem, no sentinel dir found.
			return "", fmt.Errorf("%w (starting at %s)", ErrRepoRootNotFound, path)
		}
	}
	return candidateRoot, nil
}

// reconstructRepoRoot takes 'taskpath' from [Task.File] (format:
// reponame/taskdir/taskfile) and 'pldir', the absolute directory of the pipeline. It
// returns the repository root (repoRoot) and the task path without the reponame
// (relTaskPath).
// See also [Pipeline.RelDir], [findRepoRoot].
func reconstructRepoRoot(pldir, taskpath string) (repoRoot, relTaskPath string) {
	// taskpath has format: reponame/taskdir/taskfile (taskdir can be multiple segments)
	// remove reponame at beginning of path
	_, relTaskPath, _ = strings.Cut(taskpath, "/")
	// remove taskfile at end of path
	dir := filepath.Dir(relTaskPath)
	var root string
	for {
		var found bool
		root, found = strings.CutSuffix(pldir, dir)
		if found {
			break
		}
		dir = filepath.Dir(dir)
		if dir == "." {
			break
		}
	}
	return path.Clean(root), relTaskPath
}

type pipelineObject struct {
	// Resources is the Concourse "resources" object.
	// Use [Pipeline.AddResource] to add to it.
	Resources []resources.Resource `json:"resources,omitempty"`
	// Jobs is the Concourse "jobs" object.
	// Use [Pipeline.AddJob] to add to it.
	Jobs []Job `json:"jobs,omitempty"`
}

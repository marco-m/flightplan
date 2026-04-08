// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/marco-m/flightplan/internal/goof"
)

var (
	ErrDuplicateJob       = errors.New("pipeline cannot have duplicate job")
	ErrDuplicateResource  = errors.New("pipeline cannot have duplicate resource")
	ErrNoJobs             = errors.New("pipeline must have at least one job")
	ErrEmptyPipelineName  = errors.New("pipeline name cannot be empty")
	ErrEmptyJobName       = errors.New("job name cannot be empty")
	ErrMissingNewPipeline = errors.New("must use NewPipeline to create a pipeline")
	ErrEmptyResourceName  = errors.New("resource name cannot be empty")
	ErrSetResourceType    = errors.New("resource type cannot be set (will be set by Source.Type)")
	ErrSystem             = errors.New("unexpected")
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
	po pipelineObject
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
	pl.cwd = cwd
	if pl.name == "" {
		pl.errs = append(pl.errs, goof.Wrap("NewPipeline: %w", ErrEmptyPipelineName))
	}
	// TODO ADD OTHER CHECKS: should be only 0-9 A-Z a-z _-
	// For sure not / anywhere
	// no . as first character
	// no . as last character

	parseCommandLine(pl, args)

	return pl
}

// AddResource adds resource 'res' to the pipeline. Any error will be returned by
// [Pipeline.Render].
func (pl *Pipeline) AddResource(res Resource) ResourceHandle {
	if res.Name == "" {
		pl.errs = append(pl.errs,
			goof.Wrap("AddResource: %w", ErrEmptyResourceName))
		return ""
	}
	if res.Type != "" {
		pl.errs = append(pl.errs,
			goof.Wrap("AddResource: %w: %q", ErrSetResourceType, res.Type))
		return ""
	}
	res.Type = res.Source.Type()
	for _, r := range pl.po.Resources {
		if r.Name == res.Name {
			pl.errs = append(pl.errs,
				goof.Wrap("AddResource: %w: %q", ErrDuplicateResource, r.Name))
			return ""
		}
	}
	pl.po.Resources = append(pl.po.Resources, res)
	return ResourceHandle(res.Name)
}

func (pl *Pipeline) AddJob(job Job) JobHandle {
	if job.Name == "" {
		pl.errs = append(pl.errs, goof.Wrap("AddJob: %w", ErrEmptyJobName))
		return ""
	}
	for _, j := range pl.po.Jobs {
		if j.Name == job.Name {
			pl.errs = append(pl.errs, goof.Wrap("AddJob: %w: %q", ErrDuplicateJob, j.Name))
		}
	}
	pl.po.Jobs = append(pl.po.Jobs, job)
	return JobHandle(job.Name)
}

// Render generates the pipeline and writes it by default to the same directory where
// the program is running. The directory can be changed to a relative or absolute path
// with the command-line option --directory.
// Render retuns all the errors collected during the usage of the flightplan API.
func (pl *Pipeline) Render() error {
	if !pl.newPipeline {
		// prepend
		pl.errs = slices.Insert(pl.errs, 0, goof.Wrap("render: %w", ErrMissingNewPipeline))
	}
	if len(pl.po.Jobs) == 0 {
		pl.errs = append(pl.errs, goof.Wrap("render: %w", ErrNoJobs))
	}

	if err := pl.Errors(); err != nil {
		return err
	}

	dstPath := pl.Path()
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
	if err := enc.Encode(pl.po); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote pipeline to %s\n", dstPath)
	return nil
}

// Path returns the path of the rendered pipeline.
// Client code doesn't need to call this function.
func (pl *Pipeline) Path() string {
	name := pl.name
	if filepath.Ext(name) == "" {
		name += ".json"
	}
	if filepath.IsAbs(pl.dir) {
		// Do not consider the current working directory.
		return filepath.Join(pl.dir, name)
	}
	return filepath.Join(pl.cwd, pl.dir, name)
}

// Errors returns all the errors so far, joined into a single error.
// Does not drain the errors: [Pipeline.Render] will still return all of them.
// Client code doesn't need to call this function.
func (pl *Pipeline) Errors() error {
	return errors.Join(pl.errs...)
}

// Resource returns a copy of the [Resource] associated with 'handle'.
// Client code doesn't need to call this function.
func (pl *Pipeline) Resource(handle ResourceHandle) (res Resource, found bool) {
	for _, r := range pl.po.Resources {
		if ResourceHandle(r.Name) == handle {
			return r, true
		}
	}
	return Resource{}, false
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

type pipelineObject struct {
	// Resources is the Concourse "resources" object.
	// Use [Pipeline.AddResource] to add to it.
	Resources []Resource `json:"resources,omitempty"`
	// Jobs is the Concourse "jobs" object.
	// Use [Pipeline.AddJob] to add to it.
	Jobs []Job `json:"jobs,omitempty"`
}

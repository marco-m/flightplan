// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"flag"
	"path"
	"path/filepath"
	"strings"
	"testing"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/flightplan/internal/testhelpers"
	"github.com/marco-m/flightplan/resources"

	"github.com/marco-m/rosina/assert"
)

var update = flag.Bool("golden.update", false,
	"update the golden files for this package; use with -run to update a single test")

func makeTestJobInlineTask() plan.Job {
	return plan.Job{
		Name: "bake-pizza",
		Plan: []plan.Step{
			plan.Task{
				Task: "knead",
				Config: &plan.TaskConfig{
					Platform: "linux",
					ImageResource: &resources.AnonymousResource{
						Source: resources.RegistryImage{Repository: "alpine"},
					},
					Run: plan.TaskCommand{Path: "echo", Args: []string{"Pizza"}},
				},
			},
		},
	}
}

func makeTestJobExternalTask(pl *plan.Pipeline, repo *resources.Handle, jobName, taskName, taskDir string) plan.Job {
	taskPath := path.Join(pl.RelDir(repo), taskDir, taskName+".json")
	return plan.Job{
		Name: jobName,
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: taskName,
				File: taskPath,
				FileConfig: &plan.TaskConfig{
					Platform: "linux",
					ImageResource: &resources.AnonymousResource{
						Source: resources.RegistryImage{Repository: "alpine"},
					},
					Run: plan.TaskCommand{Path: "echo", Args: []string{taskName}},
				},
			},
		},
	}
}

func TestPipelinePath(t *testing.T) {
	test := func(desc, name string, args []string, wantSuffix string) {
		t.Helper()
		pipeline := plan.NewPipeline(name, args)
		assert.True(t, strings.HasSuffix(filepath.ToSlash(pipeline.Path()), wantSuffix),
			"Path "+desc)
		assert.True(t, filepath.IsAbs(pipeline.Path()), "Path is absolute")
	}

	test("default dir", "dummy", nil, "/flightplan/dummy.json")
	test("cli override name", "dummy", []string{"-name=mango"},
		"/flightplan/mango.json")
	test("cli override relative dir", "dummy", []string{"-directory=berry"},
		"/flightplan/berry/dummy.json")
	test("cli override absolute dir", "dummy", []string{"-directory=/mango"},
		"/mango/dummy.json")
}

func TestMustUseNewPipeline(t *testing.T) {
	pipeline := plan.Pipeline{}
	err := pipeline.Render()
	assert.ErrorIs(t, err, plan.ErrMissingNewPipeline, "Render")
}

func TestEmptyPipelineIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("empty", nil)
	err := pipeline.Render()
	assert.ErrorIs(t, err, plan.ErrNoJobs, "Render")
}

func TestPipelineWithoutNameIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("", nil)
	err := pipeline.Render()
	assert.ErrorIs(t, err, plan.ErrEmptyPipelineName, "Render")
}

func TestPipelineWithOneJobIsValid(t *testing.T) {
	dir := t.TempDir()
	pipeline := plan.NewPipeline("one-job", []string{"--directory", dir})
	pipeline.AddJob(makeTestJobInlineTask())

	err := pipeline.Render()

	assert.NoError(t, err, "Render")
	testhelpers.AssertRenderedEqualsGolden(t, pipeline.Path(), "testdata/one-job.json", *update)
}

func TestPipelineNamedImageResource(t *testing.T) {
	dir := t.TempDir()
	pl := plan.NewPipeline("named-image-resource", []string{"--directory", dir})
	repo := pl.AddResource(resources.Resource{
		Name:   "flightplan.git",
		Source: resources.Git{Uri: "https://github.com/marco-m/flightplan.git"},
	})
	alpineImage := pl.AddResource(resources.Resource{
		Name:   "alpine.image",
		Source: resources.RegistryImage{Repository: "alpine"},
	})

	pl.AddJob(plan.Job{
		Name: "the-job",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Get{Get: alpineImage},
			plan.Task{
				Task:  "the-task",
				Image: alpineImage,
				Config: &plan.TaskConfig{
					Platform: "linux",
					Run: plan.TaskCommand{
						Path: "sh", Args: []string{"-c", `echo "hello"`},
					},
				},
			},
		},
	})
	err := pl.Render()

	// FIXME expect error missing name for image resource!!!
	// FIXME expect error if no ImageResource and no Image
	// FIXME named image BUT missing GET of said image!
	// FXIME both named and anon image: anon image is still written to task (as Concourse does)
	assert.NoError(t, err, "Render")
	testhelpers.AssertRenderedEqualsGolden(t, pl.Path(), "testdata/named-image-resource.json", *update)
}

func TestPipelineJobGetAndPutResource(t *testing.T) {
	dir := t.TempDir()
	pipeline := plan.NewPipeline("get-and-put", []string{"--directory", dir})
	repo := pipeline.AddResource(resources.Resource{
		Name:   "flightplan.git",
		Source: resources.Git{Uri: "https://github.com/marco-m/flightplan.git"},
	})
	artifacts := pipeline.AddResource(resources.Resource{
		Name: "artifacts.s3",
		Source: resources.S3{
			Bucket:          "concourse",
			Regexp:          "builds/simple-s3/gift-(.*)",
			Endpoint:        "((s3-endpoint))",
			RegionName:      "((s3-region))",
			AccessKeyID:     "((s3-access-key))",
			SecretAccessKey: "((s3-secret-key))",
			UsePathStyle:    true,
		},
	})
	alpineImage := resources.AnonymousResource{
		Source: resources.RegistryImage{Repository: "alpine"},
	}
	const makeGift = `
set -ex
VERSION=$(date +%Y%m%d%H%M%S)
GIFT=gift/gift-$VERSION
echo "hello" > $GIFT
`
	kneadPizza := pipeline.AddJob(plan.Job{
		Name: "knead-pizza",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: "prepare-dough",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &alpineImage,
					Outputs:       []plan.TaskOutput{{Name: "gift"}},
					Run: plan.TaskCommand{
						Path: "sh", Args: []string{"-c", makeGift},
					},
				},
			},
			plan.Put{
				Put:    artifacts,
				Params: resources.S3PutParams{File: "gift/gift-*"},
			},
		},
	})
	pipeline.AddJob(plan.Job{
		Name: "bake-pizza",
		Plan: []plan.Step{
			plan.Get{Get: repo, Passed: []plan.JobHandle{kneadPizza}, Trigger: true},
			plan.Get{Get: artifacts, Passed: []plan.JobHandle{kneadPizza}, Trigger: true},
			plan.Task{
				Task: "insert-in-hoven",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &alpineImage,
					Run:           plan.TaskCommand{Path: "env"},
				},
			},
		},
	})

	err := pipeline.Render()

	assert.NoError(t, err, "Render")
	testhelpers.AssertRenderedEqualsGolden(t, pipeline.Path(), "testdata/get-and-put.json", *update)
}

func TestPipelineExternalTaskfileFailure(t *testing.T) {
	dir := t.TempDir()
	pl := plan.NewPipeline("external-taskfile", []string{"--directory", dir})
	repo := pl.AddResource(resources.Resource{
		Name:   "banana.git",
		Source: resources.Git{Uri: "https://example.org/flightplan.git"},
	})
	pl.AddJob(makeTestJobExternalTask(pl, repo, "job1", "mango", "tasks"))
	pl.AddJob(makeTestJobExternalTask(pl, repo, "job2", "mango", "tasks"))
	err := pl.Errors()

	assert.ErrorIs(t, err, plan.ErrDuplicateExtTaskFile, "Errors")
}

func TestPipelineCreateOneExternalTaskfileSuccess(t *testing.T) {
	dir := t.TempDir()
	testhelpers.MakeFakeGitRepo(t, dir)
	name := "with-taskfile"
	pl := plan.NewPipeline(name, []string{"--directory", dir})
	repo := pl.AddResource(resources.Resource{
		Name:   "banana.git",
		Source: resources.Git{Uri: "https://example.org/flightplan.git"},
	})

	pl.AddJob(makeTestJobExternalTask(pl, repo, "job-1", "task-1", "tasks"))
	err := pl.Render()
	assert.NoError(t, err, "Render")

	testhelpers.AssertRenderedEqualsGolden(t, pl.Path(), "testdata/with-taskfile.json", *update)
	testhelpers.AssertRenderedEqualsGolden(t, filepath.Join(dir, "tasks/task-1.json"),
		"testdata/task-1.json", *update)
}

func TestPipelineCreateTwoExternalTaskfileSuccess(t *testing.T) {
	dir := t.TempDir()
	testhelpers.MakeFakeGitRepo(t, dir)
	name := "with-two-taskfiles"
	pl := plan.NewPipeline(name, []string{"--directory", dir})
	repo := pl.AddResource(resources.Resource{
		Name:   "banana.git",
		Source: resources.Git{Uri: "https://example.org/flightplan.git"},
	})

	pl.AddJob(makeTestJobExternalTask(pl, repo, "job-1", "task-1", "tasks"))
	pl.AddJob(makeTestJobExternalTask(pl, repo, "job-2", "task-2", "tasks"))
	err := pl.Render()
	assert.NoError(t, err, "Render")

	testhelpers.AssertRenderedEqualsGolden(t, pl.Path(), "testdata/with-two-taskfiles.json", *update)
	testhelpers.AssertRenderedEqualsGolden(t, filepath.Join(dir, "tasks/task-1.json"),
		"testdata/task-1.json", *update)
	testhelpers.AssertRenderedEqualsGolden(t, filepath.Join(dir, "tasks/task-2.json"),
		"testdata/task-2.json", *update)
}

func TestPipelineReuseExternalTaskfileSuccess(t *testing.T) {
	dir := t.TempDir()
	testhelpers.MakeFakeGitRepo(t, dir)
	name := "with-two-taskfiles"
	pl := plan.NewPipeline(name, []string{"--directory", dir})
	repo := pl.AddResource(resources.Resource{
		Name:   "banana.git",
		Source: resources.Git{Uri: "https://example.org/flightplan.git"},
	})

	// External task file, with FileConfig.
	job1 := makeTestJobExternalTask(pl, repo, "job-1", "task-1", "tasks")
	taskFile1 := job1.Plan[1].(plan.Task).File
	pl.AddJob(job1)

	// External task file, nil FileConfig, same File path than job1: that is, it reuses
	// the task file of job 1.
	job2 := plan.Job{
		Name: "job-2",
		Plan: []plan.Step{plan.Task{Task: "task-2", File: taskFile1}},
	}
	pl.AddJob(job2)

	err := pl.Render()
	assert.NoError(t, err, "Render")

	testhelpers.AssertRenderedEqualsGolden(t, pl.Path(), "testdata/reuse-taskfiles.json", *update)
	testhelpers.AssertRenderedEqualsGolden(t, filepath.Join(dir, "tasks/task-1.json"),
		"testdata/task-1.json", *update)
	// There is no "tasks/task-2.json" !
}

func TestPipelineGetStepValidateFailure(t *testing.T) {
	test := func(desc string, resHandle *resources.Handle) {
		dir := t.TempDir()
		pipeline := plan.NewPipeline("get-failure", []string{"--directory", dir})
		pipeline.AddJob(plan.Job{
			Name: "knead-pizza",
			Plan: []plan.Step{plan.Get{Get: resHandle}},
		})
		err := pipeline.Render()

		assert.ErrorIs(t, err, plan.ErrGetValidation, desc)
	}

	test("nil handle", nil)
	test("missing name", &resources.Handle{})
	test("non-existing name", &resources.Handle{Resource: resources.Resource{Name: "banana"}})
}

func TestPipelinePutStepValidateFailure(t *testing.T) {
	test := func(desc string, resHandle *resources.Handle) {
		dir := t.TempDir()
		pipeline := plan.NewPipeline("get-failure", []string{"--directory", dir})
		pipeline.AddJob(plan.Job{
			Name: "knead-pizza",
			Plan: []plan.Step{plan.Put{Put: resHandle}},
		})
		err := pipeline.Render()

		assert.ErrorIs(t, err, plan.ErrPutValidation, desc)
	}

	test("nil handle", nil)
	test("missing name", &resources.Handle{})
	test("non-existing name", &resources.Handle{Resource: resources.Resource{Name: "banana"}})
}

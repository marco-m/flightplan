// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"flag"
	"path/filepath"
	"testing"

	plan "github.com/marco-m/flightplan"

	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/check"
	"github.com/marco-m/rosina/golden"
)

var update = flag.Bool("golden.update", false,
	"update the golden files for this package; use with -run to update a single test")

var job = plan.Job{
	Name: "bake-pizza",
	Plan: []plan.Step{
		plan.Task{
			Task: "knead",
			Config: &plan.TaskConfig{
				Platform: "linux",
				ImageResource: plan.AnonymousResource{
					Source: plan.RegistryImageSource{Repository: "alpine"},
				},
				Run: plan.TaskCommand{
					Path: "echo",
					Args: []string{"Pizza", "Margherita"},
				},
			},
		},
	},
}

func TestPipelinePath(t *testing.T) {
	test := func(desc, name string, args []string, wantSuffix string) {
		t.Helper()
		pipeline := plan.NewPipeline(name, args)
		check.Contains(t, filepath.ToSlash(pipeline.Path()), wantSuffix, "Path "+desc)
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
	pipeline.AddJob(job)

	err := pipeline.Render()

	assert.NoError(t, err, "Render")
	assertRenderedEqualsGolden(t, pipeline.Path(), "testdata/one-job.json", *update)
}

func TestPipelineJobWithoutNameIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("job-missing-name", nil)

	pipeline.AddJob(plan.Job{})
	err := pipeline.Render()

	assert.ErrorIs(t, err, plan.ErrEmptyJobName, "Render")
}

func TestPipelineCannotAddDuplicateJob(t *testing.T) {
	pipeline := plan.NewPipeline("duplicate-job", nil)

	pipeline.AddJob(job)
	pipeline.AddJob(job)

	err := pipeline.Render()

	assert.ErrorIs(t, err, plan.ErrDuplicateJob, "Render")
}

// test lookup via resource handle
// test lookup via job handle
// consider exporting Pipeline.Errors() and Resource(handle) and Job(handle)
// so that I can remove the pipeline private tests.

// test with override from commandline

//
// Helpers.
//

func assertRenderedEqualsGolden(t *testing.T, havePath, goldenPath string, update bool) {
	t.Helper()
	if diff := golden.DiffFiles(t, havePath, goldenPath, update); diff != "" {
		t.Errorf("Render: mismatch:\n%s", diff)
		// hack, make it possible to detect whitespace changes
		// haveBytes, err := os.ReadFile(havePath)
		// assert.NoError(t, err, "reading havePath")
		// goldenBytes, err := os.ReadFile(goldenPath)
		// assert.NoError(t, err, "reading goldenPath")
		// t.Errorf("have\n%s", hex.Dump(haveBytes))
		// t.Errorf("golden\n%s", hex.Dump(goldenBytes))
	}
}

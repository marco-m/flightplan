// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"flag"
	"path/filepath"
	"testing"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/flightplan/resources"

	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/check"
	"github.com/marco-m/rosina/golden"
)

var update = flag.Bool("golden.update", false,
	"update the golden files for this package; use with -run to update a single test")

var resource = resources.Resource{
	Name: "flightplan.git",
	Source: resources.Git{
		Uri:    "https://github.com/marco-m/flightplan.git",
		Branch: "master",
	},
}

func makeTestJob() plan.Job {
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
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"Pizza", "Margherita"},
					},
				},
			},
		},
	}
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
	pipeline.AddJob(makeTestJob())

	err := pipeline.Render()

	assert.NoError(t, err, "Render")
	assertRenderedEqualsGolden(t, pipeline.Path(), "testdata/one-job.json", *update)
}

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

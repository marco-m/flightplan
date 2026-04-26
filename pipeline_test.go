// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/flightplan/resources"

	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/golden"
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

func TestPipelinePath(t *testing.T) {
	test := func(desc, name string, args []string, wantSuffix string) {
		t.Helper()
		pipeline := plan.NewPipeline(name, args)
		assert.True(t, strings.HasSuffix(filepath.ToSlash(pipeline.Path()), wantSuffix), "Path "+desc)
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
	assertRenderedEqualsGolden(t, pipeline.Path(), "testdata/one-job.json", *update)
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
			plan.Get{Get: artifacts, Passed: []plan.JobHandle{kneadPizza}},
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
	assertRenderedEqualsGolden(t, pipeline.Path(), "testdata/get-and-put.json", *update)
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

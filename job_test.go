// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"testing"

	plan "github.com/marco-m/flightplan"

	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/check"
)

func TestLookupUnknownJob(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)
	assert.NoError(t, pipeline.Errors(), "Errors")

	_, found := pipeline.Job(plan.JobHandle("non-existing"))
	assert.False(t, found, "job found")
}

func TestAddJobsAndFindThemWithHandles(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	handle1 := pipeline.AddJob(plan.Job{
		Name: "job1",
	})
	handle2 := pipeline.AddJob(plan.Job{
		Name: "job2",
	})

	assert.NoError(t, pipeline.Errors(), "Errors")

	job1, found1 := pipeline.Job(handle1)
	assert.True(t, found1, "job1 found")
	check.Equal(t, job1.Name, "job1", "job1.Name")

	job2, found2 := pipeline.Job(handle2)
	assert.True(t, found2, "job2 found")
	check.Equal(t, job2.Name, "job2", "job2.Name")
}

func TestPipelineJobWithoutNameIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		// Name: "foo" <== missing
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrEmptyJobName, "Errors")
}

func TestJobTask(t *testing.T) {
	test := func(desc string, taskStep plan.Task, wantErr error) {
		t.Helper()
		pipeline := plan.NewPipeline("pizza", nil)
		pipeline.AddJob(plan.Job{Name: "banana", Plan: []plan.Step{taskStep}})
		assert.ErrorIs(t, pipeline.Errors(), wantErr, desc)
	}

	test("TaskMustHaveConfigOrFile",
		plan.Task{Task: "mango"},
		plan.ErrTaskNoConfigNoFile)
	test("TaskCannotHaveBothConfigAndFile",
		plan.Task{
			Task:   "mango",
			Config: &plan.TaskConfig{},
			File:   "banana",
		},
		plan.ErrTaskBothConfigAndFile)
}

func TestAddJobTaskConfigCannotHaveImageType(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	job := makeTestJobInlineTask()
	job.Plan[0].(plan.Task).Config.ImageResource.Type = "this-will-fail"

	pipeline.AddJob(job)

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrSetImageResourceType, "Errors")
	assert.ErrorContains(t, pipeline.Errors(), "this-will-fail", "Errors")
}

func TestAddJobTaskConfigMustHaveImageResource(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	job := makeTestJobInlineTask()
	job.Plan[0].(plan.Task).Config.ImageResource = nil

	pipeline.AddJob(job)

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrMissingImageResource, "Errors")
}

func TestAddJobTaskConfigMustHaveImageResourceSource(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	job := makeTestJobInlineTask()
	job.Plan[0].(plan.Task).Config.ImageResource.Source = nil

	pipeline.AddJob(job)

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrImageResourceSource, "Errors")
}

func TestAddJobTaskNeedsName(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		Name: "banana",
		Plan: []plan.Step{plan.Task{}},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrTaskNoName, "Errors")
}

func TestPipelineCannotAddDuplicateJob(t *testing.T) {
	pipeline := plan.NewPipeline("duplicate-job", nil)

	job := makeTestJobInlineTask()

	pipeline.AddJob(job)
	pipeline.AddJob(job)

	err := pipeline.Render()

	assert.ErrorIs(t, err, plan.ErrDuplicateJobName, "Render")
	assert.ErrorContains(t, err, job.Name, "job.Name")
}

// TODO
func TestJobMustHaveOrImageResourceOrTaskImageButNotBoth(t *testing.T) {
}

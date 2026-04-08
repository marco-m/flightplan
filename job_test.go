// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"testing"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/flightplan/resources"

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

func TestAddJobNeedsName(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		// Name: "foo" <== missing
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrEmptyJobName, "Errors")
}

func TestAddJobTaskMustHaveConfigOrFile(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		Name: "banana",
		Plan: []plan.Step{
			plan.Task{Task: "mango"},
		},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrTaskNoConfigNoFile, "Errors")
}

func TestAddJobTaskCannotHaveBothConfigAndFile(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		Name: "banana",
		Plan: []plan.Step{
			plan.Task{
				Task:   "mango",
				Config: &plan.TaskConfig{},
				File:   "banana",
			},
		},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrTaskBothConfigAndFile, "Errors")
}

func TestAddJobTaskConfigCannotHaveImageType(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		Name: "banana",
		Plan: []plan.Step{
			plan.Task{
				Task: "mango",
				Config: &plan.TaskConfig{
					Platform: "",
					ImageResource: resources.AnonymousResource{
						Type: "this-will-fail",
					},
				},
			},
		},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrSetImageType, "Errors")
	assert.ErrorContains(t, pipeline.Errors(), "this-will-fail", "Errors")
}

func TestAddJobTaskNeedsName(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddJob(plan.Job{
		Name: "banana",
		Plan: []plan.Step{
			plan.Task{},
		},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrTaskNoName, "Errors")
}

// func TestPipelineCannotAddDuplicateResource(t *testing.T) {
// 	pipeline := plan.NewPipeline("duplicate-resource", nil)

// 	resource := plan.Resource{
// 		Name:   "resource.git",
// 		Source: plan.GitSource{Uri: "https://github.com/marco-m/resource.git"},
// 	}

// 	pipeline.AddResource(resource)
// 	pipeline.AddResource(resource)

// 	err := pipeline.Render()

// 	assert.ErrorIs(t, err, plan.ErrDuplicateResource, "Render")
// }

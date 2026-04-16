// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// Example of building a two jobs pipeline with flightplan.
package main

import (
	"fmt"
	"os"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/flightplan/resources"
)

func main() {
	if err := buildPipeline(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func buildPipeline() error {
	pipeline := plan.NewPipeline("two-jobs", os.Args[1:])

	repo := pipeline.AddResource(resources.Resource{
		Name: "flightplan.git",
		// AddResource will set field Type using the method Type() of [plan.GitSource].
		Source: resources.Git{
			Uri:    "https://github.com/marco-m/flightplan.git",
			Branch: "master",
			Paths:  []string{"ci/*"},
		},
	})

	s3 := pipeline.AddResource(resources.Resource{
		Name: "artifacts.s3",
		// AddResource will set field Type using the method Type() of [plan.S3Source].
		Source: resources.S3{
			// FIXME
		},
	})

	golangImage := resources.AnonymousResource{
		Source: resources.RegistryImage{Repository: "golang"},
	}

	kneadPizza := pipeline.AddJob(plan.Job{
		Name: "knead-pizza",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: "prepare-dough",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &golangImage,
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"ciccio"},
					},
				},
			},
			plan.Task{
				Task: "let-dough-rise",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &golangImage,
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"bello"},
					},
				},
			},
			plan.Put{Resource: s3},
		},
	})

	pipeline.AddJob(plan.Job{
		Name: "bake-pizza",
		Plan: []plan.Step{
			plan.Get{
				Get:     repo,
				Passed:  []plan.JobHandle{kneadPizza},
				Trigger: true,
			},
			plan.Get{
				Get:    s3,
				Passed: []plan.JobHandle{kneadPizza},
			},
			plan.Task{
				Task: "put-in-hoven",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &golangImage,
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"hot", "hot", "hot"},
					},
				},
			},
		},
	})

	return pipeline.Render()
}

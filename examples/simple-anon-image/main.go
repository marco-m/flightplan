// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// There are two ways to pass an OCI image to a Concourse task:
// 1. With an anonymous resource passed directly Task.Config.ImageResource
// 2. With a named resource passed to Task.Image
//
// This file shows option 1.
package main

import (
	"fmt"
	"os"

	plan "github.com/marco-m/flightplan/pkg/flightplan"
	"github.com/marco-m/flightplan/pkg/resources"
)

func main() {
	if err := buildPipeline(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func buildPipeline() error {
	pipeline := plan.NewPipeline("simple-anon-image", os.Args[1:])

	repo := pipeline.AddResource(resources.Resource{
		Name: "flightplan.git",
		// AddResource will set field Type using the method Type() of [plan.GitSource].
		Source: resources.Git{
			Uri:    "https://github.com/marco-m/flightplan.git",
			Branch: "master",
		},
	})

	alpineImage := resources.AnonymousResource{
		Source: resources.RegistryImage{Repository: "alpine"},
	}

	pipeline.AddJob(plan.Job{
		Name: "knead-pizza",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: "prepare-dough",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &alpineImage,
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
					ImageResource: &alpineImage,
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"bello"},
					},
				},
			},
		},
	})

	return pipeline.Render()
}

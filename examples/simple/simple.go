// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// Example of building a simple pipeline with flightplan.
package main

import (
	"fmt"
	"os"

	plan "github.com/marco-m/flightplan"
)

func main() {
	if err := buildPipeline(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func buildPipeline() error {
	pipeline := plan.NewPipeline("simple", os.Args[1:])

	golangImage := plan.AnonymousResource{
		Source: plan.RegistryImageSource{Repository: "golang"},
	}

	pipeline.AddJob(plan.Job{
		Name: "knead-pizza",
		Plan: []plan.Step{
			plan.Task{
				Task: "prepare-dough",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: golangImage,
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
					ImageResource: golangImage,
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

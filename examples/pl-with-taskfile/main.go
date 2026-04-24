// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// Example of a pipeline with one job and creation of external task file with flightplan.
package main

import (
	"fmt"
	"os"
	"path"

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
	pipeline := plan.NewPipeline("pl-with-task-file", os.Args[1:])

	repo := pipeline.AddResource(resources.Resource{
		Name: "flightplan.git",
		Source: resources.Git{
			Uri:    "https://github.com/marco-m/flightplan.git",
			Branch: "devel",
		},
	})

	alpineImage := resources.AnonymousResource{
		Source: resources.RegistryImage{Repository: "alpine"},
	}

	tasksDir := path.Join(pipeline.RelDir(repo), "tasks")
	pipeline.AddJob(plan.Job{
		Name: "mango",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: "make-gift",
				// flighplan will create this task file with the contents of FileConfig.
				File: path.Join(tasksDir, "make-gift.json"),
				FileConfig: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &alpineImage,
					Inputs:        []plan.TaskInput{{Name: repo.Name}},
					Outputs:       []plan.TaskOutput{{Name: "gift"}},
					Run: plan.TaskCommand{
						Path: path.Join(tasksDir, "make-gift.sh"),
					},
				},
			},
		},
	})

	return pipeline.Render()
}

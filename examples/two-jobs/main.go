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
		},
	})

	s3 := pipeline.AddResource(resources.Resource{
		Name: "artifacts.s3",
		// AddResource will set field Type using the method Type() of [plan.S3Source].
		Source: resources.S3{
			Bucket: "concourse",
			// convention: builds/<pipeline-name>/<versioned-package-name>
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

	kneadPizza := pipeline.AddJob(plan.Job{
		Name: "knead-pizza",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: "prepare-dough",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &alpineImage,
					Outputs: []plan.TaskOutput{
						{Name: "gift"},
					},
					Run: plan.TaskCommand{
						Path: "sh",
						Args: []string{
							"-c",
							`set -e
set -x
VERSION=$(date +%Y%m%d%H%M%S)
GIFT=gift/gift-$VERSION
echo "hello" > $GIFT
`,
						},
					},
				},
			},
			plan.Put{
				Put: s3,
				Params: resources.S3PutParams{
					File: "gift/gift-*",
				},
			},
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
				Task: "place-in-hoven",
				Config: &plan.TaskConfig{
					Platform:      "linux",
					ImageResource: &alpineImage,
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

package main

import (
	"fmt"
	"io"
	"os"

	plan "github.com/marco-m/flightplan/flightplan"
)

func main() {
	if err := mainErr(os.Stdout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainErr(out io.Writer) error {
	pipeline := plan.NewPipeline()

	repo := pipeline.AddResource(plan.Resource{
		Name: "repo.git",
		Type: "git",
		Source: plan.GitSource{
			Uri:    "https://github.com/marco-m/flightplan.git",
			Branch: "master",
			Paths:  []string{"ci/*"},
		},
	})

	golangImage := plan.AnonymousResource{
		Type:   "registry-image",
		Source: plan.RegistryImageSource{Repository: "golang"},
	}

	pipeline.AddJob(plan.Job{
		Name: "knead-pizza",
		Plan: []plan.Step{
			plan.Get{Get: repo, Trigger: true},
			plan.Task{
				Task: "prepare-dough",
				Config: plan.TaskConfig{
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
				Config: plan.TaskConfig{
					Platform:      "linux",
					ImageResource: golangImage,
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"bello"},
					},
				},
			},
			// fp.Put{Resource: s3},
		},
	})

	// _, err = pipeline.AddJob("bake-pizza", fp.AddJobArgs{
	// 	Steps: []fp.Step{
	// 		fp.Get{Resource: repo, Passed: []*fp.Job{kneadPizzaJob}},
	// 		fp.Get{Resource: s3},
	// 		fp.Task{},
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }

	return pipeline.Render(out)
}

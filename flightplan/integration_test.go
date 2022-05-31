package flightplan_test

import (
	"errors"
	"io"
	"os"
	"testing"

	plan "github.com/marco-m/flightplan/flightplan"
	"github.com/marco-m/flightplan/internal/testhelpers/assert"
)

var update = os.Getenv("GOLDEN_UPDATE") == "true"

func TestEmptyPipelineFails(t *testing.T) {
	pipeline := plan.NewPipeline()

	err := pipeline.Render(io.Discard)

	want := errors.New("FIXME")
	assert.ErrorIs(t, err, want, "Render")
}

func TestMinimalPipelineWithResourceOnlyIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline()

	pipeline.AddResource(plan.Resource{
		Name: "repo.git",
		Type: "git",
		Source: plan.GitSource{
			Uri:    "https://github.com/marco-m/flightplan.git",
			Branch: "master",
		},
	})

	err := pipeline.Render(io.Discard)

	want := errors.New("pipeline must contain at least one job")
	assert.ErrorIs(t, err, want, "Render")
}

func TestMinimalPipelineWithJobOnlyIsvalid(t *testing.T) {
	pipeline := plan.NewPipeline()

	pipeline.AddJob(plan.Job{
		Name: "bake-pizza",
		Plan: []plan.Step{
			plan.Task{
				Task: "knead",
				Config: plan.TaskConfig{
					Platform: "linux",
					ImageResource: plan.AnonymousResource{
						Type:   "registry-image",
						Source: plan.RegistryImageSource{Repository: "alpine"},
					},
					Run: plan.TaskCommand{
						Path: "echo",
						Args: []string{"Pizza Margherita"},
					},
				},
			},
		},
	})

	assert.RenderedEqualsGolden(t, pipeline, "testdata/minimal-job.json", update)
}

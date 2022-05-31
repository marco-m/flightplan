package flightplan

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/marco-m/flightplan/internal/goof"
)

type Pipeline struct {
	Resources []Resource `json:"resources,omitempty"`
	Jobs      []Job      `json:"jobs,omitempty"`
	errs      []error
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (pl *Pipeline) AddResource(res Resource) ResourceHandle {
	if res.Name == "" {
		pl.errs = append(pl.errs, goof.NewErr("AddResource: res.Name cannot be empty"))
		return ""
	}
	for i, r := range pl.Resources {
		if r.Name == res.Name {
			pl.errs = append(pl.errs,
				goof.NewErr("AddResource: a resource with name %q already exists at index %d",
					r.Name, i))
			return ""
		}
	}
	pl.Resources = append(pl.Resources, res)
	return ResourceHandle(res.Name)
}

func (pl *Pipeline) AddJob(job Job) JobHandle {
	if job.Name == "" {
		pl.errs = append(pl.errs, goof.NewErr("AddJob: job.Name cannot be empty"))
		return ""
	}
	for i, j := range pl.Jobs {
		if j.Name == job.Name {
			err := goof.NewErr("AddJob: a job with name %q already exists at index %d", j.Name, i)
			pl.errs = append(pl.errs, err)
		}
	}
	pl.Jobs = append(pl.Jobs, job)
	return JobHandle(job.Name)
}

func (pl *Pipeline) Render(out io.Writer) error {
	err := errors.Join(pl.errs...)
	if err != nil {
		return err
	}
	// for _, res := range pl.resources {
	// 	fmt.Fprintf(out, "%v\n", res.Name())
	// }
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(pl); err != nil {
		return err
	}
	return nil
}

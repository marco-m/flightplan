// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"fmt"

	"github.com/marco-m/flightplan/resources"
)

// Step is a step in a [Job.Plan]. You compose the plan using the following concrete
// steps: [Task], [Get], [Put], [SetPipeline], [InParallel], [Do], [Try], [Loadvar].
type Step interface {
	// Confirm that the struct is actually a [Step] (sealed interface).
	step()
	// Validate verifies that a step implementation is valid, by being aware of the
	// step type and of the referenced resources.
	// For example: a put step to a resource is aware of the resource type (s3, git, ...)
	// and validates the resource-specific "params" object.
	Validate(pl *Pipeline) error
}

// Get implements [Step].
// Fetches a version of a resource.
// See https://concourse-ci.org/docs/steps/get/
type Get struct {
	// Required. Get is the resource to get from.
	Get *resources.Handle `json:"get"`

	// Optional

	// Resource ResourceHandle `json:"resource,omitzero"` // I don't understand this one
	Passed  []JobHandle         `json:"passed,omitzero"`
	Trigger bool                `json:"trigger,omitzero"`
	Params  resources.GetParams `json:"params,omitzero"`
	Version string              `json:"version,omitzero"`
}

func (Get) step() {}

func (get Get) Validate(pl *Pipeline) error {
	if get.Get == nil {
		return fmt.Errorf("%w: field Get cannot be empty", ErrGetValidation)
	}
	for _, res := range pl.po.Resources {
		if res.Name == get.Get.Name {
			return nil
		}
	}
	return fmt.Errorf("%w: field Get: unknown resource %q", ErrGetValidation, get.Get.Name)
}

// Put implements [Step].
// Pushes to the given resource.
// See https://concourse-ci.org/docs/steps/put/
type Put struct {
	// Required. Put is the resource to put to.
	Put *resources.Handle `json:"put"`
	// Params represents an arbitrary configuration to pass to the resource.
	// Refer to the resource type's documentation to see what it supports.
	// Params map[string]any `json:"params,omitzero"`
	Params resources.PutParams `json:"params,omitzero"`
}

func (Put) step() {}

func (put Put) Validate(pl *Pipeline) error {
	if put.Put == nil {
		return fmt.Errorf("%w: field Put cannot be empty", ErrPutValidation)
	}
	found := false
	for _, res := range pl.po.Resources {
		if res.Name == put.Put.Name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%w: field Put: unknown resource %q",
			ErrPutValidation, put.Put.Name)
	}
	if err := put.Put.Source.ValidatePut(put.Params); err != nil {
		return fmt.Errorf("Put %s: %w", put.Put.Name, err)
	}
	return nil
}

// Task implements [Step].
// See https://concourse-ci.org/docs/steps/task/
type Task struct {
	// Required. The name of the task. Shown in the UI.
	Task string `json:"task,omitzero"`
	// The [TaskConfig] to execute, inline in the pipeline.
	// A [Task] must contain a [Task.Config] or a [Task.File] but not both.
	Config *TaskConfig `json:"config,omitzero"`
	// Path to a YAML file containing a [TaskConfig].
	// The first segment in the path should refer to another source from the [Job.Plan],
	// and the rest of the path is relative to that source.
	// A [Task] must contain a [Task.File] or a [Task.Config] but not both.
	File string `json:"file,omitzero"`
}

func (Task) step() {}

func (task Task) Validate(pl *Pipeline) error {
	if task.Task == "" {
		return ErrTaskNoName
	}
	if (task.Config != nil) && task.File != "" {
		return ErrTaskBothConfigAndFile
	}
	if task.Config != nil {
		if task.Config.ImageResource == nil {
			return ErrImageResource
		}
		imgRes := task.Config.ImageResource
		if imgRes.Type != "" {
			return fmt.Errorf("%w: %s", ErrSetImageResourceType, imgRes.Type)
		}
		if task.Config.ImageResource.Source == nil {
			return ErrImageResourceSource
		}
		task.Config.ImageResource.Type = task.Config.ImageResource.Source.Type()
		return nil
	}
	if task.File != "" {
		// TODO
		return nil
	}
	return ErrTaskNoConfigNoFile
}

type TaskConfig struct {
	// Required
	Platform string `json:"platform,omitzero"`

	// Optional. The container image to run with. Prefer instead [Task.Image].
	ImageResource *resources.AnonymousResource `json:"image_resource,omitzero"`
	// Optional.
	Inputs []TaskInput `json:"inputs,omitzero,omitempty"`
	// Optional.
	Outputs []TaskOutput `json:"outputs,omitzero,omitempty"`
	// Optional.
	Run TaskCommand `json:"run,omitzero"`
	// Optional.
	RootfsUri string `json:"rootfs_uri,omitzero"`
}

type TaskInput struct{}

// https://concourse-ci.org/docs/tasks/#output-schema
type TaskOutput struct {
	Name string `json:"name"`
	Path string `json:"path,omitzero"`
}

type TaskCommand struct {
	Path string   `json:"path,omitzero"`
	Args []string `json:"args,omitzero"`
	Dir  string   `json:"dir,omitzero"`
	User string   `json:"user,omitzero"`
}

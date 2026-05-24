// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"fmt"

	"github.com/marco-m/flightplan/pkg/resources"
)

// Step is a step in a [Job.Plan]. You compose the plan using the following concrete
// steps: [Task], [Get], [Put], [SetPipeline], [InParallel], [Do], [Try], [Loadvar].
type Step interface {
	// Confirm that the struct is actually a [Step] (sealed interface).
	isStep()
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

func (Get) isStep() {}

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

func (Put) isStep() {}

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

// See https://concourse-ci.org/docs/steps/task/#container_limits-schema
type ContainerLimits struct {
	// The maximum amount of CPU available to the task container, measured in shares.
	// 0 means unlimited.
	CPU int `json:"cpu"`
	// The maximum amount of memory available to the task container, measured in bytes.
	// 0 means unlimited.
	Memory int `json:"memory"`
}

// A Task step executes a [TaskConfig].
// When a task completes, the artifacts specified by [TaskConfig.Outputs] will be
// registered in the build's artifact namespace. This allows subsequent [Task] steps
// and [Put] steps to access the result of a task.
// See https://concourse-ci.org/docs/steps/task/
type Task struct {
	// Required. The name of the task. Shown in the UI.
	Task string `json:"task,omitzero"`
	// Config is the [TaskConfig] to execute, inline in the pipeline.
	// A [Task] step must contain a [Task.Config] or a [Task.File] but not both.
	Config *TaskConfig `json:"config,omitzero"`
	// File is the path to a YAML file containing the [TaskConfig] to execute.
	// Flightplan extension: if [Task.FileConfig] is not empty, then [Task.File] will be
	// created by flightplan, containing [Task.FileConfig].
	// The first segment in the path should refer to a previous get in the [Job.Plan],
	// A [Task] step must contain a [Task.File] or a [Task.Config] but not both.
	File string `json:"file,omitzero"`
	// Flightplan extension. Equivalent of [Task.Config], will be rendered into
	// [Task.File].
	FileConfig *TaskConfig `json:"-"`
	// Image specifies an artifact source containing an image to use for the task.
	// This overrides any [TaskConfig.ImageResource] present in the task configuration.
	// This is very useful when part of your pipeline involves building an image,
	// possibly with dependencies pre-baked. You can then propagate that image through
	// the rest of your pipeline, guaranteeing that the correct version (and thus a
	// consistent set of dependencies) is used throughout your pipeline.
	Image *resources.Handle `json:"image,omitzero"`
	// Default false. If set to true, the task will run with escalated capabilities
	// available on the task's platform.
	// WARNING Setting privileged: true is a gaping security hole; use wisely
	// and only if necessary.
	Privileged bool `json:"privileged,omitzero"`
	// A map of template variables to pass to an external task. Not to be confused with
	// task step params, which provides environment variables to the task.
	// This is to be used with external tasks defined in task step file.
	Vars map[string]any `json:"vars,omitzero"`
	// A map of task environment variable parameters to set, overriding those configured
	// in the task's config or file.
	Params map[string]any `json:"params,omitzero"`
	// CPU and memory limits to enforce on the task container.
	// These values will override any limits set for concourse web.
	// These values will also override any configuration set on a task's config container_limits
	ContainerLimits ContainerLimits `json:"container_limits,omitzero"`
	// Default false. If set to true, the task will have no outbound network access.
	// WARNING: Works only for Linux IF the container runtime of the worker is containerd.
	// If has no effect for: Linux with a different container runtime, macOS and Windows.
	Hermetic bool `json:"hermetic,omitzero"`
	// A map from task input names to concrete names in the build plan.
	// This allows a task with generic input names to be used multiple times in the same
	// plan, mapping its inputs to specific resources within the plan.
	InputMapping map[string]string `json:"input_mapping,omitzero"`
	// A map from task output names to concrete names to register in the build plan.
	// This allows a task with generic output names to be used multiple times in the same plan.
	OutputMapping map[string]string `json:"output_mapping,omitzero"`
}

func (Task) isStep() {}

func (task Task) Validate(pl *Pipeline) error {
	if task.Task == "" {
		return ErrTaskNoName
	}
	if (task.Config != nil) && task.File != "" {
		return ErrTaskBothConfigAndFile
	}
	if task.Config != nil {
		anonImgRes := task.Config.ImageResource
		if anonImgRes == nil {
			if task.Image == nil {
				return ErrMissingImageResource
			}
			return nil
		}
		if anonImgRes.Type != "" {
			return fmt.Errorf("%w: %s", ErrSetImageResourceType, anonImgRes.Type)
		}
		if anonImgRes.Source == nil {
			return ErrImageResourceSource
		}
		anonImgRes.Type = anonImgRes.Source.Type()
		return nil
	}
	if task.File != "" {
		// TODO
		return nil
	}
	return ErrTaskNoConfigNoFile
}

type externalTask struct {
	File       string
	FileConfig *TaskConfig
}

func (task Task) Process(extTasks []*externalTask) (*externalTask, error) {
	if task.FileConfig == nil {
		return nil, nil
	}
	for _, et := range extTasks {
		if et.File == task.File {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateExtTaskFile, task.File)
		}
	}
	return &externalTask{
		File:       task.File,
		FileConfig: task.FileConfig,
	}, nil
}

// A TaskConfig represents a Concourse task, the smallest configurable unit in a pipeline.
// A task can be thought of as a function from [TaskConfig.Inputs] to [TaskConfig.Outputs]
// that can either succeed or fail.
// See https://concourse-ci.org/docs/tasks/
type TaskConfig struct {
	// Required
	Platform string `json:"platform,omitzero"`
	// Optional. The container image to run with. Overridden by [Task.Image].
	// Prefer instead [Task.Image].
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

// The set of artifacts used by task, determining which artifacts will be available in the
// current directory when the task runs. These are satisfied by get steps or
// task-config.outputs of a previous task. These can also be provided by -i with fly execute.
// If any required inputs are missing at run-time, then the task will error immediately.
type TaskInput struct {
	// Required. The name of the input.
	Name string `json:"name"`
	// Optional. The path where the input will be placed. If not specified, the input's
	// name is used.
	Path string `json:"path,omitzero"`
	// Optional. Default false. If true, then the input is not required by the task.
	// The task may run even if this input is missing.
	Optional bool `json:"optional,omitzero"`
}

// The artifacts produced by the task.
// Each output configures a directory to make available to later steps in the build plan.
// The directory will be automatically created before the task runs, and the task should
// place any artifacts it wants to export in the directory.
// See https://concourse-ci.org/docs/tasks/#output-schema
type TaskOutput struct {
	// Required. The name of the output. The contents under path will be made available
	// to the rest of the plan under this name.
	Name string `json:"name"`
	// Optional. The path to a directory where the output will be taken from. If not
	// specified, the output's name is used.
	Path string `json:"path,omitzero"`
}

type TaskCommand struct {
	Path string   `json:"path,omitzero"`
	Args []string `json:"args,omitzero"`
	Dir  string   `json:"dir,omitzero"`
	User string   `json:"user,omitzero"`
}

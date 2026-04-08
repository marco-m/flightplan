// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import "github.com/marco-m/flightplan/resources"

// Step is a step in a [Job.Plan]. You compose the plan using the following concrete
// steps: [Task], [Get], [Put], [SetPipeline], [InParallel], [Do], [Try], [Loadvar].
type Step interface {
	// Confirm that the struct is actually a [Step] (sealed interface).
	step()
}

// Fetches a version of a resource.
// See https://concourse-ci.org/docs/steps/get/
type Get struct {
	// Required

	Get resources.ResourceHandle `json:"get,omitzero"`

	// Optional

	// Resource ResourceHandle `json:"resource,omitzero"` // I don't understand this one
	Passed  []JobHandle       `json:"passed,omitzero"`
	Trigger bool              `json:"trigger,omitzero"`
	Params  map[string]string `json:"params,omitzero"`
	Version string            `json:"version,omitzero"`
}

func (Get) step() {}

type Put struct {
	Resource resources.ResourceHandle
}

func (Put) step() {}

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

type TaskConfig struct {
	// Required
	Platform string `json:"platform,omitzero"`

	// Optional. The container image to run with. Prefer instead [Task.Image].
	ImageResource resources.AnonymousResource `json:"image_resource,omitzero"`
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

type TaskOutput struct{}

type TaskCommand struct {
	Path string   `json:"path,omitzero"`
	Args []string `json:"args,omitzero"`
	Dir  string   `json:"dir,omitzero"`
	User string   `json:"user,omitzero"`
}

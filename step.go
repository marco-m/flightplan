// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

type Step interface {
	// Confirm that the struct is actually a [Step] (sealed interface).
	step()
}

type Task struct {
	// Required
	Task string `json:"task,omitzero"`

	// Optional
	Config TaskConfig `json:"config,omitzero"`
}

func (Task) step() {}

type TaskConfig struct {
	// Required
	Platform string `json:"platform,omitzero"`

	// Optional. The container image to run with. Prefer instead [Task.Image].
	ImageResource AnonymousResource `json:"image_resource,omitzero"`
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

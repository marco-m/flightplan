package flightplan

type Step interface{}

// Fetches a version of a resource.
// See https://concourse-ci.org/docs/steps/get/
type Get struct {
	// Required

	Get ResourceHandle `json:"get,omitzero"`

	// Optional

	// Resource ResourceHandle `json:"resource,omitzero"` // I don't understand this one
	Passed  []JobHandle       `json:"passed,omitzero"`
	Trigger bool              `json:"trigger,omitzero"`
	Params  map[string]string `json:"params,omitzero"`
	Version string            `json:"version,omitzero"`
}

type Put struct {
	Resource ResourceHandle
}

type Task struct {
	// Required

	Task string `json:"task,omitzero"`

	// Optional

	Config TaskConfig `json:"config,omitzero"`
	// If true, will generate an embedded task-config. If false (the default), will
	// generate a separate task file. A separate task file is preferred because it
	// enables to use fly execute.
	Embedded bool `json:"-"`
	// A container image to run, to be fetched with [Get]. Prefer this instead of the
	// anonymous [TaskConfig.ImageResource].
	Image RegistryImage `json:"image,omitzero"`
}

type TaskConfig struct {
	// Required

	Platform string `json:"platform,omitzero"`

	// Optional

	// The container image to run with. Prefer instead [Task.Image].
	ImageResource AnonymousResource `json:"image_resource,omitzero"`
	Inputs        []TaskInput       `json:"inputs,omitzero,omitempty"`
	Outputs       []TaskOutput      `json:"outputs,omitzero,omitempty"`
	// caches [cache]
	// params env-vars
	Run       TaskCommand `json:"run,omitzero"`
	RootfsUri string      `json:"rootfs_uri,omitzero"`
	// container_limits container_limits
}

type TaskInput struct{}

type TaskOutput struct{}

type TaskCommand struct {
	Path string   `json:"path,omitzero"`
	Args []string `json:"args,omitzero"`
	Dir  string   `json:"dir,omitzero"`
	User string   `json:"user,omitzero"`
}

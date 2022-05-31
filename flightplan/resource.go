package flightplan

type Resource struct {
	// Required

	Name   string `json:"name,omitzero"`
	Type   string `json:"type,omitzero"`
	Source Source `json:"source,omitzero,omitempty"` // Only field whose shape depends on the resource type.

	// Optional

	OldName              string   `json:"old_name,omitzero"`
	Icon                 string   `json:"icon,omitzero"`
	Version              string   `json:"version,omitzero"`
	CheckEvery           string   `json:"check_every,omitzero"`
	CheckTimeout         string   `json:"check_timeout,omitzero"`
	ExposeBuildCreatedBy bool     `json:"expose_build_created_by,omitzero"`
	Tags                 []string `json:"tags,omitzero,omitempty"`
	Public               bool     `json:"public,omitzero"`
	WebhookToken         string   `json:"webhook_token,omitzero"`
}

type AnonymousResource struct {
	// Required

	// Type is usually "registry-image", see [RegistryImageSource]
	Type   string `json:"type,omitzero"`
	Source Source `json:"source,omitempty"`

	// Optional

	Params  Source            `json:"params,omitempty"`
	Version map[string]string `json:"version,omitempty"`
}

// The source object in a Concourse [Resource] or [AnonymousResource], represented as a Go struct.
type Source interface {
	// Confirm that the struct is actually a Source (sort of sealed interface)
	Source()
	// The resource type.
	Type() string
}

// A resource name must be unique per pipeline, otherwhise it could not be resolved
// unambiguously as a get or put step. Thus, we can use its name as handle: returned
// by [Pipeline.AddResource] and required by [Pipeline.AddJob].
type ResourceHandle string

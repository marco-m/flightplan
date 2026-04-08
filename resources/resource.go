// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

type Resource struct {
	// Required

	// Name is the resource name. You must set it.
	Name string `json:"name,omitzero"`
	// Type is the resource type. Leave it alone; will be set by [Source.Type].
	Type string `json:"type,omitzero"`
	// The contents of Source are specific to the resource type.
	Source Source `json:"source,omitzero,omitempty"`

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
	// Required. Type is usually "registry-image", see [RegistryImageSource]
	Type string `json:"type,omitzero"`
	// Required.
	Source Source `json:"source,omitempty"`

	// Optional.
	Params Params `json:"params,omitempty"`
	// Optional.
	Version map[string]string `json:"version,omitempty"`
}

// Source is the "source" object in a Concourse [Resource] or [AnonymousResource].
type Source interface {
	// Confirm that the struct is actually a [Source].
	// Exported to allow custom resources to be defined outside of the flightplan module.
	Source()
	// The mandatory resource type, used by [Pipeline.AddResource] to set field
	// [Resource.Type] of the outer resorce. Setting [Resource.Type] directly will be
	// considered an error.
	Type() string
}

// Params is the "params" object in a Concourse [Resource] or [AnonymousResource].
type Params interface {
	// Confirm that the struct is actually a [Params].
	// Exported to allow custom resources to be defined outside of the flightplan module.
	Params()
}

// A resource name must be unique per pipeline, otherwhise it could not be resolved
// unambiguously as a get or put step. Thus, we can use its name as handle: returned
// by [Pipeline.AddResource] and required by [Pipeline.AddJob].
type ResourceHandle string

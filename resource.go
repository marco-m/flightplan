// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

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

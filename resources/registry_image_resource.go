// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

// See https://github.com/concourse/registry-image-resource
type RegistryImage struct {
	// Required. The URI of the image repository, e.g. alpine or ghcr.io/package/image.
	// Defaults to checking docker.io if no hostname is provided in the URI.
	Repository string `json:"repository"`

	// Optional
}

var _ Source = (*RegistryImage)(nil)

func (ris RegistryImage) Source() {}

func (ris RegistryImage) Type() string { return "registry-image" }

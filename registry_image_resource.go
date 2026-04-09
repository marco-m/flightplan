// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

type RegistryImage struct{}

//
//
//

// See https://github.com/concourse/registry-image-resource
type RegistryImageSource struct {
	// Required. The URI of the image repository, e.g. alpine or ghcr.io/package/image.
	// Defaults to checking docker.io if no hostname is provided in the URI.
	Repository string `json:"repository"`

	// Optional
}

var _ Source = (*RegistryImageSource)(nil)

func (ris RegistryImageSource) Source() {}

func (ris RegistryImageSource) Type() string { return "registry-image" }

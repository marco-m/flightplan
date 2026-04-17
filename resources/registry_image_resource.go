// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

import "fmt"

// See https://github.com/concourse/registry-image-resource
type RegistryImage struct {
	// Required. The URI of the image repository, e.g. alpine or ghcr.io/package/image.
	// Defaults to checking docker.io if no hostname is provided in the URI.
	Repository string `json:"repository"`

	// Optional
}

var _ Source = (*RegistryImage)(nil)

func (rim RegistryImage) Source() {}

func (rim RegistryImage) Type() string { return "registry-image" }

func (rim RegistryImage) Validate() error {
	if rim.Repository == "" {
		return fmt.Errorf("RegistryImage.Repository cannot be empty")
	}
	return nil
}

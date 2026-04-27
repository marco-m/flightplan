// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

import (
	"errors"
	"fmt"
)

var (
	ErrRegistryImageRepositoryEmpty       = errors.New("RegistryImage.Repository cannot be empty")
	ErrRegistryImageParamsWrongType       = errors.New("Params: wrong type")
	ErrRegistryImagePutParamsEmpty        = errors.New("Put field Params cannot be empty")
	ErrRegistryImagePutParamsMissingImage = errors.New("Params: field Image cannot be empty")
)

var _ Source = (*RegistryImage)(nil)

// Supports checking, fetching, and pushing of images to OCI (Docker) registries.
// See https://github.com/concourse/registry-image-resource
type RegistryImage struct {
	// Required. The URI of the image repository, e.g. alpine or ghcr.io/package/image.
	// Defaults to checking docker.io if no hostname is provided in the URI.
	Repository string `json:"repository"`

	// Optional
}

func (rim RegistryImage) IsSource() {}

func (rim RegistryImage) Type() string { return "registry-image" }

func (rim RegistryImage) Validate() error {
	if rim.Repository == "" {
		return ErrRegistryImageRepositoryEmpty
	}
	return nil
}

func (rim RegistryImage) ValidateGet(params GetParams) error {
	if params == nil {
		// All RegistryImage get params are optional.
		return nil
	}
	if p, ok := params.(RegistryImageGetParams); ok {
		return p.Validate()
	}
	return fmt.Errorf("%w: have: %T; want: %T",
		ErrRegistryImageParamsWrongType, params, RegistryImageGetParams{})
}

func (rim RegistryImage) ValidatePut(params PutParams) error {
	if params == nil {
		// Some RegistryImage put params are required.
		return ErrRegistryImagePutParamsEmpty
	}
	if p, ok := params.(RegistryImagePutParams); ok {
		return p.Validate()
	}
	return fmt.Errorf("%w: have: %T; want: %T",
		ErrRegistryImageParamsWrongType, params, RegistryImagePutParams{})
}

// See https://github.com/concourse/registry-image-resource#get-step-params
type RegistryImageGetParams struct {
	// Optional. Default: rootfs. The format of the image to fetch.
	Format string `json:"format,omitzero"`
	// Optional.
	Platform string `json:"platform,omitzero"`
	// Optional. Default: false Skip downloading the image. Useful if you want to
	// trigger a job without using the object or when running after a put step and
	// not needing to download the image you just uploaded.
	SkipDownload bool `json:"skip_download,omitzero"`
}

func (RegistryImageGetParams) IsGetParams() {}

func (rimg RegistryImageGetParams) Validate() error {
	return nil
}

// See https://github.com/concourse/registry-image-resource#put-steps-params
type RegistryImagePutParams struct {
	// Required. Can be the path to the docker image tarball (e.g. my-image/image.tar) or
	// the path to the oci image tarball (e.g. my-image/image) or the path to an OCI image
	// layout (e.g. my-image/oci). Expanded with filepath.Glob.
	Image string `json:"image"`
	// Optional. A version number to use as a tag.
	Version string `json:"version,omitzero"`
	// Optional. When set to true and version is specified, automatically bump alias
	// tags for the version.
	BumpAliases bool `json:"bump_aliases,omitzero"`
	// Optional. The path to a file with whitespace-separated list of tag values to tag
	// the image with..
	AdditionalTags string `json:"additional_tags,omitzero"`
	// Optional. A string that will be prefixed to the tags from additional_tags.
	TagPrefix string `json:"tag_prefix,omitzero"`
}

func (RegistryImagePutParams) IsPutParams() {}

func (rimp RegistryImagePutParams) Validate() error {
	if rimp.Image == "" {
		return ErrRegistryImagePutParamsMissingImage
	}
	return nil
}

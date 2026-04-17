// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

import (
	"errors"
)

// S3 is a resource for interacting with S3-compatible object storage.
// See https://github.com/concourse/s3-resource
type S3 struct {
	// Required. Buket is the name of the bucket.
	Bucket string `json:"bucket"`
	// Regexp is a slash-delimited sequence of patterns to match against the
	// sub-directories and filenames of the objects stored within the S3 bucket.
	// Exactly one among Regexp and VersionedFile must be provided.
	Regexp string `json:"regexp,omitzero"`
	// 	VersionedFile is the path to the object in the S3 bucket. Requires S3 versioning
	// to be enabled.
	// Exactly one among Regexp and VersionedFile must be provided.
	VersionedFile string `json:"versioned_file,omitzero"`

	// Optional. Endpoint is the URL or plain hostname to a S3-compatible service.
	// Endpoint string ((s3-endpoint))

	// region_name: ((s3-region))
	// access_key_id: ((s3-access-key))
	// secret_access_key: ((s3-secret-key))
	// # Needed for Minio; not needed for S3.
	// use_path_style: true
}

var _ Source = (*S3)(nil)

func (s3 S3) Source() {}

func (s3 S3) Type() string { return "s3" }

var (
	ErrS3MissingBucket              = errors.New("field Bucket cannot be empty")
	ErrS3NoRegexpNoVersionedFile    = errors.New("one of the following fields must be filled: Regexp, VersionedFile")
	ErrS3BothRegexpAndVersionedFile = errors.New("one of the fields Regexp or VersionedFile must be filled, not both")
)

func (s3 S3) Validate() error {
	var errs []error
	if s3.Bucket == "" {
		errs = append(errs, ErrS3MissingBucket)
	}
	if s3.Regexp == "" && s3.VersionedFile == "" {
		errs = append(errs, ErrS3NoRegexpNoVersionedFile)
	}
	if s3.Regexp != "" && s3.VersionedFile != "" {
		errs = append(errs, ErrS3BothRegexpAndVersionedFile)
	}
	return errors.Join(errs...)
}

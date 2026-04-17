// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

import (
	"errors"
	"fmt"
)

var (
	ErrS3MissingBucket              = errors.New("S3: field Bucket cannot be empty")
	ErrS3NoRegexpNoVersionedFile    = errors.New("S3: one of the following fields must be filled: Regexp, VersionedFile")
	ErrS3BothRegexpAndVersionedFile = errors.New("S3: one of the fields Regexp or VersionedFile must be filled, not both")
	//
	ErrS3ParamsWrongType      = errors.New("Params: wrong type")
	ErrS3PutParamsEmpty       = errors.New("Put field Params cannot be empty")
	ErrS3PutParamsMissingFile = errors.New("Params: field File cannot be empty")
)

var _ Source = (*S3)(nil)

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
	Endpoint string `json:"endpoint,omitzero"`
	// Optional. RegionName is the region the bucket is in. Defaults to us-east-1.
	RegionName string `json:"region_name,omitzero"`
	// Optional. AccessKeyID is the S3 access key.
	AccessKeyID string `json:"access_key_id,omitzero"`
	// Optional. SecretAccessKey is the S3 secret key.
	SecretAccessKey string `json:"secret_access_key,omitzero"`
	// Optional. UsePathStyle enables legacy path-style access for S3 compatible providers.
	// The default is virtual path-style.
	UsePathStyle bool `json:"use_path_style,omitzero"`
}

func (s3 S3) IsSource() {}

func (s3 S3) Type() string { return "s3" }

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

func (s3 S3) ValidateGet(params GetParams) error {
	if params == nil {
		return nil
	}
	if p, ok := params.(S3GetParams); ok {
		return p.Validate()
	}
	return fmt.Errorf("%w: have: %T; want: %T",
		ErrS3ParamsWrongType, params, S3GetParams{})
}

func (s3 S3) ValidatePut(params PutParams) error {
	if params == nil {
		return ErrS3PutParamsEmpty
	}
	if p, ok := params.(S3PutParams); ok {
		return p.Validate()
	}
	return fmt.Errorf("%w: have: %T; want: %T",
		ErrS3ParamsWrongType, params, S3PutParams{})
}

type S3GetParams struct{}

func (S3GetParams) IsGetParams() {}

func (s3g S3GetParams) Validate() error {
	return nil
}

// S3PutParams is the Params obejct of an S3 resource.
// It implements the [Params] interface.
// See https://github.com/concourse/s3-resource#out-upload-an-object-to-the-bucket
type S3PutParams struct {
	// Required. File is the path to the file to upload, provided by an output of a task.
	// If multiple files are matched by the glob, an error is raised.
	// The file which matches will be placed into the directory structure on S3 as
	// defined in regexp in the resource definition. The matching syntax is bash glob
	// expansion, so no capture groups, etc.
	File string `json:"file"`
	//  Optional. Canned ACL for the uploaded object.
	ACL string `json:"acl,omitzero"`
	//  Optional. MIME Content-Type describing the contents of the uploaded object
	ContentType string `json:"content_type,omitzero"`
}

func (S3PutParams) IsPutParams() {}

func (s3p S3PutParams) Validate() error {
	var errs []error
	if s3p.File == "" {
		errs = append(errs, ErrS3PutParamsMissingFile)
	}
	return errors.Join(errs...)
}

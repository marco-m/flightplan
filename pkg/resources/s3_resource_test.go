// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources_test

import (
	"testing"

	"github.com/marco-m/flightplan/pkg/resources"
	"github.com/marco-m/rosina/assert"
)

func TestS3Type(t *testing.T) {
	sut := resources.S3{}
	assert.Equal(t, sut.Type(), "s3", "Type")
}

func TestS3Validate(t *testing.T) {
	test := func(desc string, sut resources.S3, wantErrs []error) {
		t.Helper()
		err := sut.Validate()
		for _, want := range wantErrs {
			assert.ErrorIs(t, err, want, desc)
		}
	}

	test("zero resource", resources.S3{},
		[]error{
			resources.ErrS3MissingBucket,
			resources.ErrS3NoRegexpNoVersionedFile,
		})
	test("incompatible fields",
		resources.S3{
			Regexp:        "mango",
			VersionedFile: "berry",
		},
		[]error{resources.ErrS3BothRegexpAndVersionedFile})
}

func TestS3ValidateGetSuccess(t *testing.T) {
	test := func(desc string, res resources.GetParams) {
		t.Helper()
		s3 := resources.S3{}
		err := s3.ValidateGet(res)
		assert.NoError(t, err, desc)
	}

	test("nil", nil)
	test("not empty", resources.S3GetParams{SkipDownload: true})
}

func TestS3ValidateGetFailure(t *testing.T) {
	test := func(desc string, res resources.GetParams, want error) {
		t.Helper()
		s3 := resources.S3{}
		err := s3.ValidateGet(res)
		assert.ErrorIs(t, err, want, desc)
	}

	test("empty", resources.S3GetParams{}, resources.ErrS3GetParamsEmpty)
	test("wrong type", resources.GitGetParams{}, resources.ErrS3ParamsWrongType)
}

func TestS3ValidatePutSuccess(t *testing.T) {
	test := func(desc string, res resources.PutParams) {
		t.Helper()
		s3 := resources.S3{}
		err := s3.ValidatePut(res)
		assert.NoError(t, err, desc)
	}

	test("field File set", resources.S3PutParams{File: "mango"})
}

func TestS3ValidatePutFailure(t *testing.T) {
	test := func(desc string, res resources.PutParams, want error) {
		t.Helper()
		s3 := resources.S3{}
		err := s3.ValidatePut(res)
		assert.ErrorIs(t, err, want, desc)
	}

	test("empty", nil, resources.ErrS3PutParamsEmpty)
	test("wrong type", resources.GitPutParams{}, resources.ErrS3ParamsWrongType)
	test("missing File", resources.S3PutParams{}, resources.ErrS3PutParamsMissingFile)
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources_test

import (
	"testing"

	"github.com/marco-m/flightplan/resources"
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
			Bucket:        "banana",
			Regexp:        "mango",
			VersionedFile: "berry",
		},
		[]error{resources.ErrS3BothRegexpAndVersionedFile})
}

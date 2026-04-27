// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources_test

import (
	"testing"

	"github.com/marco-m/flightplan/resources"
	"github.com/marco-m/rosina/assert"
)

func TestRegistryImageType(t *testing.T) {
	sut := resources.RegistryImage{}
	assert.Equal(t, sut.Type(), "registry-image", "Type")
}

func TestRegistryImageValidateSuccess(t *testing.T) {
	sut := resources.RegistryImage{Repository: "banana"}
	err := sut.Validate()
	assert.NoError(t, err, "Validate")
}

func TestRegistryImageValidateFailure(t *testing.T) {
	sut := resources.RegistryImage{}
	err := sut.Validate()
	assert.ErrorIs(t, err, resources.ErrRegistryImageRepositoryEmpty, "Validate")
}

func TestRegistryImageValidateGetSuccess(t *testing.T) {
	test := func(desc string, res resources.GetParams) {
		t.Helper()
		rim := resources.RegistryImage{}
		err := rim.ValidateGet(res)
		assert.NoError(t, err, desc)
	}

	test("nil", nil)
	test("", resources.RegistryImageGetParams{SkipDownload: true})
}

func TestRegistryImageValidateGetFailure(t *testing.T) {
	test := func(desc string, res resources.GetParams, want error) {
		t.Helper()
		rim := resources.RegistryImage{}
		err := rim.ValidateGet(res)
		assert.ErrorIs(t, err, want, desc)
	}

	test("wrong type", resources.S3GetParams{}, resources.ErrRegistryImageParamsWrongType)
}

func TestRegistryImageValidatePutSuccess(t *testing.T) {
	test := func(desc string, res resources.PutParams) {
		t.Helper()
		rim := resources.RegistryImage{}
		err := rim.ValidatePut(res)
		assert.NoError(t, err, desc)
	}

	test("valid", resources.RegistryImagePutParams{Image: "mango"})
}

func TestRegistryImageValidatePutFailure(t *testing.T) {
	test := func(desc string, res resources.PutParams, want error) {
		t.Helper()
		rim := resources.RegistryImage{}
		err := rim.ValidatePut(res)
		assert.ErrorIs(t, err, want, desc)
	}

	test("empty", nil, resources.ErrRegistryImagePutParamsEmpty)
	test("wrong type", resources.S3PutParams{}, resources.ErrRegistryImageParamsWrongType)
	test("missing Image", resources.RegistryImagePutParams{},
		resources.ErrRegistryImagePutParamsMissingImage)
}

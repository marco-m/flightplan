// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources_test

import (
	"testing"

	"github.com/marco-m/flightplan/pkg/resources"
	"github.com/marco-m/rosina/assert"
)

func TestGitType(t *testing.T) {
	sut := resources.Git{}
	assert.Equal(t, sut.Type(), "git", "Type")
}

func TestGitValidateSuccess(t *testing.T) {
	sut := resources.Git{
		Uri: "https://github.com/marco-m/flightplan.git",
	}
	err := sut.Validate()
	assert.NoError(t, err, "Validate")
}

func TestGitValidateFailure(t *testing.T) {
	sut := resources.Git{}
	err := sut.Validate()
	assert.ErrorIs(t, err, resources.ErrGitMissingUri, "Validate")
}

func TestGitValidateGetSuccess(t *testing.T) {
	test := func(desc string, res resources.GetParams) {
		t.Helper()
		git := resources.Git{}
		err := git.ValidateGet(res)
		assert.NoError(t, err, desc)
	}

	test("nil", nil)
	test("optional field", resources.GitGetParams{Depth: 2})
}

func TestGitValidateGetFailure(t *testing.T) {
	test := func(desc string, res resources.GetParams, want error) {
		t.Helper()
		git := resources.Git{}
		err := git.ValidateGet(res)
		assert.ErrorIs(t, err, want, desc)
	}

	test("empty", resources.GitGetParams{}, resources.ErrGitGetParamsEmpty)
	test("wrong type", resources.S3GetParams{}, resources.ErrGitParamsWrongType)
}

func TestGitValidatePutSuccess(t *testing.T) {
	test := func(desc string, res resources.PutParams) {
		t.Helper()
		git := resources.Git{}
		err := git.ValidatePut(res)
		assert.NoError(t, err, desc)
	}

	test("valid", resources.GitPutParams{
		// The way the Handle is constrcuted makes sense ONLY for a test!
		Repository: &resources.Handle{
			Resource: resources.Resource{Type: "git"},
		},
	})
}

func TestGitValidatePutFailure(t *testing.T) {
	test := func(desc string, res resources.PutParams, want error) {
		t.Helper()
		git := resources.Git{}
		err := git.ValidatePut(res)
		assert.ErrorIs(t, err, want, desc)
	}

	test("empty", nil, resources.ErrGitPutParamsEmpty)
	test("wrong type", resources.S3PutParams{}, resources.ErrGitParamsWrongType)
	test("wrong repo type", resources.GitPutParams{
		// This usage of the API makes sense ONLY for a test.
		Repository: &resources.Handle{Resource: resources.Resource{Type: "mango"}},
	}, resources.ErrGitPutParamsWrongRepoType)
	test("missing Repository", resources.GitPutParams{}, resources.ErrGitPutParamsMissingRepository)
}

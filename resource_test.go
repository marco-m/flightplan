// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"testing"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/flightplan/resources"

	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/check"
)

func TestLookupUnknownResource(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)
	assert.NoError(t, pipeline.Errors(), "Errors")

	_, found := pipeline.Resource(&resources.Handle{Resource: resources.Resource{Name: "foo"}})
	assert.False(t, found, "resource found")
}

func TestAddResourcesAndFindThemWithHandles(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	handle1 := pipeline.AddResource(resources.Resource{
		Name:   "res1.git",
		Source: resources.Git{Uri: "https://github.com/marco-m/res1.git"},
	})
	handle2 := pipeline.AddResource(resources.Resource{
		Name:   "res2.git",
		Source: resources.Git{Uri: "https://github.com/marco-m/res2.git"},
	})

	assert.NoError(t, pipeline.Errors(), "Errors")

	res1, found1 := pipeline.Resource(handle1)
	assert.True(t, found1, "found1")
	check.Equal(t, res1.Name, "res1.git", "res1.Name")
	check.Equal(t, res1.Type, "git", "res1.Type")

	res2, found2 := pipeline.Resource(handle2)
	assert.True(t, found2, "found2")
	check.Equal(t, res2.Name, "res2.git", "res2.Name")
	check.Equal(t, res2.Type, "git", "res2.Type")
}

func TestAddResourceKeepsNameAndSetsType(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	handle := pipeline.AddResource(resources.Resource{
		Name: "flightplan.git",
		Source: resources.Git{
			Uri: "https://github.com/marco-m/flightplan.git",
		},
	})

	assert.NoError(t, pipeline.Errors(), "Errors")
	res, found := pipeline.Resource(handle)
	assert.True(t, found, "found")
	check.Equal(t, res.Name, "flightplan.git", "res.Name")
	check.Equal(t, res.Type, "git", "res.Type")
}

func TestAddResourceNeedsName(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddResource(resources.Resource{
		// Name: "foo" <== missing
		Source: resources.Git{
			Uri: "https://github.com/marco-m/flightplan.git",
		},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrEmptyResourceName, "Errors")
}

func TestAddResourceWithTypeAlreadySetFails(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddResource(resources.Resource{
		Name: "flightplan.git",
		Type: "this-will-fail",
		Source: resources.Git{
			Uri: "https://github.com/marco-m/flightplan.git",
		},
	})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrSetResourceType, "Errors")
}

func TestAddResourceSourceCannotBeEmpty(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddResource(resources.Resource{Name: "banana"})

	assert.ErrorIs(t, pipeline.Errors(), plan.ErrMissingSource, "Errors")
}

func TestAddResourceSourceValidationFails(t *testing.T) {
	pipeline := plan.NewPipeline("pizza", nil)

	pipeline.AddResource(resources.Resource{
		Name:   "banana",
		Source: resources.S3{},
	})

	err := pipeline.Errors()
	assert.ErrorIs(t, err, plan.ErrSourceValidation, "Errors")
	assert.ErrorContains(t, err, "field Bucket cannot be empty", "Errors")
}

func TestPipelineCannotAddDuplicateResource(t *testing.T) {
	pipeline := plan.NewPipeline("duplicate-resource", nil)

	resource := resources.Resource{
		Name:   "resource.git",
		Source: resources.Git{Uri: "https://github.com/marco-m/resource.git"},
	}

	pipeline.AddResource(resource)
	pipeline.AddResource(resource)

	err := pipeline.Render()

	assert.ErrorIs(t, err, plan.ErrDuplicateResourceName, "Render")
}

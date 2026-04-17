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

func TestRegistryImageValidate(t *testing.T) {
	sut := resources.RegistryImage{}
	err := sut.Validate()
	assert.ErrorContains(t, err, "RegistryImage.Repository cannot be empty", "Validate")
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"testing"

	plan "github.com/marco-m/flightplan"
	"github.com/marco-m/rosina/assert"
)

func TestRegistryImageResourceType(t *testing.T) {
	sut := plan.RegistryImageSource{}
	assert.Equal(t, sut.Type(), "registry-image", "Type")
}

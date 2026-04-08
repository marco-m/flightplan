// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources_test

import (
	"testing"

	"github.com/marco-m/flightplan/resources"
	"github.com/marco-m/rosina/assert"
)

func TestGitResourceType(t *testing.T) {
	sut := resources.Git{}
	assert.Equal(t, sut.Type(), "git", "Type")
}

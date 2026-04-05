// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"testing"

	plan "github.com/marco-m/flightplan"

	"github.com/marco-m/rosina/assert"
)

func TestEmptyPipelineIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("empty", nil)

	err := pipeline.Render()

	assert.ErrorIs(t, err, plan.ErrEmptyPipeline, "Render")
}

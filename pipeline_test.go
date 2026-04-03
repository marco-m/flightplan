// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"errors"
	"testing"

	plan "github.com/marco-m/flightplan"
)

func TestEmptyPipelineIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("empty", nil)

	err := pipeline.Render()

	if !errors.Is(err, plan.ErrEmptyPipeline) {
		t.Fatalf("Render:\nhave: %v\nwant: %v", err, plan.ErrEmptyPipeline)
	}
}

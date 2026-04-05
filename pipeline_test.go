// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan_test

import (
	"path/filepath"
	"testing"

	plan "github.com/marco-m/flightplan"

	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/check"
)

func TestPipelinePath(t *testing.T) {
	test := func(desc, name string, args []string, wantSuffix string) {
		t.Helper()
		pipeline := plan.NewPipeline(name, args)
		check.Contains(t, filepath.ToSlash(pipeline.Path()), wantSuffix, "Path "+desc)
		assert.True(t, filepath.IsAbs(pipeline.Path()), "Path is absolute")
	}

	test("default dir", "dummy", nil, "/flightplan/dummy.json")
	test("cli override name", "dummy", []string{"-name=mango"},
		"/flightplan/mango.json")
	test("cli override relative dir", "dummy", []string{"-directory=berry"},
		"/flightplan/berry/dummy.json")
	test("cli override absolute dir", "dummy", []string{"-directory=/mango"},
		"/mango/dummy.json")
}

func TestMustUseNewPipeline(t *testing.T) {
	pipeline := plan.Pipeline{}
	err := pipeline.Render()
	assert.ErrorIs(t, err, plan.ErrMissingNewPipeline, "Render")
}

func TestEmptyPipelineIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("empty", nil)
	err := pipeline.Render()
	assert.ErrorIs(t, err, plan.ErrEmptyPipeline, "Render")
}

func TestPipelineWithoutNameIsInvalid(t *testing.T) {
	pipeline := plan.NewPipeline("", nil)
	err := pipeline.Render()
	assert.ErrorIs(t, err, plan.ErrEmptyPipelineName, "Render")
}

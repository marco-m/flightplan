// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package goof_test

import (
	"errors"
	"testing"

	"github.com/marco-m/flightplan/internal/goof/example"
	"github.com/marco-m/rosina/assert"
)

func TestAddResource(t *testing.T) {
	err := example.AddResourceWithErr("banana")
	// The test is directly simulating client code, so we want the location to refer to
	// this file and the test itself.
	want := "/flightplan/internal/goof/goof_test.go:15: AddResource: banana"
	assert.ErrorContains(t, err, want, "AddResourceWithErr")
}

func TestSimulateClientCallingAddResourceWithErr(t *testing.T) {
	err := example.SimulateClientCallingAddResourceWithErr("orange")
	// The test is validating the simulation of client code done by
	// SimulateClientCallingAddResource, so we want the location to refer to example.go
	want := "/flightplan/internal/goof/example/example.go:22: AddResource: orange"
	assert.ErrorContains(t, err, want, "SimulateClientCallingAddResourceWithErr")
}

func TestWrapSentinel(t *testing.T) {
	err := example.WrapSentinel("blueberry")

	if !errors.Is(err, example.ErrBanana) {
		t.Fatalf("error is not ErrBanana: %v", err)
	}
	want := "/flightplan/internal/goof/goof_test.go:31: WrapSentinel: blueberry: bananas are not ripe"
	assert.ErrorContains(t, err, want, "WrapSentinel")
}

func TestSimulateClientUsingPipeline(t *testing.T) {
	err := example.SimulateClientUsingPipeline()
	// The test is validating the simulation of client code done by
	// SimulateClientUsingPipeline, so we want the locations (notice the plural) to refer
	// to example.go
	assert.ErrorContains(t, err,
		"/flightplan/internal/goof/example/example.go:44: AddResource: lime\n",
		"line1")
	assert.ErrorContains(t, err,
		"/flightplan/internal/goof/example/example.go:45: AddResource: guava",
		"line2")
}

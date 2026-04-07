// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package goof_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/marco-m/flightplan/internal/goof/example"
)

func TestAddResource(t *testing.T) {
	err := example.AddResourceWithErr("banana")
	// The test is directly simulating client code, so we want the location to refer to
	// this file and the test itself.
	want := "internal/goof/goof_test.go:15: AddResource: banana"
	if have := err.Error(); have != want {
		t.Errorf("\nhave: %v\nwant: %v", have, want)
	}
}

func TestSimulateClientCallingAddResourceWithErr(t *testing.T) {
	err := example.SimulateClientCallingAddResourceWithErr("orange")
	// The test is validating the simulation of client code done by
	// SimulateClientCallingAddResource, so we want the location to refer to example.go
	want := "internal/goof/example/example.go:22: AddResource: orange"
	if have := err.Error(); have != want {
		t.Fatalf("%s:\nhave: %v\nwant: %v", "WrapErr", have, want)
	}
}

func TestWrapSentinel(t *testing.T) {
	err := example.WrapSentinel("blueberry")

	if !errors.Is(err, example.ErrBanana) {
		t.Fatalf("error is not ErrBanana: %v", err)
	}
	want := "internal/goof/goof_test.go:35: WrapSentinel: blueberry: bananas are not ripe"
	if have := err.Error(); have != want {
		t.Fatalf("%s:\nhave: %v\nwant: %v", "WrapErr", have, want)
	}
}

func TestSimulateClientUsingPipeline(t *testing.T) {
	err := example.SimulateClientUsingPipeline()
	// The test is validating the simulation of client code done by
	// SimulateClientUsingPipeline, so we want the locations (notice the plural) to refer
	// to example.go
	want := strings.TrimSpace(`
internal/goof/example/example.go:44: AddResource: lime
internal/goof/example/example.go:45: AddResource: guava
`)
	if have := err.Error(); have != want {
		t.Fatalf("\nhave:\n%v\nwant:\n%v", have, want)
	}
}

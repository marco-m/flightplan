// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// This package exists to write tests for package goof.
package example

import (
	"errors"

	"github.com/marco-m/flightplan/internal/goof"
)

var ErrBanana = errors.New("bananas are not ripe")

// Simulates a Pipeline.AddResourceWithErr that returns an error.
func AddResourceWithErr(what string) error {
	return goof.Wrap("AddResource: %s", what)
}

// Simulates client code that calls AddResource.
func SimulateClientCallingAddResourceWithErr(what string) error {
	return AddResourceWithErr(what)
}

func WrapSentinel(what string) error {
	return goof.Wrap("WrapSentinel: %s: %w", what, ErrBanana)
}

type Pipeline struct {
	errs []error
}

func (pl *Pipeline) AddResource(name string) {
	pl.errs = append(pl.errs, goof.Wrap("AddResource: %s", name))
}

func (pl *Pipeline) Render() error {
	return errors.Join(pl.errs...)
}

// This similation is higher fidelity than SimulateClientCallingAddResourceWithErr
func SimulateClientUsingPipeline() error {
	pl := Pipeline{}
	pl.AddResource("lime")
	pl.AddResource("guava")
	return pl.Render()
}

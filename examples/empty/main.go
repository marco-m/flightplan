// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// Minimal usage of flightplan. Will return an error because the pipeline is empty.
// Used in the testscript tests.
package main

import (
	"fmt"
	"os"

	plan "github.com/marco-m/flightplan"
)

func main() {
	if err := buildPipeline(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func buildPipeline() error {
	pipeline := plan.NewPipeline("minimal", os.Args[1:])
	return pipeline.Render()
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrEmptyPipeline = errors.New("pipeline cannot be empty")

// Pipeline is used to construct a Concourse pipeline. Use [NewPipeline] to instantiate.
type Pipeline struct {
	// The name of the pipeline.
	name string
	// Errors collected during construction and returned by [Pipeline.Render].
	errs []error
}

// NewPipeline parses the command-line and returns a *[Pipeline] ready to use.
//
// Param 'name' is the the name of the pipeline, used to construct its file name as
// <name.json>. Can be overridden on the command-line.
//
// Note: NewPipeline can call [os.Exit] during command-line parsing.
//
// Usage:
//
// pipeline := plan.NewPipeline("default-name", os.Args[1:])
func NewPipeline(name string, args []string) *Pipeline {
	pl := &Pipeline{name: name}

	parseCommandLine(pl, args)

	return pl
}

func (pl *Pipeline) Render() error {
	pl.errs = append(pl.errs, fmt.Errorf("render: %w", ErrEmptyPipeline))
	err := errors.Join(pl.errs...)
	if err != nil {
		return err
	}
	return nil
}

// Parses the command-line and overrides settings in 'pl' based on the parse.
// Can call directly or indirectly [os.Exit].
func parseCommandLine(pl *Pipeline, args []string) {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s -- generate a Concourse pipeline with flightplan\n\n",
			fs.Name())
		fmt.Fprintf(os.Stderr, "Usage: %s [ARGS]\n\n", fs.Name())
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}
	fs.StringVar(&pl.name, "name", pl.name, "Name of the pipeline file")
	fs.Parse(args)
	if len(fs.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "unrecognized arguments: %v\n", strings.Join(fs.Args(), " "))
		os.Exit(1)
	}
}

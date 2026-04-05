// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/marco-m/flightplan/internal/goof"
)

var (
	ErrEmptyPipeline      = errors.New("pipeline cannot be empty")
	ErrEmptyPipelineName  = errors.New("pipeline name cannot be empty")
	ErrMissingNewPipeline = errors.New("must use NewPipeline to create a pipeline")
	ErrSystem             = errors.New("unexpected")
)

// Pipeline is used to construct a Concourse pipeline. Use [NewPipeline] to instantiate.
type Pipeline struct {
	// The current working directory.
	cwd string
	// The directory set on the command-line
	dir string
	// The name of the pipeline.
	name string
	// Errors collected during construction and returned by [Pipeline.Render].
	errs        []error
	newPipeline bool // True if instantiated by NewPipeline.
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
// pipeline := plan.NewPipeline("default-name", "default-dir", os.Args[1:])
func NewPipeline(name string, args []string) *Pipeline {
	pl := &Pipeline{name: name, newPipeline: true}

	cwd, err := os.Getwd()
	if err != nil {
		pl.errs = append(pl.errs, goof.Wrap("NewPipeline: %w", ErrSystem))
	}
	pl.cwd = cwd
	if pl.name == "" {
		pl.errs = append(pl.errs, goof.Wrap("NewPipeline: %w", ErrEmptyPipelineName))
	}
	// TODO ADD OTHER CHECKS: should be only 0-9 A-Z a-z _-
	// For sure not / anywhere
	// no . as first character
	// no . as last character

	parseCommandLine(pl, args)

	return pl
}

// Render generates the pipeline and writes it by default to the same directory where
// the program is running. The directory can be changed to a relative or absolute path
// with the command-line --directory option.
// Render retuns all the errors collected during the usage of the flightplan API.
func (pl *Pipeline) Render() error {
	pl.errs = append(pl.errs, fmt.Errorf("render: %w", ErrEmptyPipeline))
	if !pl.newPipeline {
		// prepend
		pl.errs = slices.Insert(pl.errs, 0, goof.Wrap("render: %w", ErrMissingNewPipeline))
	}
	err := errors.Join(pl.errs...)
	if err != nil {
		return err
	}
	return nil
}

// Path returns the path of the rendered pipeline.
// Client code doesn't need to call this function.
func (pl *Pipeline) Path() string {
	name := pl.name
	if filepath.Ext(name) == "" {
		name += ".json"
	}
	if filepath.IsAbs(pl.dir) {
		// Do not consider the current working directory.
		return filepath.Join(pl.dir, name)
	}
	return filepath.Join(pl.cwd, pl.dir, name)
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
	fs.StringVar(&pl.dir, "directory", ".", "Directory in which to write the pipeline file")
	fs.Parse(args)
	if len(fs.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "unrecognized arguments: %v\n", strings.Join(fs.Args(), " "))
		os.Exit(1)
	}
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// On Windows, the embedded files in a txtar are NOT extracted. Seems a bug.
//go:build !windows

package flightplan_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

func TestScript(t *testing.T) {
	runScriptTests(t, "testdata/script/*.txtar")
}

func runScriptTests(t *testing.T, pattern string) {
	dir := t.TempDir()

	// The script environment variable PATH has meaning similar to PATH for a shell:
	// an executable  'foo' in PATH can be invoked in a test script with 'exec foo ...'.
	// That is, we put in PATH the systems under test (SUTs).
	env := []string{
		"PATH=" + dir,
		// "go test -cover" already sets a GOCOVERDIR; we don't want that!
		"GOCOVERDIR=" + os.Getenv("COVER_INTEGRATION"),
		// fragile but maybe good enough?
		"TESTDATA=" + filepath.Join(os.Getenv("PWD"), "testdata"),
	}

	engine := &script.Engine{
		Cmds:  scripttest.DefaultCmds(),
		Conds: scripttest.DefaultConds(),
		Quiet: !testing.Verbose(),
	}
	engine.Conds["concourse"] = script.BoolCondition("env var FLIGHTPLAN_CONCOURSE is set",
		os.Getenv("FLIGHTPLAN_CONCOURSE") != "")

	// Make an executable found in the host PATH available to the test script for
	// direct invocation (will also show in the help output). Contrast above with
	// putting a custom-built executable with goBuild into env["PATH"].
	engine.Cmds["ls"] = script.Program("ls", nil, 100*time.Millisecond)
	engine.Cmds["fly"] = script.Program("fly", nil, 100*time.Millisecond)

	ctx := context.Background()
	scripttest.Test(t, ctx, engine, env, pattern)
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// On Windows, the embedded files in a txtar are NOT extracted. Seems a bug.
//go:build !windows

package flightplan_test

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

func TestScriptExamples(t *testing.T) {
	runScriptTests(t, "testdata/script/*.txtar")
}

func runScriptTests(t *testing.T, pattern string) {
	dir := t.TempDir()
	if _, err := goBuild("examples-empty", "./examples/empty", dir); err != nil {
		t.Fatal(err)
	}
	// The script environment variable PATH has meaning similar to PATH for a shell:
	// an executable  'foo' in PATH can be invoked in a test script with 'exec foo ...'.
	// That is, we put in PATH the systems under test (SUTs).
	env := []string{"PATH=" + dir}

	engine := &script.Engine{
		Cmds:  scripttest.DefaultCmds(),
		Conds: scripttest.DefaultConds(),
		Quiet: !testing.Verbose(),
	}

	// How to make executables found in the host PATH available to the test script:
	engine.Cmds["ls"] = script.Program("ls", nil, 100*time.Millisecond)

	ctx := context.Background()
	scripttest.Test(t, ctx, engine, env, pattern)
}

// Build Go code in 'src', give it 'name', put it in 'dir'. Return path 'dir/name'.
func goBuild(name, src, dir string) (string, error) {
	outPath := filepath.Join(dir, name)
	// staticcheck: due to the file's build constraints, runtime.GOOS will never equal "windows"
	// if runtime.GOOS == "windows" {
	// 	outPath += ".exe"
	// }
	out, err := exec.Command("go", "build", "-o", outPath, src).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("building %s: %s\n%s", src, err, string(out))
	}
	return outPath, err
}

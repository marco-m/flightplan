// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

// On Windows, the embedded files in a txtar are NOT extracted. Seems a bug.
//go:build !windows

package flightplan_test

import (
	"context"
	"fmt"
	"os"
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
	if _, err := goBuild("examples-simple", "./examples/simple", dir); err != nil {
		t.Fatal(err)
	}
	// The script environment variable PATH has meaning similar to PATH for a shell:
	// an executable  'foo' in PATH can be invoked in a test script with 'exec foo ...'.
	// That is, we put in PATH the systems under test (SUTs).
	env := []string{
		"PATH=" + dir,
		// "go test -cover" already sets a GOCOVERDIR; we don't want that!
		"GOCOVERDIR=" + os.Getenv("COVER_INTEGRATION"),
		// "HOST_PWD=" + os.Getenv("PWD"),
	}

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

// Change directory to 'srcDir' and build the Go code there, call the executable 'name',
// put it in 'dstDir'. Return path 'dstDir/name'.
// If the tests are invoked from our build script (as opposed to a manual "go test"),
// also instrument for code coverage.
// Changind directory to 'srcDir' allows to handle the case of multiple Go modules in
// the same repository.
func goBuild(name, srcDir, dstDir string) (string, error) {
	outPath := filepath.Join(dstDir, name)
	// staticcheck: due to the file's build constraints, runtime.GOOS will never equal "windows"
	// if runtime.GOOS == "windows" {
	// 	outPath += ".exe"
	// }

	args := []string{"build", "-C", srcDir, "-o", outPath}
	// -cover: code coverage for integration testing; see https://go.dev/doc/build-cover
	// We add --cover only if running from our build script (detected by the presence of
	// the COVER_INTEGRATION env var) to avoid the message:
	//     warning: GOCOVERDIR not set, no coverage data emitted
	// when running the tests directly from the shell.
	if os.Getenv("COVER_INTEGRATION") != "" {
		args = append(args, "-cover")
	}
	out, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("building %s: %s\n%s", srcDir, err, string(out))
	}
	return outPath, err
}

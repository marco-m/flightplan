// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package testhelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-m/rosina/golden"
)

func AssertRenderedEqualsGolden(t *testing.T, havePath, goldenPath string, update bool) {
	t.Helper()
	if diff := golden.DiffFiles(t, havePath, goldenPath, update); diff != "" {
		t.Errorf("Render: mismatch:\n%s", diff)
	}
}

// makeFakeGitRepo creates an empty directory ".git" below 'dir'.
// This is enough for findRepoRoot() to succeed.
// 'dir' must have been created by t.TempDir().
func MakeFakeGitRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatalf("MakeFakeGitRepo: %s", err)
	}
}

func MakeDirAll(t *testing.T, dir string, mkdir bool) {
	t.Helper()
	if !mkdir {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MakeDirAll: %s", err)
	}
}

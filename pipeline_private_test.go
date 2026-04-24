// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

import (
	"path/filepath"
	"testing"

	"github.com/marco-m/flightplan/internal/testhelpers"
	"github.com/marco-m/flightplan/resources"
	"github.com/marco-m/rosina/assert"
	"github.com/marco-m/rosina/check"
)

func TestPipelineRelDirSuccess(t *testing.T) {
	test := func(desc string, pldir string, mkdir bool, want string) {
		t.Helper()
		tmpDir := t.TempDir()
		testhelpers.MakeFakeGitRepo(t, tmpDir)
		dir := filepath.Join(tmpDir, pldir)
		testhelpers.MakeDirAll(t, dir, mkdir)
		pl := NewPipeline("dummy", []string{"-directory", dir})
		repo := pl.AddResource(resources.Resource{
			Name:   "banana.git",
			Source: resources.Git{Uri: "https://example.org/mango.git"},
		})
		reldir := pl.RelDir(repo)
		check.Equal(t, reldir, want, desc)
		assert.NoError(t, pl.Errors(), desc)
	}

	test("pipeline at root of repository", ".", true, "banana.git")
	test("pipeline one level below root, dir existing", "examples", true,
		"banana.git/examples")
	test("pipeline one level below root, dir not existing", "non-existing", false,
		"banana.git/non-existing")
}

func TestPipelineRelDirFailure(t *testing.T) {
	test := func(desc string, args []string, wantErr error) {
		t.Helper()
		pl := NewPipeline("dummy", args)
		repo := pl.AddResource(resources.Resource{
			Name:   "banana.git",
			Source: resources.Git{Uri: "https://example.org/mango.git"},
		})
		pl.RelDir(repo)
		assert.ErrorIs(t, pl.Errors(), wantErr, desc)
	}

	test("cannot create pipeline directory", []string{`--directory=/cannot-write`},
		ErrCreatePipelineDir)

	// FIXME I have a mess of path.Foo and filepath.Foo. Now, _normally_ I should use
	// filepath, BUT the paths in the pipeline, the ones passed to [Task.File] and
	// [TaskCommand.Path], I think MUST be like URLs, so with forward slash "/" and so
	// with path.Foo, not filepath.Foo. I need to test what happens on Windows and just
	// by passing a path separate with back slash "\".
}

func TestReconstructRepoRoot(t *testing.T) {
	test := func(desc, pldir, taskpath, wantRoot, wantPath string) {
		t.Helper()
		root, path := reconstructRepoRoot(pldir, taskpath)

		check.Equal(t, root, wantRoot, desc)
		check.Equal(t, path, wantPath, desc)
	}

	repoRoot := "/home/user/repo"
	test("pl at root, task with pl", repoRoot, "mango.git/ta.json",
		repoRoot, "ta.json")
	test("pl at root, task below pl", repoRoot, "mango.git/ci/tasks/ta.json",
		repoRoot, "ci/tasks/ta.json")
	pldir := repoRoot + "/ci"
	test("pl below root, task with pl", pldir, "mango.git/ci/ta.json",
		repoRoot, "ci/ta.json")
	test("pl below root, task below pl", pldir, "mango.git/ci/tasks/ta.json",
		repoRoot, "ci/tasks/ta.json")
}

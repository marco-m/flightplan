// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package internal_test

import (
	"runtime"
	"testing"

	"github.com/marco-m/flightplan/internal"
	"github.com/marco-m/rosina/assert"
)

func TestTrimModuleFound(t *testing.T) {
	_, myPath, _, _ := runtime.Caller(0)
	trimmed := internal.TrimModule(myPath)
	assert.Equal(t, trimmed, "internal/internal_test.go", "TrimModule")
}

func TestTrimModuleNotFoundReturnsReasonable(t *testing.T) {
	myPath := "/abs/foo/bar.go"
	trimmed := internal.TrimModule(myPath)
	assert.Equal(t, trimmed, myPath, "TrimModule")
}

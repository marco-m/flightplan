// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package goof

import (
	"fmt"
	"runtime"
)

// Wrap creates an error and wraps it with file and line number, in such a way that
// the user can click on the printed error in an editor or in a terminal and jump to
// the corresponding file and line number of the invocation, exactly as if it was
// a compiler error.
//
// Usage:
//
//	return goof.Wrap("some text: %w", err)
//
// and then printing it:
//
//	some/dir/file.go:32: some text: the error
func Wrap(format string, a ...any) error {
	// FXIME For performance reasons, this should insteadd call Callers and store the
	// program counter, and resolve it only within Error().
	// See https://pkg.go.dev/log/slog@go1.26.1#example-package-Wrapping
	_, filePath, line, _ := runtime.Caller(2)
	err := fmt.Errorf(format, a...)
	return fmt.Errorf("%s:%d: %w", filePath, line, err)
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package internal

import (
	"path"
	"runtime/debug"
	"strings"
)

// TrimModule makes the absolute filePath relative to the Go module it belongs to.
// Example: /home/alice/src/mymodule/foo/bar.go -> foo/bar.go
func TrimModule(filePath string) string {
	info, _ := debug.ReadBuildInfo()
	base := path.Base(info.Main.Path)
	idx := strings.Index(filePath, base)
	if idx < 0 {
		return filePath
	}
	return filePath[idx+len(base)+1:]
}

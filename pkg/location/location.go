package location

import (
	"path/filepath"
	"runtime"
)

// Directory returns the directory of the Go source code calling this function.
// Note: this can be different from the directory of the Go executable.
func Directory() (directory string, ok bool) {
	_, file, _, ok := runtime.Caller(2)
	return filepath.Dir(file), ok
}

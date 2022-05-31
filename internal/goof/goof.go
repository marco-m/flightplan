package goof

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
)

func NewErr(format string, a ...any) error {
	_, file, line, _ := runtime.Caller(2)
	err := fmt.Errorf(format, a...)
	return fmt.Errorf("%s:%d: %w", file, line, err)
}

func WrapErr(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s:%d: %w", file, line, err)
}

func X() {
	info, _ := debug.ReadBuildInfo()
	fmt.Println(info.Main.Path)
	base := path.Base(info.Main.Path)
	fmt.Println(base)
	x := "/Users/mmolteni/src/hw/flightplan/internal/goof/goof_test.go:17"
	idx := strings.Index(x, base)
	if idx < 0 {
		// return x
	}
	// return x[idx:]
	fmt.Println("/"+x[idx:]+":", "aspide")
}

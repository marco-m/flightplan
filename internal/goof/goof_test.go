package goof_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/marco-m/flightplan/internal/goof"
)

func failing() error {
	return fmt.Errorf("banana")
}

func failingWithLocation() error {
	if err := failing(); err != nil {
		return goof.WrapErr(err)
	}
	return nil
}

func failingWithLocation2() error {
	if err := failing(); err != nil {
		return goof.WrapErr(err)
	}
	return nil
}

func TestWrapErr(t *testing.T) {
	err := failingWithLocation()
	goof.X()

	if have, want := err.Error(), "banana"; have != want {
		t.Fatalf("%s:\nhave: %v\nwant: %v", "WrapErr", have, want)
	}
}

func TestWrapErrMultiple(t *testing.T) {
	err := failingWithLocation()
	err2 := failingWithLocation2()
	err3 := errors.Join(err, err2)

	if have, want := err3.Error(), "banana"; have != want {
		t.Fatalf("%s:\nhave:\n%v\nwant: %v", "WrapErr", have, want)
	}
}

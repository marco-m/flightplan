package assert

import (
	"bytes"
	"errors"
	"testing"

	plan "github.com/marco-m/flightplan/flightplan"
	"github.com/marco-m/flightplan/internal/testhelpers/golden"
)

func NoError(t *testing.T, err error, desc string) {
	t.Helper()

	if err != nil {
		t.Fatalf("%s:\nunexpected error: %v (%T)", desc, err, err)
	}
}

func ErrorIs(t *testing.T, have, want error, desc string) {
	if have == nil {
		t.Fatalf("%s:\nhave: <no error>\nwant: %v (%T)", desc, want, want)
	}
	if !errors.Is(have, want) {
		t.Fatalf("%s:\nhave: %v (%T)\nwant: %v (%T)", desc, have, have, want, want)
	}
}

func RenderedEqualsGolden(t *testing.T, pipeline *plan.Pipeline, goldenPath string, update bool) {
	t.Helper()
	var out bytes.Buffer

	err := pipeline.Render(&out)
	if err != nil {
		t.Fatalf("%s:\nerror: %v", "Render", err)
	}
	if diff := golden.DiffText(t, out.String(), goldenPath, update); diff != "" {
		t.Fatalf("Render: mismatch:\n%s", diff)
	}
}

package golden

import (
	"io"
	"os"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/marco-m/rosina/diff"
)

//
// Code taken from my "quattro" project.
// TODO(marco) extract to "rosina"
//

// DiffText returns a text diff between the contents of the file at
// goldenPath and the string have (if update is false).
// If update is true, DiffText updates the contents of the file at goldenPath with
// have and returns an empty diff.
func DiffText(t *testing.T, have string, goldenPath string, update bool) string {
	t.Helper()
	haveStr := have
	wantStr := ReadOrUpdate(t, haveStr, goldenPath, update)
	return string(diff.TextDiffPatient("want", []byte(wantStr), "have", []byte(haveStr)))
}

// DiffAny returns a human-readable diff between the contents of the file at
// goldenPath and the text representation of have (if update is false).
// If update is true, DiffAny updates the contents of the file at goldenPath with
// the text representation of have and returns an empty diff.
func DiffAny(t *testing.T, have any, goldenPath string, update bool) string {
	t.Helper()
	haveStr := repr.String(have, repr.Indent("  "), repr.OmitEmpty(false), repr.OmitZero(false))
	wantStr := ReadOrUpdate(t, haveStr, goldenPath, update)
	return string(diff.TextDiffPatient("want", []byte(wantStr), "have", []byte(haveStr)))
}

// ReadOrUpdate returns the contents of the file at goldenPath (if update is false).
// If update is true, ReadOrUpdate updates the contents of the file at goldenPath
// with have and returns have.
func ReadOrUpdate(t *testing.T, have string, goldenPath string, update bool) string {
	t.Helper()
	fi, err := os.OpenFile(goldenPath, os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("opening file %s: %s", goldenPath, err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			t.Fatalf("closing file: %s", err)
		}
	}()

	if update {
		t.Logf("updating golden file: %s", goldenPath)
		if _, err := fi.WriteString(have); err != nil {
			t.Fatalf("writing file %s: %s", goldenPath, err)
		}
		return have
	}

	content, err := io.ReadAll(fi)
	if err != nil {
		t.Fatalf("reading file %s: %s", goldenPath, err)
	}
	return string(content)
}

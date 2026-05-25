// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package testhelpers

import (
	"encoding/json"
	"testing"

	"github.com/marco-m/rosina/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// BodyFormatHook reformats the response body (assumed to be JSON) to human-readable
// multiple lines.
func BodyFormatHook(i *cassette.Interaction) error {
	switch i.Response.Headers.Get("Content-Type") {
	case "application/json":
		var body any
		if err := json.Unmarshal([]byte(i.Response.Body), &body); err != nil {
			return err
		}
		bytes, err := json.MarshalIndent(body, "", "  ")
		if err != nil {
			return err
		}
		i.Response.Body = string(bytes)
		return nil
	default:
		return nil
	}
}

// SetupRecorder returns a *[recorder.Recorder] ready to use and a teardown function.
// Usage example:
//
//	rec, teardown := testhelpers.SetupRecorder(t, "testdata/NAME")
//	t.Cleanup(func() { teardown(t) })
func SetupRecorder(t *testing.T, name string) (*recorder.Recorder, func(t *testing.T)) {
	rec, err := recorder.New(name,
		recorder.WithMode(recorder.ModeRecordOnce),
		recorder.WithHook(BodyFormatHook, recorder.AfterCaptureHook),
		recorder.WithSkipRequestLatency(true),
	)
	assert.NoError(t, err, "recorder.New")
	return rec,
		func(t *testing.T) {
			// Teardown recorder. Needed to save the cassette on first run.
			assert.NoError(t, rec.Stop(), "rec.Stop")
		}
}

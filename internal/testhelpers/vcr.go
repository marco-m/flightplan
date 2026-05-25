// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package testhelpers

import (
	"encoding/json"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
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

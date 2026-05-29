// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"fmt"
	"net/http"
)

func responseError(statusCode int, body []byte, readErr error) error {
	if readErr != nil {
		return fmt.Errorf("unexpected response\nstatus: %s\nbody: %s\nread error: %s",
			http.StatusText(statusCode), string(body), readErr)
	}
	return fmt.Errorf("unexpected response\nstatus: %s\nbody: %s",
		http.StatusText(statusCode), string(body))
}

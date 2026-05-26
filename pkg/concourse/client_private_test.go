// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/marco-m/rosina/assert"
)

func TestRedaction(t *testing.T) {
	client, err := NewClient(ClientArgs{ServerURL: "the-server"})
	assert.NoError(t, err, "NewClient")
	client.token = "the-secret"

	t.Run("Stringer", func(t *testing.T) {
		have := fmt.Sprint(client)
		assert.Contains(t, have, "token=[redacted]", "")
		assert.NotContains(t, have, client.token, "")
	})
	t.Run("GoString", func(t *testing.T) {
		have := fmt.Sprintf("%#v", client)
		assert.Contains(t, have, "token=[redacted]", "")
		assert.NotContains(t, have, client.token, "")
	})
	t.Run("Slog", func(t *testing.T) {
		var bld strings.Builder
		log := slog.New(slog.NewTextHandler(&bld, nil))
		log.Info("client-created", "client", client)
		have := bld.String()
		assert.Contains(t, have, "token=[redacted]", "")
		assert.NotContains(t, have, client.token, "")
	})
}

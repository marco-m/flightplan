// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse_test

import (
	"context"
	"testing"

	"github.com/marco-m/rosina/assert"

	"github.com/marco-m/flightplan/internal/testhelpers"
	"github.com/marco-m/flightplan/pkg/concourse"
)

func TestClient_GetInfo(t *testing.T) {
	// Arrange recorder.
	rec, teardown := testhelpers.SetupRecorder(t, "testdata/getinfo")
	t.Cleanup(func() { teardown(t) })

	// Arrange SUT.
	client, err := concourse.NewClient(concourse.ClientArgs{
		ServerURL:  "http://localhost:8080",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	have, err := client.GetInfo(ctx)

	// Assert.
	assert.NoError(t, err, "GetInfo")
	want := concourse.Info{
		Version:       "8.2.1",
		WorkerVersion: "2.5",
		FeatureFlags: map[string]bool{
			"build_rerun":            false,
			"cache_streamed_volumes": true,
			"global_resources":       false,
			"resource_causality":     true,
		},
		ExternalURL: "http://localhost:8080",
		ClusterName: "dev",
	}
	assert.DeepEqual(t, have, want, "GetInfo")
}

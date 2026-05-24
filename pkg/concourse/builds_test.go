// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse_test

import (
	"context"
	"testing"
	"time"

	"github.com/marco-m/rosina/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/marco-m/flightplan/pkg/concourse"
)

var unixEpoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestClient_ListPipelineBuilds(t *testing.T) {
	// Arrange recorder.
	rec, err := recorder.New("testdata/list-pipeline-builds-short",
		recorder.WithSkipRequestLatency(true),
	)
	assert.NoError(t, err, "recorder.New")
	// Teardown recorder. Needed to save the cassette on first run.
	t.Cleanup(func() {
		assert.NoError(t, rec.Stop(), "rec.Stop")
	})

	// Arrange SUT.
	const team = "main"
	const pipeline = "concourse"
	concourseClient, err := concourse.NewClient(concourse.Client{
		Server:     "https://ci.concourse-ci.org",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	have, err := concourseClient.ListPipelineBuilds(ctx, team, pipeline)

	// Assert.
	assert.NoError(t, err, "ListPipelineBuilds")

	want := []concourse.Build{
		{
			ID:           214953806,
			TeamName:     "main",
			Name:         "575",
			Status:       "started",
			APIURL:       "/api/v1/builds/214953806",
			JobName:      "dev-image",
			PipelineID:   24,
			PipelineName: "concourse",
			StartTime:    time.Date(2023, time.October, 16, 14, 7, 45, 0, time.UTC),
			EndTime:      unixEpoch,
			ReapTime:     unixEpoch,
		},
		{
			ID:           214953805,
			TeamName:     "main",
			Name:         "402",
			Status:       "succeeded",
			APIURL:       "/api/v1/builds/214953805",
			JobName:      "quickstart-smoke",
			PipelineID:   24,
			PipelineName: "concourse",
			StartTime:    time.Date(2023, time.October, 16, 14, 7, 45, 0, time.UTC),
			EndTime:      time.Date(2023, time.October, 16, 14, 11, 55, 0, time.UTC),
			ReapTime:     unixEpoch,
		},
	}
	assert.DeepEqual(t, have, want, "ListPipelineBuilds")
}

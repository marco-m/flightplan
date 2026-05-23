// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/marco-m/flightplan/pkg/concourse"
	"github.com/marco-m/rosina/assert"
)

func TestClient_ListPipelineBuilds(t *testing.T) {
	// Arrange recorder.
	rec, err := recorder.New("testdata/list-pipeline-builds-short",
		recorder.WithMode(recorder.ModeRecordOnce),
		recorder.WithRealTransport(http.DefaultTransport),
		recorder.WithSkipRequestLatency(true),
	)
	assert.NoError(t, err, "recorder.New")

	// Arrange SUT.
	const teamName = "main"
	const pipelineName = "concourse"
	concourseClient, err := concourse.NewClient(concourse.Client{
		Server:     "https://ci.concourse-ci.org",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	have, err := concourseClient.ListPipelineBuilds(ctx, teamName, pipelineName)

	// Teardown
	// Needed to actually save the cassette on first run.
	assert.NoError(t, rec.Stop(), "rec.Stop")

	// Assert.
	assert.NoError(t, err, "ListPipelineBuilds")

	want := []concourse.Build{
		{
			Id:           214953806,
			TeamName:     "main",
			Name:         "575",
			Status:       "started",
			ApiUrl:       "/api/v1/builds/214953806",
			JobName:      "dev-image",
			PipelineId:   24,
			PipelineName: "concourse",
			StartTime:    time.Date(2023, time.October, 16, 14, 7, 45, 0, time.UTC),
			EndTime:      time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
			CreatedBy:    "",
		},
		{
			Id:           214953805,
			TeamName:     "main",
			Name:         "402",
			Status:       "succeeded",
			ApiUrl:       "/api/v1/builds/214953805",
			JobName:      "quickstart-smoke",
			PipelineId:   24,
			PipelineName: "concourse",
			StartTime:    time.Date(2023, time.October, 16, 14, 7, 45, 0, time.UTC),
			EndTime:      time.Date(2023, time.October, 16, 14, 11, 55, 0, time.UTC),
			CreatedBy:    "",
		},
	}
	assert.DeepEqual(t, have, want, "ListPipelineBuilds")
}

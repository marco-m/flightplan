// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse_test

import (
	"context"
	"testing"
	"time"

	"github.com/marco-m/rosina/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/marco-m/flightplan/internal/testhelpers"
	"github.com/marco-m/flightplan/pkg/concourse"
)

var unixEpoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestClient_ListPipelineBuilds(t *testing.T) {
	// Arrange recorder.
	rec, err := recorder.New("testdata/list-pipeline-builds-short",
		recorder.WithHook(testhelpers.BodyFormatHook, recorder.AfterCaptureHook),
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
	const limit = 2
	concourseClient, err := concourse.NewClient(concourse.ClientArgs{
		ServerURL:  "https://ci.concourse-ci.org",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	have, err := concourseClient.ListPipelineBuilds(ctx, team, pipeline, limit)

	// Assert.
	assert.NoError(t, err, "ListPipelineBuilds")
	want := []concourse.Build{
		{
			ID:           562340543,
			TeamName:     "main",
			Name:         "2668",
			Status:       concourse.StatusSucceeded,
			APIURL:       "/api/v1/builds/562340543",
			JobName:      "bosh-check-props",
			PipelineID:   24,
			PipelineName: "concourse",
			StartTime:    time.Date(2026, time.May, 24, 23, 29, 58, 0, time.UTC),
			EndTime:      time.Date(2026, time.May, 24, 23, 30, 46, 0, time.UTC),
			ReapTime:     unixEpoch,
		},
		{
			ID:           562340123,
			TeamName:     "main",
			Name:         "1627",
			Status:       concourse.StatusSucceeded,
			APIURL:       "/api/v1/builds/562340123",
			JobName:      "k8s-check-helm-params",
			PipelineID:   24,
			PipelineName: "concourse",
			StartTime:    time.Date(2026, time.May, 24, 23, 25, 38, 0, time.UTC),
			EndTime:      time.Date(2026, time.May, 24, 23, 27, 32, 0, time.UTC),
			ReapTime:     unixEpoch,
		},
	}
	assert.DeepEqual(t, have, want, "ListPipelineBuilds")
}

// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// BuildStatus taken from https://github.com/concourse/concourse/blob/master/atc/build.go
// Go doesn't have real enums; this is the simplest less bad we can do.
// Other approaches are possible, but I am not convinced of the trade-offs.
type BuildStatus string

const (
	StatusStarted   BuildStatus = "started"
	StatusPending   BuildStatus = "pending"
	StatusSucceeded BuildStatus = "succeeded"
	StatusFailed    BuildStatus = "failed"
	StatusErrored   BuildStatus = "errored"
	StatusAborted   BuildStatus = "aborted"
	//
	StatusInvalid BuildStatus = "invalid"
)

func ParseStatus(s string) BuildStatus {
	switch s {
	case "started":
		return StatusStarted
	case "pending":
		return StatusPending
	case "succeeded":
		return StatusSucceeded
	case "failed":
		return StatusFailed
	case "errored":
		return StatusErrored
	case "aborted":
		return StatusAborted
	default:
		return StatusInvalid
	}
}

func (status BuildStatus) String() string {
	return string(status)
}

type Build struct {
	ID                   int           `json:"id"`
	TeamName             string        `json:"team_name"`
	Name                 string        `json:"name"`
	Status               BuildStatus   `json:"status"`
	APIURL               string        `json:"api_url"`
	Comment              string        `json:"comment,omitempty"`
	JobName              string        `json:"job_name,omitempty"`
	ResourceName         string        `json:"resource_name,omitempty"`
	PipelineID           int           `json:"pipeline_id,omitempty"`
	PipelineName         string        `json:"pipeline_name,omitempty"`
	PipelineInstanceVars InstanceVars  `json:"pipeline_instance_vars,omitempty"`
	StartTime            time.Time     `json:"start_time,omitempty"` // on the wire: int64
	EndTime              time.Time     `json:"end_time,omitempty"`   // on the wire: int64
	ReapTime             time.Time     `json:"reap_time,omitempty"`  // on the wire: int64
	RerunNumber          int           `json:"rerun_number,omitempty"`
	RerunOf              *RerunOfBuild `json:"rerun_of,omitempty"`
	CreatedBy            *string       `json:"created_by,omitempty"`
}

type RerunOfBuild struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UnmarshalJSON is normally not needed; function [json.Unmarshal] already knows
// what to do. In this case, we override the default behavior because we want to
// transparently parse time-related fields, which on the wire are encoded as int64
// (Unix time), to the Go-native [time.Time].
//
// For a detailed explanation, see
// - https://eli.thegreenplace.net/2019/go-json-cookbook/
// - https://choly.ca/post/go-json-marshalling/
func (bld *Build) UnmarshalJSON(data []byte) error {
	type Alias Build // Avoid infinite loop calling UnmarshalJSON.
	aux := &struct {
		StartTime int64 `json:"start_time,omitempty"`
		EndTime   int64 `json:"end_time,omitempty"`
		ReapTime  int64 `json:"reap_time,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(bld),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	bld.StartTime = time.Unix(aux.StartTime, 0).UTC()
	bld.EndTime = time.Unix(aux.EndTime, 0).UTC()
	bld.ReapTime = time.Unix(aux.ReapTime, 0).UTC()

	return nil
}

func (bld *Build) MarshalJSON() ([]byte, error) {
	type Alias Build // Avoid infinite loop calling MarshalJSON.
	return json.Marshal(&struct {
		StartTime int64 `json:"start_time,omitempty"`
		EndTime   int64 `json:"end_time,omitempty"`
		ReapTime  int64 `json:"reap_time,omitempty"`
		*Alias
	}{
		StartTime: bld.StartTime.Unix(),
		EndTime:   bld.EndTime.Unix(),
		ReapTime:  bld.ReapTime.Unix(),
		Alias:     (*Alias)(bld),
	})
}

// ListPipelineBuilds returns the last builds for pipeline 'pipelineName'
// in team 'teamName'.
//
// From https://github.com/concourse/concourse/blob/master/atc/routes.go
// Path: /api/v1/teams/:team_name/pipelines/:pipeline_name/builds?limit=N
// Method: GET
func (cl *Client) ListPipelineBuilds(ctx context.Context, team, pipeline string,
	limit int,
) ([]Build, error) {
	uri := cl.serverURL.JoinPath("/api/v1/teams/", team, "/pipelines", pipeline, "/builds")
	values := url.Values{}
	values.Set("limit", strconv.Itoa(limit))
	body, err := get(ctx, cl.httpClient, uri.String(), values)
	if err != nil {
		return nil, fmt.Errorf("ListPipelineBuilds: url: %s: %s", uri, err)
	}
	var builds []Build
	if err := json.Unmarshal(body, &builds); err != nil {
		return nil, err
	}
	return builds, nil
}

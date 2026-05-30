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

func TestClient_PasswordLogin(t *testing.T) {
	// Arrange recorder.
	rec, teardown := testhelpers.SetupRecorder(t, "testdata/password-login")
	t.Cleanup(func() { teardown(t) })

	// Arrange SUT.
	client, err := concourse.NewClient(concourse.ClientArgs{
		ServerURL:  "http://localhost:8080",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	const username = "main"
	const password = "main"
	err = client.PasswordLogin(ctx, username, password)

	// Assert.
	assert.NoError(t, err, "PasswordLogin")
	// NOTE verification of the fact that the token is set in Client will be done
	// indirectly by subsequent tests.
}

func TestClient_GetUserInfoWithoutLoginFails(t *testing.T) {
	// Arrange recorder.
	rec, teardown := testhelpers.SetupRecorder(t, "testdata/get-userinfo-without-login")
	t.Cleanup(func() { teardown(t) })

	// Arrange SUT.
	client, err := concourse.NewClient(concourse.ClientArgs{
		ServerURL:  "http://localhost:8080",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	_, err = client.GetUserInfo(ctx)

	// Assert.
	assert.ErrorIs(t, err, concourse.ErrUnauthorized, "GetUserInfo")
}

func TestClient_GetUserInfoAfterLoginPasses(t *testing.T) {
	// Arrange recorder.
	rec, teardown := testhelpers.SetupRecorder(t, "testdata/get-userinfo-after-login")
	t.Cleanup(func() { teardown(t) })

	// Arrange SUT.
	client, err := concourse.NewClient(concourse.ClientArgs{
		ServerURL:  "http://localhost:1234",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	const username = "main"
	const password = "main"
	err = client.PasswordLogin(ctx, username, password)
	assert.NoError(t, err, "PasswordLogin")
	have, err := client.GetUserInfo(ctx)

	// Assert.
	assert.NoError(t, err, "GetUserInfo")
	want := concourse.UserInfo{
		Sub:           "CgRtYWluEgVsb2NhbA",
		Name:          "main",
		UserId:        "main",
		UserName:      "",
		Email:         "main",
		IsAdmin:       true,
		IsSystem:      false,
		Teams:         map[string][]string{"main": {"owner"}},
		Connector:     "local",
		DisplayUserId: "main",
	}
	assert.DeepEqual(t, have, want, "GetUserInfo")
}

func TestClient_ListTeamsWithoutLoginFails(t *testing.T) {
	// Arrange recorder.
	rec, teardown := testhelpers.SetupRecorder(t, "testdata/list-teams-without-login")
	t.Cleanup(func() { teardown(t) })

	// Arrange SUT.
	client, err := concourse.NewClient(concourse.ClientArgs{
		ServerURL:  "http://localhost:8080",
		HttpClient: rec.GetDefaultClient(),
	})
	assert.NoError(t, err, "concourse.NewClient")
	ctx := context.Background()

	// Act.
	_, err = client.ListTeams(ctx)

	// Assert.
	assert.ErrorIs(t, err, concourse.ErrUnauthorized, "ListTeams")
}

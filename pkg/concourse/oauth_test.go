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

func TestClient_GetOauthTokenFromPassword(t *testing.T) {
	// Arrange recorder.
	rec, teardown := testhelpers.SetupRecorder(t, "testdata/oauth-password")
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
	have, err := client.GetOauthTokenFromPassword(ctx, username, password)

	// Assert.
	assert.NoError(t, err, "GetOauthTokenFromPassword")
	want := "bU6YRVVJYUe3DpOeQ6sM5UQGa2meJRdqAAAAAA"
	assert.Equal(t, have.AccessToken, want, "GetOauthTokenFromPassword")
}

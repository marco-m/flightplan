// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Info struct {
	Version       string          `json:"version"`
	WorkerVersion string          `json:"worker_version"`
	FeatureFlags  map[string]bool `json:"feature_flags"`
	ExternalURL   string          `json:"external_url"`
	ClusterName   string          `json:"cluster_name"`
}

// GetInfo returns [Info] about the Concourse server.
//
// From https://github.com/concourse/concourse/blob/master/atc/routes.go
// Path: /api/v1/info
// Method: GET
func (cl *Client) GetInfo(ctx context.Context) (Info, error) {
	const op = "GetInfo"
	uri := cl.serverURL.JoinPath("/api/v1/info")

	resp, err := get(ctx, cl.httpClient, uri.String(), nil, nil)
	if err != nil {
		return Info{}, fmt.Errorf("%s: url: %s: %s", op, uri, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return Info{}, fmt.Errorf("%s: %s", op, responseError(resp.StatusCode, body, err))
	}
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return Info{}, fmt.Errorf("%s: %s", op, err)
	}
	return info, nil
}

// PasswordLogin performs a username and password "login", by asking for an Oauth access
// token and storing it in Client, so that subsequent operations will be authenticated
// with the granted token.
//
// Written by observing the traffic of "fly login -p password -u username -t team"
//
// Flow:
//
//	>> GET http://localhost/api/v1/info
//	<< 200 OK
//
//	>> POST http://localhost/sky/issuer/token
//	<< 200 OK
//
//	>> GET http://localhost/api/v1/user
//	<< 200 OK
//
//	>> GET http://localhost/api/v1/teams
//	<< 200 OK
func (cl *Client) PasswordLogin(ctx context.Context, username, password string) error {
	const op = "PasswordLogin"

	// This API call is done by "fly login". Fly looks at the reported version and
	// decideds wether to continue or not. It is not needed in order to get authenticated.
	_, err := cl.GetInfo(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// The needed API call.
	grant, err := cl.GetOauthTokenFromPassword(ctx, username, password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	cl.token = grant.AccessToken

	return nil
}

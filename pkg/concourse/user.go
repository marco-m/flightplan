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

	//
	// The following API calls are done by "fly login".
	// They are not needed in orded to remain authenticated.
	//

	_, err = cl.GetUserInfo(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type UserInfo struct {
	Sub           string              `json:"sub"`
	Name          string              `json:"name"`
	UserId        string              `json:"user_id"`
	UserName      string              `json:"user_name"`
	Email         string              `json:"email"`
	IsAdmin       bool                `json:"is_admin"`
	IsSystem      bool                `json:"is_system"`
	Teams         map[string][]string `json:"teams"`
	Connector     string              `json:"connector"`
	DisplayUserId string              `json:"display_user_id"`
}

// GetUserInfo returns [UserInfo] about the user associated to the bearer token of the
// request.
//
// From https://github.com/concourse/concourse/blob/master/atc/routes.go
// Path: /api/v1/user
// Method: GET
func (cl *Client) GetUserInfo(ctx context.Context) (UserInfo, error) {
	const op = "GetUserInfo"
	uri := cl.serverURL.JoinPath("/api/v1/user")
	header := make(http.Header)
	header.Set("Authorization", "Bearer "+cl.token)

	resp, err := get(ctx, cl.httpClient, uri.String(), header, nil)
	if err != nil {
		return UserInfo{}, fmt.Errorf("%s: %s", op, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var userInfo UserInfo
		err = json.NewDecoder(resp.Body).Decode(&userInfo)
		if err != nil {
			return UserInfo{}, fmt.Errorf("%s: %s", op, err)
		}
		return userInfo, nil
	case http.StatusUnauthorized:
		body, _ := io.ReadAll(resp.Body)
		return UserInfo{}, fmt.Errorf("%s: %w (%s)", op, ErrUnauthorized, string(body))
	default:
		body, err := io.ReadAll(resp.Body)
		return UserInfo{}, fmt.Errorf("%s: %s", op, responseError(resp.StatusCode, body, err))

	}
}

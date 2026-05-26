// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TokenGrant struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

// Request a OAuth password-based grant.
// NOTE This function is low-level and stateless. Do not use unless you have a specific
// reason. Use instead the stateful [Client.PasswordLogin]
// Returns the access token (to be used in header Authorization as Bearer).
// See https://oauth.net/2/grant-types/password/
//
// Path http://localhost:8080/sky/issuer/token
// Method POST
// Not present in atc/routes.go, see instead skymarshall/token.
//
// Written by observing the traffic of "fly login -p password -u username"
func (cl *Client) GetOauthTokenFromPassword(ctx context.Context, username, password string,
) (TokenGrant, error) {
	const op = "GetOauthTokenFromPassword"
	uri := cl.serverURL.JoinPath("/sky/issuer/token")

	header := make(http.Header)
	// From concourse/fly/commands/login.go
	const clientID = "fly"
	const clientSecret = "Zmx5"
	const src = clientID + ":" + clientSecret
	// Equivalent of Request.SetBasicAuth(url.QueryEscape(clientID), url.QueryEscape(clientSecret))
	credentials := base64.StdEncoding.EncodeToString([]byte(src))
	header.Set("Authorization", "Basic "+credentials)

	values := make(url.Values)
	values.Set("grant_type", "password")
	values.Set("scope", "openid profile email federated:id groups")
	values.Set("username", username)
	values.Set("password", password)

	resp, err := postForm(ctx, cl.httpClient, uri.String(), header, values)
	if err != nil {
		return TokenGrant{}, fmt.Errorf("%s: url: %s: %s", op, uri, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return TokenGrant{}, fmt.Errorf("%s: %s", op, responseError(resp.StatusCode, body, err))
	}
	var grant TokenGrant
	if err := json.NewDecoder(resp.Body).Decode(&grant); err != nil {
		return TokenGrant{}, fmt.Errorf("%s: %s", op, err)
	}
	return grant, nil
}

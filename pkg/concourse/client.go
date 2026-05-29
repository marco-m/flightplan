// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is a minimal client for the Concourse HTTP API.
// Use [NewClient] to instantiate.
type Client struct {
	serverURL  *url.URL
	httpClient *http.Client
}

// Arguments to [NewClient].
type ClientArgs struct {
	ServerURL  string       // Mandatory.
	HttpClient *http.Client // Optional; overridable in tests.
}

// NewClient instantiates a Concourse [Client].
func NewClient(args ClientArgs) (*Client, error) {
	client := &Client{}

	uri, err := url.Parse(args.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("concourse.NewClient: %s", err)
	}
	client.serverURL = uri

	client.httpClient = args.HttpClient
	if client.httpClient == nil {
		client.httpClient = &http.Client{}
	}

	return client, nil
}

func get(ctx context.Context, hclient *http.Client, theUrl string, values url.Values,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, theUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("get: new request: %s", err)
	}
	req.URL.RawQuery = values.Encode()
	resp, err := hclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: do: %s", err)
	}
	defer resp.Body.Close()

	body, errBody := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if errBody != nil {
			return nil, fmt.Errorf("get: status code: %d (%s)", resp.StatusCode, errBody)
		}
		return nil, fmt.Errorf("get: status code: %d (%s)", resp.StatusCode, string(body))
	}
	if errBody != nil {
		return body, fmt.Errorf("get: read body: %s", errBody)
	}

	return body, nil
}

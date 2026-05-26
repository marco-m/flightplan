// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

// Client is a minimal client for the Concourse HTTP API.
// Use [NewClient] to instantiate.
type Client struct {
	serverURL  *url.URL
	httpClient *http.Client
	token      string `json:"-"` // Sensitive!
}

// LogValue redacts sensitive fields when Client is passed to package slog.
func (cl Client) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Server", cl.serverURL.String()),
		slog.String("HttpClient", fmt.Sprint(cl.httpClient)),
		slog.String("token", "[redacted]"),
	)
}

// String redacts sensitive fields when Client is printed with a string verb.
// See https://pkg.go.dev/fmt#Stringer.
func (cl Client) String() string {
	return cl.LogValue().String()
}

// GoString redacts sensitive fields when Client is printed with verb "%#v".
// See https://pkg.go.dev/fmt#GoStringer
func (cl Client) GoString() string {
	return cl.String()
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

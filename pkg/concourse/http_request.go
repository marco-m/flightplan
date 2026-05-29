// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package concourse

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func get(
	ctx context.Context,
	hclient *http.Client,
	url string,
	header http.Header,
	values url.Values,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get: new request: %s", err)
	}
	req.Header.Set("User-Agent", "ghibli")
	for k, v := range header {
		req.Header[k] = v
	}
	req.URL.RawQuery = values.Encode()
	return hclient.Do(req)
}

func postForm(
	ctx context.Context,
	hclient *http.Client,
	url string,
	header http.Header,
	data url.Values,
) (*http.Response, error) {
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return post(ctx, hclient, url, header, strings.NewReader(data.Encode()))
}

func post(
	ctx context.Context,
	hclient *http.Client,
	url string,
	header http.Header,
	body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("post: new request: %s", err)
	}
	req.Header.Set("User-Agent", "ghibli")
	for k, v := range header {
		req.Header[k] = v
	}
	return hclient.Do(req)
}

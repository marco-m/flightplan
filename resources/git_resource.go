// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

import (
	"errors"
	"fmt"
)

var ErrGitParamsEmpty = errors.New("Git: field Params cannot be empty")

//

var _ Source = (*Git)(nil)

// Git gets and puts commits in a git repository.
// See https://github.com/concourse/git-resource for details.
type Git struct {
	// Required. The location of the repository.
	Uri string `json:"uri"`
	// Optional. The branch to track.
	Branch string `json:"branch,omitzero"`
	// Optional. If specified (as a list of glob patterns), only changes to the
	// specified files will yield new versions from `check`.
	Paths []string `json:"paths,omitzero"`
	// Optional. The HTTPS proxy that will be used to tunnel SSH-based git commands.
	HttpsTunnel GitTunnel `json:"https_tunnel,omitzero"`
}

type GitTunnel struct {
	// Required. The host name or IP of the proxy server.
	ProxyHost string `json:"proxy_host,omitzero"`
	//  Required. The proxy server's listening port.
	ProxyPort int `json:"proxy_port,omitzero"`
	// Optional. If the proxy requires authentication, use this username.
	ProxyUser string `json:"proxy_user,omitzero"`
	// Optional. If the proxy requires authenticate, use this password.
	ProxyPassword string `json:"proxy_password,omitzero"`
}

func (git Git) IsSource() {}

func (git Git) Type() string { return "git" }

func (git Git) Validate() error {
	if git.Uri == "" {
		return fmt.Errorf("Git.Uri cannot be empty")
	}
	return nil
}

func (git Git) ValidateGet(params GetParams) error {
	if params == nil {
		return ErrGitParamsEmpty
	}
	return params.(GitGetParams).Validate()
}

func (git Git) ValidatePut(params PutParams) error {
	if params == nil {
		return ErrGitParamsEmpty
	}
	return params.(GitPutParams).Validate()
}

type GitGetParams struct{}

func (GitGetParams) IsGetParams() {}

func (GitGetParams) Validate() error {
	return nil
}

type GitPutParams struct{}

func (GitPutParams) IsPutParams() {}

func (GitPutParams) Validate() error {
	return nil
}

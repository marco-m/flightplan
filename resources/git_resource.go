// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package resources

import (
	"errors"
	"fmt"
	"reflect"
)

//lint:file-ignore ST1005 Capitalized error strings are OK in this case.
var (
	ErrGitMissingUri = errors.New("Git.Uri cannot be empty")
	//
	ErrGitParamsWrongType            = errors.New("Params: wrong type")
	ErrGitGetParamsEmpty             = errors.New("Get field Params provided but empty (remove or fill)")
	ErrGitPutParamsEmpty             = errors.New("Put field Params cannot be empty")
	ErrGitPutParamsMissingRepository = errors.New("Params: field Repository cannot be empty")
	ErrGitPutParamsWrongRepoType     = errors.New("Params: field Repository: wrong type")
)

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
		return ErrGitMissingUri
	}
	return nil
}

func (git Git) ValidateGet(params GetParams) error {
	if params == nil {
		// All Git get params are optional.
		return nil
	}
	if p, ok := params.(GitGetParams); ok {
		return p.Validate()
	}
	return fmt.Errorf("%w: have: %T; want: %T",
		ErrGitParamsWrongType, params, GitGetParams{})
}

func (git Git) ValidatePut(params PutParams) error {
	if params == nil {
		// Some Git put params are required.
		return ErrGitPutParamsEmpty
	}
	if p, ok := params.(GitPutParams); ok {
		return p.Validate()
	}
	return fmt.Errorf("%w: have: %T; want: %T",
		ErrGitParamsWrongType, params, GitPutParams{})
}

func (git Git) SentinelDir() string {
	return ".git"
}

type GitGetParams struct {
	// Optional. If a positive integer is given, shallow clone the repository using the
	// --depth option.
	Depth int `json:"depth,omitzero"`
	// Optional. Additional branches to fetch and make available as a local branch.
	Fetch []string `json:"fetch,omitzero"`
	// Optional	If true all tags in the repository will be fetched. If false no tags
	// will be fetched. Overrides the fetch_tags source configuration.
	FetchTags bool `json:"fetch_tags,omitzero"`
	// Optional. Enable debugging output. Secrets may not be redacted.
	Debug bool `json:"debug,omitzero"`
}

func (GitGetParams) IsGetParams() {}

func (gg GitGetParams) Validate() error {
	// Since GitGetParams is not comparable because it contains []string, we must use
	// reflection.
	if reflect.ValueOf(gg).IsZero() {
		return ErrGitGetParamsEmpty
	}
	return nil
}

type GitPutParams struct {
	// Required. The [Handle] of a [Git] repo to push to the source [Git] repo.
	Repository *Handle `json:"repository"`
	// Optional. Enable debugging output. Secrets may not be redacted.
	Debug bool `json:"debug,omitzero"`
}

func (GitPutParams) IsPutParams() {}

func (gp GitPutParams) Validate() error {
	if gp.Repository == nil {
		return ErrGitPutParamsMissingRepository
	}

	if have, want := gp.Repository.Type, (Git{}).Type(); have != want {
		return fmt.Errorf("%w: have: %q; want: %q", ErrGitPutParamsWrongRepoType, have, want)
	}
	return nil
}

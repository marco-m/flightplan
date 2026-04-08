// Copyright 2026 The Flightplan Authors. All rights reserved.
// Use of this source code is governed by the MIT license; see the LICENSE file.

package flightplan

// GitSource gets and puts commits in a git repository.
// See https://github.com/concourse/git-resource for details.
type GitSource struct {
	// Required. The location of the repository.
	Uri string `json:"uri,omitzero"`
	// Optional. The branch to track.
	Branch string `json:"branch,omitzero"`
	// Optional. If specified (as a list of glob patterns), only changes to the
	// specified files will yield new versions from `check`.
	Paths []string `json:"paths,omitzero"`
	// Optional. The HTTPS proxy that will be used to tunnel SSH-based git commands.
	HttpsTunnel GitSourceTunnel `json:"https_tunnel,omitzero"`
}

type GitSourceTunnel struct {
	// Required. The host name or IP of the proxy server.
	ProxyHost string `json:"proxy_host,omitzero"`
	//  Required. The proxy server's listening port.
	ProxyPort int `json:"proxy_port,omitzero"`
	// Optional. If the proxy requires authentication, use this username.
	ProxyUser string `json:"proxy_user,omitzero"`
	// Optional. If the proxy requires authenticate, use this password.
	ProxyPassword string `json:"proxy_password,omitzero"`
}

var _ Source = (*GitSource)(nil)

func (src GitSource) Source() {}

func (src GitSource) Type() string { return "git" }

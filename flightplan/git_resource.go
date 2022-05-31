package flightplan

// This approach cannot work :-/ Because I can do:
// bar.Name = "caca"
// but not (what I would like to do)
//
//	bar := fp.GitResource2{
//			Name: "caca",
//
// I must do instead the absurd
//
//	bar := fp.GitResource2{
//	  ResourceObject: fp.ResourceObject{
//	    Name: "caca"
type GitResource2 struct {
	Resource
	GitSource
}

var _ Source = (*GitSource)(nil)

type GitSource struct {
	Uri         string          `json:"uri,omitzero"`
	Branch      string          `json:"branch,omitzero"`
	Paths       []string        `json:"paths,omitzero"`
	HttpsTunnel GitSourceTunnel `json:"https_tunnel,omitzero"`
}

func (src GitSource) Source() {}

func (src GitSource) Type() string { return "git" }

type GitSourceTunnel struct {
	ProxyHost     string `json:"proxy_host,omitzero"`
	ProxyPort     int    `json:"proxy_port,omitzero"`
	ProxyUser     string `json:"proxy_user,omitzero"`
	ProxyPassword string `json:"proxy_password,omitzero"`
}

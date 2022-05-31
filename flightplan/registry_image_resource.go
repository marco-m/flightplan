package flightplan

type RegistryImage struct{}

var _ Source = (*RegistryImageSource)(nil)

// See https://github.com/concourse/registry-image-resource
type RegistryImageSource struct {
	// Required

	// The URI of the image repository, e.g. alpine or ghcr.io/package/image. Defaults to
	// checking docker.io if no hostname is provided in the URI.
	Repository string `json:"repository"`

	// Optional
}

func (src RegistryImageSource) Source() {}

func (src RegistryImageSource) Type() string { return "registry-image" }

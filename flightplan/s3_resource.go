package flightplan

type S3Resource struct {
	Name_ string `json:"name"`
}

// var _ Resource = (*S3Resource)(nil)

func NewS3Resource(name string) *S3Resource {
	return &S3Resource{Name_: name}
}

func (s3r S3Resource) Name() string {
	return s3r.Name_
}

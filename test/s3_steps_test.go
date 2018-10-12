package fujilane

import (
	"github.com/nerde/fuji-lane-back/flservices"
)

type fakeS3 struct {
	flservices.S3Service
}

// DeleteFile overrides the default behaviour, which would actually trigger a call to S3
func (f *fakeS3) DeleteFile(_ string) error {
	return nil
}

func newFakeS3(original flservices.S3Service) *fakeS3 {
	return &fakeS3{original}
}

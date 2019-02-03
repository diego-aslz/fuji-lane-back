package flservices

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nerde/fuji-lane-back/flconfig"
)

// S3 is an abstraction that provides easy-to-use functions to access Amazon S3
type S3 struct {
	service *s3.S3
}

// GenerateURLToUploadPublicFile generates a pre-signed URL to upload public files to Amazon S3. To test it, you can
// use the following CURL command:
//
//   curl -H 'x-amz-acl: public-read' -X PUT -F file=@<path/to/file> <presigned_url>
func (s S3) GenerateURLToUploadPublicFile(key, cType string, size int) (string, error) {
	req, _ := s.service.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(flconfig.Config.AWSBucket),
		Key:           aws.String(key),
		ACL:           aws.String(s3.ObjectCannedACLPublicRead),
		ContentLength: aws.Int64(int64(size)),
		ContentType:   aws.String(cType),
	})

	return req.Presign(1 * time.Hour)
}

// DeleteFile removes a file stored in S3
func (s S3) DeleteFile(key string) error {
	_, err := s.service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(flconfig.Config.AWSBucket),
		Key:    aws.String(key),
	})

	return err
}

// S3Service provides the functionality the system needs from AWS S3
type S3Service interface {
	GenerateURLToUploadPublicFile(string, string, int) (string, error)
	DeleteFile(string) error
}

// NewS3 returns a new S3Service based on the env vars configuration
func NewS3() (S3Service, error) {
	config := &aws.Config{
		Region:      aws.String(flconfig.Config.AWSRegion),
		Credentials: credentials.NewEnvCredentials(),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &S3{s3.New(sess)}, nil
}

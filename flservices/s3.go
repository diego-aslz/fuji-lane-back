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
	config *aws.Config
}

// GenerateURLToUploadPublicFile generates a pre-signed URL to upload public files to Amazon S3. To test it, you can
// use the following CURL command:
//
//   curl -H 'x-amz-acl: public-read' -X PUT -F file=@<path/to/file> <presigned_url>
func (s S3) GenerateURLToUploadPublicFile(key string) (string, error) {
	sess, err := session.NewSession(s.config)

	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(flconfig.Config.AWSBucket),
		Key:    aws.String("public/" + key),
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	})

	return req.Presign(1 * time.Hour)
}

// NewS3 creates a new instance for the S3 service with configuration
func NewS3() *S3 {
	return &S3{&aws.Config{
		Region:      aws.String(flconfig.Config.AWSRegion),
		Credentials: credentials.NewEnvCredentials(),
	}}
}

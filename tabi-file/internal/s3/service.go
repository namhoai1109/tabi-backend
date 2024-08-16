package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var Metadata = map[string]*string{
	"x-amz-storage-class": aws.String("STANDARD_IA"),
}

// Config represents the configuration
type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Debug           bool
}

// New initialize s3 service
func New(cfg Config) *Service {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		panic(err)
	}
	return &Service{
		cfg:       cfg,
		s3:        s3.New(session),
		s3manager: s3manager.NewUploader(session),
	}
}

// Service represents the s3 service
type Service struct {
	cfg       Config
	s3        *s3.S3
	s3manager *s3manager.Uploader
}

package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	timeOut     = 15
	timeDefault = 15
)

// PreparePresignedURL to prepare the S3 request so a signature can be generated
func (s *Service) PreparePresignedURL(key, bucketName string, expireTime int) (*string, error) {
	putReq, _ := s.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if expireTime == 0 {
		expireTime = timeDefault
	}

	// Create the pre-signed url with an expiry
	presignedURL, err := putReq.Presign(time.Duration(expireTime) * time.Minute)
	if err != nil {
		return nil, err
	}

	// Display the pre-signed url
	return &presignedURL, nil
}

// GetSignedURL to get object presign
func (s *Service) GetSignedURL(key, bucketName string, expireTime int) (*string, error) {
	resp, _ := s.s3.GetObjectRequest(
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
	if expireTime == 0 {
		expireTime = timeDefault
	}
	url, err := resp.Presign(time.Duration(expireTime) * time.Minute)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

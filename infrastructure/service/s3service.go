package service

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

// S3Service ...
type S3Service struct {
	s3         *s3.S3
	s3Uploader *s3manager.Uploader
	BucketName string
	Prefix     string
}

// NewS3Service ...
func NewS3Service(opt *setting.BackupStorage) (S3Service, error) {
	creds := credentials.NewStaticCredentials(opt.AccessKey, opt.SecretKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(opt.Region),
		Endpoint:    aws.String(opt.CustomEndpoint),
	}))

	s3Client := s3.New(sess)
	s3Uploader := s3manager.NewUploader(sess)

	return S3Service{
		s3:         s3Client,
		s3Uploader: s3Uploader,
		BucketName: opt.BucketName,
	}, nil
}

// Add ...
func (s *S3Service) Add(key string, file []byte) error {
	param := &s3manager.UploadInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	}
	_, err := s.s3Uploader.Upload(param)
	if err != nil {
		panic(err)
	}

	return nil
}

// Delete ...
func (s *S3Service) Delete(key string) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}

	_, err := s.s3.DeleteObject(params)
	if err != nil {
		return err
	}
	return nil
}

// List ...
func (s *S3Service) List() ([]*s3.Object, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(s.Prefix),
	}

	output, err := s.s3.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}
	return output.Contents, nil
}

// TouchBucket ...
func (s *S3Service) TouchBucket() error {
	if s.ExistBucket() {
		return nil
	}

	return s.MakeBucket()
}

// ExistBucket ...
func (s *S3Service) ExistBucket() bool {
	param := &s3.HeadBucketInput{
		Bucket: aws.String(s.BucketName),
	}
	_, err := s.s3.HeadBucket(param)
	if err != nil {
		return false
	}
	return true
}

// MakeBucket ...
func (s *S3Service) MakeBucket() error {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(s.BucketName),
	}

	_, err := s.s3.CreateBucket(params)
	if err != nil {
		return err
	}
	return nil
}

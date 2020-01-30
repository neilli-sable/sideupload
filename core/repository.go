package core

import "github.com/aws/aws-sdk-go/service/s3"

// S3Repository ...
type S3Repository interface {
	Add(key string, body []byte) error
	List() ([]*s3.Object, error)
	Delete(file string) error
}

// ArchiveRepository ...
type ArchiveRepository interface {
	Compress(path string) (body []byte, err error)
}

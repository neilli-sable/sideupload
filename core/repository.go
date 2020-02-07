package core

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

//go:generate mockgen -source=./repository.go -destination=mock_core/repository_mock.go -package=mock_core

// S3Repository ...
type S3Repository interface {
	Add(key string, body []byte) error
	List() ([]*s3.Object, error)
	Delete(file string) error
	IsOld(*s3.Object) bool
}

// ArchiveRepository ...
type ArchiveRepository interface {
	Timestamp() time.Time
	Compress(path string) (body []byte, err error)
}

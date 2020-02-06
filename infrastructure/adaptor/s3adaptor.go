package adaptor

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/neilli-sable/sideupload/infrastructure/service"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

// NewS3Adaptor コンストラクタ
func NewS3Adaptor(opt *setting.Setting) (*S3Adaptor, error) {
	s3Service, err := service.NewS3Service(&opt.Sideupload.BackupStorage)
	if err != nil {
		return nil, err
	}
	return &S3Adaptor{
		s3Service: s3Service,
		opt:       opt,
	}, nil
}

// S3Adaptor 参照操作アダプター
type S3Adaptor struct {
	opt       *setting.Setting
	s3Service service.S3Service
}

// Add ...
func (ad *S3Adaptor) Add(key string, body []byte) error {
	err := ad.s3Service.TouchBucket()
	if err != nil {
		return err
	}
	return ad.s3Service.Add(key, body)
}

// Delete ...
func (ad *S3Adaptor) Delete(key string) error {
	err := ad.s3Service.TouchBucket()
	if err != nil {
		return err
	}
	return ad.s3Service.Delete(key)
}

// List ...
func (ad *S3Adaptor) List() ([]*s3.Object, error) {
	err := ad.s3Service.TouchBucket()
	if err != nil {
		return nil, err
	}
	return ad.s3Service.List()
}

// IsOld ...
func (ad *S3Adaptor) IsOld(object *s3.Object) bool {
	deleteDate := time.Now().AddDate(0, 0, -ad.opt.Sideupload.StorageDays)
	if object.LastModified.Before(deleteDate) {
		return true
	}
	return false
}

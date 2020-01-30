package adaptor

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/neilli-sable/sideupload/infrastructure/service"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

// NewS3Adaptor コンストラクタ
func NewS3Adaptor(setting *setting.Setting) (*S3Adaptor, error) {
	s3Service, err := service.NewS3Service(&setting.BackupStorage)
	if err != nil {
		return nil, err
	}
	return &S3Adaptor{
		s3Service: s3Service,
	}, nil
}

// S3Adaptor 参照操作アダプター
type S3Adaptor struct {
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

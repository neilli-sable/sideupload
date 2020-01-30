package adaptor

// ユースケースにアダプターを注入します

import (
	"github.com/neilli-sable/sideupload/core"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

// UsecaseFactory ユースケースの生成
func UsecaseFactory(setting *setting.Setting) core.Usecase {
	s3Adaptor, err := NewS3Adaptor(setting)
	if err != nil {
		panic(err)
	}
	archiveAdaptor, err := NewArchiveAdaptor(setting)
	if err != nil {
		panic(err)
	}

	return core.Usecase{
		S3:      s3Adaptor,
		Archive: archiveAdaptor,
	}
}

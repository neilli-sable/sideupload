package adaptor

import (
	"github.com/neilli-sable/sideupload/infrastructure/service"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

// NewArchiveAdaptor コンストラクタ
func NewArchiveAdaptor(setting *setting.Setting) (*ArchiveAdaptor, error) {
	archiveService, err := service.NewArchiveService()
	if err != nil {
		return nil, err
	}
	return &ArchiveAdaptor{
		archiveService: archiveService,
	}, nil
}

// ArchiveAdaptor 参照操作アダプター
type ArchiveAdaptor struct {
	archiveService service.ArchiveService
}

// Compress ...
func (ad *ArchiveAdaptor) Compress(path string) (body []byte, err error) {
	return ad.archiveService.Compress(path)
}

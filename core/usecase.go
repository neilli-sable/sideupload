package core

import (
	"io/ioutil"
	"path/filepath"
)

// Usecase ...
type Usecase struct {
	S3      S3Repository
	Archive ArchiveRepository
}

// Close ...
func (usecase *Usecase) Close() {
	// no operation, currently
}

// GetTargets ...
func (usecase *Usecase) GetTargets(dir string) (paths []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths, nil
}

// CompressTargets ...
func (usecase *Usecase) CompressTargets(paths []string) (archives map[string][]byte, err error) {
	archives = map[string][]byte{}
	timestamp := usecase.Archive.Timestamp().Format("20060102150405")

	for _, p := range paths {
		body, err := usecase.Archive.Compress(p)
		if err != nil {
			return nil, err
		}

		key := filepath.Base(p) + "_" + timestamp + ".tar.gz"

		archives[key] = body
	}
	return archives, nil
}

// UploadArchives ...
func (usecase *Usecase) UploadArchives(archives map[string][]byte) error {
	for key, body := range archives {
		err := usecase.S3.Add(key, body)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteOldArchives ...
func (usecase *Usecase) DeleteOldArchives() (int, error) {
	list, err := usecase.S3.List()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range list {
		if usecase.S3.IsOld(item) {
			usecase.S3.Delete(*item.Key)
			count++
		}
	}
	return count, nil
}

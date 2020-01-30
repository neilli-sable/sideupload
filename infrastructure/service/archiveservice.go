package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// ArchiveService ...
type ArchiveService struct {
}

// NewArchiveService ...
func NewArchiveService() (ArchiveService, error) {
	return ArchiveService{}, nil
}

// Compress ...
func (s *ArchiveService) Compress(path string) (body []byte, err error) {
	var buf bytes.Buffer
	err = compress(path, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func compress(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	if err := tw.Close(); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}
	//
	return nil
}

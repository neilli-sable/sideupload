package service

import (
	"testing"
)

func TestCompress(t *testing.T) {
	path := sampleDir + "storage_object/object1.json"
	service := &ArchiveService{}

	_, err := service.Compress(path)
	if err != nil {
		t.Fatal(err)
	}
}

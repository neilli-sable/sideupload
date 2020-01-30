package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

var testOpt = &setting.BackupStorage{
	CustomEndpoint: "http://localhost:9000",
	Region:         "dummy",
	BucketName:     "/check",
	Prefix:         "",
	AccessKey:      "minio",
	SecretKey:      "minio_good_secretkey",
}

const sampleDir = "../../samples/"

func TestS3AddListDelete(t *testing.T) {
	key := "storage_object/object1.json"

	localPath := sampleDir + key
	file, err := os.Open(localPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	service, err := NewS3Service(testOpt)
	if err != nil {
		panic(err)
	}
	err = service.TouchBucket()
	if err != nil {
		panic(err)
	}

	err = service.Add(key, bytes)
	if err != nil {
		panic(err)
	}

	_, err = service.List()
	if err != nil {
		panic(err)
	}

	err = service.Delete(key)
	if err != nil {
		panic(err)
	}
}

func TestTouchBuchet(t *testing.T) {
	service, err := NewS3Service(testOpt)
	if err != nil {
		panic(err)
	}
	_ = service.TouchBucket()
}

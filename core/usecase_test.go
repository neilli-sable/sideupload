package core

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	gomock "github.com/golang/mock/gomock"
	"github.com/neilli-sable/sideupload/core/mock_core"
)

func TestGetTargets(t *testing.T) {
	cases := []struct {
		label    string
		inDir    string
		outPaths []string
		isError  bool
	}{
		{
			label: "ディレクトリを入力すると、その直下のファイルパスを配列で取得する",
			inDir: "../samples/storage_object",
			outPaths: []string{
				"../samples/storage_object/object1.json",
			},
			isError: false,
		},
	}

	usecase := Usecase{}

	for _, c := range cases {
		actual, err := usecase.GetTargets(c.inDir)
		if c.isError && err == nil {
			t.Log(c.label)
			t.Fatal(fmt.Errorf("error not occured at [%s]", c.label))
		}
		if !c.isError && err != nil {
			t.Log(c.label)
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, c.outPaths) {
			t.Log(c.label)
			t.Fatal(fmt.Errorf("wrong output \nactual:%v\nexpected:%v", actual, c.outPaths))
		}
	}
}

func TestCompressTargets(t *testing.T) {
	baseDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	con := gomock.NewController(t)
	cases := []struct {
		label       string
		mock        func() ArchiveRepository
		inPaths     []string
		outArchives map[string][]byte
		isError     bool
	}{
		{
			label: "対象のファイル・ディレクトリを圧縮する",
			mock: func() ArchiveRepository {
				archiveMock := mock_core.NewMockArchiveRepository(con)
				archiveMock.EXPECT().Timestamp().
					Return(baseDate)
				archiveMock.EXPECT().Compress("../samples/storage_object/object1.json").
					Return([]byte("compressed body"), nil)
				return archiveMock
			},
			inPaths: []string{"../samples/storage_object/object1.json"},
			outArchives: map[string][]byte{
				"object1.json_20200102000000.tar.gz": []byte("compressed body"),
			},
			isError: false,
		},
		{
			label: "Compressでエラーが発生した場合、ユースケースもエラーを返す",
			mock: func() ArchiveRepository {
				archiveMock := mock_core.NewMockArchiveRepository(con)
				archiveMock.EXPECT().Timestamp().
					Return(baseDate)
				archiveMock.EXPECT().Compress("../samples/storage_object/object1.json").
					Return(nil, errors.New("compress error"))
				return archiveMock
			},
			inPaths:     []string{"../samples/storage_object/object1.json"},
			outArchives: nil,
			isError:     true,
		},
	}

	for _, c := range cases {
		usecase := Usecase{
			Archive: c.mock(),
		}
		actual, err := usecase.CompressTargets(c.inPaths)

		if c.isError && err == nil {
			t.Log(c.label)
			t.Fatal(fmt.Errorf("error not occured at [%s]", c.label))
		}
		if !c.isError && err != nil {
			t.Log(c.label)
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, c.outArchives) {
			t.Log(c.label)
			t.Fatal(fmt.Errorf("wrong output \nactual:%v\nexpected:%v", actual, c.outArchives))
		}
	}
}

func TestUploadArchives(t *testing.T) {
	con := gomock.NewController(t)
	cases := []struct {
		label      string
		mock       func() S3Repository
		inArchives map[string][]byte
		isError    bool
	}{
		{
			label: "圧縮ファイルのアップロードを行う",
			mock: func() S3Repository {
				s3Mock := mock_core.NewMockS3Repository(con)
				s3Mock.EXPECT().Add("object1.json_20200102000000.tar.gz", []byte("compressed body")).
					Return(nil)
				return s3Mock
			},
			inArchives: map[string][]byte{
				"object1.json_20200102000000.tar.gz": []byte("compressed body"),
			},
			isError: false,
		},
		{
			label: "アップロードでエラーが発生した場合、ユースケースもエラーを返す",
			mock: func() S3Repository {
				s3Mock := mock_core.NewMockS3Repository(con)
				s3Mock.EXPECT().Add("object1.json_20200102000000.tar.gz", []byte("compressed body")).
					Return(errors.New("upload error"))
				return s3Mock
			},
			inArchives: map[string][]byte{
				"object1.json_20200102000000.tar.gz": []byte("compressed body"),
			},
			isError: true,
		},
	}

	for _, c := range cases {
		usecase := Usecase{
			S3: c.mock(),
		}
		err := usecase.UploadArchives(c.inArchives)

		if c.isError && err == nil {
			t.Log(c.label)
			t.Fatal(fmt.Errorf("error not occured at [%s]", c.label))
		}
		if !c.isError && err != nil {
			t.Log(c.label)
			t.Fatal(err)
		}
	}
}

func TestDeleteOldArchives(t *testing.T) {
	baseDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	baseArchive := &s3.Object{
		Key:          aws.String("object1.json_20200102000000.tar.gz"),
		LastModified: &baseDate,
	}

	con := gomock.NewController(t)
	cases := []struct {
		label      string
		mock       func() S3Repository
		inArchives map[string][]byte
		isError    bool
	}{
		{
			label: "古いファイルの削除を行う",
			mock: func() S3Repository {
				s3Mock := mock_core.NewMockS3Repository(con)
				s3Mock.EXPECT().List().Return([]*s3.Object{
					baseArchive,
				}, nil)
				s3Mock.EXPECT().IsOld(baseArchive).Return((true))
				return s3Mock
			},
			isError: false,
		},
		// {
		// 	label: "アップロードでエラーが発生した場合、ユースケースもエラーを返す",
		// 	mock: func() S3Repository {
		// 		s3Mock := mock_core.NewMockS3Repository(con)
		// 		s3Mock.EXPECT().Add("object1.json_20200102000000.tar.gz", []byte("compressed body")).
		// 			Return(errors.New("upload error"))
		// 		return s3Mock
		// 	},
		// 	inArchives: map[string][]byte{
		// 		"object1.json_20200102000000.tar.gz": []byte("compressed body"),
		// 	},
		// 	isError: true,
		// },
	}

	for _, c := range cases {
		usecase := Usecase{
			S3: c.mock(),
		}
		_, err := usecase.DeleteOldArchives()

		if c.isError && err == nil {
			t.Log(c.label)
			t.Fatal(fmt.Errorf("error not occured at [%s]", c.label))
		}
		if !c.isError && err != nil {
			t.Log(c.label)
			t.Fatal(err)
		}
	}
}

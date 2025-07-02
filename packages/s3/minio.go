package s3

import (
	"context"
	"github.com/amupxm/go-video-concat/internal/logger"

	"github.com/amupxm/go-video-concat/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorageStruct struct {
	Client *minio.Client
}

var ObjectStorage = ObjectStorageStruct{}

func (object *ObjectStorageStruct) Connect(cfg *config.Config) {

	m, err := minio.New(cfg.Storage.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.MinioUser, cfg.Storage.MinioPassword, ""),
		Secure: false,
	})
	if err != nil {
		logger.Log.Fatal(err)
	}
	ObjectStorage.Client = m
}

func InitBuckets(buckets []string) {
	for _, bucket := range buckets {
		s, err := ObjectStorage.Client.BucketExists(context.Background(), bucket)
		if err != nil {
			logger.Log.Fatal(err)
		}
		if !s {
			if err := ObjectStorage.Client.MakeBucket(context.Background(),
				bucket,
				minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false}); err != nil {
				logger.Log.Fatal(err)
			}
		}
	}
}

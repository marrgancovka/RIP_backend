package myminio

import (
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

const (
	BucketName = "spacey"
	MinioHost  = "192.168.1.3:9000"
)

func NewMinioClient(logger *logrus.Logger) *minio.Client {
	accessKeyID := "minio"
	secretAccessKey := "minio124"
	minioClient, err := minio.New(MinioHost, accessKeyID, secretAccessKey, false)

	if err != nil {
		logger.Fatalf("error1: %s", err)
	}

	location := "us-east-1"
	err = minioClient.MakeBucket(BucketName, location)
	if err != nil {
		exists, err2 := minioClient.BucketExists(BucketName)
		if err2 == nil && exists {
			logger.Infof("We already own %s", BucketName)
		} else {
			logger.Fatalf("error2: %s", err2)
		}
	} else {
		logger.Infof("Successfully created %s\n", BucketName)
	}

	return minioClient
}

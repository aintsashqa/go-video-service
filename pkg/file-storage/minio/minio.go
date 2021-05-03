package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	EnableSSL       bool
}

func NewClient(conf Config) (*minio.Client, error) {
	return minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyId, conf.SecretAccessKey, ""),
		Secure: conf.EnableSSL,
	})
}

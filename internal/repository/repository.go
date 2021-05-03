package repository

import (
	"github.com/aintsashqa/go-video-service/internal/config"
	"github.com/aintsashqa/go-video-service/internal/module/video"
	minioimpl "github.com/aintsashqa/go-video-service/internal/repository/minio"
	"github.com/aintsashqa/go-video-service/pkg/file-storage/minio"
)

type Container struct {
	VideoRepository video.Repository
}

func NewContainer(conf config.Config) (*Container, error) {
	client, err := minio.NewClient(minio.Config{
		Endpoint:        conf.MinioConfig.Endpoint,
		AccessKeyId:     conf.MinioConfig.AccessKeyId,
		SecretAccessKey: conf.MinioConfig.SecretAccessKey,
		EnableSSL:       conf.MinioConfig.EnableSSL,
	})
	if err != nil {
		return nil, err
	}

	videoRepository := minioimpl.NewVideoRepository(client)

	container := &Container{
		VideoRepository: videoRepository,
	}
	return container, nil
}

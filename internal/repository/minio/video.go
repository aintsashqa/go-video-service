package minio

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/aintsashqa/go-video-service/internal/module/video"
	"github.com/minio/minio-go/v7"
	uuid "github.com/satori/go.uuid"
)

const (
	VideoRepositoryBucketName                     string = "video-repository-bucket"
	VideoRepositoryUserMetadataOriginalNameKey    string = "Originalname"
	VideoRepositoryUserMetadataCreatedAtKey       string = "Createdat"
	VideoRepositoryUserMetadataCreatedAtFormatKey string = "Createdatformat"
	VideoRepositoryUserMetadataContentType        string = "Contenttype"
)

var (
	_ video.Repository = &repository{}
)

type repository struct {
	Client *minio.Client
}

func NewVideoRepository(client *minio.Client) *repository {
	return &repository{
		Client: client,
	}
}

func (r *repository) Create(ctx context.Context, v *video.Video) error {
	if !isBucketExists(ctx, r.Client, VideoRepositoryBucketName) {
		if err := r.Client.MakeBucket(ctx, VideoRepositoryBucketName, minio.MakeBucketOptions{}); err != nil {
			log.Printf("Repository create error: %s: %s", video.ErrRepositoryCannotCreate, err)
			return video.ErrRepositoryCannotCreate
		}
	}

	reader := bytes.NewBuffer(v.Bytes)
	timeFormat := time.RFC3339
	if _, err := r.Client.PutObject(ctx, VideoRepositoryBucketName, v.Uuid.String(), reader, v.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			VideoRepositoryUserMetadataOriginalNameKey:    v.OriginalName,
			VideoRepositoryUserMetadataCreatedAtKey:       v.CreatedAt.Format(timeFormat),
			VideoRepositoryUserMetadataCreatedAtFormatKey: timeFormat,
			VideoRepositoryUserMetadataContentType:        v.ContentType,
		},
	}); err != nil {
		log.Printf("Repository create error: %s: %s", video.ErrRepositoryCannotCreate, err)
		return video.ErrRepositoryCannotCreate
	}

	return nil
}

func (r *repository) Find(ctx context.Context, id uuid.UUID) (*video.Video, error) {
	bytes, objectInfo, err := findObject(ctx, r.Client, VideoRepositoryBucketName, id.String())
	if err != nil {
		log.Printf("Repository find error: %s: %s", video.ErrRepositoryNotFound, err)
		return nil, video.ErrRepositoryNotFound
	}

	createdAt, err := time.Parse(
		objectInfo.UserMetadata[VideoRepositoryUserMetadataCreatedAtFormatKey],
		objectInfo.UserMetadata[VideoRepositoryUserMetadataCreatedAtKey])
	if err != nil {
		log.Printf("Repository find error: Cannot parse time: %s", err)
		return nil, err
	}

	v := video.Video{
		Uuid:         id,
		OriginalName: objectInfo.UserMetadata[VideoRepositoryUserMetadataOriginalNameKey],
		Size:         objectInfo.Size,
		ContentType:  objectInfo.UserMetadata[VideoRepositoryUserMetadataContentType],
		Bytes:        bytes,
		CreatedAt:    createdAt,
	}

	return &v, nil
}

func (r *repository) Remove(ctx context.Context, id uuid.UUID) error {
	err := r.Client.RemoveObject(ctx, VideoRepositoryBucketName, id.String(), minio.RemoveObjectOptions{})

	errResponse := minio.ToErrorResponse(err)
	if errResponse.Code == "NoSuchKey" {
		log.Printf("Repository remove error: %s: %s", video.ErrRepositoryNotFound, err)
		return video.ErrRepositoryNotFound
	}

	return err
}

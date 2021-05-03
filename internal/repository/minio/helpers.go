package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

func isBucketExists(ctx context.Context, client *minio.Client, bucketName string) bool {
	exists, err := client.BucketExists(ctx, bucketName)
	return err == nil && exists
}

func findObject(ctx context.Context, client *minio.Client, bucketName string, objectName string) ([]byte, minio.ObjectInfo, error) {
	object, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return []byte{}, minio.ObjectInfo{}, err
	}
	defer object.Close()

	objectInfo, err := object.Stat()
	if err != nil {
		return []byte{}, minio.ObjectInfo{}, err
	}

	buffer := make([]byte, objectInfo.Size)
	_, err = object.Read(buffer)
	if err != nil && err != io.EOF {
		return []byte{}, minio.ObjectInfo{}, err
	}

	return buffer, objectInfo, nil
}

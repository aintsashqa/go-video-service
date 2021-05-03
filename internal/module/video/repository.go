package video

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrRepositoryCannotCreate error = errors.New("Cannot create video file")
	ErrRepositoryNotFound     error = errors.New("Video file not found")
)

type Repository interface {
	Create(context.Context, *Video) error
	Find(context.Context, uuid.UUID) (*Video, error)
	Remove(context.Context, uuid.UUID) error
}

package video

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	ContentTypePattern string = "^video/\\w+$"
	ContentChunkSize   int    = 1 << 20
)

var (
	ErrCommandBadUuidValue           error = errors.New("Bad uuid value")
	ErrCommandCannotReadBytes        error = errors.New("Cannot read slice of bytes")
	ErrCommandInvalidContentType     error = errors.New("Invalid content type, expected `video/*`")
	ErrCommandCannotMatchContentType error = errors.New("Cannot match content type")
)

type CreateCommandDeps struct {
	Context    context.Context
	Repository Repository
}

type CreateCommandArgs struct {
	OriginalName string
	Size         int64
	Reader       io.Reader
}

func CreateCommand(deps CreateCommandDeps, args CreateCommandArgs) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(deps.Context, time.Second*5)
	defer cancel()

	bytes, err := ioutil.ReadAll(args.Reader)
	if err != nil {
		return uuid.Nil, ErrCommandCannotReadBytes
	}

	contentType := http.DetectContentType(bytes)
	success, err := regexp.MatchString(ContentTypePattern, contentType)
	if err != nil {
		return uuid.Nil, ErrCommandCannotMatchContentType
	}
	if !success {
		return uuid.Nil, ErrCommandInvalidContentType
	}

	v := Video{
		Uuid:         uuid.NewV4(),
		OriginalName: strings.TrimSpace(args.OriginalName),
		Size:         args.Size,
		ContentType:  contentType,
		Bytes:        bytes,
		CreatedAt:    time.Now(),
	}

	if err := deps.Repository.Create(ctx, &v); err != nil {
		return uuid.Nil, err
	}

	return v.Uuid, nil
}

type FindCommandDeps struct {
	Context    context.Context
	Repository Repository
}

type FindCommandArgs struct {
	Uuid uuid.UUID
}

func FindCommand(deps FindCommandDeps, args FindCommandArgs) (*Video, error) {
	ctx, cancel := context.WithTimeout(deps.Context, time.Second*5)
	defer cancel()

	if args.Uuid == uuid.Nil {
		return nil, ErrCommandBadUuidValue
	}

	v, err := deps.Repository.Find(ctx, args.Uuid)
	return v, err
}

type ReadLimitCommandDeps struct {
	Context    context.Context
	Repository Repository
}

type ReadLimitCommandArgs struct {
	Uuid       uuid.UUID
	StartIndex int
}

func ReadLimitCommand(deps ReadLimitCommandDeps, args ReadLimitCommandArgs) (*Limit, error) {
	ctx, cancel := context.WithTimeout(deps.Context, time.Second*5)
	defer cancel()

	if args.Uuid == uuid.Nil {
		return nil, ErrCommandBadUuidValue
	}

	v, err := deps.Repository.Find(ctx, args.Uuid)
	if err != nil {
		return nil, err
	}

	chunkSize := ContentChunkSize

	endIndex := int(math.Min(
		float64(args.StartIndex+ContentChunkSize),
		float64(v.Size-1),
	))

	if size := endIndex - args.StartIndex; size < chunkSize {
		chunkSize = size
	}

	reader := bytes.NewReader(v.Bytes)
	chunkBytes := make([]byte, chunkSize)

	_, err = reader.ReadAt(chunkBytes, int64(args.StartIndex))
	if err != nil && err != io.EOF {
		return nil, ErrCommandCannotReadBytes
	}

	l := Limit{
		Size:     v.Size,
		Bytes:    chunkBytes,
		EndIndex: endIndex,
	}

	return &l, nil
}

type RemoveCommandDeps struct {
	Context    context.Context
	Repository Repository
}

type RemoveCommandArgs struct {
	Uuid uuid.UUID
}

func RemoveCommand(deps RemoveCommandDeps, args RemoveCommandArgs) error {
	ctx, cancel := context.WithTimeout(deps.Context, time.Second*5)
	defer cancel()

	if args.Uuid == uuid.Nil {
		return ErrCommandBadUuidValue
	}

	err := deps.Repository.Remove(ctx, args.Uuid)
	return err
}

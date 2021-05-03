package video

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Video struct {
	Uuid         uuid.UUID `json:"uuid"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	Bytes        []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Limit struct {
	Size     int64
	Bytes    []byte
	EndIndex int
}

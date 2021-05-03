package http

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/aintsashqa/go-video-service/internal/module/video"
	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"
)

const (
	HeaderXOriginalName string = "X-Original-Name"
)

func HandleVideoError(w http.ResponseWriter, err error) {
	switch err {
	case video.ErrCommandBadUuidValue, video.ErrCommandInvalidContentType:
		ApiErrorResponse(w, http.StatusBadRequest, err)
		break

	case video.ErrRepositoryNotFound:
		ApiErrorResponse(w, http.StatusNotFound, err)
		break

	default:
		ApiErrorResponse(w, http.StatusInternalServerError, err)
	}
}

type UploadVideoResponse struct {
	Uuid  string `json:"uuid"`
	Links []Link `json:"links"`
}

func (h *Handler) UploadVideoAction(w http.ResponseWriter, r *http.Request) {
	originNameHeader := r.Header.Get(HeaderXOriginalName)
	size := r.ContentLength

	deps := video.CreateCommandDeps{
		Context:    r.Context(),
		Repository: h.Repository.VideoRepository,
	}

	args := video.CreateCommandArgs{
		OriginalName: originNameHeader,
		Size:         size,
		Reader:       r.Body,
	}

	result, err := video.CreateCommand(deps, args)
	if err != nil {
		log.Printf("Failure to create video: %s", err)
		HandleVideoError(w, err)
		return
	}

	response := UploadVideoResponse{
		Uuid: result.String(),
		Links: []Link{
			{Method: http.MethodGet, Url: fmt.Sprintf("/api/video/%s", result)},
			{Method: http.MethodGet, Url: fmt.Sprintf("/api/video/%s/stream", result)},
			{Method: http.MethodDelete, Url: fmt.Sprintf("/api/video/%s", result)},
		},
	}

	ApiResponse(w, http.StatusCreated, response)
}

type VideoResponse struct {
	*video.Video
	Links []Link `json:"links"`
}

func (h *Handler) FindVideoAction(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(chi.URLParam(r, "uuid"))

	deps := video.FindCommandDeps{
		Context:    r.Context(),
		Repository: h.Repository.VideoRepository,
	}

	args := video.FindCommandArgs{
		Uuid: id,
	}

	result, err := video.FindCommand(deps, args)
	if err != nil {
		log.Printf("Failure to find video: %s", err)
		HandleVideoError(w, err)
		return
	}

	response := VideoResponse{
		Video: result,
		Links: []Link{
			{Method: http.MethodGet, Url: fmt.Sprintf("/api/video/%s/stream", result.Uuid)},
			{Method: http.MethodDelete, Url: fmt.Sprintf("/api/video/%s", result.Uuid)},
		},
	}
	ApiResponse(w, http.StatusOK, response)
}

func (h *Handler) ReadBytesVideoAction(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(chi.URLParam(r, "uuid"))

	rangeHeader := r.Header.Get(HeaderRange)

	re := regexp.MustCompile("\\D")
	value := re.ReplaceAllString(rangeHeader, "")
	startIndex, _ := strconv.Atoi(value)

	deps := video.ReadLimitCommandDeps{
		Context:    r.Context(),
		Repository: h.Repository.VideoRepository,
	}

	args := video.ReadLimitCommandArgs{
		Uuid:       id,
		StartIndex: startIndex,
	}

	result, err := video.ReadLimitCommand(deps, args)
	if err != nil {
		log.Printf("Failure to read video: %s", err)
		HandleVideoError(w, err)
		return
	}

	ApiStreamResponse(w, startIndex, result.EndIndex, int(result.Size), result.Bytes)
}

func (h *Handler) RemoveVideoAction(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(chi.URLParam(r, "uuid"))

	deps := video.RemoveCommandDeps{
		Context:    r.Context(),
		Repository: h.Repository.VideoRepository,
	}

	args := video.RemoveCommandArgs{
		Uuid: id,
	}

	if err := video.RemoveCommand(deps, args); err != nil {
		log.Printf("Failure to remove video: %s", err)
		HandleVideoError(w, err)
		return
	}

	EmptyResponse(w)
}

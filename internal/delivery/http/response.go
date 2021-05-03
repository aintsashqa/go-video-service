package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	HeaderRange         string = "Range"
	HeaderAcceptRanges  string = "Accept-Ranges"
	HeaderContentType   string = "Content-Type"
	HeaderContentRange  string = "Content-Range"
	HeaderContentLength string = "Content-Length"

	ContentTypeJson   string = "application/json"
	ContentTypeStream string = "application/octet-stream"
)

type Link struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Error struct {
	Message string `json:"message"`
}

func Response(w http.ResponseWriter, headers map[string]string, statusCode int, bytes []byte) {
	if headers == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func EmptyResponse(w http.ResponseWriter) {
	headers := map[string]string{
		HeaderContentType: ContentTypeJson,
	}

	Response(w, headers, http.StatusNoContent, []byte{})
}
func ApiResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	headers := map[string]string{
		HeaderContentType: ContentTypeJson,
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Print(err)
		Response(w, nil, 0, []byte{})
		return
	}

	Response(w, headers, statusCode, bytes)
}

func ApiErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	payload := Error{Message: err.Error()}

	ApiResponse(w, statusCode, payload)
}

func ApiStreamResponse(w http.ResponseWriter, startIndex int, nextIndex, size int, bytes []byte) {
	length := nextIndex - startIndex
	statusCode := http.StatusPartialContent
	if length == 0 {
		statusCode = http.StatusOK
	}

	contentRangeHeader := fmt.Sprintf("bytes %d-%d/%d", startIndex, nextIndex, size)
	contentLengthHeader := fmt.Sprintf("%d", length)

	headers := map[string]string{
		HeaderAcceptRanges:  "bytes",
		HeaderContentType:   ContentTypeStream,
		HeaderContentRange:  contentRangeHeader,
		HeaderContentLength: contentLengthHeader,
	}

	Response(w, headers, statusCode, bytes)
}

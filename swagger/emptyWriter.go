package swagger

import (
	"net/http"
)

type EmptyWriter struct{}

func (e *EmptyWriter) Header() http.Header {
	return make(map[string][]string)
}

func (e *EmptyWriter) Write(data []byte) (int, error) {
	return 0, nil
}

func (e *EmptyWriter) WriteHeader(statusCode int) {
}

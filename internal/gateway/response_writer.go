package gateway

import (
	"bytes"
	"net/http"
)

// responseWriter is a custom http.ResponseWriter that captures the response body
type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

// Write captures the response body while writing it to the underlying ResponseWriter
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// getBody returns the captured response body
func (w *responseWriter) getBody() []byte {
	return w.body.Bytes()
}

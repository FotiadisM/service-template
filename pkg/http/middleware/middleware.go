package middleware

import (
	"bytes"
	"net/http"
)

type WrappedResponseWriter struct {
	Body       bytes.Buffer
	StatusCode int

	http.ResponseWriter
}

func (w *WrappedResponseWriter) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return w.ResponseWriter.Write(buf)
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

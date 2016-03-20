package encoding_test

import (
	"bytes"
	"net/http"
)

type embedMime struct {
	mime string
}

func (em *embedMime) GetMime() string {
	if em == nil || em.mime == "" {
		return "application/json"
	}

	return em.mime
}

func (em *embedMime) SetMime(mime string) {
	em.mime = mime
}

type request struct {
	*embedMime
	Str  string      `json:"str" xml:"str"`
	Num  float64     `json:"num" xml:"num"`
	Bool bool        `json:"bool" xml:"bool"`
	Null interface{} `json:"null" xml:"null"`
}

type responseWriter struct {
	*bytes.Buffer
	headers http.Header
}

func (rw *responseWriter) Header() http.Header {
	return rw.headers
}

func (rw *responseWriter) WriteHeader(status int) {
}

func createResponseWriter(buf *bytes.Buffer) http.ResponseWriter {
	return &responseWriter{
		Buffer:  buf,
		headers: make(http.Header),
	}
}

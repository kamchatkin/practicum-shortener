package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var respContentTypes4Gzip = []string{
	"application/json",
	"text/html",
}

var isCompressedContentType map[string]bool

func init() {
	for _, contentType := range respContentTypes4Gzip {
		isCompressedContentType[contentType] = true
	}
}

// WithGzipped Жмем и разжимаем запросы и ответы
func WithGzipped(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respWriter := w
		isGzipSupported := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if isGzipSupported {
			cw := newCompressWriter(w)
			respWriter = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isReceivedGzip := strings.Contains(contentEncoding, "gzip")
		if isReceivedGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(respWriter, r)
	}
}

type (
	compressWriter struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}

	compressReader struct {
		r  io.ReadCloser
		zr *gzip.Reader
	}
)

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	if _, ok := isCompressedContentType[cw.w.Header().Get("Content-Type")]; !ok {
		return cw.w.Write(b)
	}

	return cw.zw.Write(b)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.zw.Close()
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (int, error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}

	return cr.zr.Close()
}

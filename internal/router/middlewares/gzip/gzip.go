package gzip

import (
	//"bytes"
	"compress/gzip"
	//"fmt"
	"io"
	"net/http"
	"strings"
)

var respContentTypes4Gzip = []string{
	"application/json",
	"text/html",
}

var isCompressedContentType = map[string]bool{}

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

		isReceivedGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
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

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w      http.ResponseWriter
	zw     *gzip.Writer
	wBytes int64
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	isCompressedType := false
	ct := c.w.Header().Get("Content-Type")
	cts := strings.Split(ct, ";")
	ct = strings.TrimSpace(cts[0])
	_, isCompressedType = isCompressedContentType[ct]

	if len(p) == 0 || !isCompressedType {
		return c.w.Write(p)
	}

	b, err := c.zw.Write(p)
	c.wBytes += int64(b)
	return b, err
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if c.wBytes == 0 {
		return nil
	}

	return c.zw.Close()
}

// -----------------

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
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

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

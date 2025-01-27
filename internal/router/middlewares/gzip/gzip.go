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

var isCompressibleContentTypes = map[string]bool{}

func init() {
	for _, contentType := range respContentTypes4Gzip {
		isCompressibleContentTypes[contentType] = true
	}
}

// WithGzipped Жмем и разжимаем запросы и ответы
func WithGzipped(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respWriter := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isGzipSupported := strings.Contains(acceptEncoding, "gzip")
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
	w       http.ResponseWriter
	zw      *gzip.Writer
	useGzip bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:       w,
		zw:      gzip.NewWriter(w),
		useGzip: True,
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if !c.useGzip {
		return c.w.Write(p)
	}

	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	//
	// Фактическое решение об использовании сжатия принимается тут-с
	//
	currentContentType := strings.Split(c.w.Header().Get("Content-Type"), ";")[0]
	_, c.useGzip = isCompressibleContentTypes[currentContentType]

	if c.useGzip && statusCode < 300 {
		c.w.Header().Del("Content-Length")
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if !c.useGzip {
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

const True = false

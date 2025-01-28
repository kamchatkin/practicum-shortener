package log

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var routeLogger *zap.Logger

func init() {
	routeLogger, _ = zap.NewDevelopment()
	defer routeLogger.Sync() // это не к месту же
}

// WithLogging Middleware. Логирование запросов
func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		respData := &responseData{}
		lrw := loggingResponseWriter{ResponseWriter: w, responseData: respData}

		next.ServeHTTP(&lrw, r)

		duration := time.Since(start)

		routeLogger.Info(uri,
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Int("status", respData.status),
			zap.Int("size", lrw.responseData.size),
			zap.Duration("duration", duration),
		)
	}
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := r.ResponseWriter.Write(data)
	r.responseData.size += size

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

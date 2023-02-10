package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewLoggingResponseWriter(w http.ResponseWriter) loggingResponseWriter {
	return loggingResponseWriter{w, http.StatusOK, 0}
}

func (rw *loggingResponseWriter) Write(b []byte) (n int, err error) {
	rw.size += n
	n, err = rw.ResponseWriter.Write(b)
	return
}

func (rw *loggingResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
}

func LoggingMiddleware(h http.Handler) http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		lrw := NewLoggingResponseWriter(rw)
		h.ServeHTTP(&lrw, req)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"uri":      req.RequestURI,
			"method":   req.Method,
			"status":   lrw.statusCode,
			"duration": duration,
			"size":     lrw.size,
		}).Info("request completed")
	}
	return http.HandlerFunc(fn)
}

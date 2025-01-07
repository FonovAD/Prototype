package api

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *server) WriteMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)
		s.metricMonitor.IncRequestsTotal(req.Method, lrw.statusCode)
		s.metricMonitor.IncRequestLatency(req.Method, req.URL.Path, time.Since(start).Seconds())
	})
}

func (s *server) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		var body string
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			body = string(bodyBytes)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next.ServeHTTP(lrw, req)
		s.logger.LogRequest(req.Method, req.URL.Path, lrw.statusCode, time.Since(start).Seconds(), body, nil)
	})
}

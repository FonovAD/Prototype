package api

import (
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

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		}
		next.ServeHTTP(w, r)
	})
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

package main

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

// recovery middleware
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered", "err", err, "path", request.URL.Path)
				http.Error(writer, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

// logging middleware
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: writer}

		next.ServeHTTP(rec, request)

		slog.Info("request",
			"method", request.Method,
			"path", request.URL.Path,
			"pattern", request.Pattern,
			"status", rec.status,
			"dur_ms", time.Since(start).Milliseconds(),
			"remoteAddress", request.RemoteAddr,
			"origin", request.Header.Get("Origin"),
		)
	})
}

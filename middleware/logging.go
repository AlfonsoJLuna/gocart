package middleware

import (
    "log"
    "time"
    "net/http"
)

// Wraper for http.ResponseWriter, so we can log status code and size of the response
type loggingResponseWriter struct {
    http.ResponseWriter
    status int
    size   int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
    w.status = code
    w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.size += n // Track number of bytes written
    return n, err
}

func WithLogging(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap ResponseWriter
        lw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

        next(lw, r)

        log.Printf("Received %s %s from %s, answered %d %dB in %s",
            r.Method,
            r.URL.Path,
            r.RemoteAddr,
            lw.status,
            lw.size,
            time.Since(start),
        )
    }
}

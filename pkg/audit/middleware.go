package audit

import (
	"net/http"
	"time"
)

// ResponseWriter wrapper to capture HTTP response
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // default 200
	}
}

// WriteHeader captures status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures data writing
func (rw *ResponseWriter) Write(data []byte) (int, error) {
	written, err := rw.ResponseWriter.Write(data)
	rw.written += written
	return written, err
}

// StatusCode returns captured status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Middleware audit logging middleware
func Middleware(auditLogger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer
			wrappedWriter := NewResponseWriter(w)

			// Process request
			next.ServeHTTP(wrappedWriter, r)

			// Calculate duration
			duration := time.Since(start)

			// Log API call
			details := map[string]interface{}{
				"content_length": r.ContentLength,
				"bytes_written":  wrappedWriter.written,
			}

			// Add Content-Type if present
			if contentType := r.Header.Get("Content-Type"); contentType != "" {
				details["content_type"] = contentType
			}

			auditLogger.LogAPICall(r, wrappedWriter.StatusCode(), duration, details)
		})
	}
}

// MiddlewareFunc middleware as function
func MiddlewareFunc(auditLogger *Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer
		wrappedWriter := NewResponseWriter(w)

		// Process request
		next.ServeHTTP(wrappedWriter, r)

		// Calculate duration
		duration := time.Since(start)

		// Log API call
		details := map[string]interface{}{
			"content_length": r.ContentLength,
			"bytes_written":  wrappedWriter.written,
		}

		// Add Content-Type if present
		if contentType := r.Header.Get("Content-Type"); contentType != "" {
			details["content_type"] = contentType
		}

		auditLogger.LogAPICall(r, wrappedWriter.StatusCode(), duration, details)
	}
}

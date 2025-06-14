package audit

import (
	"net/http"
	"time"
)

// ResponseWriter HTTP response'u yakalamak için wrapper
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

// NewResponseWriter yeni bir ResponseWriter oluşturur
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // varsayılan 200
	}
}

// WriteHeader status code'u yakalar
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write data yazımını yakalar
func (rw *ResponseWriter) Write(data []byte) (int, error) {
	written, err := rw.ResponseWriter.Write(data)
	rw.written += written
	return written, err
}

// StatusCode yakalanan status code'u döner
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Middleware audit logging middleware'i
func Middleware(auditLogger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Response writer'ı wrap et
			wrappedWriter := NewResponseWriter(w)

			// Request'i işle
			next.ServeHTTP(wrappedWriter, r)

			// Süreyi hesapla
			duration := time.Since(start)

			// API çağrısını logla
			details := map[string]interface{}{
				"content_length": r.ContentLength,
				"bytes_written":  wrappedWriter.written,
			}

			// Content-Type varsa ekle
			if contentType := r.Header.Get("Content-Type"); contentType != "" {
				details["content_type"] = contentType
			}

			auditLogger.LogAPICall(r, wrappedWriter.StatusCode(), duration, details)
		})
	}
}

// MiddlewareFunc fonksiyon olarak middleware
func MiddlewareFunc(auditLogger *Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Response writer'ı wrap et
		wrappedWriter := NewResponseWriter(w)

		// Request'i işle
		next.ServeHTTP(wrappedWriter, r)

		// Süreyi hesapla
		duration := time.Since(start)

		// API çağrısını logla
		details := map[string]interface{}{
			"content_length": r.ContentLength,
			"bytes_written":  wrappedWriter.written,
		}

		// Content-Type varsa ekle
		if contentType := r.Header.Get("Content-Type"); contentType != "" {
			details["content_type"] = contentType
		}

		auditLogger.LogAPICall(r, wrappedWriter.StatusCode(), duration, details)
	}
}

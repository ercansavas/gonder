package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// EventType audit event türlerini tanımlar
type EventType string

const (
	EventTypeAPICall     EventType = "api_call"
	EventTypeMessageSent EventType = "message_sent"
	EventTypeError       EventType = "error"
	EventTypeStartup     EventType = "startup"
	EventTypeShutdown    EventType = "shutdown"
	EventTypeHealthCheck EventType = "health_check"
)

// AuditEvent sistem olaylarını temsil eder
type AuditEvent struct {
	Timestamp  time.Time   `json:"timestamp"`
	EventType  EventType   `json:"event_type"`
	UserID     string      `json:"user_id,omitempty"`
	SessionID  string      `json:"session_id,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
	Method     string      `json:"method,omitempty"`
	Path       string      `json:"path,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
	Duration   string      `json:"duration,omitempty"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Error      string      `json:"error,omitempty"`
	RemoteAddr string      `json:"remote_addr,omitempty"`
	UserAgent  string      `json:"user_agent,omitempty"`
}

// Logger audit logger
type Logger struct {
	logger *log.Logger
}

// New yeni bir audit logger oluşturur
func New() *Logger {
	logger := log.New(os.Stdout, "[AUDIT] ", 0)
	return &Logger{
		logger: logger,
	}
}

// LogEvent bir audit event'i loglar
func (l *Logger) LogEvent(event AuditEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// JSON formatında serialize et
	jsonData, err := json.Marshal(event)
	if err != nil {
		l.logger.Printf("AUDIT LOG ERROR: %v", err)
		return
	}

	// Console'a yaz
	l.logger.Println(string(jsonData))
}

// LogAPICall API çağrısını loglar
func (l *Logger) LogAPICall(r *http.Request, statusCode int, duration time.Duration, details interface{}) {
	event := AuditEvent{
		EventType:  EventTypeAPICall,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: statusCode,
		Duration:   duration.String(),
		Message:    fmt.Sprintf("%s %s - %d", r.Method, r.URL.Path, statusCode),
		Details:    details,
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
	}

	// Query parameters varsa ekle
	if r.URL.RawQuery != "" {
		event.Details = map[string]interface{}{
			"query_params": r.URL.RawQuery,
			"details":      details,
		}
	}

	l.LogEvent(event)
}

// LogMessageSent mesaj gönderimi loglar
func (l *Logger) LogMessageSent(recipient, messageType, messageID string, success bool, details interface{}) {
	message := fmt.Sprintf("Mesaj gönderildi: %s -> %s (ID: %s)", messageType, recipient, messageID)
	if !success {
		message = fmt.Sprintf("Mesaj gönderme başarısız: %s -> %s", messageType, recipient)
	}

	event := AuditEvent{
		EventType: EventTypeMessageSent,
		Message:   message,
		Details: map[string]interface{}{
			"recipient":    recipient,
			"message_type": messageType,
			"message_id":   messageID,
			"success":      success,
			"extra":        details,
		},
	}

	l.LogEvent(event)
}

// LogError hata durumunu loglar
func (l *Logger) LogError(err error, context string, details interface{}) {
	event := AuditEvent{
		EventType: EventTypeError,
		Message:   fmt.Sprintf("Error in %s: %v", context, err),
		Error:     err.Error(),
		Details:   details,
	}

	l.LogEvent(event)
}

// LogStartup uygulama başlama durumunu loglar
func (l *Logger) LogStartup(port string, details interface{}) {
	event := AuditEvent{
		EventType: EventTypeStartup,
		Message:   fmt.Sprintf("Gonder uygulaması başlatıldı - Port: %s", port),
		Details:   details,
	}

	l.LogEvent(event)
}

// LogHealthCheck health check durumunu loglar
func (l *Logger) LogHealthCheck(status string, details interface{}) {
	event := AuditEvent{
		EventType: EventTypeHealthCheck,
		Message:   fmt.Sprintf("Health check: %s", status),
		Details:   details,
	}

	l.LogEvent(event)
}

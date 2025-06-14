package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// EventType defines audit event types
type EventType string

const (
	EventTypeAPICall     EventType = "api_call"
	EventTypeMessageSent EventType = "message_sent"
	EventTypeError       EventType = "error"
	EventTypeStartup     EventType = "startup"
	EventTypeShutdown    EventType = "shutdown"
	EventTypeHealthCheck EventType = "health_check"
)

// AuditEvent represents system events
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

// New creates a new audit logger
func New() *Logger {
	logger := log.New(os.Stdout, "[AUDIT] ", 0)
	return &Logger{
		logger: logger,
	}
}

// LogEvent logs an audit event
func (l *Logger) LogEvent(event AuditEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Serialize to JSON format
	jsonData, err := json.Marshal(event)
	if err != nil {
		l.logger.Printf("AUDIT LOG ERROR: %v", err)
		return
	}

	// Write to console
	l.logger.Println(string(jsonData))
}

// LogAPICall logs API calls
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

	// Add query parameters if present
	if r.URL.RawQuery != "" {
		event.Details = map[string]interface{}{
			"query_params": r.URL.RawQuery,
			"details":      details,
		}
	}

	l.LogEvent(event)
}

// LogMessageSent logs message sending
func (l *Logger) LogMessageSent(recipient, messageType, messageID string, success bool, details interface{}) {
	message := fmt.Sprintf("Message sent: %s -> %s (ID: %s)", messageType, recipient, messageID)
	if !success {
		message = fmt.Sprintf("Message sending failed: %s -> %s", messageType, recipient)
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

// LogError logs error conditions
func (l *Logger) LogError(err error, context string, details interface{}) {
	event := AuditEvent{
		EventType: EventTypeError,
		Message:   fmt.Sprintf("Error in %s: %v", context, err),
		Error:     err.Error(),
		Details:   details,
	}

	l.LogEvent(event)
}

// LogStartup logs application startup
func (l *Logger) LogStartup(port string, details interface{}) {
	event := AuditEvent{
		EventType: EventTypeStartup,
		Message:   fmt.Sprintf("Gonder application started - Port: %s", port),
		Details:   details,
	}

	l.LogEvent(event)
}

// LogHealthCheck logs health check status
func (l *Logger) LogHealthCheck(status string, details interface{}) {
	event := AuditEvent{
		EventType: EventTypeHealthCheck,
		Message:   fmt.Sprintf("Health check: %s", status),
		Details:   details,
	}

	l.LogEvent(event)
}

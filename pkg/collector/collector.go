package collector

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gonder/pkg/audit"
)

// LogSource defines log source types
type LogSource string

const (
	SourceSyslog     LogSource = "syslog"
	SourceNginx      LogSource = "nginx"
	SourceApache     LogSource = "apache"
	SourceDocker     LogSource = "docker"
	SourceKubernetes LogSource = "kubernetes"
	SourceCustom     LogSource = "custom"
)

// LogLevel defines log levels
type LogLevel string

const (
	LevelDebug   LogLevel = "debug"
	LevelInfo    LogLevel = "info"
	LevelWarn    LogLevel = "warn"
	LevelError   LogLevel = "error"
	LevelFatal   LogLevel = "fatal"
	LevelUnknown LogLevel = "unknown"
)

// SystemLog represents a system log entry
type SystemLog struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Source      LogSource              `json:"source"`
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	Host        string                 `json:"host,omitempty"`
	Service     string                 `json:"service,omitempty"`
	PID         int                    `json:"pid,omitempty"`
	User        string                 `json:"user,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Path        string                 `json:"path,omitempty"`
	StatusCode  int                    `json:"status_code,omitempty"`
	RawLog      string                 `json:"raw_log"`
	ParsedData  map[string]interface{} `json:"parsed_data,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	CollectedAt time.Time              `json:"collected_at"`
}

// LogCollector manages the log collection system
type LogCollector struct {
	auditLogger *audit.Logger
	parsers     map[LogSource]*LogParser
	sources     []LogSourceConfig
	running     bool
}

// LogSourceConfig log source configuration
type LogSourceConfig struct {
	Name     string    `json:"name"`
	Source   LogSource `json:"source"`
	Path     string    `json:"path"`
	Pattern  string    `json:"pattern,omitempty"`
	Enabled  bool      `json:"enabled"`
	Tags     []string  `json:"tags,omitempty"`
	Interval int       `json:"interval"` // seconds
}

// LogParser log parser
type LogParser struct {
	Source  LogSource
	Pattern *regexp.Regexp
	Fields  []string
}

// New creates a new log collector
func New(auditLogger *audit.Logger) *LogCollector {
	collector := &LogCollector{
		auditLogger: auditLogger,
		parsers:     make(map[LogSource]*LogParser),
		running:     false,
	}

	// Add default parsers
	collector.initDefaultParsers()

	// Add default log sources
	collector.initDefaultSources()

	return collector
}

// initDefaultParsers initializes default log parsers
func (lc *LogCollector) initDefaultParsers() {
	// Syslog parser
	syslogPattern := regexp.MustCompile(`^(\w+\s+\d+\s+\d+:\d+:\d+)\s+(\S+)\s+(\S+)(\[\d+\])?\s*:\s*(.*)$`)
	lc.parsers[SourceSyslog] = &LogParser{
		Source:  SourceSyslog,
		Pattern: syslogPattern,
		Fields:  []string{"timestamp", "host", "service", "pid", "message"},
	}

	// Nginx access log parser
	nginxPattern := regexp.MustCompile(`^(\S+)\s+-\s+\S+\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+\S+"\s+(\d+)\s+(\d+)\s+"[^"]*"\s+"([^"]*)"`)
	lc.parsers[SourceNginx] = &LogParser{
		Source:  SourceNginx,
		Pattern: nginxPattern,
		Fields:  []string{"ip", "timestamp", "method", "path", "status", "size", "user_agent"},
	}

	// Docker log parser
	dockerPattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z)\s+(.*)$`)
	lc.parsers[SourceDocker] = &LogParser{
		Source:  SourceDocker,
		Pattern: dockerPattern,
		Fields:  []string{"timestamp", "message"},
	}
}

// initDefaultSources initializes default log sources
func (lc *LogCollector) initDefaultSources() {
	// Get working directory
	workDir, _ := os.Getwd()

	lc.sources = []LogSourceConfig{
		{
			Name:     "test_syslog",
			Source:   SourceSyslog,
			Path:     filepath.Join(workDir, "test_logs/syslog"),
			Enabled:  true,
			Tags:     []string{"system", "syslog", "test"},
			Interval: 3,
		},
		{
			Name:     "test_auth",
			Source:   SourceSyslog,
			Path:     filepath.Join(workDir, "test_logs/auth.log"),
			Enabled:  true,
			Tags:     []string{"security", "auth", "test"},
			Interval: 3,
		},
		{
			Name:     "system_syslog",
			Source:   SourceSyslog,
			Path:     "/var/log/syslog",
			Enabled:  false, // disabled by default (may not exist in test environment)
			Tags:     []string{"system", "syslog"},
			Interval: 5,
		},
		{
			Name:     "system_messages",
			Source:   SourceSyslog,
			Path:     "/var/log/messages",
			Enabled:  false, // disabled by default
			Tags:     []string{"system", "messages"},
			Interval: 5,
		},
		{
			Name:     "nginx_access",
			Source:   SourceNginx,
			Path:     "/var/log/nginx/access.log",
			Enabled:  false, // disabled by default
			Tags:     []string{"web", "nginx", "access"},
			Interval: 5,
		},
		{
			Name:     "auth_log",
			Source:   SourceSyslog,
			Path:     "/var/log/auth.log",
			Enabled:  false, // disabled by default
			Tags:     []string{"security", "auth"},
			Interval: 5,
		},
	}
}

// Start begins the log collection process
func (lc *LogCollector) Start() error {
	if lc.running {
		return fmt.Errorf("log collector already running")
	}

	lc.running = true
	lc.auditLogger.LogEvent(audit.AuditEvent{
		EventType: "log_collector_start",
		Message:   "System log collection started",
		Details: map[string]interface{}{
			"sources_count": len(lc.sources),
			"enabled_sources": func() []string {
				var enabled []string
				for _, source := range lc.sources {
					if source.Enabled {
						enabled = append(enabled, source.Name)
					}
				}
				return enabled
			}(),
		},
	})

	// Start goroutine for each enabled source
	for _, source := range lc.sources {
		if source.Enabled {
			go lc.collectFromSource(source)
		}
	}

	return nil
}

// Stop stops the log collection process
func (lc *LogCollector) Stop() {
	lc.running = false
	lc.auditLogger.LogEvent(audit.AuditEvent{
		EventType: "log_collector_stop",
		Message:   "System log collection stopped",
	})
}

// collectFromSource collects logs from a specific source
func (lc *LogCollector) collectFromSource(config LogSourceConfig) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	var lastPosition int64 = 0

	for lc.running {
		select {
		case <-ticker.C:
			// Check log file
			if _, err := os.Stat(config.Path); os.IsNotExist(err) {
				// File doesn't exist, continue
				continue
			}

			// Open file
			file, err := os.Open(config.Path)
			if err != nil {
				lc.auditLogger.LogError(err, fmt.Sprintf("Failed to open log file: %s", config.Path), map[string]interface{}{
					"source": config.Name,
					"path":   config.Path,
				})
				continue
			}

			// Get file info
			fileInfo, err := file.Stat()
			if err != nil {
				file.Close()
				continue
			}

			// If file is smaller than last position, file might have been rotated
			if fileInfo.Size() < lastPosition {
				lastPosition = 0
			}

			// Seek to last position
			if _, err := file.Seek(lastPosition, 0); err != nil {
				file.Close()
				continue
			}

			// Read new lines
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if systemLog := lc.parseLogLine(line, config); systemLog != nil {
					lc.processSystemLog(*systemLog)
				}
			}

			// Save new position
			newPosition, _ := file.Seek(0, 1)
			lastPosition = newPosition

			file.Close()
		}
	}
}

// parseLogLine parses a log line based on source type
func (lc *LogCollector) parseLogLine(line string, config LogSourceConfig) *SystemLog {
	if strings.TrimSpace(line) == "" {
		return nil
	}

	systemLog := &SystemLog{
		ID:          fmt.Sprintf("log_%d%06d", time.Now().Unix(), time.Now().Nanosecond()/1000),
		Source:      config.Source,
		RawLog:      line,
		Tags:        config.Tags,
		CollectedAt: time.Now(),
		ParsedData:  make(map[string]interface{}),
	}

	// If no parser exists, save as raw log
	parser, exists := lc.parsers[config.Source]
	if !exists {
		systemLog.Timestamp = time.Now()
		systemLog.Message = line
		systemLog.Level = lc.detectLogLevel(line)
		return systemLog
	}

	// Parse with regex
	matches := parser.Pattern.FindStringSubmatch(line)
	if matches == nil {
		// If parsing fails, save as raw log
		systemLog.Timestamp = time.Now()
		systemLog.Message = line
		systemLog.Level = lc.detectLogLevel(line)
		return systemLog
	}

	// Convert parsed data to SystemLog
	systemLog.Timestamp = time.Now() // default

	// Map data to parser fields
	for i, field := range parser.Fields {
		if i+1 < len(matches) {
			value := matches[i+1]
			systemLog.ParsedData[field] = value

			// Copy special fields to system's corresponding fields
			switch field {
			case "timestamp":
				if ts, err := lc.parseTimestamp(value); err == nil {
					systemLog.Timestamp = ts
				}
			case "message":
				systemLog.Message = value
				systemLog.Level = lc.detectLogLevel(value)
			case "host":
				systemLog.Host = value
			case "service":
				systemLog.Service = value
			case "ip":
				systemLog.IP = value
			case "method":
				systemLog.Method = value
			case "path":
				systemLog.Path = value
			case "status":
				if statusCode, err := parseStatusCode(value); err == nil {
					systemLog.StatusCode = statusCode
				}
			}
		}
	}

	return systemLog
}

// processSystemLog processes a system log
func (lc *LogCollector) processSystemLog(log SystemLog) {
	// Write to console in structured format
	jsonData, err := json.Marshal(log)
	if err != nil {
		lc.auditLogger.LogError(err, "Failed to marshal system log", map[string]interface{}{
			"log_id": log.ID,
		})
		return
	}

	// Print to console with SYSTEM_LOG prefix
	fmt.Printf("[SYSTEM_LOG] %s\n", string(jsonData))

	// Additional processing can be added here
	// - Database insertion
	// - Alert checking
	// - Metric collection
	// - External system integration
}

// parseTimestamp parses timestamp from log
func (lc *LogCollector) parseTimestamp(ts string) (time.Time, error) {
	// Common timestamp formats
	formats := []string{
		"Jan 2 15:04:05",
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"02/Jan/2006:15:04:05 -0700",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			// If year is missing, use current year
			if t.Year() == 0 {
				now := time.Now()
				t = t.AddDate(now.Year(), 0, 0)
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", ts)
}

// detectLogLevel detects log level from message
func (lc *LogCollector) detectLogLevel(message string) LogLevel {
	lower := strings.ToLower(message)

	if strings.Contains(lower, "fatal") || strings.Contains(lower, "panic") {
		return LevelFatal
	}
	if strings.Contains(lower, "error") || strings.Contains(lower, "err") {
		return LevelError
	}
	if strings.Contains(lower, "warn") || strings.Contains(lower, "warning") {
		return LevelWarn
	}
	if strings.Contains(lower, "debug") {
		return LevelDebug
	}
	if strings.Contains(lower, "info") {
		return LevelInfo
	}

	return LevelUnknown
}

// parseStatusCode parses HTTP status code
func parseStatusCode(s string) (int, error) {
	var statusCode int
	_, err := fmt.Sscanf(s, "%d", &statusCode)
	return statusCode, err
}

// GetSources returns all log sources
func (lc *LogCollector) GetSources() []LogSourceConfig {
	return lc.sources
}

// IsRunning returns whether collector is running
func (lc *LogCollector) IsRunning() bool {
	return lc.running
}

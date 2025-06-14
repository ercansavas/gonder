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

// LogSource log kaynağı tiplerini tanımlar
type LogSource string

const (
	SourceSyslog     LogSource = "syslog"
	SourceNginx      LogSource = "nginx"
	SourceApache     LogSource = "apache"
	SourceDocker     LogSource = "docker"
	SourceKubernetes LogSource = "kubernetes"
	SourceCustom     LogSource = "custom"
)

// LogLevel log seviyelerini tanımlar
type LogLevel string

const (
	LevelDebug   LogLevel = "debug"
	LevelInfo    LogLevel = "info"
	LevelWarn    LogLevel = "warn"
	LevelError   LogLevel = "error"
	LevelFatal   LogLevel = "fatal"
	LevelUnknown LogLevel = "unknown"
)

// SystemLog sistem log entry'sini temsil eder
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

// LogCollector log toplama sistemini yönetir
type LogCollector struct {
	auditLogger *audit.Logger
	parsers     map[LogSource]*LogParser
	sources     []LogSourceConfig
	running     bool
}

// LogSourceConfig log kaynağı konfigürasyonu
type LogSourceConfig struct {
	Name     string    `json:"name"`
	Source   LogSource `json:"source"`
	Path     string    `json:"path"`
	Pattern  string    `json:"pattern,omitempty"`
	Enabled  bool      `json:"enabled"`
	Tags     []string  `json:"tags,omitempty"`
	Interval int       `json:"interval"` // saniye
}

// LogParser log parse edici
type LogParser struct {
	Source  LogSource
	Pattern *regexp.Regexp
	Fields  []string
}

// New yeni bir log collector oluşturur
func New(auditLogger *audit.Logger) *LogCollector {
	collector := &LogCollector{
		auditLogger: auditLogger,
		parsers:     make(map[LogSource]*LogParser),
		running:     false,
	}

	// Varsayılan parser'ları ekle
	collector.initDefaultParsers()

	// Varsayılan log kaynaklarını ekle
	collector.initDefaultSources()

	return collector
}

// initDefaultParsers varsayılan log parser'larını başlatır
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

// initDefaultSources varsayılan log kaynaklarını başlatır
func (lc *LogCollector) initDefaultSources() {
	// Çalışma dizinini al
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
			Enabled:  false, // varsayılan olarak kapalı (test ortamında olmayabilir)
			Tags:     []string{"system", "syslog"},
			Interval: 5,
		},
		{
			Name:     "system_messages",
			Source:   SourceSyslog,
			Path:     "/var/log/messages",
			Enabled:  false, // varsayılan olarak kapalı
			Tags:     []string{"system", "messages"},
			Interval: 5,
		},
		{
			Name:     "nginx_access",
			Source:   SourceNginx,
			Path:     "/var/log/nginx/access.log",
			Enabled:  false, // varsayılan olarak kapalı
			Tags:     []string{"web", "nginx", "access"},
			Interval: 5,
		},
		{
			Name:     "auth_log",
			Source:   SourceSyslog,
			Path:     "/var/log/auth.log",
			Enabled:  false, // varsayılan olarak kapalı
			Tags:     []string{"security", "auth"},
			Interval: 5,
		},
	}
}

// Start log toplama işlemini başlatır
func (lc *LogCollector) Start() error {
	if lc.running {
		return fmt.Errorf("log collector already running")
	}

	lc.running = true
	lc.auditLogger.LogEvent(audit.AuditEvent{
		EventType: "log_collector_start",
		Message:   "Sistem log toplama başlatıldı",
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

	// Her enabled source için goroutine başlat
	for _, source := range lc.sources {
		if source.Enabled {
			go lc.collectFromSource(source)
		}
	}

	return nil
}

// Stop log toplama işlemini durdurur
func (lc *LogCollector) Stop() {
	lc.running = false
	lc.auditLogger.LogEvent(audit.AuditEvent{
		EventType: "log_collector_stop",
		Message:   "Sistem log toplama durduruldu",
	})
}

// collectFromSource belirli bir kaynaktan log toplar
func (lc *LogCollector) collectFromSource(config LogSourceConfig) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	var lastPosition int64 = 0

	for lc.running {
		select {
		case <-ticker.C:
			// Log dosyasını kontrol et
			if _, err := os.Stat(config.Path); os.IsNotExist(err) {
				// Dosya yoksa devam et
				continue
			}

			// Dosyayı aç
			file, err := os.Open(config.Path)
			if err != nil {
				lc.auditLogger.LogError(err, fmt.Sprintf("log file open error: %s", config.Path), map[string]interface{}{
					"source": config.Name,
					"path":   config.Path,
				})
				continue
			}

			// Son pozisyondan itibaren oku
			file.Seek(lastPosition, 0)
			scanner := bufio.NewScanner(file)

			lineCount := 0
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == "" {
					continue
				}

				// Log'u parse et ve işle
				systemLog := lc.parseLogLine(line, config)
				if systemLog != nil {
					lc.processSystemLog(*systemLog)
					lineCount++
				}
			}

			// Yeni pozisyonu kaydet
			if stat, err := file.Stat(); err == nil {
				lastPosition = stat.Size()
			}

			file.Close()

			if lineCount > 0 {
				lc.auditLogger.LogEvent(audit.AuditEvent{
					EventType: "logs_collected",
					Message:   fmt.Sprintf("%d log entry toplandı: %s", lineCount, config.Name),
					Details: map[string]interface{}{
						"source":     config.Name,
						"path":       config.Path,
						"line_count": lineCount,
					},
				})
			}
		}
	}
}

// parseLogLine log satırını parse eder
func (lc *LogCollector) parseLogLine(line string, config LogSourceConfig) *SystemLog {
	parser, exists := lc.parsers[config.Source]
	if !exists {
		// Parser yoksa raw log olarak kaydet
		return &SystemLog{
			ID:          fmt.Sprintf("log_%d", time.Now().UnixNano()),
			Timestamp:   time.Now(),
			Source:      config.Source,
			Level:       LevelUnknown,
			Message:     line,
			RawLog:      line,
			Tags:        config.Tags,
			CollectedAt: time.Now(),
		}
	}

	matches := parser.Pattern.FindStringSubmatch(line)
	if matches == nil {
		// Parse edilemezse raw log olarak kaydet
		return &SystemLog{
			ID:          fmt.Sprintf("log_%d", time.Now().UnixNano()),
			Timestamp:   time.Now(),
			Source:      config.Source,
			Level:       LevelUnknown,
			Message:     line,
			RawLog:      line,
			Tags:        config.Tags,
			CollectedAt: time.Now(),
		}
	}

	// Parse edilmiş veriyi SystemLog'a dönüştür
	systemLog := &SystemLog{
		ID:          fmt.Sprintf("log_%d", time.Now().UnixNano()),
		Source:      config.Source,
		RawLog:      line,
		Tags:        config.Tags,
		CollectedAt: time.Now(),
		ParsedData:  make(map[string]interface{}),
	}

	// Parser field'larına göre veriyi map'le
	for i, field := range parser.Fields {
		if i+1 < len(matches) {
			value := matches[i+1]
			systemLog.ParsedData[field] = value

			// Özel field'ları sistemin ilgili alanlarına kopyala
			switch field {
			case "timestamp":
				if ts, err := lc.parseTimestamp(value); err == nil {
					systemLog.Timestamp = ts
				} else {
					systemLog.Timestamp = time.Now()
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

// processSystemLog sistem log'unu işler
func (lc *LogCollector) processSystemLog(log SystemLog) {
	// Console'a structured format olarak yaz
	jsonData, err := json.Marshal(log)
	if err != nil {
		lc.auditLogger.LogError(err, "SystemLog JSON marshal error", log)
		return
	}

	fmt.Printf("[SYSTEM_LOG] %s\n", string(jsonData))

	// Kritik log seviyelerini ayrıca audit log'a yaz
	if log.Level == LevelError || log.Level == LevelFatal {
		lc.auditLogger.LogEvent(audit.AuditEvent{
			EventType: "critical_system_log",
			Message:   fmt.Sprintf("Kritik sistem log tespit edildi: %s", log.Message),
			Details: map[string]interface{}{
				"log_id":      log.ID,
				"source":      log.Source,
				"level":       log.Level,
				"host":        log.Host,
				"service":     log.Service,
				"raw_log":     log.RawLog,
				"parsed_data": log.ParsedData,
			},
		})
	}
}

// Yardımcı fonksiyonlar
func (lc *LogCollector) parseTimestamp(ts string) (time.Time, error) {
	// Çeşitli timestamp formatlarını dene
	formats := []string{
		"Jan 2 15:04:05",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
		time.RFC3339,
		time.RFC822,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			// Yıl yoksa şu anki yılı ekle
			if t.Year() == 0 {
				t = t.AddDate(time.Now().Year(), 0, 0)
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", ts)
}

func (lc *LogCollector) detectLogLevel(message string) LogLevel {
	lowerMsg := strings.ToLower(message)

	if strings.Contains(lowerMsg, "error") || strings.Contains(lowerMsg, "err") {
		return LevelError
	}
	if strings.Contains(lowerMsg, "warn") || strings.Contains(lowerMsg, "warning") {
		return LevelWarn
	}
	if strings.Contains(lowerMsg, "fatal") || strings.Contains(lowerMsg, "panic") {
		return LevelFatal
	}
	if strings.Contains(lowerMsg, "debug") {
		return LevelDebug
	}
	if strings.Contains(lowerMsg, "info") {
		return LevelInfo
	}

	return LevelInfo // varsayılan
}

func parseStatusCode(s string) (int, error) {
	var code int
	_, err := fmt.Sscanf(s, "%d", &code)
	return code, err
}

// GetSources aktif log kaynaklarını döner
func (lc *LogCollector) GetSources() []LogSourceConfig {
	return lc.sources
}

// IsRunning collector'ın çalışıp çalışmadığını döner
func (lc *LogCollector) IsRunning() bool {
	return lc.running
}

# Gonder - Technical Design Document

## 1. Architecture Overview


### 1.1 High-Level Architecture

```
┌─────────────────────┐
│   Configuration     │
│   Manager          │
└─────────┬───────────┘
          │
┌─────────▼───────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Log Collector     │    │   Log Parser    │    │  Elasticsearch  │
│   - File Watcher    │───▶│   - Audit Parser│───▶│    Client       │
│   - Batch Reader    │    │   - JSON Output │    │   - Bulk Insert │
│   - Buffer Manager  │    │   - Field Extract│    │   - Retry Logic │
└─────────────────────┘    └─────────────────┘    └─────────────────┘
          │                          │                        │
          ▼                          ▼                        ▼
┌─────────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Health Monitor    │    │   Metrics       │    │   Error Handler │
│   - Status Check    │    │   - Counters    │    │   - Log Failures│
│   - HTTP Endpoint   │    │   - Performance │    │   - Retry Queue │
└─────────────────────┘    └─────────────────┘    └─────────────────┘
```

### 1.2 Component Responsibilities

#### Log Collector
- Monitor log files for changes using `fsnotify`
- Read existing log files on startup
- Maintain a buffer for efficient processing
- Handle file rotation and new file detection

#### Log Parser
- Parse audit log entries using regex patterns
- Extract structured data (timestamp, user, action, etc.)
- Convert to standardized JSON format
- Apply data enrichment and filtering

#### Elasticsearch Client
- Manage connection pool to Elasticsearch
- Batch operations for performance
- Handle authentication and TLS
- Implement retry logic for failed operations

#### Configuration Manager
- Load and validate configuration from YAML
- Support environment variable overrides
- Hot-reload configuration changes
- Provide configuration to other components

## 2. Data Flow

### 2.1 Log Processing Pipeline

```
[Log File] → [Collector] → [Buffer] → [Parser] → [Formatter] → [ES Client] → [Elasticsearch]
     │            │           │          │           │            │              │
     │            │           │          │           │            │              ▼
     │            │           │          │           │            │         [Index Created]
     │            │           │          │           │            │              │
     │            │           │          │           │            ▼              │
     │            │           │          │           │       [Retry Queue]       │
     │            │           │          │           │            │              │
     │            │           │          │           ▼            │              │
     │            │           │          │      [Error Log]       │              │
     │            │           │          │                        │              │
     │            │           │          ▼                        │              │
     │            │           │     [Metrics Update]              │              │
     │            │           │                                   │              │
     │            │           ▼                                   │              │
     │            │      [Buffer Full?]                          │              │
     │            │           │                                   │              │
     │            ▼           │                                   │              │
     │      [File Watch]      │                                   │              │
     │            │           │                                   │              │
     ▼            │           │                                   │              │
[New Content]    │           │                                   │              │
     │            │           │                                   │              │
     └────────────┴───────────┴───────────────────────────────────┴──────────────┘
```

### 2.2 Error Handling Flow

```
[Processing Error]
        │
        ▼
   [Retry Attempt]
        │
        ├─[Success]──────────────▶ [Continue Processing]
        │
        ├─[Temporary Failure]────▶ [Add to Retry Queue]
        │                              │
        │                              ▼
        │                         [Exponential Backoff]
        │                              │
        │                              ▼
        │                         [Retry After Delay]
        │
        └─[Permanent Failure]────▶ [Log Error + Skip]
```

## 3. Module Design

### 3.1 Package Structure

```
gonder/
├── cmd/
│   └── gonder/
│       └── main.go
├── internal/
│   ├── collector/
│   │   ├── collector.go
│   │   ├── filewatcher.go
│   │   └── buffer.go
│   ├── parser/
│   │   ├── parser.go
│   │   ├── audit.go
│   │   └── common.go
│   ├── elasticsearch/
│   │   ├── client.go
│   │   ├── bulk.go
│   │   └── retry.go
│   ├── config/
│   │   ├── config.go
│   │   └── validation.go
│   ├── health/
│   │   ├── health.go
│   │   └── server.go
│   └── metrics/
│       ├── metrics.go
│       └── prometheus.go
├── pkg/
│   └── types/
│       └── log.go
├── configs/
│   ├── config.yaml
│   └── config.example.yaml
├── scripts/
│   ├── install.sh
│   └── systemd/
│       └── gonder.service
└── docs/
    ├── PRD.md
    ├── Technical-Design.md
    └── API.md
```

### 3.2 Core Interfaces

```go
// pkg/types/log.go
type LogEntry struct {
    Timestamp   time.Time              `json:"@timestamp"`
    Source      string                 `json:"source"`
    Level       string                 `json:"level"`
    Message     string                 `json:"message"`
    User        string                 `json:"user,omitempty"`
    SourceIP    string                 `json:"source_ip,omitempty"`
    Action      string                 `json:"action,omitempty"`
    Resource    string                 `json:"resource,omitempty"`
    Result      string                 `json:"result,omitempty"`
    PID         int                    `json:"pid,omitempty"`
    Host        string                 `json:"host"`
    Tags        []string               `json:"tags,omitempty"`
    Fields      map[string]interface{} `json:"fields,omitempty"`
}

type LogCollector interface {
    Start(ctx context.Context) error
    Stop() error
    Subscribe() <-chan LogEntry
}

type LogParser interface {
    Parse(line string) (*LogEntry, error)
    CanParse(line string) bool
}

type ElasticsearchClient interface {
    BulkIndex(ctx context.Context, entries []LogEntry) error
    Health(ctx context.Context) error
    CreateIndex(ctx context.Context, name string) error
}
```

## 4. Configuration Design

### 4.1 Configuration Structure

```yaml
# config.yaml
gonder:
  # Elasticsearch Configuration
  elasticsearch:
    hosts:
      - "https://elasticsearch:9200"
    username: "elastic"
    password: "${ELASTICSEARCH_PASSWORD}"
    index_prefix: "gonder-logs"
    batch_size: 1000
    flush_interval: "30s"
    timeout: "10s"
    retry_max: 3
    retry_delay: "1s"
    
  # Log Sources Configuration
  log_sources:
    audit:
      enabled: true
      paths:
        - "/var/log/audit/audit.log"
        - "/var/log/auth.log"
      parser: "audit"
      poll_interval: "1s"
      start_from: "end"  # start, end, beginning
      
    syslog:
      enabled: false
      paths:
        - "/var/log/syslog"
      parser: "syslog"
      poll_interval: "1s"
      
  # Processing Configuration
  processing:
    buffer_size: 10000
    workers: 4
    batch_timeout: "5s"
    max_line_length: 65536
    
  # Output Configuration
  output:
    format: "json"
    timestamp_format: "2006-01-02T15:04:05Z07:00"
    add_host: true
    add_tags: ["gonder", "audit"]
    
  # Health Check Configuration
  health:
    enabled: true
    bind_address: "0.0.0.0:8080"
    path: "/health"
    
  # Logging Configuration
  logging:
    level: "info"  # debug, info, warn, error
    format: "json"  # json, text
    output: "stdout"  # stdout, stderr, file
    file: "/var/log/gonder/gonder.log"
```

### 4.2 Environment Variables

```bash
# Elasticsearch
GONDER_ELASTICSEARCH_HOSTS="https://es1:9200,https://es2:9200"
GONDER_ELASTICSEARCH_USERNAME="elastic"
GONDER_ELASTICSEARCH_PASSWORD="secret"

# Processing
GONDER_PROCESSING_WORKERS="8"
GONDER_PROCESSING_BUFFER_SIZE="20000"

# Logging
GONDER_LOG_LEVEL="debug"
GONDER_LOG_OUTPUT="file"
GONDER_LOG_FILE="/var/log/gonder/gonder.log"
```

## 5. Data Processing

### 5.1 Audit Log Parsing

Linux audit log example:
```
type=USER_AUTH msg=audit(1639123456.789:12345): pid=1234 uid=0 auid=1000 ses=1 subj=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023 msg='op=PAM:authentication grantors=pam_unix acct="user" exe="/usr/sbin/sshd" hostname=192.168.1.100 addr=192.168.1.100 terminal=ssh res=success'
```

Parsed JSON output:
```json
{
  "@timestamp": "2021-12-10T10:30:56.789Z",
  "source": "audit",
  "level": "info",
  "message": "PAM authentication successful",
  "user": "user",
  "source_ip": "192.168.1.100",
  "action": "authentication",
  "resource": "ssh",
  "result": "success",
  "pid": 1234,
  "host": "server01",
  "tags": ["gonder", "audit", "authentication"],
  "fields": {
    "audit_type": "USER_AUTH",
    "audit_id": 12345,
    "uid": 0,
    "auid": 1000,
    "session": 1,
    "executable": "/usr/sbin/sshd",
    "hostname": "192.168.1.100",
    "terminal": "ssh",
    "grantors": "pam_unix"
  }
}
```

### 5.2 Parser Implementation

```go
// internal/parser/audit.go
type AuditParser struct {
    hostname string
}

func (p *AuditParser) Parse(line string) (*types.LogEntry, error) {
    // Extract basic audit fields
    auditType := extractAuditType(line)
    timestamp := extractTimestamp(line)
    message := extractMessage(line)
    
    entry := &types.LogEntry{
        Timestamp: timestamp,
        Source:    "audit",
        Level:     "info",
        Message:   message,
        Host:      p.hostname,
        Tags:      []string{"gonder", "audit"},
        Fields:    make(map[string]interface{}),
    }
    
    // Parse audit-specific fields
    if fields := parseAuditFields(line); fields != nil {
        entry.Fields = fields
        
        // Extract common fields
        if user, ok := fields["acct"].(string); ok {
            entry.User = user
        }
        if ip, ok := fields["addr"].(string); ok {
            entry.SourceIP = ip
        }
        if result, ok := fields["res"].(string); ok {
            entry.Result = result
        }
    }
    
    return entry, nil
}

func extractAuditType(line string) string {
    re := regexp.MustCompile(`type=(\w+)`)
    matches := re.FindStringSubmatch(line)
    if len(matches) > 1 {
        return matches[1]
    }
    return "UNKNOWN"
}

func extractTimestamp(line string) time.Time {
    re := regexp.MustCompile(`audit\((\d+\.\d+):\d+\)`)
    matches := re.FindStringSubmatch(line)
    if len(matches) > 1 {
        if timestamp, err := strconv.ParseFloat(matches[1], 64); err == nil {
            return time.Unix(int64(timestamp), int64((timestamp-math.Floor(timestamp))*1e9))
        }
    }
    return time.Now()
}
```

## 6. Performance Considerations

### 6.1 Buffering Strategy

- **Ring Buffer**: Fixed-size circular buffer to prevent memory overflow
- **Batch Processing**: Group logs into batches for efficient Elasticsearch operations
- **Backpressure**: Slow down reading when buffer is full
- **Memory Limits**: Configurable maximum memory usage

### 6.2 Concurrency Model

```go
// Collector goroutines
for i := 0; i < numWorkers; i++ {
    go func() {
        for entry := range logChannel {
            processLog(entry)
        }
    }()
}

// Elasticsearch sender
go func() {
    ticker := time.NewTicker(flushInterval)
    batch := make([]LogEntry, 0, batchSize)
    
    for {
        select {
        case entry := <-processedChannel:
            batch = append(batch, entry)
            if len(batch) >= batchSize {
                sendBatch(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                sendBatch(batch)
                batch = batch[:0]
            }
        }
    }
}()
```

### 6.3 Resource Management

- **Memory Pools**: Reuse objects to reduce GC pressure
- **Connection Pooling**: Maintain persistent connections to Elasticsearch
- **Graceful Shutdown**: Flush buffers before termination
- **Resource Limits**: Configurable CPU and memory constraints

## 7. Monitoring and Observability

### 7.1 Metrics

```go
type Metrics struct {
    LogsProcessed     prometheus.Counter
    LogsErrored       prometheus.Counter
    ProcessingLatency prometheus.Histogram
    BufferSize        prometheus.Gauge
    ESConnectionStatus prometheus.Gauge
}
```

### 7.2 Health Checks

```go
type HealthCheck struct {
    ElasticsearchHealthy bool   `json:"elasticsearch_healthy"`
    LogSourcesAccessible bool   `json:"log_sources_accessible"`
    BufferUtilization    float64 `json:"buffer_utilization"`
    LastProcessedLog     string  `json:"last_processed_log"`
    Uptime              string  `json:"uptime"`
}
```

## 8. Security Implementation

### 8.1 Elasticsearch Security

- **TLS/SSL**: Encrypted communication with certificate validation
- **Authentication**: Username/password or API key authentication
- **Index Security**: Template-based index creation with proper mappings

### 8.2 File System Security

- **Permission Checks**: Verify read access to log files
- **Safe Path Handling**: Prevent directory traversal attacks
- **Privilege Dropping**: Run with minimal required permissions

## 9. Deployment Architecture

### 9.1 Systemd Service

```ini
[Unit]
Description=Gonder Log Collector
After=network.target

[Service]
Type=simple
User=gonder
Group=gonder
ExecStart=/usr/local/bin/gonder -config /etc/gonder/config.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### 9.2 Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o gonder cmd/gonder/main.go

FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/gonder .
COPY --from=builder /app/configs/config.yaml .
EXPOSE 8080
CMD ["./gonder", "-config", "config.yaml"]
```

## 10. Testing Strategy

### 10.1 Unit Tests
- Parser functionality with various log formats
- Configuration validation
- Elasticsearch client operations

### 10.2 Integration Tests
- End-to-end log processing pipeline
- Elasticsearch integration
- File watching and rotation handling

### 10.3 Performance Tests
- High-volume log processing
- Memory usage under load
- Elasticsearch bulk operation performance

---

**Document Version:** 1.0  
**Last Updated:** 2024-12-19  
**Author:** Technical Team 

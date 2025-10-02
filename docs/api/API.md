# Gonder - API Documentation & Implementation Guide

## 1. Health Check API

### 1.1 Health Endpoint

**Endpoint:** `GET /health`

**Description:** Application health status check

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-12-19T10:30:45Z",
  "version": "1.0.0",
  "uptime": "2h45m30s",
  "checks": {
    "elasticsearch": {
      "status": "healthy",
      "response_time": "15ms",
      "cluster_name": "elasticsearch",
      "cluster_status": "green"
    },
    "log_sources": {
      "status": "healthy",
      "accessible_files": [
        "/var/log/audit/audit.log",
        "/var/log/auth.log"
      ],
      "inaccessible_files": []
    },
    "buffer": {
      "status": "healthy",
      "utilization": 0.45,
      "current_size": 4500,
      "max_size": 10000
    }
  },
  "metrics": {
    "logs_processed_total": 125430,
    "logs_failed_total": 12,
    "processing_rate_per_second": 150.5,
    "last_processed_timestamp": "2024-12-19T10:30:44Z"
  }
}
```

**Status Codes:**
- `200 OK`: Service is healthy
- `503 Service Unavailable`: Service is unhealthy

### 1.2 Metrics Endpoint

**Endpoint:** `GET /metrics`

**Description:** Prometheus-compatible metrics

**Response:**
```
# HELP gonder_logs_processed_total Total number of logs processed
# TYPE gonder_logs_processed_total counter
gonder_logs_processed_total 125430

# HELP gonder_logs_failed_total Total number of failed log processing attempts
# TYPE gonder_logs_failed_total counter
gonder_logs_failed_total 12

# HELP gonder_processing_duration_seconds Time spent processing logs
# TYPE gonder_processing_duration_seconds histogram
gonder_processing_duration_seconds_bucket{le="0.001"} 1000
gonder_processing_duration_seconds_bucket{le="0.01"} 5000
gonder_processing_duration_seconds_bucket{le="0.1"} 8000
gonder_processing_duration_seconds_bucket{le="1"} 8500
gonder_processing_duration_seconds_bucket{le="+Inf"} 8550

# HELP gonder_buffer_size_current Current buffer size
# TYPE gonder_buffer_size_current gauge
gonder_buffer_size_current 4500

# HELP gonder_elasticsearch_connection_status Elasticsearch connection status (1=connected, 0=disconnected)
# TYPE gonder_elasticsearch_connection_status gauge
gonder_elasticsearch_connection_status 1
```

## 2. Configuration API

### 2.1 Configuration Reload

**Endpoint:** `POST /config/reload`

**Description:** Reload configuration without restarting

**Response:**
```json
{
  "status": "success",
  "message": "Configuration reloaded successfully",
  "timestamp": "2024-12-19T10:30:45Z",
  "changes": [
    "elasticsearch.batch_size: 1000 -> 1500",
    "processing.workers: 4 -> 6"
  ]
}
```

### 2.2 Get Current Configuration

**Endpoint:** `GET /config`

**Description:** Get current configuration (sensitive data masked)

**Response:**
```json
{
  "gonder": {
    "elasticsearch": {
      "hosts": ["https://elasticsearch:9200"],
      "username": "elastic",
      "password": "***MASKED***",
      "index_prefix": "gonder-logs",
      "batch_size": 1000,
      "flush_interval": "30s"
    },
    "log_sources": {
      "audit": {
        "enabled": true,
        "paths": ["/var/log/audit/audit.log"],
        "parser": "audit"
      }
    },
    "processing": {
      "buffer_size": 10000,
      "workers": 4
    }
  }
}
```

## 3. Implementation Phases

### Phase 1: Core Implementation (Week 1-2)

#### 3.1 Project Structure Setup

```bash
# Initialize Go module
go mod init github.com/yourusername/gonder

# Create directory structure
mkdir -p {cmd/gonder,internal/{collector,parser,elasticsearch,config,health,metrics},pkg/types,configs,scripts/systemd,docs}
```

#### 3.2 Basic Types and Interfaces

```go
// pkg/types/log.go
package types

import "time"

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

type Config struct {
    Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
    LogSources    LogSourcesConfig    `yaml:"log_sources"`
    Processing    ProcessingConfig    `yaml:"processing"`
    Output        OutputConfig        `yaml:"output"`
    Health        HealthConfig        `yaml:"health"`
    Logging       LoggingConfig       `yaml:"logging"`
}

type ElasticsearchConfig struct {
    Hosts         []string `yaml:"hosts"`
    Username      string   `yaml:"username"`
    Password      string   `yaml:"password"`
    IndexPrefix   string   `yaml:"index_prefix"`
    BatchSize     int      `yaml:"batch_size"`
    FlushInterval string   `yaml:"flush_interval"`
    Timeout       string   `yaml:"timeout"`
    RetryMax      int      `yaml:"retry_max"`
    RetryDelay    string   `yaml:"retry_delay"`
}
```

#### 3.3 Configuration Management

```go
// internal/config/config.go
package config

import (
    "fmt"
    "os"
    "strings"
    
    "gopkg.in/yaml.v3"
    "github.com/yourusername/gonder/pkg/types"
)

func LoadConfig(path string) (*types.Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    // Replace environment variables
    content := string(data)
    content = os.ExpandEnv(content)
    
    var config types.Config
    if err := yaml.Unmarshal([]byte(content), &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return &config, nil
}

func validateConfig(config *types.Config) error {
    if len(config.Elasticsearch.Hosts) == 0 {
        return fmt.Errorf("elasticsearch hosts cannot be empty")
    }
    
    if config.Processing.BufferSize <= 0 {
        return fmt.Errorf("buffer size must be positive")
    }
    
    if config.Processing.Workers <= 0 {
        return fmt.Errorf("workers count must be positive")
    }
    
    return nil
}
```

### Phase 2: Log Collection (Week 2-3)

#### 3.4 File Collector Implementation

```go
// internal/collector/collector.go
package collector

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    
    "github.com/fsnotify/fsnotify"
    "github.com/yourusername/gonder/pkg/types"
)

type Collector struct {
    config   *types.Config
    watcher  *fsnotify.Watcher
    logChan  chan types.LogEntry
    stopChan chan struct{}
    wg       sync.WaitGroup
}

func New(config *types.Config) (*Collector, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, fmt.Errorf("failed to create file watcher: %w", err)
    }
    
    return &Collector{
        config:   config,
        watcher:  watcher,
        logChan:  make(chan types.LogEntry, config.Processing.BufferSize),
        stopChan: make(chan struct{}),
    }, nil
}

func (c *Collector) Start(ctx context.Context) error {
    // Add log files to watcher
    for _, source := range c.config.LogSources {
        if !source.Enabled {
            continue
        }
        
        for _, path := range source.Paths {
            if err := c.addPath(path); err != nil {
                return fmt.Errorf("failed to add path %s: %w", path, err)
            }
        }
    }
    
    c.wg.Add(1)
    go c.watchFiles(ctx)
    
    return nil
}

func (c *Collector) Stop() error {
    close(c.stopChan)
    c.wg.Wait()
    return c.watcher.Close()
}

func (c *Collector) Subscribe() <-chan types.LogEntry {
    return c.logChan
}

func (c *Collector) watchFiles(ctx context.Context) {
    defer c.wg.Done()
    
    for {
        select {
        case event := <-c.watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                c.handleFileChange(event.Name)
            }
        case err := <-c.watcher.Errors:
            fmt.Printf("watcher error: %v\n", err)
        case <-c.stopChan:
            return
        case <-ctx.Done():
            return
        }
    }
}

func (c *Collector) addPath(path string) error {
    // Check if file exists and is readable
    if _, err := os.Stat(path); err != nil {
        return fmt.Errorf("cannot access file %s: %w", path, err)
    }
    
    // Add directory to watcher (for file rotation detection)
    dir := filepath.Dir(path)
    return c.watcher.Add(dir)
}
```

### Phase 3: Log Parsing (Week 3)

#### 3.5 Audit Log Parser

```go
// internal/parser/audit.go
package parser

import (
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"
    
    "github.com/yourusername/gonder/pkg/types"
)

type AuditParser struct {
    hostname string
    patterns map[string]*regexp.Regexp
}

func NewAuditParser() *AuditParser {
    hostname, _ := os.Hostname()
    
    return &AuditParser{
        hostname: hostname,
        patterns: map[string]*regexp.Regexp{
            "type":      regexp.MustCompile(`type=(\w+)`),
            "timestamp": regexp.MustCompile(`audit\((\d+\.\d+):\d+\)`),
            "pid":       regexp.MustCompile(`pid=(\d+)`),
            "uid":       regexp.MustCompile(`uid=(\d+)`),
            "user":      regexp.MustCompile(`acct="([^"]+)"`),
            "addr":      regexp.MustCompile(`addr=(\d+\.\d+\.\d+\.\d+)`),
            "result":    regexp.MustCompile(`res=(\w+)`),
            "exe":       regexp.MustCompile(`exe="([^"]+)"`),
        },
    }
}

func (p *AuditParser) Parse(line string) (*types.LogEntry, error) {
    if !p.CanParse(line) {
        return nil, fmt.Errorf("cannot parse line: %s", line)
    }
    
    entry := &types.LogEntry{
        Source:    "audit",
        Level:     "info",
        Host:      p.hostname,
        Tags:      []string{"gonder", "audit"},
        Fields:    make(map[string]interface{}),
        Timestamp: time.Now(), // default fallback
    }
    
    // Extract audit type
    if matches := p.patterns["type"].FindStringSubmatch(line); len(matches) > 1 {
        entry.Fields["audit_type"] = matches[1]
        entry.Action = strings.ToLower(matches[1])
    }
    
    // Extract timestamp
    if matches := p.patterns["timestamp"].FindStringSubmatch(line); len(matches) > 1 {
        if timestamp, err := strconv.ParseFloat(matches[1], 64); err == nil {
            entry.Timestamp = time.Unix(int64(timestamp), 0)
        }
    }
    
    // Extract PID
    if matches := p.patterns["pid"].FindStringSubmatch(line); len(matches) > 1 {
        if pid, err := strconv.Atoi(matches[1]); err == nil {
            entry.PID = pid
        }
    }
    
    // Extract user
    if matches := p.patterns["user"].FindStringSubmatch(line); len(matches) > 1 {
        entry.User = matches[1]
    }
    
    // Extract source IP
    if matches := p.patterns["addr"].FindStringSubmatch(line); len(matches) > 1 {
        entry.SourceIP = matches[1]
    }
    
    // Extract result
    if matches := p.patterns["result"].FindStringSubmatch(line); len(matches) > 1 {
        entry.Result = matches[1]
        if entry.Result == "success" {
            entry.Level = "info"
        } else {
            entry.Level = "warning"
        }
    }
    
    // Extract executable
    if matches := p.patterns["exe"].FindStringSubmatch(line); len(matches) > 1 {
        entry.Resource = matches[1]
        entry.Fields["executable"] = matches[1]
    }
    
    // Generate message
    entry.Message = p.generateMessage(entry)
    
    return entry, nil
}

func (p *AuditParser) CanParse(line string) bool {
    return strings.Contains(line, "type=") && strings.Contains(line, "audit(")
}

func (p *AuditParser) generateMessage(entry *types.LogEntry) string {
    if entry.User != "" && entry.Action != "" {
        return fmt.Sprintf("User %s performed %s", entry.User, entry.Action)
    }
    if entry.Action != "" {
        return fmt.Sprintf("Audit event: %s", entry.Action)
    }
    return "Audit log entry"
}
```

### Phase 4: Elasticsearch Integration (Week 4)

#### 3.6 Elasticsearch Client

```go
// internal/elasticsearch/client.go
package elasticsearch

import (
    "bytes"
    "context"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"
    
    "github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
    "github.com/yourusername/gonder/pkg/types"
)

type Client struct {
    es     *elasticsearch.Client
    config *types.ElasticsearchConfig
}

func NewClient(config *types.ElasticsearchConfig) (*Client, error) {
    cfg := elasticsearch.Config{
        Addresses: config.Hosts,
        Username:  config.Username,
        Password:  config.Password,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
            },
        },
    }
    
    es, err := elasticsearch.NewClient(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
    }
    
    return &Client{
        es:     es,
        config: config,
    }, nil
}

func (c *Client) BulkIndex(ctx context.Context, entries []types.LogEntry) error {
    if len(entries) == 0 {
        return nil
    }
    
    var buf bytes.Buffer
    
    for _, entry := range entries {
        // Index action
        indexName := fmt.Sprintf("%s-%s", c.config.IndexPrefix, 
            entry.Timestamp.Format("2006.01.02"))
        
        action := map[string]interface{}{
            "index": map[string]interface{}{
                "_index": indexName,
            },
        }
        
        actionBytes, _ := json.Marshal(action)
        buf.Write(actionBytes)
        buf.WriteByte('\n')
        
        // Document
        docBytes, err := json.Marshal(entry)
        if err != nil {
            return fmt.Errorf("failed to marshal log entry: %w", err)
        }
        
        buf.Write(docBytes)
        buf.WriteByte('\n')
    }
    
    req := esapi.BulkRequest{
        Body:    strings.NewReader(buf.String()),
        Refresh: "false",
    }
    
    res, err := req.Do(ctx, c.es)
    if err != nil {
        return fmt.Errorf("bulk request failed: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("bulk request error: %s", res.Status())
    }
    
    return nil
}

func (c *Client) Health(ctx context.Context) error {
    res, err := c.es.Cluster.Health()
    if err != nil {
        return fmt.Errorf("health check failed: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("cluster unhealthy: %s", res.Status())
    }
    
    return nil
}

func (c *Client) CreateIndex(ctx context.Context, name string) error {
    mapping := `{
        "mappings": {
            "properties": {
                "@timestamp": { "type": "date" },
                "source": { "type": "keyword" },
                "level": { "type": "keyword" },
                "message": { "type": "text" },
                "user": { "type": "keyword" },
                "source_ip": { "type": "ip" },
                "action": { "type": "keyword" },
                "resource": { "type": "keyword" },
                "result": { "type": "keyword" },
                "pid": { "type": "integer" },
                "host": { "type": "keyword" },
                "tags": { "type": "keyword" }
            }
        }
    }`
    
    req := esapi.IndicesCreateRequest{
        Index: name,
        Body:  strings.NewReader(mapping),
    }
    
    res, err := req.Do(ctx, c.es)
    if err != nil {
        return fmt.Errorf("create index failed: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() && res.StatusCode != 400 { // 400 = already exists
        return fmt.Errorf("create index error: %s", res.Status())
    }
    
    return nil
}
```

## 4. Main Application

### 4.1 Main Function

```go
// cmd/gonder/main.go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    
    "github.com/yourusername/gonder/internal/collector"
    "github.com/yourusername/gonder/internal/config"
    "github.com/yourusername/gonder/internal/elasticsearch"
    "github.com/yourusername/gonder/internal/health"
    "github.com/yourusername/gonder/internal/parser"
)

func main() {
    var configPath = flag.String("config", "config.yaml", "Path to configuration file")
    flag.Parse()
    
    // Load configuration
    cfg, err := config.LoadConfig(*configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Create context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Initialize components
    collector, err := collector.New(cfg)
    if err != nil {
        log.Fatalf("Failed to create collector: %v", err)
    }
    
    parser := parser.NewAuditParser()
    
    esClient, err := elasticsearch.NewClient(&cfg.Elasticsearch)
    if err != nil {
        log.Fatalf("Failed to create Elasticsearch client: %v", err)
    }
    
    // Start health server
    healthServer := health.NewServer(cfg.Health.BindAddress, esClient, collector)
    go healthServer.Start()
    
    // Start log processing pipeline
    var wg sync.WaitGroup
    
    // Start collector
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := collector.Start(ctx); err != nil {
            log.Printf("Collector error: %v", err)
        }
    }()
    
    // Start processor
    wg.Add(1)
    go func() {
        defer wg.Done()
        processLogs(ctx, collector.Subscribe(), parser, esClient, cfg)
    }()
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    log.Println("Shutting down...")
    
    cancel()
    collector.Stop()
    wg.Wait()
    
    log.Println("Shutdown complete")
}

func processLogs(ctx context.Context, logChan <-chan types.LogEntry, 
                parser *parser.AuditParser, esClient *elasticsearch.Client, 
                cfg *types.Config) {
    
    batch := make([]types.LogEntry, 0, cfg.Elasticsearch.BatchSize)
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case entry := <-logChan:
            batch = append(batch, entry)
            
            if len(batch) >= cfg.Elasticsearch.BatchSize {
                if err := esClient.BulkIndex(ctx, batch); err != nil {
                    log.Printf("Failed to index batch: %v", err)
                }
                batch = batch[:0]
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                if err := esClient.BulkIndex(ctx, batch); err != nil {
                    log.Printf("Failed to index batch: %v", err)
                }
                batch = batch[:0]
            }
            
        case <-ctx.Done():
            // Final flush
            if len(batch) > 0 {
                if err := esClient.BulkIndex(ctx, batch); err != nil {
                    log.Printf("Failed to index final batch: %v", err)
                }
            }
            return
        }
    }
}
```

## 5. Deployment & Build

### 5.1 Go Dependencies

```bash
# go.mod
module github.com/yourusername/gonder

go 1.21

require (
    github.com/elastic/go-elasticsearch/v8 v8.11.0
    github.com/fsnotify/fsnotify v1.7.0
    github.com/prometheus/client_golang v1.17.0
    github.com/sirupsen/logrus v1.9.3
    gopkg.in/yaml.v3 v3.0.1
)
```

### 5.2 Build Script

```bash
#!/bin/bash
# scripts/build.sh

set -e

VERSION=${VERSION:-"1.0.0"}
BINARY_NAME="gonder"

echo "Building ${BINARY_NAME} version ${VERSION}..."

# Build for Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o dist/${BINARY_NAME}-linux-amd64 cmd/gonder/main.go

echo "Build complete!"
```

### 5.3 Installation Script

```bash
#!/bin/bash
# scripts/install.sh

set -e

BINARY_NAME="gonder"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/gonder"
LOG_DIR="/var/log/gonder"
SERVICE_DIR="/etc/systemd/system"

echo "Installing ${BINARY_NAME}..."

# Create user
if ! id -u gonder > /dev/null 2>&1; then
    useradd -r -s /bin/false gonder
fi

# Create directories
mkdir -p ${CONFIG_DIR} ${LOG_DIR}
chown gonder:gonder ${LOG_DIR}

# Install binary
cp dist/${BINARY_NAME}-linux-amd64 ${INSTALL_DIR}/${BINARY_NAME}
chmod +x ${INSTALL_DIR}/${BINARY_NAME}

# Install config
cp configs/config.yaml ${CONFIG_DIR}/
chown gonder:gonder ${CONFIG_DIR}/config.yaml

# Install systemd service
cp scripts/systemd/gonder.service ${SERVICE_DIR}/
systemctl daemon-reload
systemctl enable gonder

echo "Installation complete!"
echo "Start with: systemctl start gonder"
```

---

**Document Version:** 1.0  
**Last Updated:** 2024-12-19  
**Author:** Development Team 
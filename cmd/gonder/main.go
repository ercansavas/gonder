package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gonder/internal/config"
	"gonder/pkg/audit"
	"gonder/pkg/collector"
	"gonder/pkg/handler"
)

func main() {
	fmt.Println("🚀 Gonder - System Log Collection Service starting...")

	// Start audit logger
	auditLogger := audit.New()

	// Load configuration
	cfg := config.Load()

	// Start log collector
	logCollector := collector.New(auditLogger)

	// Startup audit log
	auditLogger.LogStartup(cfg.Port, map[string]interface{}{
		"host":      cfg.Host,
		"log_level": cfg.LogLevel,
		"version":   "2.0.0",
		"purpose":   "system_log_collection",
		"features": []string{
			"system_log_collection",
			"audit_logging",
			"real_time_monitoring",
			"log_parsing",
			"structured_output",
		},
	})

	// Start handlers
	h := handler.New(auditLogger)
	logHandler := handler.NewLogHandler(logCollector)

	// Define routes - wrap with audit middleware
	http.HandleFunc("/", audit.MiddlewareFunc(auditLogger, h.Home))
	http.HandleFunc("/api/health", audit.MiddlewareFunc(auditLogger, h.Health))

	// Log management endpoints
	http.HandleFunc("/api/logs/status", audit.MiddlewareFunc(auditLogger, logHandler.GetStatus))
	http.HandleFunc("/api/logs/sources", audit.MiddlewareFunc(auditLogger, logHandler.GetSources))
	http.HandleFunc("/api/logs/start", audit.MiddlewareFunc(auditLogger, logHandler.StartCollector))
	http.HandleFunc("/api/logs/stop", audit.MiddlewareFunc(auditLogger, logHandler.StopCollector))

	// Backward compatibility (deprecated)
	http.HandleFunc("/api/send", audit.MiddlewareFunc(auditLogger, h.Send))

	// Auto-start log collector
	fmt.Println("🔧 Starting system log collector...")
	if err := logCollector.Start(); err != nil {
		auditLogger.LogError(err, "Log collector startup error", nil)
		fmt.Printf("⚠️ Log collector could not be started: %v\n", err)
	} else {
		fmt.Println("✅ System log collector started successfully")
	}

	// Signal handler for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n🛑 Shutdown signal received, starting clean shutdown process...")

		// Stop log collector
		logCollector.Stop()

		// Shutdown audit log
		auditLogger.LogEvent(audit.AuditEvent{
			EventType: "system_shutdown",
			Message:   "System is shutting down cleanly",
		})

		os.Exit(0)
	}()

	// Start server
	fmt.Printf("🌐 Server running on port %s\n", cfg.Port)
	fmt.Println("📋 Endpoints:")
	fmt.Println("  GET  /                    - Home page")
	fmt.Println("  GET  /api/health          - System health check")
	fmt.Println("  GET  /api/logs/status     - Log collector status")
	fmt.Println("  GET  /api/logs/sources    - List log sources")
	fmt.Println("  POST /api/logs/start      - Start log collector")
	fmt.Println("  POST /api/logs/stop       - Stop log collector")
	fmt.Println("  POST /api/send            - [DEPRECATED] Send message")
	fmt.Println("📊 System log collection active - Logs are written to console")
	fmt.Println("🔍 Monitored log files:")

	// Show active log sources
	sources := logCollector.GetSources()
	for _, source := range sources {
		if source.Enabled {
			fmt.Printf("  ✅ %s (%s) - %s\n", source.Name, source.Source, source.Path)
		} else {
			fmt.Printf("  ❌ %s (%s) - %s [DISABLED]\n", source.Name, source.Source, source.Path)
		}
	}

	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

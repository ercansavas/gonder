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
	fmt.Println("🚀 Gonder - Sistem Log Toplama Servisi başlatılıyor...")

	// Audit logger'ı başlat
	auditLogger := audit.New()

	// Konfigürasyon yükle
	cfg := config.Load()

	// Log collector'ı başlat
	logCollector := collector.New(auditLogger)

	// Startup audit log
	auditLogger.LogStartup(cfg.Port, map[string]interface{}{
		"host":      cfg.Host,
		"log_level": cfg.LogLevel,
		"version":   "2.0.0",
		"purpose":   "sistem_log_toplama",
		"features": []string{
			"sistem_log_toplama",
			"audit_logging",
			"real_time_monitoring",
			"log_parsing",
			"structured_output",
		},
	})

	// Handler'ları başlat
	h := handler.New(auditLogger)
	logHandler := handler.NewLogHandler(logCollector)

	// Routes tanımla - audit middleware ile wrap et
	http.HandleFunc("/", audit.MiddlewareFunc(auditLogger, h.Home))
	http.HandleFunc("/api/health", audit.MiddlewareFunc(auditLogger, h.Health))

	// Log management endpoints
	http.HandleFunc("/api/logs/status", audit.MiddlewareFunc(auditLogger, logHandler.GetStatus))
	http.HandleFunc("/api/logs/sources", audit.MiddlewareFunc(auditLogger, logHandler.GetSources))
	http.HandleFunc("/api/logs/start", audit.MiddlewareFunc(auditLogger, logHandler.StartCollector))
	http.HandleFunc("/api/logs/stop", audit.MiddlewareFunc(auditLogger, logHandler.StopCollector))

	// Backward compatibility (deprecated)
	http.HandleFunc("/api/send", audit.MiddlewareFunc(auditLogger, h.Send))

	// Log collector'ı otomatik başlat
	fmt.Println("🔧 Sistem log collector başlatılıyor...")
	if err := logCollector.Start(); err != nil {
		auditLogger.LogError(err, "Log collector başlatma hatası", nil)
		fmt.Printf("⚠️ Log collector başlatılamadı: %v\n", err)
	} else {
		fmt.Println("✅ Sistem log collector başarıyla başlatıldı")
	}

	// Graceful shutdown için signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n🛑 Shutdown sinyali alındı, temiz kapatma işlemi başlatılıyor...")

		// Log collector'ı durdur
		logCollector.Stop()

		// Shutdown audit log
		auditLogger.LogEvent(audit.AuditEvent{
			EventType: "system_shutdown",
			Message:   "Sistem temiz şekilde kapatılıyor",
		})

		os.Exit(0)
	}()

	// Sunucuyu başlat
	fmt.Printf("🌐 Sunucu %s portunda çalışıyor\n", cfg.Port)
	fmt.Println("📋 Endpoints:")
	fmt.Println("  GET  /                    - Ana sayfa")
	fmt.Println("  GET  /api/health          - Sistem sağlık kontrolü")
	fmt.Println("  GET  /api/logs/status     - Log collector durumu")
	fmt.Println("  GET  /api/logs/sources    - Log kaynaklarını listele")
	fmt.Println("  POST /api/logs/start      - Log collector'ı başlat")
	fmt.Println("  POST /api/logs/stop       - Log collector'ı durdur")
	fmt.Println("  POST /api/send            - [DEPRECATED] Mesaj gönder")
	fmt.Println("📊 Sistem log toplama aktif - Loglar console'a yazılıyor")
	fmt.Println("🔍 Takip edilen log dosyaları:")

	// Aktif log kaynaklarını göster
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

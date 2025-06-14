package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gonder/pkg/audit"
)

// Handler HTTP handler'larını içerir
type Handler struct {
	auditLogger *audit.Logger
}

// New yeni bir handler örneği oluşturur
func New(auditLogger *audit.Logger) *Handler {
	return &Handler{
		auditLogger: auditLogger,
	}
}

// Home ana sayfa handler'ı
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gonder - Sistem Log Toplama Servisi</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #2c3e50; margin-bottom: 10px; }
        .header p { color: #7f8c8d; font-size: 16px; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin: 20px 0; }
        .card { background: #f8f9fa; padding: 20px; border-radius: 8px; border-left: 4px solid #3498db; }
        .card h3 { margin-top: 0; color: #2c3e50; }
        .endpoint { background: white; margin: 10px 0; padding: 15px; border-radius: 5px; border: 1px solid #ddd; }
        .method { display: inline-block; padding: 4px 8px; border-radius: 4px; color: white; font-weight: bold; font-size: 12px; }
        .get { background: #27ae60; }
        .post { background: #e74c3c; }
        .feature-list { list-style: none; padding: 0; }
        .feature-list li { padding: 8px 0; border-bottom: 1px solid #eee; }
        .feature-list li:before { content: "✅ "; color: #27ae60; font-weight: bold; }
        .status-indicator { display: inline-block; width: 10px; height: 10px; border-radius: 50%; margin-right: 8px; }
        .status-active { background: #27ae60; }
        .status-inactive { background: #e74c3c; }
        .btn { background: #3498db; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none; display: inline-block; }
        .btn:hover { background: #2980b9; }
        .btn-danger { background: #e74c3c; }
        .btn-danger:hover { background: #c0392b; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Gonder - Sistem Log Toplama Servisi</h1>
            <p>Gerçek zamanlı sistem log toplama, parsing ve monitoring platformu</p>
        </div>
        
        <div class="grid">
            <div class="card">
                <h3>📊 Log Toplama Özellikleri</h3>
                <ul class="feature-list">
                    <li>Syslog toplama ve parsing</li>
                    <li>Nginx/Apache access log'ları</li>
                    <li>Docker container log'ları</li>
                    <li>Authentication log'ları</li>
                    <li>Real-time log monitoring</li>
                    <li>Structured JSON output</li>
                    <li>Kritik log alerting</li>
                </ul>
            </div>
            
            <div class="card">
                <h3>⚙️ Sistem Durumu</h3>
                <p><span class="status-indicator status-active"></span><strong>Log Collector:</strong> Aktif</p>
                <p><span class="status-indicator status-active"></span><strong>Audit Logger:</strong> Aktif</p>
                <p><span class="status-indicator status-active"></span><strong>API Server:</strong> Çalışıyor</p>
                <br>
                <a href="/api/logs/start" class="btn">Log Collector Başlat</a>
                <a href="/api/logs/stop" class="btn btn-danger">Log Collector Durdur</a>
            </div>
        </div>

        <h2>📋 API Endpoints</h2>
        
        <div class="endpoint">
            <span class="method get">GET</span> <strong>/</strong> - Ana sayfa
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span> <strong>/api/logs/status</strong> - Log collector durumu
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span> <strong>/api/logs/sources</strong> - Log kaynaklarını listele
        </div>
        
        <div class="endpoint">
            <span class="method post">POST</span> <strong>/api/logs/start</strong> - Log collector'ı başlat
        </div>
        
        <div class="endpoint">
            <span class="method post">POST</span> <strong>/api/logs/stop</strong> - Log collector'ı durdur
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span> <strong>/api/health</strong> - Sistem sağlık kontrolü
        </div>

        <h2>🧪 Test Komutları</h2>
        <div class="card">
            <h3>Log Collector Yönetimi</h3>
            <pre>
# Log collector durumunu kontrol et
curl http://localhost:8080/api/logs/status

# Log collector'ı başlat
curl -X POST http://localhost:8080/api/logs/start

# Log kaynaklarını listele
curl http://localhost:8080/api/logs/sources

# Log collector'ı durdur
curl -X POST http://localhost:8080/api/logs/stop
            </pre>
        </div>
        
        <h2>📈 Log Formatları</h2>
        <div class="card">
            <h3>Sistem Log Örneği</h3>
            <pre>[SYSTEM_LOG] {
  "id": "log_1749941868123456789",
  "timestamp": "2025-06-15T01:57:48+03:00",
  "source": "syslog",
  "level": "info",
  "message": "systemd[1]: Started nginx.service",
  "host": "server01",
  "service": "systemd",
  "raw_log": "Jun 15 01:57:48 server01 systemd[1]: Started nginx.service",
  "parsed_data": {
    "timestamp": "Jun 15 01:57:48",
    "host": "server01", 
    "service": "systemd",
    "pid": "1",
    "message": "Started nginx.service"
  },
  "tags": ["system", "syslog"],
  "collected_at": "2025-06-15T01:57:48+03:00"
}</pre>
        </div>
        
        <h2>🎯 Desteklenen Log Kaynakları</h2>
        <div class="grid">
            <div class="card">
                <h3>Sistem Log'ları</h3>
                <ul>
                    <li>/var/log/syslog</li>
                    <li>/var/log/messages</li>
                    <li>/var/log/auth.log</li>
                    <li>/var/log/kern.log</li>
                </ul>
            </div>
            <div class="card">
                <h3>Uygulama Log'ları</h3>
                <ul>
                    <li>Nginx access/error log</li>
                    <li>Apache access/error log</li>
                    <li>Docker container logs</li>
                    <li>Custom application logs</li>
                </ul>
            </div>
        </div>
    </div>
</body>
</html>`
	fmt.Fprint(w, html)
}

// SendRequest mesaj gönderme isteği (eski)
type SendRequest struct {
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
	Type      string `json:"type,omitempty"` // email, sms, etc.
}

// SendResponse mesaj gönderme yanıtı (eski)
type SendResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ID        string `json:"id,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Send mesaj gönderme handler'ı (eski, backward compatibility için)
func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Hata audit log
		h.auditLogger.LogError(err, "JSON decode error in Send endpoint", map[string]interface{}{
			"path":        r.URL.Path,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
		})
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basit validasyon
	if req.Message == "" {
		h.auditLogger.LogError(fmt.Errorf("message field is empty"), "Validation error in Send endpoint", map[string]interface{}{
			"request": req,
		})
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	if req.Recipient == "" {
		h.auditLogger.LogError(fmt.Errorf("recipient field is empty"), "Validation error in Send endpoint", map[string]interface{}{
			"request": req,
		})
		http.Error(w, "Recipient is required", http.StatusBadRequest)
		return
	}

	// Mesaj tipi varsayılan
	if req.Type == "" {
		req.Type = "email"
	}

	// Mesaj ID oluştur
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())

	// Eski işlevsellik - artık log toplama odaklı olduğumuz için deprecated
	fmt.Printf("📤 [DEPRECATED] Mesaj gönderiliyor: %s -> %s\n", req.Message, req.Recipient)

	// Mesaj gönderimi audit log
	h.auditLogger.LogMessageSent(req.Recipient, req.Type, messageID, true, map[string]interface{}{
		"message_length": len(req.Message),
		"message_preview": func() string {
			if len(req.Message) > 50 {
				return req.Message[:50] + "..."
			}
			return req.Message
		}(),
		"deprecated": true,
		"note":       "Bu endpoint artık deprecated. Sistem log toplama odaklı çalışıyoruz.",
	})

	response := SendResponse{
		Success:   true,
		Message:   "Mesaj gönderildi (deprecated feature)",
		ID:        messageID,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthResponse sağlık kontrolü yanıtı
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime"`
	Purpose   string `json:"purpose"`
}

// Health sağlık kontrolü handler'ı
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	// Health check audit log
	h.auditLogger.LogHealthCheck("healthy", map[string]interface{}{
		"version":    "2.0.0",
		"purpose":    "sistem_log_toplama",
		"check_time": time.Now().Format(time.RFC3339),
	})

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "2.0.0",
		Uptime:    "N/A", // Burada gerçek uptime hesaplanabilir
		Purpose:   "Sistem Log Toplama Servisi",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package handler

import (
	"encoding/json"
	"net/http"

	"gonder/pkg/collector"
)

// LogHandler log toplama ile ilgili handler'ları içerir
type LogHandler struct {
	collector *collector.LogCollector
}

// NewLogHandler yeni bir log handler oluşturur
func NewLogHandler(collector *collector.LogCollector) *LogHandler {
	return &LogHandler{
		collector: collector,
	}
}

// GetSources log kaynaklarını döner
func (lh *LogHandler) GetSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sources := lh.collector.GetSources()

	response := map[string]interface{}{
		"success": true,
		"data":    sources,
		"count":   len(sources),
		"running": lh.collector.IsRunning(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StartCollector log collector'ı başlatır
func (lh *LogHandler) StartCollector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if lh.collector.IsRunning() {
		response := map[string]interface{}{
			"success": false,
			"message": "Log collector zaten çalışıyor",
			"running": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	err := lh.collector.Start()
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"message": "Log collector başlatılamadı: " + err.Error(),
			"running": false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Log collector başarıyla başlatıldı",
		"running": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopCollector log collector'ı durdurur
func (lh *LogHandler) StopCollector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !lh.collector.IsRunning() {
		response := map[string]interface{}{
			"success": false,
			"message": "Log collector zaten durmuş durumda",
			"running": false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	lh.collector.Stop()

	response := map[string]interface{}{
		"success": true,
		"message": "Log collector başarıyla durduruldu",
		"running": false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetStatus log collector durumunu döner
func (lh *LogHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sources := lh.collector.GetSources()
	enabledCount := 0
	for _, source := range sources {
		if source.Enabled {
			enabledCount++
		}
	}

	response := map[string]interface{}{
		"success": true,
		"status": map[string]interface{}{
			"running":         lh.collector.IsRunning(),
			"total_sources":   len(sources),
			"enabled_sources": enabledCount,
			"sources":         sources,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

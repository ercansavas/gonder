package handler

import (
	"encoding/json"
	"net/http"

	"gonder/pkg/collector"
)

// LogHandler contains handlers for log collection
type LogHandler struct {
	collector *collector.LogCollector
}

// NewLogHandler creates a new log handler
func NewLogHandler(collector *collector.LogCollector) *LogHandler {
	return &LogHandler{
		collector: collector,
	}
}

// GetSources returns log sources
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

// StartCollector starts the log collector
func (lh *LogHandler) StartCollector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if lh.collector.IsRunning() {
		response := map[string]interface{}{
			"success": false,
			"message": "Log collector is already running",
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
			"message": "Log collector could not be started: " + err.Error(),
			"running": false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Log collector started successfully",
		"running": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopCollector stops the log collector
func (lh *LogHandler) StopCollector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !lh.collector.IsRunning() {
		response := map[string]interface{}{
			"success": false,
			"message": "Log collector is already stopped",
			"running": false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	lh.collector.Stop()

	response := map[string]interface{}{
		"success": true,
		"message": "Log collector stopped successfully",
		"running": false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetStatus returns log collector status
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

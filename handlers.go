package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// setNoCacheHeaders instructs browsers not to hold onto stale assets during development
func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, must-revalidate, max-age=0")
}
// enableCORS handles cross-origin headers and preflight OPTIONS requests
func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// respondJSON helper writes typed JSON responses
func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
// HandleStatus returns the system and Gemini AI readiness status
func HandleStatus(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)
		if enableCORS(w, r) {
			return
		}
		respondJSON(w, http.StatusOK, StatusResponse{
			Status:       "ok",
			HasGeminiKey: cfg.HasValidKey(),
			Model:        cfg.Model,
			Port:         cfg.Port,
		})
	}
}

// HandleChat processes chat messages and returns responses from Gemini or local fallback
func HandleChat(cfg *Config, client *GeminiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)
		if enableCORS(w, r) {
			return
		}

		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
			return
		}

		msg := strings.TrimSpace(req.Message)
		if msg == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Message cannot be empty"})
			return
		}

		// Re-check config key dynamically in case user updated .env without restart
		loadEnvFile(".env")
		loadEnvFile(filepath.Join("..", ".env"))
		apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
		hasKey := apiKey != "" && apiKey != "your_gemini_api_key_here"

		if !hasKey {
			// Serve profile knowledge base fallback
			reply := GetFallbackReply(msg)
			respondJSON(w, http.StatusOK, ChatResponse{
				Reply:        reply,
				Model:        "profile-knowledge-base",
				HasGeminiKey: false,
			})
			return
		}

		// Call Gemini API
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		reply, err := client.GenerateReply(ctx, msg, req.History)
		if err != nil {
			log.Printf("Gemini API call failed: %v", err)
			fallbackReply := GetFallbackReply(msg)
			respondJSON(w, http.StatusOK, ChatResponse{
				Reply:        fallbackReply + "\n\n*(Note: Gemini API temporarily unavailable: " + err.Error() + ")*",
				Model:        cfg.Model,
				HasGeminiKey: true,
				Error:        err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, ChatResponse{
			Reply:        reply,
			Model:        cfg.Model,
			HasGeminiKey: true,
		})
	}
}

// HandleStatic serves static assets from cv_frontend directory with no-cache headers
func HandleStatic(frontendDir string) http.Handler {
	fs := http.FileServer(http.Dir(frontendDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)
		fs.ServeHTTP(w, r)
	})
}

// HandleIndex serves index.html from cv_frontend
func HandleIndex(frontendDir string) http.HandlerFunc {
	indexPath := filepath.Join(frontendDir, "index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)
		if r.URL.Path != "/" {
			// Serve asset if it exists in frontendDir
			target := filepath.Join(frontendDir, r.URL.Path)
			if _, err := os.Stat(target); err == nil {
				http.ServeFile(w, r, target)
				return
			}
		}
		http.ServeFile(w, r, indexPath)
	}
}

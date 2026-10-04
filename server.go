package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Server coordinates HTTP routing and external API integration
type Server struct {
	Config       *Config
	GeminiClient *GeminiClient
	Mux          *http.ServeMux
}

// NewServer initializes a new Server instance
func NewServer(cfg *Config) *Server {
	geminiClient := NewGeminiClient(cfg.APIKey, cfg.Model, cfg.SecondaryModels)
	s := &Server{
		Config:       cfg,
		GeminiClient: geminiClient,
		Mux:          http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

// setupRoutes registers all HTTP handlers for APIs and frontend assets
func (s *Server) setupRoutes() {
	// API routes
	s.Mux.HandleFunc("/api/status", HandleStatus(s.Config))
	s.Mux.HandleFunc("/api/chat", HandleChat(s.Config, s.GeminiClient))

	// Frontend static assets
	staticHandler := HandleStatic(s.Config.FrontendDir)
	s.Mux.Handle("/css/", staticHandler)
	s.Mux.Handle("/js/", staticHandler)
	s.Mux.Handle("/partials/", staticHandler)
	s.Mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Root index handler
	s.Mux.HandleFunc("/", HandleIndex(s.Config.FrontendDir))
}

// Start launches the HTTP listener on the configured port
// corsAndCleanMiddleware ensures CORS headers on every response and eliminates double-slash redirects
func corsAndCleanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always attach CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS immediately with 200 OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Clean multiple consecutive slashes from URL path to prevent 301 redirects
		for strings.Contains(r.URL.Path, "//") {
			r.URL.Path = strings.ReplaceAll(r.URL.Path, "//", "/")
		}

		next.ServeHTTP(w, r)
	})
}

// Start launches the HTTP listener on the configured port
func (s *Server) Start() error {
	addr := ":" + s.Config.Port
	printBanner(s.Config)
	return http.ListenAndServe(addr, corsAndCleanMiddleware(s.Mux))
}

// printBanner displays the startup status
func printBanner(cfg *Config) {
	fmt.Printf("====================================================\n")
	fmt.Printf("Mohammad Wildan Kosasih - Interactive CV & AI Twin\n")
	fmt.Printf("Server listening on http://localhost:%s\n", cfg.Port)
	fmt.Printf("Frontend directory: %s\n", cfg.FrontendDir)
	fmt.Printf("Using Gemini Model: %s\n", cfg.Model)
	if len(cfg.SecondaryModels) > 0 {
		fmt.Printf("Secondary Models  : %s\n", strings.Join(cfg.SecondaryModels, ", "))
	}
	if cfg.HasValidKey() {
		fmt.Printf("Gemini API Key: Configured\n")
	} else {
		fmt.Printf("Gemini API Key: Not configured (using knowledge-base fallback)\n")
		fmt.Printf("To configure: add GEMINI_API_KEY to your .env file\n")
	}
	fmt.Printf("====================================================\n")
}

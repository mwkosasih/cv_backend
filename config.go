package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultPort        = "8081"
	DefaultGeminiModel = "gemini-3.5-flash-lite"
	DefaultFrontendDir = "./cv_frontend"
)

// Config holds runtime configuration for the CV backend service
type Config struct {
	Port            string
	Model           string
	SecondaryModels []string
	APIKey          string
	FrontendDir     string
}

// LoadConfig initializes configuration by reading .env and environment variables
func LoadConfig() *Config {
	// Look for .env in current directory or parent directory
	loadEnvFile(".env")
	loadEnvFile(filepath.Join("cv_backend", ".env"))
	if os.Getenv("GEMINI_API_KEY") == "" {
		loadEnvFile(filepath.Join("..", ".env"))
	}

	port := getEnvWithDefault("PORT", DefaultPort)
	model := getEnvWithDefault("GEMINI_MODEL", DefaultGeminiModel)
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))

	// Load optional secondary models (GEMINI_MODEL_2, GEMINI_MODEL_3, GEMINI_MODEL_4)
	var secondaryModels []string
	for _, envKey := range []string{"GEMINI_MODEL_2", "GEMINI_MODEL_3", "GEMINI_MODEL_4"} {
		if val := strings.TrimSpace(os.Getenv(envKey)); val != "" {
			secondaryModels = append(secondaryModels, val)
		}
	}

	frontendDir := getEnvWithDefault("FRONTEND_DIR", DefaultFrontendDir)
	// Check if frontend directory exists relative to current path
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		// Fallback to ../cv_frontend if run from cv_backend directory
		if _, err := os.Stat("../cv_frontend"); err == nil {
			frontendDir = "../cv_frontend"
		}
	}

	return &Config{
		Port:            port,
		Model:           model,
		SecondaryModels: secondaryModels,
		APIKey:          apiKey,
		FrontendDir:     frontendDir,
	}
}

// HasValidKey checks if a valid Gemini API key is configured
func (c *Config) HasValidKey() bool {
	return c.APIKey != "" && c.APIKey != "your_gemini_api_key_here"
}

// loadEnvFile parses key=value pairs from a .env file without external dependencies
func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	err = scanner.Err()
	if err != nil {
		return
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = trimQuotes(val)

		if key != "" && os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	if err := scanner.Err(); err != nil {
		return
	}
}

// trimQuotes strips optional leading and trailing quotes from values
func trimQuotes(val string) string {
	if len(val) >= 2 {
		if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
			return val[1 : len(val)-1]
		}
	}
	return val
}

// getEnvWithDefault retrieves an environment variable or returns a fallback value
func getEnvWithDefault(key, fallback string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	return val
}

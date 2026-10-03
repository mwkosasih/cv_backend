package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ChatMessage represents a single message in the conversational history
type ChatMessage struct {
	Role string `json:"role"` // "user" or "model"
	Text string `json:"text"`
}

// ChatRequest is the incoming payload for /api/chat
type ChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history,omitempty"`
}

// ChatResponse is the outgoing payload for /api/chat
type ChatResponse struct {
	Reply        string `json:"reply"`
	Model        string `json:"model"`
	HasGeminiKey bool   `json:"has_gemini_key"`
	Error        string `json:"error,omitempty"`
}

// StatusResponse is the payload for /api/status
type StatusResponse struct {
	Status       string `json:"status"`
	HasGeminiKey bool   `json:"has_gemini_key"`
	Model        string `json:"model"`
	Port         string `json:"port"`
}

// Internal Gemini API structures
type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiSystemInstruction struct {
	Parts []geminiPart `json:"parts"`
}

type geminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiRequestBody struct {
	SystemInstruction *geminiSystemInstruction `json:"system_instruction,omitempty"`
	Contents          []geminiContent          `json:"contents"`
	GenerationConfig  *geminiGenerationConfig  `json:"generationConfig,omitempty"`
}

type geminiCandidate struct {
	Content struct {
		Parts []geminiPart `json:"parts"`
		Role  string       `json:"role"`
	} `json:"content"`
	FinishReason string `json:"finishReason"`
}

type geminiResponseBody struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// GeminiClient handles communication with Google's Generative Language API
type GeminiClient struct {
	apiKey          string
	model           string
	secondaryModels []string
	httpClient      *http.Client
}

// NewGeminiClient instantiates a new GeminiClient
func NewGeminiClient(apiKey, model string, secondaryModels []string) *GeminiClient {
	return &GeminiClient{
		apiKey:          apiKey,
		model:           model,
		secondaryModels: secondaryModels,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

// GenerateReply calls Gemini with fallback model support
func (c *GeminiClient) GenerateReply(ctx context.Context, userMsg string, history []ChatMessage) (string, error) {
	modelsToTry := []string{c.model}
	for _, m := range c.secondaryModels {
		m = strings.TrimSpace(m)
		if m != "" && m != c.model {
			modelsToTry = append(modelsToTry, m)
		}
	}
	var lastErr error
	for _, m := range modelsToTry {
		reply, err := c.tryCallModel(ctx, m, userMsg, history)
		if err == nil && reply != "" {
			return reply, nil
		}
		lastErr = err
		log.Printf("Gemini call for model %s failed: %v", m, err)
	}

	return "", lastErr
}

// tryCallModel sends a single HTTP request to Gemini API endpoint
func (c *GeminiClient) tryCallModel(ctx context.Context, model, userMsg string, history []ChatMessage) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, c.apiKey)

	payloadBytes, err := buildGeminiPayload(userMsg, history)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	return parseGeminiResponse(bodyBytes)
}

// buildGeminiPayload constructs the JSON request body
func buildGeminiPayload(userMsg string, history []ChatMessage) ([]byte, error) {
	var contents []geminiContent

	// Add recent conversational history (up to last 10 messages)
	startIdx := 0
	if len(history) > 10 {
		startIdx = len(history) - 10
	}
	for _, h := range history[startIdx:] {
		role := "user"
		if h.Role == "model" || h.Role == "assistant" {
			role = "model"
		}
		trimmed := strings.TrimSpace(h.Text)
		if trimmed != "" {
			contents = append(contents, geminiContent{
				Role:  role,
				Parts: []geminiPart{{Text: trimmed}},
			})
		}
	}

	// Add the current user message
	contents = append(contents, geminiContent{
		Role:  "user",
		Parts: []geminiPart{{Text: strings.TrimSpace(userMsg)}},
	})

	reqPayload := geminiRequestBody{
		SystemInstruction: &geminiSystemInstruction{
			Parts: []geminiPart{{Text: GetSystemPrompt()}},
		},
		Contents: contents,
		GenerationConfig: &geminiGenerationConfig{
			Temperature:     0.7,
			MaxOutputTokens: 1000,
		},
	}

	return json.Marshal(reqPayload)
}

// parseGeminiResponse decodes the response body from Gemini
func parseGeminiResponse(bodyBytes []byte) (string, error) {
	var geminiResp geminiResponseBody
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return "", fmt.Errorf("unmarshal gemini response: %w (body: %s)", err, string(bodyBytes))
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("gemini api error (code %d): %s", geminiResp.Error.Code, geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("empty response from gemini: %s", string(bodyBytes))
}

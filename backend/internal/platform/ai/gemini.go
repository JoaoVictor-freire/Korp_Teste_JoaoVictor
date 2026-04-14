package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"korp_backend/internal/platform/config"
)

var ErrGeminiNotConfigured = errors.New("GEMINI_API_KEY is not configured")

type GeminiClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type geminiGenerateContentRequest struct {
	SystemInstruction geminiContent          `json:"system_instruction"`
	Contents          []geminiContent        `json:"contents"`
	GenerationConfig  geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiGenerateContentResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewGeminiClientFromEnv() *GeminiClient {
	timeout := time.Duration(config.GetEnvAsInt("GEMINI_TIMEOUT_MS", 12000)) * time.Millisecond

	return &GeminiClient{
		apiKey: config.GetEnvTrimmed("GEMINI_API_KEY", ""),
		model:  config.GetEnvTrimmed("GEMINI_MODEL", "gemini-2.5-flash"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *GeminiClient) GenerateOperationalInsights(ctx context.Context, prompt string) (string, string, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return "", "", ErrGeminiNotConfigured
	}

	requestBody := geminiGenerateContentRequest{
		SystemInstruction: geminiContent{
			Parts: []geminiPart{
				{
					Text: "Voce e um assistente de operacoes para um sistema de estoque e faturamento. Responda sempre em portugues do Brasil, seja objetivo, nao invente dados e deixe claro quando uma observacao for apenas uma inferencia.",
				},
			},
		},
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:     0.4,
			MaxOutputTokens: 700,
		},
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var parsed geminiGenerateContentResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return "", "", err
	}

	if response.StatusCode >= http.StatusBadRequest {
		message := "failed to generate insights with Gemini"
		if parsed.Error != nil && strings.TrimSpace(parsed.Error.Message) != "" {
			message = parsed.Error.Message
		}
		return "", "", errors.New(message)
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", "", errors.New("Gemini returned an empty response")
	}

	textParts := make([]string, 0, len(parsed.Candidates[0].Content.Parts))
	for _, part := range parsed.Candidates[0].Content.Parts {
		if strings.TrimSpace(part.Text) != "" {
			textParts = append(textParts, strings.TrimSpace(part.Text))
		}
	}

	if len(textParts) == 0 {
		return "", "", errors.New("Gemini returned an empty response")
	}

	return strings.Join(textParts, "\n\n"), c.model, nil
}

package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"korp_backend/internal/platform/config"

	"google.golang.org/genai"
)

var ErrGeminiNotConfigured = errors.New("GEMINI_API_KEY is not configured")

type GeminiClient struct {
	apiKey       string
	model        string
	httpClient   *http.Client
	sdkClient    *genai.Client
	initErr      error
	timeoutValue time.Duration
}

type GeneratedOperationalInsights struct {
	Overview           string
	Alerts             []string
	Actions            []string
	BillingNotes       []string
	BuyRecommendations []BuyRecommendation
	SearchQueries      []string
	Sources            []GroundingSource
}

type BuyRecommendation struct {
	Name          string
	Category      string
	Reason        string
	MarketSignal  string
	StockRelation string
}

type GroundingSource struct {
	Title string
	URI   string
}

type geminiStructuredInsights struct {
	Overview           string                    `json:"overview"`
	Alerts             []string                  `json:"alerts"`
	Actions            []string                  `json:"actions"`
	BillingNotes       []string                  `json:"billing_notes"`
	BuyRecommendations []geminiRecommendationDTO `json:"buy_recommendations"`
}

type geminiRecommendationDTO struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Reason        string `json:"reason"`
	MarketSignal  string `json:"market_signal"`
	StockRelation string `json:"stock_relation"`
}

func NewGeminiClientFromEnv() *GeminiClient {
	timeout := time.Duration(config.GetEnvAsInt("GEMINI_TIMEOUT_MS", 30000)) * time.Millisecond
	apiKey := config.GetEnvTrimmed("GEMINI_API_KEY", "")
	model := config.GetEnvTrimmed("GEMINI_MODEL", "gemini-2.5-flash")
	httpClient := &http.Client{Timeout: timeout}

	client := &GeminiClient{
		apiKey:       apiKey,
		model:        model,
		httpClient:   httpClient,
		timeoutValue: timeout,
	}

	if apiKey == "" {
		return client
	}

	sdkClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:     apiKey,
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: httpClient,
	})
	if err != nil {
		client.initErr = err
		return client
	}

	client.sdkClient = sdkClient
	return client
}

func (c *GeminiClient) GenerateOperationalInsights(ctx context.Context, prompt string) (GeneratedOperationalInsights, string, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return GeneratedOperationalInsights{}, "", ErrGeminiNotConfigured
	}
	if c.initErr != nil {
		return GeneratedOperationalInsights{}, "", c.initErr
	}

	groundedText, groundingMetadata, err := c.generateGroundedContext(ctx, prompt)
	if err != nil {
		return GeneratedOperationalInsights{}, "", err
	}

	jsonPrompt := buildStructuredPrompt(prompt, groundedText)
	structuredText, err := c.generateStructuredJSON(ctx, jsonPrompt)
	if err != nil {
		return GeneratedOperationalInsights{}, "", err
	}

	structured, err := parseStructuredInsights(structuredText)
	if err != nil {
		return GeneratedOperationalInsights{}, "", err
	}

	result := GeneratedOperationalInsights{
		Overview:           strings.TrimSpace(structured.Overview),
		Alerts:             filterNonEmpty(structured.Alerts),
		Actions:            filterNonEmpty(structured.Actions),
		BillingNotes:       filterNonEmpty(structured.BillingNotes),
		BuyRecommendations: make([]BuyRecommendation, 0, len(structured.BuyRecommendations)),
	}

	for _, item := range structured.BuyRecommendations {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		result.BuyRecommendations = append(result.BuyRecommendations, BuyRecommendation{
			Name:          strings.TrimSpace(item.Name),
			Category:      strings.TrimSpace(item.Category),
			Reason:        strings.TrimSpace(item.Reason),
			MarketSignal:  strings.TrimSpace(item.MarketSignal),
			StockRelation: strings.TrimSpace(item.StockRelation),
		})
	}

	result.SearchQueries, result.Sources = extractGroundingMetadata(groundingMetadata)

	return result, c.model, nil
}

func (c *GeminiClient) generateGroundedContext(ctx context.Context, prompt string) (string, *genai.GroundingMetadata, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			"Voce pesquisa tendencias de mercado no Brasil para um sistema de estoque e faturamento. Resuma fatos de forma curta e objetiva.",
			"system",
		),
		Temperature:     genai.Ptr[float32](0.2),
		MaxOutputTokens: 500,
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		},
	}

	response, err := c.sdkClient.Models.GenerateContent(ctx, c.model, genai.Text(prompt), config)
	if err != nil {
		return "", nil, err
	}
	if len(response.Candidates) == 0 || response.Candidates[0].Content == nil {
		return "", nil, errors.New("Gemini returned an empty grounded response")
	}

	return response.Text(), response.Candidates[0].GroundingMetadata, nil
}

func (c *GeminiClient) generateStructuredJSON(ctx context.Context, prompt string) (string, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			"Responda somente com um unico JSON valido. Nao use markdown. Nao use comentarios. Nao escreva nada fora do JSON.",
			"system",
		),
		Temperature:      genai.Ptr[float32](0.1),
		MaxOutputTokens:  5000,
		ResponseMIMEType: "application/json",
		ResponseSchema:   insightsResponseSchema(),
		CandidateCount:   1,
	}

	response, err := c.sdkClient.Models.GenerateContent(ctx, c.model, genai.Text(prompt), config)
	if err != nil {
		return "", err
	}
	if len(response.Candidates) == 0 || response.Candidates[0].Content == nil {
		return "", errors.New("Gemini returned an empty structured response")
	}

	return response.Text(), nil
}

func insightsResponseSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"overview": {
				Type: genai.TypeString,
			},
			"alerts": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"actions": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"billing_notes": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"buy_recommendations": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name":           {Type: genai.TypeString},
						"category":       {Type: genai.TypeString},
						"reason":         {Type: genai.TypeString},
						"market_signal":  {Type: genai.TypeString},
						"stock_relation": {Type: genai.TypeString},
					},
					Required: []string{"name", "category", "reason", "market_signal", "stock_relation"},
				},
			},
		},
		Required: []string{"overview", "alerts", "actions", "billing_notes", "buy_recommendations"},
	}
}

func buildStructuredPrompt(originalPrompt string, groundedText string) string {
	return fmt.Sprintf(
		`Use o contexto abaixo e responda apenas com JSON valido.

FORMATO:
{
  "overview": "string",
  "alerts": ["string"],
  "actions": ["string"],
  "billing_notes": ["string"],
  "buy_recommendations": [
    {
      "name": "string",
      "category": "string",
      "reason": "string",
      "market_signal": "string",
      "stock_relation": "string"
    }
  ]
}

REGRAS:
- Seja extremamente objetivo.
- Nao use quebras de linha dentro de strings JSON.
- "overview" deve ter no maximo 1 frase curta.
- "alerts" deve ter no maximo 3 itens.
- "actions" deve ter no maximo 3 itens.
- "billing_notes" deve ter no maximo 2 itens.
- "buy_recommendations" deve ter exatamente 6 itens.
- Cada item das listas deve ser curto.
- "reason", "market_signal" e "stock_relation" devem ser curtos.
- Nao invente dados.

CONTEXTO INTERNO:
%s

CONTEXTO PESQUISADO NA WEB:
%s`,
		originalPrompt,
		strings.TrimSpace(groundedText),
	)
}

func parseStructuredInsights(text string) (geminiStructuredInsights, error) {
	var structured geminiStructuredInsights
	normalized := normalizeGeminiJSON(text)
	if err := json.Unmarshal([]byte(normalized), &structured); err != nil {
		return geminiStructuredInsights{}, err
	}

	return structured, nil
}

func filterNonEmpty(values []string) []string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			items = append(items, trimmed)
		}
	}

	return items
}

func normalizeGeminiJSON(text string) string {
	return sanitizeJSONStringLiterals(extractJSONObject(text))
}

func extractJSONObject(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)
	return extractFirstJSONObject(trimmed)
}

func extractFirstJSONObject(value string) string {
	start := strings.Index(value, "{")
	if start < 0 {
		return strings.TrimSpace(value)
	}

	depth := 0
	inString := false
	escaped := false

	for i := start; i < len(value); i++ {
		ch := value[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return strings.TrimSpace(value[start : i+1])
			}
		}
	}

	return strings.TrimSpace(value[start:])
}

func sanitizeJSONStringLiterals(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))

	inString := false
	escaped := false

	for i := 0; i < len(value); i++ {
		ch := value[i]
		if inString {
			if escaped {
				builder.WriteByte(ch)
				escaped = false
				continue
			}
			switch ch {
			case '\\':
				builder.WriteByte(ch)
				escaped = true
			case '"':
				builder.WriteByte(ch)
				inString = false
			case '\n', '\r', '\t':
				builder.WriteByte(' ')
			default:
				builder.WriteByte(ch)
			}
			continue
		}

		if ch == '"' {
			inString = true
		}
		builder.WriteByte(ch)
	}

	return builder.String()
}

func countSearchQueries(metadata *genai.GroundingMetadata) int {
	if metadata == nil {
		return 0
	}
	return len(metadata.WebSearchQueries)
}

func countGroundingChunks(metadata *genai.GroundingMetadata) int {
	if metadata == nil {
		return 0
	}
	return len(metadata.GroundingChunks)
}

func extractGroundingMetadata(metadata *genai.GroundingMetadata) ([]string, []GroundingSource) {
	if metadata == nil {
		return nil, nil
	}

	queries := filterNonEmpty(metadata.WebSearchQueries)
	sources := make([]GroundingSource, 0, len(metadata.GroundingChunks))
	seen := make(map[string]struct{})
	for _, chunk := range metadata.GroundingChunks {
		if chunk == nil || chunk.Web == nil {
			continue
		}
		uri := strings.TrimSpace(chunk.Web.URI)
		if uri == "" {
			continue
		}
		if _, exists := seen[uri]; exists {
			continue
		}
		seen[uri] = struct{}{}
		sources = append(sources, GroundingSource{
			Title: strings.TrimSpace(chunk.Web.Title),
			URI:   uri,
		})
	}

	return queries, sources
}

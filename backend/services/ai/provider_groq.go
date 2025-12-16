package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GroqProvider implements AIProvider for Groq
// Groq uses OpenAI-compatible API
type GroqProvider struct {
	apiKey      string
	modelName   string
	maxTokens   int
	temperature float32
	httpClient  *http.Client
}

// NewGroqProvider creates a new Groq provider
func NewGroqProvider(config ProviderConfig) (*GroqProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Groq API key is required")
	}

	modelName := config.Model
	if modelName == "" {
		modelName = "llama-3.3-70b-versatile"
	}

	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 150 // Reduced default for faster response
	}

	temp := config.Temperature
	if temp == 0 {
		temp = 0.7
	}

	// Create HTTP client with timeout for faster failure detection
	httpClient := &http.Client{
		Timeout: 8 * time.Second, // Groq is fast, 8 second timeout should be enough
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	return &GroqProvider{
		apiKey:      config.APIKey,
		modelName:   modelName,
		maxTokens:   maxTokens,
		temperature: temp,
		httpClient:  httpClient,
	}, nil
}

func (g *GroqProvider) GetProviderName() string {
	return "groq"
}

func (g *GroqProvider) GetModelName() string {
	return g.modelName
}

func (g *GroqProvider) GenerateResponse(ctx context.Context, systemPrompt, userMessage, contextInfo string) (*AIResponse, error) {
	fullContext := systemPrompt
	if contextInfo != "" {
		fullContext += "\n\nCONTEXT:\n" + contextInfo
	}
	fullContext += "\n\nIMPORTANT RULES:\n- Respond in Indonesian language\n- Be helpful, friendly, and professional\n- Keep responses concise (max 300 characters)\n- If you don't know something, say so honestly"

	reqBody := openAIRequest{
		Model: g.modelName,
		Messages: []openAIMessage{
			{Role: "system", Content: fullContext},
			{Role: "user", Content: userMessage},
		},
		MaxTokens:   g.maxTokens,
		Temperature: g.temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var groqResp openAIResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if groqResp.Error != nil {
		return nil, fmt.Errorf("Groq API error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Groq")
	}

	responseText := groqResp.Choices[0].Message.Content
	inputTokens := groqResp.Usage.PromptTokens
	outputTokens := groqResp.Usage.CompletionTokens

	// Calculate cost based on model (Groq pricing)
	var inputCost, outputCost float64
	switch g.modelName {
	case "llama-3.3-70b-versatile":
		inputCost = float64(inputTokens) / 1_000_000 * 0.59
		outputCost = float64(outputTokens) / 1_000_000 * 0.79
	case "llama-3.1-8b-instant":
		inputCost = float64(inputTokens) / 1_000_000 * 0.05
		outputCost = float64(outputTokens) / 1_000_000 * 0.08
	case "mixtral-8x7b-32768":
		inputCost = float64(inputTokens) / 1_000_000 * 0.24
		outputCost = float64(outputTokens) / 1_000_000 * 0.24
	default:
		inputCost = float64(inputTokens) / 1_000_000 * 0.20
		outputCost = float64(outputTokens) / 1_000_000 * 0.20
	}

	confidence := 0.90
	if groqResp.Choices[0].FinishReason != "stop" {
		confidence = 0.75
	}

	return &AIResponse{
		Response:     responseText,
		Confidence:   confidence,
		TokensUsed:   inputTokens + outputTokens,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		CostUSD:      inputCost + outputCost,
		Model:        g.modelName,
	}, nil
}

func (g *GroqProvider) Close() {
	// HTTP client doesn't need explicit cleanup
}


package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AnthropicProvider implements AIProvider for Anthropic Claude
type AnthropicProvider struct {
	apiKey      string
	modelName   string
	maxTokens   int
	temperature float32
	httpClient  *http.Client
}

type anthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	System      string             `json:"system,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
	Temperature float32            `json:"temperature,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(config ProviderConfig) (*AnthropicProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	modelName := config.Model
	if modelName == "" {
		modelName = "claude-3-5-haiku-latest"
	}

	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 500
	}

	temp := config.Temperature
	if temp == 0 {
		temp = 0.7
	}

	return &AnthropicProvider{
		apiKey:      config.APIKey,
		modelName:   modelName,
		maxTokens:   maxTokens,
		temperature: temp,
		httpClient:  &http.Client{},
	}, nil
}

func (a *AnthropicProvider) GetProviderName() string {
	return "anthropic"
}

func (a *AnthropicProvider) GetModelName() string {
	return a.modelName
}

func (a *AnthropicProvider) GenerateResponse(ctx context.Context, systemPrompt, userMessage, contextInfo string) (*AIResponse, error) {
	fullSystem := systemPrompt
	if contextInfo != "" {
		fullSystem += "\n\nCONTEXT:\n" + contextInfo
	}
	fullSystem += "\n\nIMPORTANT RULES:\n- Respond in Indonesian language\n- Be helpful, friendly, and professional\n- Keep responses concise (max 300 characters)\n- If you don't know something, say so honestly"

	reqBody := anthropicRequest{
		Model:     a.modelName,
		MaxTokens: a.maxTokens,
		System:    fullSystem,
		Messages: []anthropicMessage{
			{Role: "user", Content: userMessage},
		},
		Temperature: a.temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if anthropicResp.Error != nil {
		return nil, fmt.Errorf("Anthropic API error: %s", anthropicResp.Error.Message)
	}

	if len(anthropicResp.Content) == 0 {
		return nil, fmt.Errorf("no response from Anthropic")
	}

	responseText := ""
	for _, content := range anthropicResp.Content {
		if content.Type == "text" {
			responseText += content.Text
		}
	}

	inputTokens := anthropicResp.Usage.InputTokens
	outputTokens := anthropicResp.Usage.OutputTokens

	// Calculate cost based on model (Anthropic pricing)
	var inputCost, outputCost float64
	switch a.modelName {
	case "claude-3-5-haiku-latest":
		inputCost = float64(inputTokens) / 1_000_000 * 0.80
		outputCost = float64(outputTokens) / 1_000_000 * 4.0
	case "claude-3-5-sonnet-latest":
		inputCost = float64(inputTokens) / 1_000_000 * 3.0
		outputCost = float64(outputTokens) / 1_000_000 * 15.0
	default:
		inputCost = float64(inputTokens) / 1_000_000 * 0.80
		outputCost = float64(outputTokens) / 1_000_000 * 4.0
	}

	confidence := 0.90
	if anthropicResp.StopReason != "end_turn" {
		confidence = 0.75
	}

	return &AIResponse{
		Response:     responseText,
		Confidence:   confidence,
		TokensUsed:   inputTokens + outputTokens,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		CostUSD:      inputCost + outputCost,
		Model:        a.modelName,
	}, nil
}

func (a *AnthropicProvider) Close() {
	// HTTP client doesn't need explicit cleanup
}


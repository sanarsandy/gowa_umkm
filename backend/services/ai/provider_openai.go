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

// OpenAIProvider implements AIProvider for OpenAI
type OpenAIProvider struct {
	apiKey      string
	modelName   string
	maxTokens   int
	temperature float32
	httpClient  *http.Client
}

type openAIRequest struct {
	Model       string           `json:"model"`
	Messages    []openAIMessage  `json:"messages"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Temperature float32          `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	modelName := config.Model
	if modelName == "" {
		modelName = "gpt-4o-mini"
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
		Timeout: 10 * time.Second, // 10 second timeout for API calls
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	return &OpenAIProvider{
		apiKey:      config.APIKey,
		modelName:   modelName,
		maxTokens:   maxTokens,
		temperature: temp,
		httpClient:  httpClient,
	}, nil
}

func (o *OpenAIProvider) GetProviderName() string {
	return "openai"
}

func (o *OpenAIProvider) GetModelName() string {
	return o.modelName
}

func (o *OpenAIProvider) GenerateResponse(ctx context.Context, systemPrompt, userMessage, contextInfo string) (*AIResponse, error) {
	fullContext := systemPrompt
	if contextInfo != "" {
		fullContext += "\n\nCONTEXT:\n" + contextInfo
	}
	fullContext += "\n\nIMPORTANT RULES:\n- Respond in Indonesian language\n- Be helpful, friendly, and professional\n- Keep responses concise (max 300 characters)\n- If you don't know something, say so honestly"

	reqBody := openAIRequest{
		Model: o.modelName,
		Messages: []openAIMessage{
			{Role: "system", Content: fullContext},
			{Role: "user", Content: userMessage},
		},
		MaxTokens:   o.maxTokens,
		Temperature: o.temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseText := openAIResp.Choices[0].Message.Content
	inputTokens := openAIResp.Usage.PromptTokens
	outputTokens := openAIResp.Usage.CompletionTokens

	// Calculate cost based on model
	var inputCost, outputCost float64
	switch o.modelName {
	case "gpt-4o-mini":
		inputCost = float64(inputTokens) / 1_000_000 * 0.15
		outputCost = float64(outputTokens) / 1_000_000 * 0.60
	case "gpt-4o":
		inputCost = float64(inputTokens) / 1_000_000 * 2.50
		outputCost = float64(outputTokens) / 1_000_000 * 10.0
	default:
		inputCost = float64(inputTokens) / 1_000_000 * 0.50
		outputCost = float64(outputTokens) / 1_000_000 * 1.50
	}

	confidence := 0.90
	if openAIResp.Choices[0].FinishReason != "stop" {
		confidence = 0.75
	}

	return &AIResponse{
		Response:     responseText,
		Confidence:   confidence,
		TokensUsed:   inputTokens + outputTokens,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		CostUSD:      inputCost + outputCost,
		Model:        o.modelName,
	}, nil
}

func (o *OpenAIProvider) Close() {
	// HTTP client doesn't need explicit cleanup
}


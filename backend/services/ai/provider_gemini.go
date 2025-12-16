package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiProvider implements AIProvider for Google Gemini
type GeminiProvider struct {
	client    *genai.Client
	model     *genai.GenerativeModel
	modelName string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(config ProviderConfig) (*GeminiProvider, error) {
	ctx := context.Background()
	
	if config.APIKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	modelName := config.Model
	if modelName == "" {
		// Use faster model by default (gemini-1.5-flash is faster than 2.0-flash)
		modelName = "gemini-1.5-flash"
	}

	model := client.GenerativeModel(modelName)

	// Configure model parameters
	temp := config.Temperature
	if temp == 0 {
		temp = 0.7
	}
	model.SetTemperature(temp)
	model.SetTopK(40)
	model.SetTopP(0.95)
	
	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		// Reduce default max tokens for faster response (150 is enough for concise replies)
		maxTokens = 150
	}
	model.SetMaxOutputTokens(int32(maxTokens))

	// Safety settings
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockOnlyHigh},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockOnlyHigh},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockOnlyHigh},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockOnlyHigh},
	}

	return &GeminiProvider{
		client:    client,
		model:     model,
		modelName: modelName,
	}, nil
}

func (g *GeminiProvider) GetProviderName() string {
	return "gemini"
}

func (g *GeminiProvider) GetModelName() string {
	return g.modelName
}

func (g *GeminiProvider) GenerateResponse(ctx context.Context, systemPrompt, userMessage, contextInfo string) (*AIResponse, error) {
	// Add timeout to context (8 seconds for Gemini)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	fullPrompt := buildPrompt(systemPrompt, contextInfo, userMessage)

	resp, err := g.model.GenerateContent(ctxWithTimeout, genai.Text(fullPrompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response candidates returned")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response content")
	}

	responseText := extractTextFromParts(candidate.Content.Parts)
	inputTokens := estimateTokens(fullPrompt)
	outputTokens := estimateTokens(responseText)
	totalTokens := inputTokens + outputTokens

	// Gemini Flash pricing
	inputCost := float64(inputTokens) / 1_000_000 * 0.075
	outputCost := float64(outputTokens) / 1_000_000 * 0.30
	totalCost := inputCost + outputCost

	confidence := calculateConfidence(candidate)

	return &AIResponse{
		Response:     responseText,
		Confidence:   confidence,
		TokensUsed:   totalTokens,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		CostUSD:      totalCost,
		Model:        g.modelName,
	}, nil
}

func (g *GeminiProvider) Close() {
	if g.client != nil {
		g.client.Close()
	}
}

// Helper functions
func extractTextFromParts(parts []genai.Part) string {
	var sb strings.Builder
	for _, part := range parts {
		if text, ok := part.(genai.Text); ok {
			sb.WriteString(string(text))
		}
	}
	return sb.String()
}

func calculateConfidence(candidate *genai.Candidate) float64 {
	baseConfidence := 0.85
	switch candidate.FinishReason {
	case genai.FinishReasonStop:
		baseConfidence = 0.90
	case genai.FinishReasonMaxTokens:
		baseConfidence = 0.75
	case genai.FinishReasonSafety:
		baseConfidence = 0.50
	case genai.FinishReasonRecitation:
		baseConfidence = 0.60
	case genai.FinishReasonOther:
		baseConfidence = 0.70
	}
	return baseConfidence
}


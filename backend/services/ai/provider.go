package ai

import (
	"context"
)

// AIResponse is the response from an AI provider
type AIResponse struct {
	Response     string
	Confidence   float64
	TokensUsed   int
	InputTokens  int
	OutputTokens int
	CostUSD      float64
	Model        string
}

// AIProvider defines the interface for AI providers
type AIProvider interface {
	// GenerateResponse generates a response from the AI
	GenerateResponse(ctx context.Context, systemPrompt, userMessage, contextInfo string) (*AIResponse, error)
	// GetProviderName returns the provider name
	GetProviderName() string
	// GetModelName returns the current model name
	GetModelName() string
	// Close cleans up resources
	Close()
}

// ProviderConfig contains configuration for creating a provider
type ProviderConfig struct {
	Provider   string
	APIKey     string
	Model      string
	MaxTokens  int
	Temperature float32
}

// ProviderInfo contains information about an AI provider
type ProviderInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Models      []ModelInfo `json:"models"`
	RequiresKey bool     `json:"requires_key"`
	FreeAvailable bool   `json:"free_available"`
}

// ModelInfo contains information about a model
type ModelInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	InputCost   float64 `json:"input_cost"`  // per 1M tokens
	OutputCost  float64 `json:"output_cost"` // per 1M tokens
}

// GetAvailableProviders returns list of available AI providers
func GetAvailableProviders() []ProviderInfo {
	return []ProviderInfo{
		{
			ID:          "gemini",
			Name:        "Google Gemini",
			Description: "Google's multimodal AI model, great for general use",
			RequiresKey: true,
			FreeAvailable: true,
			Models: []ModelInfo{
				{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash", Description: "Fast and efficient", InputCost: 0.075, OutputCost: 0.30},
				{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Description: "Latest flash model", InputCost: 0.075, OutputCost: 0.30},
				{ID: "gemini-2.0-flash-lite", Name: "Gemini 2.0 Flash Lite", Description: "Lightweight version", InputCost: 0.0375, OutputCost: 0.15},
			},
		},
		{
			ID:          "openai",
			Name:        "OpenAI (ChatGPT)",
			Description: "OpenAI's powerful GPT models",
			RequiresKey: true,
			FreeAvailable: false,
			Models: []ModelInfo{
				{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Description: "Fast and affordable", InputCost: 0.15, OutputCost: 0.60},
				{ID: "gpt-4o", Name: "GPT-4o", Description: "Most capable model", InputCost: 2.50, OutputCost: 10.0},
				{ID: "gpt-3.5-turbo", Name: "GPT-3.5 Turbo", Description: "Legacy fast model", InputCost: 0.50, OutputCost: 1.50},
			},
		},
		{
			ID:          "groq",
			Name:        "Groq",
			Description: "Ultra-fast inference with open models",
			RequiresKey: true,
			FreeAvailable: true,
			Models: []ModelInfo{
				{ID: "llama-3.3-70b-versatile", Name: "Llama 3.3 70B", Description: "Most capable open model", InputCost: 0.59, OutputCost: 0.79},
				{ID: "llama-3.1-8b-instant", Name: "Llama 3.1 8B", Description: "Fast and efficient", InputCost: 0.05, OutputCost: 0.08},
				{ID: "mixtral-8x7b-32768", Name: "Mixtral 8x7B", Description: "MoE architecture", InputCost: 0.24, OutputCost: 0.24},
				{ID: "gemma2-9b-it", Name: "Gemma 2 9B", Description: "Google's open model", InputCost: 0.20, OutputCost: 0.20},
			},
		},
		{
			ID:          "anthropic",
			Name:        "Anthropic (Claude)",
			Description: "Claude AI - safe and helpful",
			RequiresKey: true,
			FreeAvailable: false,
			Models: []ModelInfo{
				{ID: "claude-3-5-haiku-latest", Name: "Claude 3.5 Haiku", Description: "Fast and affordable", InputCost: 0.80, OutputCost: 4.0},
				{ID: "claude-3-5-sonnet-latest", Name: "Claude 3.5 Sonnet", Description: "Balanced performance", InputCost: 3.0, OutputCost: 15.0},
			},
		},
	}
}

// GetProviderModels returns models for a specific provider
func GetProviderModels(providerID string) []ModelInfo {
	providers := GetAvailableProviders()
	for _, p := range providers {
		if p.ID == providerID {
			return p.Models
		}
	}
	return nil
}

// CreateProvider creates an AI provider based on config
func CreateProvider(config ProviderConfig) (AIProvider, error) {
	switch config.Provider {
	case "gemini":
		return NewGeminiProvider(config)
	case "openai":
		return NewOpenAIProvider(config)
	case "groq":
		return NewGroqProvider(config)
	case "anthropic":
		return NewAnthropicProvider(config)
	default:
		return NewGeminiProvider(config) // default to Gemini
	}
}


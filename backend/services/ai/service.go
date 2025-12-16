package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// AIService provides AI-powered auto-reply functionality
type AIService struct {
	defaultProvider AIProvider
	defaultAPIKey   string
}

// AutoReplyRequest contains the request data for auto-reply
type AutoReplyRequest struct {
	TenantID        string
	CustomerMessage string
	CustomerID      string
	SystemPrompt    string
	BusinessContext string
	KnowledgeBase   []Knowledge
	// Provider config from database
	Provider  string
	APIKey    string
	Model     string
	MaxTokens int
}

// AutoReplyResponse contains the AI response and metadata
type AutoReplyResponse struct {
	Response         string
	Confidence       float64
	DetectedIntent   string
	KnowledgeUsed    []string
	Attachments      []Knowledge // Knowledge entries that have media
	ShouldEscalate   bool
	EscalationReason string
	TokensUsed       int
	InputTokens      int
	OutputTokens     int
	CostUSD          float64
	Model            string
	Provider         string
	ResponseTimeMs   int64
}

// Knowledge represents a knowledge base entry
type Knowledge struct {
	ID        string
	Title     string
	Content   string
	Category  string
	Priority  int
	MediaURL  string
	MediaType string
}

// NewAIService creates a new AI service
func NewAIService() (*AIService, error) {
	// Try to create a default Gemini provider with system API key
	defaultAPIKey := os.Getenv("GEMINI_API_KEY")
	defaultModel := os.Getenv("GEMINI_MODEL")
	if defaultModel == "" {
		defaultModel = "gemini-2.0-flash"
	}

	var defaultProvider AIProvider
	if defaultAPIKey != "" {
		provider, err := NewGeminiProvider(ProviderConfig{
			APIKey: defaultAPIKey,
			Model:  defaultModel,
		})
		if err == nil {
			defaultProvider = provider
		}
	}

	return &AIService{
		defaultProvider: defaultProvider,
		defaultAPIKey:   defaultAPIKey,
	}, nil
}

// GenerateAutoReply generates an AI-powered auto-reply
func (s *AIService) GenerateAutoReply(ctx context.Context, req AutoReplyRequest) (*AutoReplyResponse, error) {
	startTime := time.Now()

	// Determine which provider to use
	var provider AIProvider
	var err error

	// If user provided API key, use it (even if use_system_key is true, if APIKey is provided, use it)
	if req.APIKey != "" && req.Provider != "" {
		// Use user's custom API key
		provider, err = CreateProvider(ProviderConfig{
			Provider:    req.Provider,
			APIKey:      req.APIKey,
			Model:       req.Model,
			MaxTokens:   req.MaxTokens,
			Temperature: 0.7,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create provider with user API key: %w", err)
		}
		defer provider.Close()
	} else if s.defaultProvider != nil {
		// Use system default provider (when use_system_key = true and no user API key)
		// If provider is specified and doesn't match, we can still use default but log warning
		// Or we could try to create provider with system key for requested provider
		if req.Provider != "" && req.Provider != s.defaultProvider.GetProviderName() {
			// Try to get system API key for requested provider
			envKey := fmt.Sprintf("%s_API_KEY", strings.ToUpper(req.Provider))
			systemAPIKey := os.Getenv(envKey)
			if systemAPIKey != "" {
				// Create provider with system API key for requested provider
				provider, err = CreateProvider(ProviderConfig{
					Provider:    req.Provider,
					APIKey:      systemAPIKey,
					Model:       req.Model,
					MaxTokens:   req.MaxTokens,
					Temperature: 0.7,
				})
				if err != nil {
					// Fallback to default provider
					fmt.Printf("[AI] Warning: Failed to create %s provider, using default: %v\n", req.Provider, err)
					provider = s.defaultProvider
				} else {
					defer provider.Close()
				}
			} else {
				// No system key for requested provider, use default
				fmt.Printf("[AI] Warning: Provider %s requested but no %s_API_KEY found, using default provider %s\n", 
					req.Provider, envKey, s.defaultProvider.GetProviderName())
				provider = s.defaultProvider
			}
		} else {
			provider = s.defaultProvider
		}
	} else {
		// No API key available
		if req.Provider != "" {
			return nil, fmt.Errorf("no API key available for provider %s. Please configure API key in AI settings or set %s_API_KEY environment variable", 
				req.Provider, strings.ToUpper(req.Provider))
		}
		return nil, fmt.Errorf("no AI provider available. Please configure an API key (GEMINI_API_KEY environment variable or user API key in settings)")
	}

	// Build context from knowledge base
	contextInfo := s.buildContext(req.BusinessContext, req.KnowledgeBase)

	// Generate response
	aiResp, err := provider.GenerateResponse(ctx, req.SystemPrompt, req.CustomerMessage, contextInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Detect intent
	intent := s.detectIntent(req.CustomerMessage)

	// Determine if should escalate
	shouldEscalate := false
	escalationReason := ""

	if aiResp.Confidence < 0.80 {
		shouldEscalate = true
		escalationReason = fmt.Sprintf("Low confidence: %.2f%%", aiResp.Confidence*100)
	}

	responseTimeMs := time.Since(startTime).Milliseconds()

	knowledgeUsed := make([]string, 0)
	attachments := make([]Knowledge, 0)
	for _, k := range req.KnowledgeBase {
		knowledgeUsed = append(knowledgeUsed, k.ID)
		// If knowledge has media, add to attachments
		if k.MediaURL != "" {
			attachments = append(attachments, k)
		}
	}

	return &AutoReplyResponse{
		Response:         aiResp.Response,
		Confidence:       aiResp.Confidence,
		DetectedIntent:   intent,
		KnowledgeUsed:    knowledgeUsed,
		Attachments:      attachments,
		ShouldEscalate:   shouldEscalate,
		EscalationReason: escalationReason,
		TokensUsed:       aiResp.TokensUsed,
		InputTokens:      aiResp.InputTokens,
		OutputTokens:     aiResp.OutputTokens,
		CostUSD:          aiResp.CostUSD,
		Model:            aiResp.Model,
		Provider:         provider.GetProviderName(),
		ResponseTimeMs:   responseTimeMs,
	}, nil
}

// TestConnection tests connection to an AI provider
func (s *AIService) TestConnection(ctx context.Context, providerName, apiKey, model string) error {
	provider, err := CreateProvider(ProviderConfig{
		Provider:    providerName,
		APIKey:      apiKey,
		Model:       model,
		MaxTokens:   50,
		Temperature: 0.5,
	})
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	defer provider.Close()

	_, err = provider.GenerateResponse(ctx, "You are a test assistant.", "Say 'OK' if you can hear me.", "")
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

// buildContext builds the context string from business info and knowledge base
func (s *AIService) buildContext(businessContext string, knowledge []Knowledge) string {
	var sb strings.Builder

	if businessContext != "" {
		sb.WriteString(businessContext)
		sb.WriteString("\n\n")
	}

	if len(knowledge) > 0 {
		sb.WriteString("RELEVANT INFORMATION:\n")
		for _, k := range knowledge {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", k.Title, k.Content))
		}
	}

	return sb.String()
}

// detectIntent detects the intent from customer message
func (s *AIService) detectIntent(message string) string {
	msg := strings.ToLower(message)

	if containsAny(msg, []string{"harga", "berapa", "price", "cost", "biaya"}) {
		return "price_inquiry"
	}
	if containsAny(msg, []string{"lokasi", "alamat", "dimana", "where", "tempat"}) {
		return "location_inquiry"
	}
	if containsAny(msg, []string{"jam", "buka", "tutup", "kapan", "hours", "operasional"}) {
		return "hours_inquiry"
	}
	if containsAny(msg, []string{"ada", "ready", "stock", "tersedia", "available"}) {
		return "availability_inquiry"
	}
	if containsAny(msg, []string{"pesan", "order", "beli", "buy", "mau"}) {
		return "order_intent"
	}
	if containsAny(msg, []string{"komplain", "kecewa", "marah", "complaint", "buruk", "jelek"}) {
		return "complaint"
	}
	if containsAny(msg, []string{"kirim", "ongkir", "shipping", "delivery", "pengiriman"}) {
		return "shipping_inquiry"
	}
	if containsAny(msg, []string{"bayar", "payment", "transfer", "cash", "pembayaran"}) {
		return "payment_inquiry"
	}

	return "general_inquiry"
}

func containsAny(str string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(str, keyword) {
			return true
		}
	}
	return false
}

// Close closes the AI service and releases resources
func (s *AIService) Close() {
	if s.defaultProvider != nil {
		s.defaultProvider.Close()
	}
}

// HasDefaultProvider returns true if a default provider is configured
func (s *AIService) HasDefaultProvider() bool {
	return s.defaultProvider != nil
}

// Helper function for prompt building
func buildPrompt(systemPrompt, contextInfo, userMessage string) string {
	var sb strings.Builder

	sb.WriteString("SYSTEM INSTRUCTIONS:\n")
	sb.WriteString(systemPrompt)
	sb.WriteString("\n\n")

	if contextInfo != "" {
		sb.WriteString("CONTEXT:\n")
		sb.WriteString(contextInfo)
		sb.WriteString("\n\n")
	}

	sb.WriteString("IMPORTANT RULES:\n")
	sb.WriteString("- Respond in Indonesian language\n")
	sb.WriteString("- Be helpful, friendly, and professional\n")
	sb.WriteString("- Keep responses concise (max 300 characters)\n")
	sb.WriteString("- If you don't know something, say so honestly\n")
	sb.WriteString("- Use the context provided to answer accurately\n")
	sb.WriteString("- Do not make up information\n\n")

	sb.WriteString("USER MESSAGE:\n")
	sb.WriteString(userMessage)

	return sb.String()
}

func estimateTokens(text string) int {
	return len(text) / 4
}

// Placeholder prompt templates
const (
	SystemPromptGeneral = `You are a helpful AI assistant for a small business.
Analyze customer messages and provide insights about their intent, sentiment, and potential as a lead.`

	SystemPromptOrderTaking = `You are an AI assistant helping to take orders for a small business.
Extract order details, quantities, and customer preferences from messages.`

	SystemPromptFAQ = `You are an AI assistant answering frequently asked questions for a small business.
Use the provided FAQ database to answer customer questions accurately.`
)

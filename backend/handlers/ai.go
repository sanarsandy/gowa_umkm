package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"net/http"
	"os"

	"gowa-backend/db"
	"gowa-backend/services/ai"

	"github.com/labstack/echo/v4"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiService *ai.AIService
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiService *ai.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

// AIConfigDB represents AI config in database
type AIConfigDB struct {
	TenantID            string  `db:"tenant_id" json:"tenant_id"`
	Enabled             bool    `db:"enabled" json:"enabled"`
	AIProvider          string  `db:"ai_provider" json:"ai_provider"`
	Model               string  `db:"model" json:"model"`
	APIKeySet           bool    `db:"api_key_set" json:"api_key_set"`
	UseSystemKey        bool    `db:"use_system_key" json:"use_system_key"`
	ConfidenceThreshold float64 `db:"confidence_threshold" json:"confidence_threshold"`
	MaxTokens           int     `db:"max_tokens" json:"max_tokens"`
	Language            string  `db:"language" json:"language"`
	SystemPrompt        string  `db:"system_prompt" json:"system_prompt"`
	BusinessName        string  `db:"business_name" json:"business_name"`
	BusinessType        string  `db:"business_type" json:"business_type"`
	BusinessHours       string  `db:"business_hours" json:"business_hours"`
	BusinessDescription string  `db:"business_description" json:"business_description"`
	BusinessAddress     string  `db:"business_address" json:"business_address"`
	PaymentMethods      string  `db:"payment_methods" json:"payment_methods"`
	EscalateLowConf     bool    `db:"escalate_low_confidence" json:"escalate_low_confidence"`
	EscalateComplaint   bool    `db:"escalate_complaint" json:"escalate_complaint"`
	EscalateOrder       bool    `db:"escalate_order" json:"escalate_order"`
	EscalateUrgent      bool    `db:"escalate_urgent" json:"escalate_urgent"`
	NotifyWhatsApp      bool    `db:"notify_whatsapp" json:"notify_whatsapp"`
	NotifyEmail         bool    `db:"notify_email" json:"notify_email"`
	TotalRequests       int     `db:"total_requests" json:"total_requests"`
	TotalTokensUsed     int64   `db:"total_tokens_used" json:"total_tokens_used"`
	TotalCostUSD        float64 `db:"total_cost_usd" json:"total_cost_usd"`
}

// GetProviders returns list of available AI providers
// GET /api/ai/providers
func (h *AIHandler) GetProviders(c echo.Context) error {
	providers := ai.GetAvailableProviders()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"providers": providers,
	})
}

// GetProviderModels returns models for a specific provider
// GET /api/ai/providers/:provider/models
func (h *AIHandler) GetProviderModels(c echo.Context) error {
	providerID := c.Param("provider")
	models := ai.GetProviderModels(providerID)
	if models == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Provider not found")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"models": models,
	})
}

// TestConnection tests connection to an AI provider
// POST /api/ai/test-connection
func (h *AIHandler) TestConnection(c echo.Context) error {
	var req struct {
		Provider string `json:"provider"`
		APIKey   string `json:"api_key"`
		Model    string `json:"model"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Provider == "" || req.APIKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Provider and API key are required")
	}

	if req.Model == "" {
		// Use default model for provider
		models := ai.GetProviderModels(req.Provider)
		if len(models) > 0 {
			req.Model = models[0].ID
		}
	}

	err := h.aiService.TestConnection(c.Request().Context(), req.Provider, req.APIKey, req.Model)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Connection successful!",
	})
}

// TestAIResponse tests the AI auto-reply functionality
// POST /api/ai/test
func (h *AIHandler) TestAIResponse(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var req struct {
		Message string `json:"message" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Message is required")
	}

	// Load AI config from database
	config, userAPIKey, err := getAIConfigWithKey(tenantID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load AI config")
	}

	// Determine API key to use
	apiKey := ""
	if !config.UseSystemKey && userAPIKey != "" {
		apiKey = userAPIKey
	}

	// Use defaults if not configured
	systemPrompt := config.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = `Kamu adalah asisten toko online yang ramah dan membantu.
Jawab pertanyaan customer dengan sopan, jelas, dan ringkas.
Jika tidak tahu jawabannya, katakan dengan jujur dan sarankan untuk menghubungi admin.`
	}

	businessContext := buildBusinessContextFromConfig(config)
	knowledgeBase, _ := getKnowledgeBaseForTenant(tenantID)

	// Generate AI response
	resp, err := h.aiService.GenerateAutoReply(c.Request().Context(), ai.AutoReplyRequest{
		TenantID:        tenantID,
		CustomerMessage: req.Message,
		SystemPrompt:    systemPrompt,
		BusinessContext: businessContext,
		KnowledgeBase:   knowledgeBase,
		Provider:        config.AIProvider,
		APIKey:          apiKey,
		Model:           config.Model,
		MaxTokens:       config.MaxTokens,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate AI response: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response":          resp.Response,
		"confidence":        resp.Confidence,
		"detected_intent":   resp.DetectedIntent,
		"should_escalate":   resp.ShouldEscalate,
		"escalation_reason": resp.EscalationReason,
		"tokens_used":       resp.TokensUsed,
		"input_tokens":      resp.InputTokens,
		"output_tokens":     resp.OutputTokens,
		"cost_usd":          resp.CostUSD,
		"model":             resp.Model,
		"provider":          resp.Provider,
		"response_time_ms":  resp.ResponseTimeMs,
	})
}

// GetAIConfig gets the AI configuration for the tenant
// GET /api/ai/config
func (h *AIHandler) GetAIConfig(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	config, _, err := getAIConfigWithKey(tenantID)
	if err == sql.ErrNoRows {
		// Return default config
		return c.JSON(http.StatusOK, AIConfigDB{
			TenantID:            tenantID,
			Enabled:             false,
			AIProvider:          "gemini",
			Model:               "gemini-2.0-flash",
			APIKeySet:           false,
			UseSystemKey:        true,
			ConfidenceThreshold: 0.80,
			MaxTokens:           200,
			Language:            "id",
			EscalateLowConf:     true,
			EscalateComplaint:   true,
			EscalateOrder:       false,
			EscalateUrgent:      true,
			NotifyWhatsApp:      true,
			NotifyEmail:         false,
		})
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get AI config")
	}

	return c.JSON(http.StatusOK, config)
}

// UpdateAIConfig updates the AI configuration for the tenant
// PUT /api/ai/config
func (h *AIHandler) UpdateAIConfig(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var req struct {
		Enabled             bool    `json:"enabled"`
		AIProvider          string  `json:"ai_provider"`
		Model               string  `json:"model"`
		APIKey              string  `json:"api_key"` // Only sent when changing
		UseSystemKey        bool    `json:"use_system_key"`
		ConfidenceThreshold float64 `json:"confidence_threshold"`
		MaxTokens           int     `json:"max_tokens"`
		Language            string  `json:"language"`
		SystemPrompt        string  `json:"system_prompt"`
		BusinessName        string  `json:"business_name"`
		BusinessType        string  `json:"business_type"`
		BusinessHours       string  `json:"business_hours"`
		BusinessDescription string  `json:"business_description"`
		BusinessAddress     string  `json:"business_address"`
		PaymentMethods      string  `json:"payment_methods"`
		EscalateLowConf     bool    `json:"escalate_low_confidence"`
		EscalateComplaint   bool    `json:"escalate_complaint"`
		EscalateOrder       bool    `json:"escalate_order"`
		EscalateUrgent      bool    `json:"escalate_urgent"`
		NotifyWhatsApp      bool    `json:"notify_whatsapp"`
		NotifyEmail         bool    `json:"notify_email"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Set defaults
	if req.AIProvider == "" {
		req.AIProvider = "gemini"
	}
	if req.ConfidenceThreshold < 0.5 || req.ConfidenceThreshold > 0.95 {
		req.ConfidenceThreshold = 0.80
	}
	if req.MaxTokens < 50 || req.MaxTokens > 2000 {
		req.MaxTokens = 200
	}

	// Encrypt API key if provided
	var encryptedKey *string
	apiKeySet := false
	if req.APIKey != "" {
		encrypted, err := encryptAPIKey(req.APIKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encrypt API key")
		}
		encryptedKey = &encrypted
		apiKeySet = true
	}

	// Upsert config
	query := `
		INSERT INTO ai_configs (
			tenant_id, enabled, ai_provider, model, confidence_threshold, max_tokens, language,
			system_prompt, business_name, business_type, business_hours,
			business_description, business_address, payment_methods,
			escalate_low_confidence, escalate_complaint, escalate_order,
			escalate_urgent, notify_whatsapp, notify_email, use_system_key
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		ON CONFLICT (tenant_id) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			ai_provider = EXCLUDED.ai_provider,
			model = EXCLUDED.model,
			confidence_threshold = EXCLUDED.confidence_threshold,
			max_tokens = EXCLUDED.max_tokens,
			language = EXCLUDED.language,
			system_prompt = EXCLUDED.system_prompt,
			business_name = EXCLUDED.business_name,
			business_type = EXCLUDED.business_type,
			business_hours = EXCLUDED.business_hours,
			business_description = EXCLUDED.business_description,
			business_address = EXCLUDED.business_address,
			payment_methods = EXCLUDED.payment_methods,
			escalate_low_confidence = EXCLUDED.escalate_low_confidence,
			escalate_complaint = EXCLUDED.escalate_complaint,
			escalate_order = EXCLUDED.escalate_order,
			escalate_urgent = EXCLUDED.escalate_urgent,
			notify_whatsapp = EXCLUDED.notify_whatsapp,
			notify_email = EXCLUDED.notify_email,
			use_system_key = EXCLUDED.use_system_key,
			updated_at = NOW()
	`

	_, err := db.DB.Exec(query,
		tenantID, req.Enabled, req.AIProvider, req.Model, req.ConfidenceThreshold, req.MaxTokens,
		req.Language, req.SystemPrompt, req.BusinessName, req.BusinessType,
		req.BusinessHours, req.BusinessDescription, req.BusinessAddress,
		req.PaymentMethods, req.EscalateLowConf, req.EscalateComplaint,
		req.EscalateOrder, req.EscalateUrgent, req.NotifyWhatsApp, req.NotifyEmail, req.UseSystemKey,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save AI config: "+err.Error())
	}

	// Update API key separately if provided
	if encryptedKey != nil {
		_, err = db.DB.Exec(`
			UPDATE ai_configs 
			SET user_api_key = $1, api_key_set = $2 
			WHERE tenant_id = $3
		`, *encryptedKey, apiKeySet, tenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save API key")
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "AI configuration updated successfully",
		"tenant_id": tenantID,
	})
}

// GetAIStats returns AI usage statistics
// GET /api/ai/stats
func (h *AIHandler) GetAIStats(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var stats struct {
		TotalRequests   int     `db:"total_requests" json:"total_requests"`
		TotalTokensUsed int64   `db:"total_tokens_used" json:"total_tokens_used"`
		TotalCostUSD    float64 `db:"total_cost_usd" json:"total_cost_usd"`
	}

	query := `
		SELECT 
			COALESCE(total_requests, 0) as total_requests,
			COALESCE(total_tokens_used, 0) as total_tokens_used,
			COALESCE(total_cost_usd, 0) as total_cost_usd
		FROM ai_configs
		WHERE tenant_id = $1
	`

	err := db.DB.Get(&stats, query, tenantID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"total_requests":    0,
			"total_tokens_used": 0,
			"total_cost_usd":    0,
		})
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get AI stats")
	}

	var logsStats struct {
		TotalConversations int     `db:"total_conversations"`
		AutoReplied        int     `db:"auto_replied"`
		Escalated          int     `db:"escalated"`
		AvgConfidence      float64 `db:"avg_confidence"`
	}

	logsQuery := `
		SELECT 
			COUNT(*) as total_conversations,
			COUNT(*) FILTER (WHERE action_taken = 'auto_replied') as auto_replied,
			COUNT(*) FILTER (WHERE action_taken = 'escalated') as escalated,
			COALESCE(AVG(confidence_score), 0) as avg_confidence
		FROM ai_conversation_logs
		WHERE tenant_id = $1
	`

	db.DB.Get(&logsStats, logsQuery, tenantID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_requests":      stats.TotalRequests,
		"total_tokens_used":   stats.TotalTokensUsed,
		"total_cost_usd":      stats.TotalCostUSD,
		"total_conversations": logsStats.TotalConversations,
		"auto_replied":        logsStats.AutoReplied,
		"escalated":           logsStats.Escalated,
		"avg_confidence":      logsStats.AvgConfidence,
	})
}

// Helper functions

func getAIConfigWithKey(tenantID string) (*AIConfigDB, string, error) {
	var config AIConfigDB
	var userAPIKey sql.NullString

	query := `
		SELECT 
			tenant_id,
			COALESCE(enabled, false) as enabled,
			COALESCE(ai_provider, 'gemini') as ai_provider,
			COALESCE(model, 'gemini-2.0-flash') as model,
			COALESCE(api_key_set, false) as api_key_set,
			COALESCE(use_system_key, true) as use_system_key,
			COALESCE(confidence_threshold, 0.80) as confidence_threshold,
			COALESCE(max_tokens, 200) as max_tokens,
			COALESCE(language, 'id') as language,
			COALESCE(system_prompt, '') as system_prompt,
			COALESCE(business_name, '') as business_name,
			COALESCE(business_type, '') as business_type,
			COALESCE(business_hours, '') as business_hours,
			COALESCE(business_description, '') as business_description,
			COALESCE(business_address, '') as business_address,
			COALESCE(payment_methods, '') as payment_methods,
			COALESCE(escalate_low_confidence, true) as escalate_low_confidence,
			COALESCE(escalate_complaint, true) as escalate_complaint,
			COALESCE(escalate_order, false) as escalate_order,
			COALESCE(escalate_urgent, true) as escalate_urgent,
			COALESCE(notify_whatsapp, true) as notify_whatsapp,
			COALESCE(notify_email, false) as notify_email,
			COALESCE(total_requests, 0) as total_requests,
			COALESCE(total_tokens_used, 0) as total_tokens_used,
			COALESCE(total_cost_usd, 0) as total_cost_usd,
			user_api_key
		FROM ai_configs
		WHERE tenant_id = $1
	`

	row := db.DB.QueryRow(query, tenantID)
	err := row.Scan(
		&config.TenantID, &config.Enabled, &config.AIProvider, &config.Model,
		&config.APIKeySet, &config.UseSystemKey, &config.ConfidenceThreshold, &config.MaxTokens,
		&config.Language, &config.SystemPrompt, &config.BusinessName, &config.BusinessType,
		&config.BusinessHours, &config.BusinessDescription, &config.BusinessAddress,
		&config.PaymentMethods, &config.EscalateLowConf, &config.EscalateComplaint,
		&config.EscalateOrder, &config.EscalateUrgent, &config.NotifyWhatsApp, &config.NotifyEmail,
		&config.TotalRequests, &config.TotalTokensUsed, &config.TotalCostUSD, &userAPIKey,
	)

	if err != nil {
		return &AIConfigDB{
			TenantID:   tenantID,
			AIProvider: "gemini",
			Model:      "gemini-2.0-flash",
		}, "", err
	}

	// Decrypt API key if exists
	decryptedKey := ""
	if userAPIKey.Valid && userAPIKey.String != "" {
		decrypted, err := decryptAPIKey(userAPIKey.String)
		if err == nil {
			decryptedKey = decrypted
		}
	}

	return &config, decryptedKey, nil
}

func getAIConfigFromDB(tenantID string) (*AIConfigDB, error) {
	config, _, err := getAIConfigWithKey(tenantID)
	return config, err
}

func buildBusinessContextFromConfig(config *AIConfigDB) string {
	var parts []string

	if config.BusinessName != "" {
		parts = append(parts, "Nama Bisnis: "+config.BusinessName)
	}
	if config.BusinessType != "" {
		parts = append(parts, "Jenis Bisnis: "+config.BusinessType)
	}
	if config.BusinessHours != "" {
		parts = append(parts, "Jam Operasional: "+config.BusinessHours)
	}
	if config.BusinessAddress != "" {
		parts = append(parts, "Alamat: "+config.BusinessAddress)
	}
	if config.PaymentMethods != "" {
		parts = append(parts, "Metode Pembayaran: "+config.PaymentMethods)
	}
	if config.BusinessDescription != "" {
		parts = append(parts, "Deskripsi: "+config.BusinessDescription)
	}

	if len(parts) == 0 {
		return ""
	}

	result := ""
	for _, p := range parts {
		result += p + "\n"
	}
	return result
}

func getKnowledgeBaseForTenant(tenantID string) ([]ai.Knowledge, error) {
	var knowledge []ai.Knowledge

	query := `
		SELECT id, title, content, COALESCE(category, '') as category, priority
		FROM knowledge_base
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY priority DESC
		LIMIT 10
	`

	rows, err := db.DB.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k ai.Knowledge
		if err := rows.Scan(&k.ID, &k.Title, &k.Content, &k.Category, &k.Priority); err != nil {
			continue
		}
		knowledge = append(knowledge, k)
	}

	return knowledge, nil
}

// Encryption functions for API keys
func getEncryptionKey() []byte {
	key := os.Getenv("API_KEY_ENCRYPTION_KEY")
	if key == "" {
		key = os.Getenv("JWT_SECRET")
	}
	if key == "" {
		key = "default-encryption-key-32bytes!!" // fallback (32 bytes)
	}
	// Ensure key is exactly 32 bytes for AES-256
	for len(key) < 32 {
		key += key
	}
	return []byte(key[:32])
}

func encryptAPIKey(plaintext string) (string, error) {
	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAPIKey(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", err
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

package workers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"gowa-backend/services/ai"
	"gowa-backend/services/redis"

	"github.com/jmoiron/sqlx"
)

// MessageWorker processes messages from Redis queue
type MessageWorker struct {
	redisClient     *redis.Client
	db              *sqlx.DB
	whatsappService WhatsAppService
	aiService       *ai.AIService
	stopChan        chan struct{}
}

// WhatsAppService interface for sending messages
type WhatsAppService interface {
	SendMessage(ctx context.Context, tenantID string, recipientJID string, message string) (string, error)
	SendMediaMessage(ctx context.Context, tenantID string, recipientJID string, mediaData []byte, mediaType string, fileName string, caption string) (string, error)
}

// NewMessageWorker creates a new message worker
func NewMessageWorker(redisClient *redis.Client, db *sqlx.DB, whatsappService WhatsAppService) *MessageWorker {
	// Try to initialize AI service
	aiService, err := ai.NewAIService()
	if err != nil {
		fmt.Printf("[Worker] Warning: AI service not available: %v\n", err)
	}

	return &MessageWorker{
		redisClient:     redisClient,
		db:              db,
		whatsappService: whatsappService,
		aiService:       aiService,
		stopChan:        make(chan struct{}),
	}
}

// Start begins processing messages from the queue
func (w *MessageWorker) Start(ctx context.Context) {
	fmt.Println("Message worker started")
	
	// Start goroutines for different queues
	go w.processAIMessages(ctx)
	go w.processBroadcastMessages(ctx)
	
	<-ctx.Done()
	fmt.Println("Message worker stopped")
}

// Stop stops the worker
func (w *MessageWorker) Stop() {
	close(w.stopChan)
}

// processAIMessages processes AI queue messages
func (w *MessageWorker) processAIMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		default:
			w.processAIMessage(ctx)
		}
	}
}

// processBroadcastMessages processes broadcast queue messages
func (w *MessageWorker) processBroadcastMessages(ctx context.Context) {
	fmt.Println("[Worker] Broadcast message processor started")
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		default:
			w.processBroadcastMessage(ctx)
		}
	}
}

// AIConfig holds the AI configuration for a tenant
type AIConfig struct {
	Enabled             bool    `db:"enabled"`
	AIProvider          string  `db:"ai_provider"`
	Model               string  `db:"model"`
	UseSystemKey        bool    `db:"use_system_key"`
	UserAPIKey          string  // Decrypted API key (not from DB directly)
	ConfidenceThreshold float64 `db:"confidence_threshold"`
	MaxTokens           int     `db:"max_tokens"`
	SystemPrompt        string  `db:"system_prompt"`
	BusinessName        string  `db:"business_name"`
	BusinessType        string  `db:"business_type"`
	BusinessHours       string  `db:"business_hours"`
	BusinessDescription string  `db:"business_description"`
	BusinessAddress     string  `db:"business_address"`
	PaymentMethods      string  `db:"payment_methods"`
	EscalateComplaint   bool    `db:"escalate_complaint"`
	EscalateOrder       bool    `db:"escalate_order"`
}

// getAIConfig loads AI configuration for a tenant
func (w *MessageWorker) getAIConfig(ctx context.Context, tenantID string) (*AIConfig, error) {
	var config AIConfig
	var encryptedAPIKey sql.NullString
	
	query := `
		SELECT 
			COALESCE(enabled, false) as enabled,
			COALESCE(ai_provider, 'gemini') as ai_provider,
			COALESCE(model, 'gemini-1.5-flash') as model,
			COALESCE(use_system_key, true) as use_system_key,
			user_api_key,
			COALESCE(confidence_threshold, 0.80) as confidence_threshold,
			COALESCE(max_tokens, 200) as max_tokens,
			COALESCE(system_prompt, '') as system_prompt,
			COALESCE(business_name, '') as business_name,
			COALESCE(business_type, '') as business_type,
			COALESCE(business_hours, '') as business_hours,
			COALESCE(business_description, '') as business_description,
			COALESCE(business_address, '') as business_address,
			COALESCE(payment_methods, '') as payment_methods,
			COALESCE(escalate_complaint, true) as escalate_complaint,
			COALESCE(escalate_order, false) as escalate_order
		FROM ai_configs
		WHERE tenant_id = $1
	`
	
	row := w.db.QueryRowContext(ctx, query, tenantID)
	err := row.Scan(
		&config.Enabled,
		&config.AIProvider,
		&config.Model,
		&config.UseSystemKey,
		&encryptedAPIKey,
		&config.ConfidenceThreshold,
		&config.MaxTokens,
		&config.SystemPrompt,
		&config.BusinessName,
		&config.BusinessType,
		&config.BusinessHours,
		&config.BusinessDescription,
		&config.BusinessAddress,
		&config.PaymentMethods,
		&config.EscalateComplaint,
		&config.EscalateOrder,
	)
	
	if err == sql.ErrNoRows {
		// Return default config
		return &AIConfig{
			Enabled:             false,
			AIProvider:          "gemini",
			Model:               "gemini-1.5-flash",
			UseSystemKey:        true,
			ConfidenceThreshold: 0.80,
			MaxTokens:           200,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Decrypt user API key if exists and not using system key
	if !config.UseSystemKey && encryptedAPIKey.Valid && encryptedAPIKey.String != "" {
		decrypted, err := decryptAPIKeyWorker(encryptedAPIKey.String)
		if err != nil {
			fmt.Printf("[Worker] Warning: Failed to decrypt user API key: %v\n", err)
		} else {
			config.UserAPIKey = decrypted
			fmt.Printf("[Worker] Using user's custom API key for tenant %s\n", tenantID)
		}
	}
	
	return &config, nil
}

// getKnowledgeBase loads active knowledge base for a tenant
func (w *MessageWorker) getKnowledgeBase(ctx context.Context, tenantID string) ([]ai.Knowledge, error) {
	var knowledge []ai.Knowledge
	
	query := `
		SELECT id, title, content, COALESCE(category, '') as category, priority,
		       COALESCE(media_url, '') as media_url, COALESCE(media_type, '') as media_type
		FROM knowledge_base
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY priority DESC
		LIMIT 10
	`
	
	rows, err := w.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var k ai.Knowledge
		if err := rows.Scan(&k.ID, &k.Title, &k.Content, &k.Category, &k.Priority, &k.MediaURL, &k.MediaType); err != nil {
			continue
		}
		knowledge = append(knowledge, k)
	}
	
	return knowledge, nil
}

// buildBusinessContext creates business context string from config
func (w *MessageWorker) buildBusinessContext(config *AIConfig) string {
	var parts []string
	
	if config.BusinessName != "" {
		parts = append(parts, fmt.Sprintf("Nama Bisnis: %s", config.BusinessName))
	}
	if config.BusinessType != "" {
		parts = append(parts, fmt.Sprintf("Jenis Bisnis: %s", config.BusinessType))
	}
	if config.BusinessHours != "" {
		parts = append(parts, fmt.Sprintf("Jam Operasional: %s", config.BusinessHours))
	}
	if config.BusinessAddress != "" {
		parts = append(parts, fmt.Sprintf("Alamat: %s", config.BusinessAddress))
	}
	if config.PaymentMethods != "" {
		parts = append(parts, fmt.Sprintf("Metode Pembayaran: %s", config.PaymentMethods))
	}
	if config.BusinessDescription != "" {
		parts = append(parts, fmt.Sprintf("Deskripsi: %s", config.BusinessDescription))
	}
	
	return strings.Join(parts, "\n")
}

// getChatHistory retrieves recent chat messages for context
func (w *MessageWorker) getChatHistory(ctx context.Context, tenantID, chatJID string) string {
	query := `
		SELECT message_text, is_from_me, timestamp
		FROM whatsapp_messages
		WHERE tenant_id = $1 AND chat_jid = $2 AND message_text IS NOT NULL AND message_text != ''
		ORDER BY timestamp DESC
		LIMIT 10
	`

	rows, err := w.db.QueryContext(ctx, query, tenantID, chatJID)
	if err != nil {
		fmt.Printf("[Worker] Failed to get chat history: %v\n", err)
		return ""
	}
	defer rows.Close()

	type historyMsg struct {
		Text     string
		IsFromMe bool
		Time     int64
	}

	var messages []historyMsg
	for rows.Next() {
		var msg historyMsg
		if err := rows.Scan(&msg.Text, &msg.IsFromMe, &msg.Time); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return ""
	}

	// Reverse to get oldest first (chronological order)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Build chat history string
	var sb strings.Builder
	sb.WriteString("RIWAYAT PERCAKAPAN SEBELUMNYA:\n")
	for _, msg := range messages {
		role := "Customer"
		if msg.IsFromMe {
			role = "Anda"
		}
		// Truncate long messages
		text := msg.Text
		if len(text) > 200 {
			text = text[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", role, text))
	}

	fmt.Printf("[Worker] Including %d messages in chat history context\n", len(messages))
	return sb.String()
}

// logAIConversation logs AI conversation to database
func (w *MessageWorker) logAIConversation(ctx context.Context, tenantID string, senderJID string, customerMessage string, response *ai.AutoReplyResponse, action string) {
	query := `
		INSERT INTO ai_conversation_logs (
			tenant_id, customer_message, ai_response, detected_intent,
			confidence_score, action_taken, escalation_reason,
			response_time_ms, tokens_used, input_tokens, output_tokens,
			cost_usd, model_used
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	
	_, err := w.db.ExecContext(ctx, query,
		tenantID,
		customerMessage,
		response.Response,
		response.DetectedIntent,
		response.Confidence,
		action,
		response.EscalationReason,
		response.ResponseTimeMs,
		response.TokensUsed,
		response.InputTokens,
		response.OutputTokens,
		response.CostUSD,
		response.Model,
	)
	
	if err != nil {
		fmt.Printf("[Worker] Failed to log AI conversation: %v\n", err)
	}
}

// updateKnowledgeUsage updates usage count for knowledge entries
func (w *MessageWorker) updateKnowledgeUsage(ctx context.Context, knowledgeIDs []string) {
	if len(knowledgeIDs) == 0 {
		return
	}
	
	for _, id := range knowledgeIDs {
		query := `UPDATE knowledge_base SET usage_count = usage_count + 1, last_used_at = NOW() WHERE id = $1`
		w.db.ExecContext(ctx, query, id)
	}
}

// updateAIUsageStats updates AI usage statistics
func (w *MessageWorker) updateAIUsageStats(ctx context.Context, tenantID string, tokensUsed int, costUSD float64) {
	query := `
		UPDATE ai_configs 
		SET total_requests = total_requests + 1,
		    total_tokens_used = total_tokens_used + $1,
		    total_cost_usd = total_cost_usd + $2,
		    updated_at = NOW()
		WHERE tenant_id = $3
	`
	w.db.ExecContext(ctx, query, tokensUsed, costUSD, tenantID)
}

// processAIMessage processes a single message from the AI queue
func (w *MessageWorker) processAIMessage(ctx context.Context) {
	// Pop message from AI queue with 5 second timeout
	payload, err := w.redisClient.PopFromAIQueue(ctx, 5*time.Second)
	if err != nil {
		// Timeout or error - just continue
		return
	}

	if payload == nil {
		// No message available
		return
	}

	fmt.Printf("[Worker] Processing AI message from tenant %s: %s\n", payload.TenantID, payload.MessageText)

	// Load AI config for tenant
	config, err := w.getAIConfig(ctx, payload.TenantID)
	if err != nil {
		fmt.Printf("[Worker] Failed to load AI config: %v\n", err)
		return
	}

	// Check if AI auto-reply is enabled
	if !config.Enabled {
		fmt.Printf("[Worker] AI auto-reply disabled for tenant %s\n", payload.TenantID)
		w.markMessageProcessed(ctx, payload, false)
		w.updateCustomerInsight(ctx, payload)
		return
	}

	// Check if AI service is available
	if w.aiService == nil {
		fmt.Printf("[Worker] AI service not available\n")
		w.markMessageProcessed(ctx, payload, false)
		w.updateCustomerInsight(ctx, payload)
		return
	}

	// Load knowledge base
	knowledgeBase, err := w.getKnowledgeBase(ctx, payload.TenantID)
	if err != nil {
		fmt.Printf("[Worker] Failed to load knowledge base: %v\n", err)
		knowledgeBase = []ai.Knowledge{}
	}

	// Build system prompt
	systemPrompt := config.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = `Kamu adalah asisten toko online yang ramah dan membantu.
Jawab pertanyaan customer dengan sopan, jelas, dan ringkas.
Jika tidak tahu jawabannya, katakan dengan jujur dan sarankan untuk menghubungi admin.
Selalu gunakan bahasa yang santun dan profesional.`
	}

	// Build business context
	businessContext := w.buildBusinessContext(config)

	// Get chat history for context
	chatHistory := w.getChatHistory(ctx, payload.TenantID, normalizeJID(payload.SenderJID))
	if chatHistory != "" {
		if businessContext != "" {
			businessContext += "\n\n"
		}
		businessContext += chatHistory
	}

	// Generate AI response
	aiReq := ai.AutoReplyRequest{
		TenantID:        payload.TenantID,
		CustomerMessage: payload.MessageText,
		CustomerID:      payload.SenderJID,
		SystemPrompt:    systemPrompt,
		BusinessContext: businessContext,
		KnowledgeBase:   knowledgeBase,
		Provider:        config.AIProvider,
		APIKey:          config.UserAPIKey, // Will be empty if using system key
		Model:           config.Model,
		MaxTokens:       config.MaxTokens,
	}

	response, err := w.aiService.GenerateAutoReply(ctx, aiReq)
	if err != nil {
		fmt.Printf("[Worker] Failed to generate AI response: %v\n", err)
		w.markMessageProcessed(ctx, payload, false)
		w.updateCustomerInsight(ctx, payload)
		return
	}

	fmt.Printf("[Worker] AI Response: confidence=%.2f, intent=%s, escalate=%v\n",
		response.Confidence, response.DetectedIntent, response.ShouldEscalate)

	// Determine action based on confidence and escalation rules
	shouldEscalate := response.ShouldEscalate
	action := "auto_replied"

	// Check escalation rules from config
	if config.EscalateComplaint && response.DetectedIntent == "complaint" {
		shouldEscalate = true
		response.EscalationReason = "Complaint detected"
	}
	if config.EscalateOrder && response.DetectedIntent == "order_intent" {
		shouldEscalate = true
		response.EscalationReason = "Order intent detected"
	}

	// Check confidence threshold
	if response.Confidence < config.ConfidenceThreshold {
		shouldEscalate = true
		response.EscalationReason = fmt.Sprintf("Low confidence: %.0f%%", response.Confidence*100)
	}

	if shouldEscalate {
		action = "escalated"
		fmt.Printf("[Worker] Message escalated: %s\n", response.EscalationReason)
		// TODO: Send notification to admin (WhatsApp/email)
	} else {
		// Send auto-reply via WhatsApp
		if w.whatsappService != nil {
			// Send text response first
			messageID, err := w.whatsappService.SendMessage(ctx, payload.TenantID, payload.SenderJID, response.Response)
			if err != nil {
				fmt.Printf("[Worker] Failed to send auto-reply: %v\n", err)
				action = "failed"
			} else {
				fmt.Printf("[Worker] Auto-reply sent! MessageID: %s\n", messageID)
			}

			// Send attachments if any
			if len(response.Attachments) > 0 {
				fmt.Printf("[Worker] Sending %d attachments\n", len(response.Attachments))
				for _, attachment := range response.Attachments {
					if attachment.MediaURL == "" {
						continue
					}

					// Fetch media data (assuming URL is accessible)
					// In a real implementation, we might need to handle local files vs remote URLs
					// For now, we'll assume it's a remote URL and we need to download it first?
					// Or if it's a local path, read it.
					// Let's implement a simple helper to get media data
					mediaData, err := w.fetchMediaData(attachment.MediaURL)
					if err != nil {
						fmt.Printf("[Worker] Failed to fetch media data: %v\n", err)
						continue
					}

					// Send media
					mediaType := attachment.MediaType
					if mediaType == "" {
						mediaType = "image" // Default to image
					}

					_, err = w.whatsappService.SendMediaMessage(ctx, payload.TenantID, payload.SenderJID, mediaData, mediaType, "attachment", attachment.Title)
					if err != nil {
						fmt.Printf("[Worker] Failed to send attachment: %v\n", err)
					} else {
						fmt.Printf("[Worker] Attachment sent: %s\n", attachment.Title)
					}
				}
			}
		}
	}

	// Log conversation
	w.logAIConversation(ctx, payload.TenantID, payload.SenderJID, payload.MessageText, response, action)

	// Update knowledge usage
	w.updateKnowledgeUsage(ctx, response.KnowledgeUsed)

	// Update AI usage stats
	w.updateAIUsageStats(ctx, payload.TenantID, response.TokensUsed, response.CostUSD)

	// Mark message as AI processed
	w.markMessageProcessed(ctx, payload, true)

	// Update customer insight
	w.updateCustomerInsight(ctx, payload)
}

// markMessageProcessed marks a message as processed
func (w *MessageWorker) markMessageProcessed(ctx context.Context, payload *redis.MessagePayload, aiProcessed bool) {
	query := `
		UPDATE whatsapp_messages 
		SET ai_processed = $1, ai_processed_at = NOW()
		WHERE tenant_id = $2 AND message_id = $3
	`
	
	_, err := w.db.ExecContext(ctx, query, aiProcessed, payload.TenantID, payload.MessageID)
	if err != nil {
		fmt.Printf("[Worker] Failed to update message status: %v\n", err)
	}
}

// processBroadcastMessage processes a single broadcast message from the queue
func (w *MessageWorker) processBroadcastMessage(ctx context.Context) {
	// Pop message from broadcast queue with 5 second timeout
	payload, err := w.redisClient.PopFromBroadcastQueue(ctx, 5*time.Second)
	if err != nil {
		// Timeout or error - just continue
		return
	}

	if payload == nil {
		// No message available
		return
	}

	fmt.Printf("[Worker] Processing broadcast message to %s: %s\n", payload.CustomerJID, payload.Message)

	// Send message via WhatsApp service
	messageID, err := w.whatsappService.SendMessage(ctx, payload.TenantID, payload.CustomerJID, payload.Message)
	
	if err != nil {
		fmt.Printf("[Worker] Error sending WhatsApp message: %v\n", err)
		// Mark as failed
		failQuery := `UPDATE broadcast_recipients SET status = 'failed', error_message = $1 WHERE id = $2`
		w.db.ExecContext(ctx, failQuery, err.Error(), payload.RecipientID)
		
		// Update broadcast failed count
		w.db.ExecContext(ctx, `UPDATE broadcasts SET failed_count = failed_count + 1, updated_at = NOW() WHERE id = $1`, payload.BroadcastID)
		return
	}

	fmt.Printf("[Worker] WhatsApp message sent successfully! MessageID: %s\n", messageID)
	
	// Update recipient status to sent with message ID
	updateQuery := `UPDATE broadcast_recipients SET status = 'sent', message_id = $1, sent_at = NOW() WHERE id = $2`
	if _, err := w.db.ExecContext(ctx, updateQuery, messageID, payload.RecipientID); err != nil {
		fmt.Printf("[Worker] Error updating recipient status: %v\n", err)
		return
	}

	// Update broadcast sent count
	countQuery := `UPDATE broadcasts SET sent_count = sent_count + 1, updated_at = NOW() WHERE id = $1`
	w.db.ExecContext(ctx, countQuery, payload.BroadcastID)

	// Check if all recipients are processed
	var pendingCount int
	w.db.Get(&pendingCount, `SELECT COUNT(*) FROM broadcast_recipients WHERE broadcast_id = $1 AND status = 'queued'`, payload.BroadcastID)
	
	if pendingCount == 0 {
		// All messages sent, mark broadcast as completed
		completeQuery := `UPDATE broadcasts SET status = 'completed', completed_at = NOW(), updated_at = NOW() WHERE id = $1`
		w.db.ExecContext(ctx, completeQuery, payload.BroadcastID)
		fmt.Printf("[Worker] Broadcast %s completed\n", payload.BroadcastID)
	}

	fmt.Printf("[Worker] Broadcast message delivered to %s (MessageID: %s)\n", payload.CustomerJID, messageID)
}

// normalizeJID removes device part from JID for consistent customer identification
// e.g., "6281234567890:46@s.whatsapp.net" -> "6281234567890@s.whatsapp.net"
func normalizeJID(jid string) string {
	// Check if JID contains device part (e.g., :46@)
	if idx := strings.Index(jid, ":"); idx != -1 {
		atIdx := strings.Index(jid, "@")
		if atIdx != -1 && idx < atIdx {
			// Remove the device part
			return jid[:idx] + jid[atIdx:]
		}
	}
	return jid
}

// extractPhoneFromJID extracts phone number from JID
// e.g., "6281234567890@s.whatsapp.net" -> "6281234567890"
func extractPhoneFromJID(jid string) string {
	// First normalize
	normalized := normalizeJID(jid)
	// Then extract phone part (before @)
	if idx := strings.Index(normalized, "@"); idx != -1 {
		return normalized[:idx]
	}
	return normalized
}

// updateCustomerInsight creates or updates customer insight record
func (w *MessageWorker) updateCustomerInsight(ctx context.Context, payload *redis.MessagePayload) {
	// Normalize the JID to ensure consistent customer identification
	normalizedJID := normalizeJID(payload.SenderJID)
	phone := extractPhoneFromJID(payload.SenderJID)
	
	fmt.Printf("Customer insight: original JID=%s, normalized=%s, phone=%s\n", 
		payload.SenderJID, normalizedJID, phone)

	query := `
		INSERT INTO customer_insights (
			tenant_id, customer_jid, customer_phone, 
			message_count, last_message_at, first_message_at,
			last_message_summary, created_at, updated_at
		) VALUES ($1, $2, $3, 1, NOW(), NOW(), $4, NOW(), NOW())
		ON CONFLICT (tenant_id, customer_jid)
		DO UPDATE SET
			message_count = customer_insights.message_count + 1,
			last_message_at = NOW(),
			last_message_summary = EXCLUDED.last_message_summary,
			updated_at = NOW()
	`

	// Truncate message for summary
	summary := payload.MessageText
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}

	_, err := w.db.ExecContext(ctx, query,
		payload.TenantID,
		normalizedJID,
		phone,
		summary,
	)

	if err != nil {
		fmt.Printf("Failed to update customer insight: %v\n", err)
	}
}

// Encryption helper functions for API keys (mirrors handlers/ai.go)
func getEncryptionKeyWorker() []byte {
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

func decryptAPIKeyWorker(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(getEncryptionKeyWorker())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

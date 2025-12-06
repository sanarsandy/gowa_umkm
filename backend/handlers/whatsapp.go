package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gowa-backend/db"
	"gowa-backend/models"
	"gowa-backend/services/redis"
	"gowa-backend/services/whatsapp"

	"github.com/labstack/echo/v4"
)

var whatsappService *whatsapp.ClientService
var redisClient *redis.Client

// InitWhatsAppService initializes the WhatsApp service
func InitWhatsAppService() {
	// Initialize Redis client
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	var err error
	redisClient, err = redis.NewClient(redis.Config{
		Host:     redisHost,
		Port:     redisPort,
		Password: "",
		DB:       0,
	})
	
	if err != nil {
		fmt.Printf("Warning: Failed to connect to Redis: %v\n", err)
		// Continue without Redis - messages won't be queued for AI processing
		redisClient = nil
	} else {
		fmt.Println("Redis client initialized successfully")
	}

	// db.DB is *sqlx.DB, we need *sql.DB for the service
	whatsappService = whatsapp.NewClientService(db.DB.DB, redisClient)

	// Auto-reconnect WhatsApp clients that were previously connected
	go autoReconnectWhatsAppClients()
}

// autoReconnectWhatsAppClients reconnects WhatsApp clients that were previously connected
func autoReconnectWhatsAppClients() {
	fmt.Println("[Auto-Reconnect] Starting auto-reconnect for previously connected WhatsApp clients...")
	
	// Wait a bit for the server to fully start
	time.Sleep(3 * time.Second)
	
	// Query database for connected devices
	rows, err := db.DB.Query(`
		SELECT tenant_id, jid 
		FROM whatsapp_devices 
		WHERE is_connected = true AND jid IS NOT NULL AND jid != ''
	`)
	if err != nil {
		fmt.Printf("[Auto-Reconnect] Failed to query connected devices: %v\n", err)
		return
	}
	defer rows.Close()

	var reconnected int
	for rows.Next() {
		var tenantID, jid string
		if err := rows.Scan(&tenantID, &jid); err != nil {
			fmt.Printf("[Auto-Reconnect] Failed to scan row: %v\n", err)
			continue
		}

		fmt.Printf("[Auto-Reconnect] Attempting to reconnect tenant %s (JID: %s)...\n", tenantID, jid)
		
		// Try to reconnect
		ctx := context.Background()
		client, err := whatsappService.Connect(ctx, tenantID)
		if err != nil {
			fmt.Printf("[Auto-Reconnect] Failed to create client for tenant %s: %v\n", tenantID, err)
			// Mark as disconnected in database
			db.DB.Exec(`UPDATE whatsapp_devices SET is_connected = false WHERE tenant_id = $1`, tenantID)
			continue
		}

		// Connect the client
		if err := client.Connect(); err != nil {
			fmt.Printf("[Auto-Reconnect] Failed to connect client for tenant %s: %v\n", tenantID, err)
			// Mark as disconnected in database
			db.DB.Exec(`UPDATE whatsapp_devices SET is_connected = false WHERE tenant_id = $1`, tenantID)
			continue
		}

		// Check if actually connected
		if client.IsConnected() {
			fmt.Printf("[Auto-Reconnect] Successfully reconnected tenant %s\n", tenantID)
			reconnected++
		} else {
			fmt.Printf("[Auto-Reconnect] Client created but not connected for tenant %s\n", tenantID)
		}
	}

	fmt.Printf("[Auto-Reconnect] Completed. Reconnected %d client(s).\n", reconnected)
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return redisClient
}

// GetWhatsAppService returns the WhatsApp service instance
func GetWhatsAppService() *whatsapp.ClientService {
	return whatsappService
}

// ConnectWhatsApp initiates a WhatsApp connection
func ConnectWhatsApp(c echo.Context) error {
	fmt.Printf("[DEBUG] ConnectWhatsApp: handler called\n")
	
	// Get tenant ID from JWT claims
	tenantID := getTenantIDFromContext(c)
	fmt.Printf("[DEBUG] ConnectWhatsApp: tenantID = '%s'\n", tenantID)
	
	if tenantID == "" {
		fmt.Printf("[DEBUG] ConnectWhatsApp: tenantID is empty, returning error\n")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found. Please create a tenant first.",
		})
	}

	ctx := c.Request().Context()

	// Check if already connected
	fmt.Printf("[DEBUG] ConnectWhatsApp: calling whatsappService.GetStatus\n")
	isConnected, jid, err := whatsappService.GetStatus(ctx, tenantID)
	if err != nil {
		fmt.Printf("[DEBUG] ConnectWhatsApp: GetStatus error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to check status: %v", err),
		})
	}
	fmt.Printf("[DEBUG] ConnectWhatsApp: GetStatus result - isConnected=%v, jid=%s\n", isConnected, jid)

	if isConnected {
		fmt.Printf("[DEBUG] ConnectWhatsApp: already connected, returning\n")
		return c.JSON(http.StatusOK, models.WhatsAppConnectionResponse{
			Status:  "already_connected",
			Message: fmt.Sprintf("Already connected as %s", jid),
		})
	}

	// Initiate connection
	fmt.Printf("[DEBUG] ConnectWhatsApp: calling whatsappService.Connect\n")
	client, err := whatsappService.Connect(ctx, tenantID)
	if err != nil {
		fmt.Printf("[DEBUG] ConnectWhatsApp: Connect error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect: %v", err),
		})
	}
	fmt.Printf("[DEBUG] ConnectWhatsApp: Connect successful, checking if logged in\n")

	// Check if already logged in
	fmt.Printf("[DEBUG] ConnectWhatsApp: calling client.IsLoggedIn\n")
	if client.IsLoggedIn() {
		fmt.Printf("[DEBUG] ConnectWhatsApp: client is already logged in\n")
		return c.JSON(http.StatusOK, models.WhatsAppConnectionResponse{
			Status:  "connected",
			Message: "Successfully connected (already logged in)",
		})
	}

	// Need to pair - return stream URL for QR codes
	fmt.Printf("[DEBUG] ConnectWhatsApp: pairing required, returning stream URL\n")
	streamURL := fmt.Sprintf("/api/whatsapp/qr/stream?tenant_id=%s", tenantID)

	return c.JSON(http.StatusOK, models.WhatsAppConnectionResponse{
		Status:    "pairing_required",
		Message:   "Please scan QR code to pair",
		StreamURL: streamURL,
	})
}

// DisconnectWhatsApp disconnects WhatsApp
func DisconnectWhatsApp(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found. Please create a tenant first.",
		})
	}

	ctx := c.Request().Context()

	if err := whatsappService.Disconnect(ctx, tenantID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to disconnect: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "disconnected",
		"message": "Successfully disconnected",
	})
}

// GetWhatsAppStatus returns the current WhatsApp connection status
func GetWhatsAppStatus(c echo.Context) error {
	fmt.Printf("[DEBUG] GetWhatsAppStatus: handler called\n")
	
	tenantID := getTenantIDFromContext(c)
	fmt.Printf("[DEBUG] GetWhatsAppStatus: tenantID = '%s'\n", tenantID)
	
	if tenantID == "" {
		fmt.Printf("[DEBUG] GetWhatsAppStatus: tenantID is empty, returning error\n")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found. Please create a tenant first.",
		})
	}

	ctx := c.Request().Context()

	isConnected, jid, err := whatsappService.GetStatus(ctx, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get status: %v", err),
		})
	}

	status := "disconnected"
	if isConnected {
		status = "connected"
	}

	return c.JSON(http.StatusOK, models.WhatsAppStatusResponse{
		IsConnected: isConnected,
		JID:         jid,
		Status:      status,
	})
}

// StreamQRCode streams QR codes via Server-Sent Events (SSE)
func StreamQRCode(c echo.Context) error {
	fmt.Printf("[DEBUG] StreamQRCode: handler called\n")
	
	// Get tenant ID from JWT claims (not from query param for security)
	tenantID := getTenantIDFromContext(c)
	fmt.Printf("[DEBUG] StreamQRCode: tenantID = '%s'\n", tenantID)
	
	if tenantID == "" {
		// Send error via SSE
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
		c.Response().WriteHeader(http.StatusOK)
		sendSSE(c, "error", map[string]string{
			"error": "Tenant not found. Please login again or create a tenant.",
		})
		return nil
	}

	// Set headers for SSE
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	ctx := c.Request().Context()

	// Get client manager
	manager := whatsapp.GetGlobalClientManager()
	client, err := manager.GetClient(tenantID)
	if err != nil {
		fmt.Printf("[DEBUG] StreamQRCode: client not found, error: %v\n", err)
		// Send error event
		sendSSE(c, "error", map[string]string{
			"error": "WhatsApp client not initialized. Please click 'Connect WhatsApp' first.",
		})
		return nil
	}

	fmt.Printf("[DEBUG] StreamQRCode: client found, getting QR channel\n")
	
	// Get QR channel with retry logic
	var qrChan <-chan whatsapp.QRChannelResult
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		qrChan, err = whatsapp.GetQRChannel(ctx, client)
		if err == nil {
			fmt.Printf("[DEBUG] StreamQRCode: QR channel obtained successfully\n")
			break
		}
		fmt.Printf("[DEBUG] StreamQRCode: QR channel attempt %d failed: %v\n", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}
	
	if err != nil {
		fmt.Printf("[DEBUG] StreamQRCode: failed to get QR channel after retries: %v\n", err)
		sendSSE(c, "error", map[string]string{
			"error": fmt.Sprintf("Failed to initialize QR code generation: %v. Please try reconnecting.", err),
		})
		return nil
	}

	// Connect the client if not already connected
	// This must be done AFTER getting the QR channel
	if !client.IsConnected() {
		fmt.Printf("[DEBUG] StreamQRCode: connecting client...\n")
		if err := client.Connect(); err != nil {
			fmt.Printf("[DEBUG] StreamQRCode: failed to connect client: %v\n", err)
			sendSSE(c, "error", map[string]string{
				"error": fmt.Sprintf("Failed to connect to WhatsApp: %v", err),
			})
			return nil
		}
		fmt.Printf("[DEBUG] StreamQRCode: client connected successfully\n")
	}

	// Stream QR codes
	fmt.Printf("[DEBUG] StreamQRCode: starting QR stream loop\n")
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[DEBUG] StreamQRCode: context cancelled\n")
			return nil
		case qrResult, ok := <-qrChan:
			if !ok {
				// Channel closed
				fmt.Printf("[DEBUG] StreamQRCode: QR channel closed\n")
				return nil
			}

			fmt.Printf("[DEBUG] StreamQRCode: received QR event: %s\n", qrResult.Event)
			
			// Send QR event
			if err := sendSSE(c, qrResult.Event, qrResult); err != nil {
				fmt.Printf("[DEBUG] StreamQRCode: error sending SSE: %v\n", err)
				return err
			}

			// If success or timeout, close the stream
			if qrResult.Event == "success" || qrResult.Event == "timeout" || qrResult.Event == "error" {
				fmt.Printf("[DEBUG] StreamQRCode: stream ending with event: %s\n", qrResult.Event)
				return nil
			}

			c.Response().Flush()
		case <-time.After(90 * time.Second):
			// Timeout after 90 seconds (increased from 60)
			fmt.Printf("[DEBUG] StreamQRCode: timeout after 90 seconds\n")
			sendSSE(c, "timeout", map[string]string{
				"error": "QR code request timeout. Please try again.",
			})
			return nil
		}
	}
}

// Helper function to send SSE events
func sendSSE(c echo.Context, event string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c.Response(), "event: %s\ndata: %s\n\n", event, string(jsonData))
	if err != nil {
		return err
	}

	c.Response().Flush()
	return nil
}

// HealthCheck returns the current WhatsApp connection health status
func HealthCheck(c echo.Context) error {
	fmt.Printf("[DEBUG] HealthCheck: handler called\n")
	
	tenantID := getTenantIDFromContext(c)
	fmt.Printf("[DEBUG] HealthCheck: tenantID = '%s'\n", tenantID)
	
	if tenantID == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "no_tenant",
			"message": "No tenant found. Please login or create a tenant.",
		})
	}

	// Check if client exists
	manager := whatsapp.GetGlobalClientManager()
	client, err := manager.GetClient(tenantID)
	if err != nil {
		fmt.Printf("[DEBUG] HealthCheck: no client found\n")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "no_client",
			"message":    "No WhatsApp client initialized",
			"tenant_id":  tenantID,
			"suggestion": "Click 'Connect WhatsApp' to initialize",
		})
	}

	// Check if connected
	if client.IsConnected() {
		jid := ""
		if client.Store != nil && client.Store.ID != nil {
			jid = client.Store.ID.String()
		}
		fmt.Printf("[DEBUG] HealthCheck: client connected, jid: %s\n", jid)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "connected",
			"message":   "WhatsApp is connected and ready",
			"jid":       jid,
			"tenant_id": tenantID,
		})
	}

	fmt.Printf("[DEBUG] HealthCheck: client exists but not connected\n")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "disconnected",
		"message":    "Client exists but not connected",
		"tenant_id":  tenantID,
		"suggestion": "Scan QR code to connect",
	})
}

// SendWhatsAppMessage sends a message to a WhatsApp recipient
func SendWhatsAppMessage(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found. Please create a tenant first.",
		})
	}

	// Parse request body
	var req struct {
		RecipientJID string `json:"recipient_jid"`
		Message      string `json:"message"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.RecipientJID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Recipient JID is required",
		})
	}

	if req.Message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Message is required",
		})
	}

	// Trim message
	if len(req.Message) > 4096 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Message too long (max 4096 characters)",
		})
	}

	ctx := c.Request().Context()

	fmt.Printf("[DEBUG] SendWhatsAppMessage: tenantID=%s, recipientJID=%s, message=%s\n", tenantID, req.RecipientJID, req.Message)

	// Send message
	messageID, err := whatsappService.SendMessage(ctx, tenantID, req.RecipientJID, req.Message)
	if err != nil {
		fmt.Printf("[DEBUG] SendWhatsAppMessage: error sending message: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	fmt.Printf("[DEBUG] SendWhatsAppMessage: message sent successfully, messageID=%s\n", messageID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"message_id": messageID,
		"status":     "sent",
	})
}

// getTenantIDFromContext extracts tenant ID from JWT claims by:
// 1. Getting user_id from JWT claims
// 2. Querying database to get tenant_id for that user
func getTenantIDFromContext(c echo.Context) string {
	// First, get user_id from JWT
	userID := getUserIDFromContext(c)
	fmt.Printf("[DEBUG] getTenantIDFromContext: userID from getUserIDFromContext = '%s'\n", userID)
	
	if userID == "" {
		fmt.Printf("[DEBUG] getTenantIDFromContext: userID is empty from JWT\n")
		return ""
	}

	// Query database to get tenant_id for this user
	var tenantID string
	query := `SELECT id FROM tenants WHERE user_id = $1 AND is_active = true LIMIT 1`
	fmt.Printf("[DEBUG] getTenantIDFromContext: executing query with user_id = '%s'\n", userID)
	err := db.DB.QueryRow(query, userID).Scan(&tenantID)
	if err != nil {
		// Tenant not found - this is OK, user might not have created tenant yet
		if err == sql.ErrNoRows {
			fmt.Printf("[DEBUG] getTenantIDFromContext: tenant not found for user_id: %s (sql.ErrNoRows)\n", userID)
			return ""
		}
		// Other database error
		fmt.Printf("[DEBUG] getTenantIDFromContext: database error for user_id %s: %v\n", userID, err)
		return ""
	}

	fmt.Printf("[DEBUG] getTenantIDFromContext: found tenant_id: %s for user_id: %s\n", tenantID, userID)
	return tenantID
}

// ClearChatMessages deletes all messages for a specific chat/customer
// DELETE /api/whatsapp/messages/:jid
func ClearChatMessages(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found",
		})
	}

	customerJIDRaw := c.Param("jid")
	if customerJIDRaw == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Customer JID is required",
		})
	}

	// URL decode the JID (@ is encoded as %40)
	customerJID, err := url.QueryUnescape(customerJIDRaw)
	if err != nil {
		customerJID = customerJIDRaw // fallback to raw value
	}

	fmt.Printf("[ClearChat] TenantID: %s, CustomerJID: %s (raw: %s)\n", tenantID, customerJID, customerJIDRaw)

	// Delete messages for this customer (both incoming and outgoing)
	query := `
		DELETE FROM whatsapp_messages 
		WHERE tenant_id = $1 AND (chat_jid = $2 OR sender_jid = $2)
	`
	
	fmt.Printf("[ClearChat] Executing query with jid: '%s'\n", customerJID)
	result, err := db.DB.Exec(query, tenantID, customerJID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to clear chat messages",
		})
	}

	rowsAffected, _ := result.RowsAffected()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":        true,
		"message":        "Chat cleared successfully",
		"deleted_count":  rowsAffected,
	})
}

// SendWhatsAppMedia sends a media message to a WhatsApp recipient
func SendWhatsAppMedia(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found. Please create a tenant first.",
		})
	}

	// Parse multipart form
	recipientJID := c.FormValue("recipient_jid")
	caption := c.FormValue("caption")
	mediaType := c.FormValue("media_type") // image, document

	if recipientJID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Recipient JID is required",
		})
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File is required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// Read file content
	mediaData, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read file",
		})
	}

	// Determine media type if not provided
	if mediaType == "" {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			mediaType = "image"
		} else if ext == ".pdf" || ext == ".doc" || ext == ".docx" {
			mediaType = "document"
		} else {
			mediaType = "document" // Default fallback
		}
	}

	ctx := c.Request().Context()

	fmt.Printf("[DEBUG] SendWhatsAppMedia: tenantID=%s, recipientJID=%s, type=%s, file=%s\n", tenantID, recipientJID, mediaType, file.Filename)

	// Send media message
	messageID, err := whatsappService.SendMediaMessage(ctx, tenantID, recipientJID, mediaData, mediaType, file.Filename, caption)
	if err != nil {
		fmt.Printf("[DEBUG] SendWhatsAppMedia: error sending media: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	fmt.Printf("[DEBUG] SendWhatsAppMedia: media sent successfully, messageID=%s\n", messageID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"message_id": messageID,
		"status":     "sent",
	})
}

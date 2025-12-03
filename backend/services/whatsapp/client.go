package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"gowa-backend/services/redis"
	"gowa-backend/services/websocket"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver for whatsmeow store
)

// ClientService handles WhatsApp client operations
type ClientService struct {
	db            *sql.DB
	clientManager *ClientManager
	redisClient   *redis.Client
	logger        waLog.Logger
}

// NewClientService creates a new client service
func NewClientService(db *sql.DB, redisClient *redis.Client) *ClientService {
	return &ClientService{
		db:            db,
		clientManager: GetGlobalClientManager(),
		redisClient:   redisClient,
		logger:        waLog.Stdout("ClientService", "INFO", true),
	}
}

// Connect initiates a WhatsApp connection for a tenant
// Note: This creates the client but does NOT connect yet
// The actual connection happens when QR stream is requested
func (s *ClientService) Connect(ctx context.Context, tenantID string) (*whatsmeow.Client, error) {
	// Check if client already exists
	existingClient, err := s.clientManager.GetClient(tenantID)
	if err == nil && existingClient != nil {
		// Client exists - check if it's connected
		if existingClient.IsConnected() {
			s.logger.Infof("[%s] Client already connected", tenantID)
			return existingClient, nil
		}
		
		// Client exists but not connected - remove it and create a new one
		s.logger.Warnf("[%s] Client exists but not connected, removing and recreating", tenantID)
		existingClient.Disconnect()
		s.clientManager.RemoveClient(tenantID)
	}

	// Create persistent store directory
	storeDir := "/app/data/whatsapp_stores"
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		s.logger.Errorf("Failed to create store directory: %v", err)
		// Continue anyway - might work if directory already exists
	}

	// Create device store with persistent path
	storePath := fmt.Sprintf("file:%s/store_%s.db?_foreign_keys=on", storeDir, tenantID)
	s.logger.Infof("Creating WhatsApp store at: %s", storePath)
	
	container, err := sqlstore.New(ctx, "sqlite3", storePath, s.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Create new client
	client := whatsmeow.NewClient(deviceStore, s.logger)
	client.AddEventHandler(s.eventHandler(tenantID))

	// Add to manager
	if err := s.clientManager.AddClient(tenantID, client); err != nil {
		return nil, fmt.Errorf("failed to add client to manager: %w", err)
	}

	// DON'T connect yet - QR channel must be set up first
	// Connection will happen in the QR stream handler
	s.logger.Infof("[%s] WhatsApp client created (not connected yet)", tenantID)
	return client, nil
}

// Disconnect disconnects a WhatsApp client
func (s *ClientService) Disconnect(ctx context.Context, tenantID string) error {
	client, err := s.clientManager.GetClient(tenantID)
	if err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	// Logout
	if client.IsConnected() {
		if err := client.Logout(ctx); err != nil {
			s.logger.Errorf("Failed to logout: %v", err)
		}
	}

	// Disconnect
	client.Disconnect()

	// Remove from manager
	if err := s.clientManager.RemoveClient(tenantID); err != nil {
		return fmt.Errorf("failed to remove client: %w", err)
	}

	// Update database
	query := `UPDATE whatsapp_devices SET is_connected = false, updated_at = NOW() WHERE tenant_id = $1`
	if _, err := s.db.ExecContext(ctx, query, tenantID); err != nil {
		s.logger.Errorf("Failed to update device status: %v", err)
	}

	return nil
}

// GetStatus returns the connection status for a tenant
func (s *ClientService) GetStatus(ctx context.Context, tenantID string) (bool, string, error) {
	client, err := s.clientManager.GetClient(tenantID)
	if err != nil {
		// Check database for last known status
		var isConnected bool
		var jid string
		query := `SELECT is_connected, COALESCE(jid, '') FROM whatsapp_devices WHERE tenant_id = $1`
		err := s.db.QueryRowContext(ctx, query, tenantID).Scan(&isConnected, &jid)
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		if err != nil {
			return false, "", fmt.Errorf("failed to query status: %w", err)
		}
		return isConnected, jid, nil
	}

	isConnected := client.IsConnected()
	jid := ""
	if client.Store != nil && client.Store.ID != nil {
		jid = client.Store.ID.String()
	}

	return isConnected, jid, nil
}

// eventHandler creates an event handler for a specific tenant
func (s *ClientService) eventHandler(tenantID string) func(interface{}) {
	return func(evt interface{}) {
		// Log all event types for debugging
		s.logger.Infof("[%s] Event received: %T", tenantID, evt)
		
		switch v := evt.(type) {
		case *events.Message:
			s.handleMessage(tenantID, v)
		case *events.Connected:
			s.handleConnected(tenantID)
		case *events.Disconnected:
			s.handleDisconnected(tenantID)
		case *events.LoggedOut:
			s.handleLoggedOut(tenantID)
		case *events.PairSuccess:
			s.handlePairSuccess(tenantID, v)
		case *events.HistorySync:
			s.handleHistorySync(tenantID, v)
		case *events.Receipt:
			// Log and process receipt events (message delivery/read status)
			s.logger.Infof("[%s] Receipt: Type=%s, MessageIDs=%v, From=%s, Chat=%s", tenantID, v.Type, v.MessageIDs, v.MessageSource.Sender, v.MessageSource.Chat)
			s.handleReceipt(tenantID, v)
		}
	}
}

// handleMessage processes incoming AND outgoing messages
func (s *ClientService) handleMessage(tenantID string, evt *events.Message) {
	s.logger.Infof("[%s] Received message from %s (IsFromMe: %v): %s", tenantID, evt.Info.Sender, evt.Info.IsFromMe, evt.Message.GetConversation())
	
	// Skip group messages for now
	if evt.Info.IsGroup {
		s.logger.Infof("[%s] Skipping group message", tenantID)
		return
	}

	// Extract message content
	messageText := ""
	messageType := "text"
	
	if evt.Message.GetConversation() != "" {
		messageText = evt.Message.GetConversation()
	} else if evt.Message.GetExtendedTextMessage() != nil {
		messageText = evt.Message.GetExtendedTextMessage().GetText()
	} else if evt.Message.GetImageMessage() != nil {
		messageType = "image"
		messageText = evt.Message.GetImageMessage().GetCaption()
	} else if evt.Message.GetVideoMessage() != nil {
		messageType = "video"
		messageText = evt.Message.GetVideoMessage().GetCaption()
	} else if evt.Message.GetDocumentMessage() != nil {
		messageType = "document"
		messageText = evt.Message.GetDocumentMessage().GetCaption()
	}

	// Normalize JIDs - remove device part for consistent customer identification
	// For outgoing messages (IsFromMe=true), chat JID is the recipient (customer)
	// For incoming messages (IsFromMe=false), sender JID is the customer
	normalizedChatJID := evt.Info.Chat.ToNonAD().String()
	var normalizedSenderJID string
	var customerJID string
	
	if evt.Info.IsFromMe {
		// Outgoing message: we are the sender, chat is the customer
		normalizedSenderJID = evt.Info.Chat.ToNonAD().String() // Use chat as "sender" for consistency
		customerJID = normalizedChatJID
		s.logger.Infof("[%s] Outgoing message to customer: %s", tenantID, customerJID)
	} else {
		// Incoming message: sender is the customer
		normalizedSenderJID = evt.Info.Sender.ToNonAD().String()
		customerJID = normalizedSenderJID
		s.logger.Infof("[%s] Incoming message from customer: %s", tenantID, customerJID)
	}
	
	s.logger.Infof("[%s] Normalized JIDs - Chat: %s, Sender: %s, Customer: %s", tenantID, normalizedChatJID, normalizedSenderJID, customerJID)

	// Store message in database
	ctx := context.Background()
	query := `
		INSERT INTO whatsapp_messages (
			tenant_id, message_id, chat_jid, sender_jid, 
			message_type, message_text, is_from_me, is_group, 
			timestamp, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (tenant_id, message_id) DO NOTHING
	`
	
	_, err := s.db.ExecContext(ctx, query,
		tenantID,
		evt.Info.ID,
		normalizedChatJID,
		customerJID, // Use customer JID as sender for consistent storage
		messageType,
		messageText,
		evt.Info.IsFromMe,
		evt.Info.IsGroup,
		evt.Info.Timestamp.Unix(),
	)
	
	if err != nil {
		s.logger.Errorf("Failed to store message: %v", err)
		return
	}

	// Push to Redis queue for AI processing (only for incoming messages)
	if s.redisClient != nil && !evt.Info.IsFromMe {
		payload := &redis.MessagePayload{
			TenantID:    tenantID,
			MessageID:   evt.Info.ID,
			SenderJID:   customerJID,
			MessageText: messageText,
			Timestamp:   evt.Info.Timestamp.Unix(),
		}

		if err := s.redisClient.PushToAIQueue(ctx, payload); err != nil {
			s.logger.Errorf("Failed to push message to Redis: %v", err)
		} else {
			s.logger.Infof("[%s] Message pushed to AI queue", tenantID)
		}
	}
	
	// Resolve JIDs to phone format for consistent matching in frontend
	resolvedCustomerJID := s.resolveJID(tenantID, customerJID)
	resolvedChatJID := s.resolveJID(tenantID, normalizedChatJID)
	
	// Broadcast new message event via WebSocket
	hub := websocket.GetHub()
	hub.BroadcastToTenant(tenantID, websocket.EventNewMessage, map[string]interface{}{
		"message_id":   evt.Info.ID,
		"sender_jid":   resolvedCustomerJID,
		"chat_jid":     resolvedChatJID,
		"message_text": messageText,
		"message_type": messageType,
		"timestamp":    evt.Info.Timestamp.Unix(),
		"is_from_me":   evt.Info.IsFromMe,
	})
	
	s.logger.Infof("[%s] Message stored and broadcasted (IsFromMe: %v, ResolvedJID: %s)", tenantID, evt.Info.IsFromMe, resolvedCustomerJID)
}

// handleConnected handles connection events
func (s *ClientService) handleConnected(tenantID string) {
	s.logger.Infof("[%s] Connected to WhatsApp", tenantID)
	
	// Update database
	ctx := context.Background()
	query := `UPDATE whatsapp_devices SET is_connected = true, last_connected_at = NOW(), updated_at = NOW() WHERE tenant_id = $1`
	if _, err := s.db.ExecContext(ctx, query, tenantID); err != nil {
		s.logger.Errorf("Failed to update connection status: %v", err)
	}
}

// handleDisconnected handles disconnection events
func (s *ClientService) handleDisconnected(tenantID string) {
	s.logger.Infof("[%s] Disconnected from WhatsApp", tenantID)
	
	// Update database
	ctx := context.Background()
	query := `UPDATE whatsapp_devices SET is_connected = false, updated_at = NOW() WHERE tenant_id = $1`
	if _, err := s.db.ExecContext(ctx, query, tenantID); err != nil {
		s.logger.Errorf("Failed to update connection status: %v", err)
	}
}

// handleLoggedOut handles logout events
func (s *ClientService) handleLoggedOut(tenantID string) {
	s.logger.Infof("[%s] Logged out from WhatsApp", tenantID)
	
	// Remove client from manager
	if err := s.clientManager.RemoveClient(tenantID); err != nil {
		s.logger.Errorf("Failed to remove client: %v", err)
	}
	
	// Update database
	ctx := context.Background()
	query := `UPDATE whatsapp_devices SET is_connected = false, updated_at = NOW() WHERE tenant_id = $1`
	if _, err := s.db.ExecContext(ctx, query, tenantID); err != nil {
		s.logger.Errorf("Failed to update connection status: %v", err)
	}
}

// pendingMappings stores MessageID -> JID for correlating Receipt events
var pendingMappings = make(map[string]map[string]string) // tenantID -> messageID -> jid

// handleReceipt handles receipt events and stores JID mappings
func (s *ClientService) handleReceipt(tenantID string, evt *events.Receipt) {
	chatJID := evt.MessageSource.Chat.String()
	
	for _, msgID := range evt.MessageIDs {
		// Initialize tenant map if needed
		if pendingMappings[tenantID] == nil {
			pendingMappings[tenantID] = make(map[string]string)
		}
		
		// Check if we have a pending JID for this message
		if existingJID, ok := pendingMappings[tenantID][msgID]; ok {
			// We found a pair! Create mapping
			if strings.Contains(existingJID, "@lid") && strings.Contains(chatJID, "@s.whatsapp.net") {
				s.storeJIDMapping(tenantID, existingJID, chatJID)
				delete(pendingMappings[tenantID], msgID)
			} else if strings.Contains(chatJID, "@lid") && strings.Contains(existingJID, "@s.whatsapp.net") {
				s.storeJIDMapping(tenantID, chatJID, existingJID)
				delete(pendingMappings[tenantID], msgID)
			}
		} else {
			// Store this JID for later correlation
			if strings.Contains(chatJID, "@lid") || strings.Contains(chatJID, "@s.whatsapp.net") {
				pendingMappings[tenantID][msgID] = chatJID
				s.logger.Infof("[%s] Stored pending JID for msgID %s: %s", tenantID, msgID, chatJID)
			}
		}
	}
}

// storeJIDMapping stores a mapping between @lid and @s.whatsapp.net JIDs
func (s *ClientService) storeJIDMapping(tenantID, lidJID, phoneJID string) {
	// Extract phone number from @s.whatsapp.net JID
	phoneNumber := strings.Split(phoneJID, "@")[0]
	// Remove device part if present (e.g., 62813:31@s.whatsapp.net -> 62813)
	if idx := strings.Index(phoneNumber, ":"); idx > 0 {
		phoneNumber = phoneNumber[:idx]
	}
	
	// Normalize phone JID (remove device part)
	normalizedPhoneJID := phoneNumber + "@s.whatsapp.net"
	
	ctx := context.Background()
	query := `
		INSERT INTO jid_mappings (tenant_id, lid_jid, phone_jid, phone_number)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (tenant_id, lid_jid) 
		DO UPDATE SET phone_jid = EXCLUDED.phone_jid, phone_number = EXCLUDED.phone_number, updated_at = NOW()
	`
	
	if _, err := s.db.ExecContext(ctx, query, tenantID, lidJID, normalizedPhoneJID, phoneNumber); err != nil {
		s.logger.Errorf("[%s] Failed to store JID mapping: %v", tenantID, err)
		return
	}
	
	s.logger.Infof("[%s] Stored JID mapping: %s -> %s (%s)", tenantID, lidJID, normalizedPhoneJID, phoneNumber)
	
	// Update existing messages with @lid JID to use phone JID
	updateQuery := `
		UPDATE whatsapp_messages 
		SET chat_jid = $1, sender_jid = $1
		WHERE tenant_id = $2 AND (chat_jid = $3 OR sender_jid = $3)
	`
	if result, err := s.db.ExecContext(ctx, updateQuery, normalizedPhoneJID, tenantID, lidJID); err != nil {
		s.logger.Errorf("[%s] Failed to update messages with new JID: %v", tenantID, err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		s.logger.Infof("[%s] Updated %d messages with new JID mapping", tenantID, rows)
		
		// Broadcast a JID update event so frontend can refresh
		hub := websocket.GetHub()
		hub.BroadcastToTenant(tenantID, websocket.EventNewMessage, map[string]interface{}{
			"type":        "jid_mapping",
			"old_jid":     lidJID,
			"new_jid":     normalizedPhoneJID,
			"message_text": "", // Empty to not create duplicate message
		})
	}
}

// resolveJID resolves a JID to the phone format (@s.whatsapp.net) if possible
func (s *ClientService) resolveJID(tenantID, jid string) string {
	// Already in phone format
	if strings.Contains(jid, "@s.whatsapp.net") {
		return jid
	}
	
	// Try to find mapping for @lid JID
	if strings.Contains(jid, "@lid") {
		ctx := context.Background()
		var phoneJID string
		query := `SELECT phone_jid FROM jid_mappings WHERE tenant_id = $1 AND lid_jid = $2`
		if err := s.db.QueryRowContext(ctx, query, tenantID, jid).Scan(&phoneJID); err == nil {
			s.logger.Infof("[%s] Resolved JID %s -> %s", tenantID, jid, phoneJID)
			return phoneJID
		}
	}
	
	return jid
}

// handlePairSuccess handles successful pairing
func (s *ClientService) handlePairSuccess(tenantID string, evt *events.PairSuccess) {
	s.logger.Infof("[%s] Successfully paired with %s", tenantID, evt.ID.String())
	
	// Save device info to database
	ctx := context.Background()
	query := `
		INSERT INTO whatsapp_devices (tenant_id, jid, is_connected, last_connected_at, platform, created_at, updated_at)
		VALUES ($1, $2, true, NOW(), 'web', NOW(), NOW())
		ON CONFLICT (tenant_id)
		DO UPDATE SET
			jid = EXCLUDED.jid,
			is_connected = true,
			last_connected_at = NOW(),
			updated_at = NOW()
	`
	
	if _, err := s.db.ExecContext(ctx, query, tenantID, evt.ID.String()); err != nil {
		s.logger.Errorf("Failed to save device info: %v", err)
	}
}

// handleHistorySync handles history sync events (messages sent from phone)
func (s *ClientService) handleHistorySync(tenantID string, evt *events.HistorySync) {
	s.logger.Infof("[%s] History sync event received: %d conversations", tenantID, len(evt.Data.Conversations))
	
	for _, conv := range evt.Data.Conversations {
		chatJID := conv.GetID()
		s.logger.Infof("[%s] Processing conversation: %s with %d messages", tenantID, chatJID, len(conv.Messages))
		
		for _, historyMsg := range conv.Messages {
			msg := historyMsg.Message
			if msg == nil || msg.Message == nil {
				continue
			}
			
			msgInfo := msg.GetMessageTimestamp()
			isFromMe := msg.GetKey().GetFromMe()
			
			// Only process recent messages (last 5 minutes) to avoid re-syncing old messages
			msgTime := int64(msgInfo)
			now := time.Now().Unix()
			if now - msgTime > 300 { // 5 minutes
				continue
			}
			
			s.logger.Infof("[%s] History sync message: IsFromMe=%v, timestamp=%d", tenantID, isFromMe, msgTime)
			
			// Extract message content
			messageText := ""
			messageType := "text"
			
			if msg.Message.GetConversation() != "" {
				messageText = msg.Message.GetConversation()
			} else if msg.Message.GetExtendedTextMessage() != nil {
				messageText = msg.Message.GetExtendedTextMessage().GetText()
			} else if msg.Message.GetImageMessage() != nil {
				messageType = "image"
				messageText = msg.Message.GetImageMessage().GetCaption()
			} else if msg.Message.GetVideoMessage() != nil {
				messageType = "video"
				messageText = msg.Message.GetVideoMessage().GetCaption()
			} else if msg.Message.GetDocumentMessage() != nil {
				messageType = "document"
				messageText = msg.Message.GetDocumentMessage().GetCaption()
			}
			
			if messageText == "" && messageType == "text" {
				continue // Skip empty text messages
			}
			
			// Parse JID
			parsedChatJID, err := types.ParseJID(chatJID)
			if err != nil {
				s.logger.Errorf("[%s] Failed to parse chat JID: %v", tenantID, err)
				continue
			}
			
			normalizedChatJID := parsedChatJID.ToNonAD().String()
			messageID := msg.GetKey().GetID()
			
			// Store message in database
			ctx := context.Background()
			query := `
				INSERT INTO whatsapp_messages (
					tenant_id, message_id, chat_jid, sender_jid, 
					message_type, message_text, is_from_me, is_group, 
					timestamp, created_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
				ON CONFLICT (tenant_id, message_id) DO NOTHING
			`
			
			_, err = s.db.ExecContext(ctx, query,
				tenantID,
				messageID,
				normalizedChatJID,
				normalizedChatJID, // For consistency with handleMessage
				messageType,
				messageText,
				isFromMe,
				parsedChatJID.Server == types.GroupServer,
				msgTime,
			)
			
			if err != nil {
				s.logger.Errorf("[%s] Failed to store history message: %v", tenantID, err)
				continue
			}
			
			// Broadcast via WebSocket
			hub := websocket.GetHub()
			hub.BroadcastToTenant(tenantID, websocket.EventNewMessage, map[string]interface{}{
				"message_id":   messageID,
				"sender_jid":   normalizedChatJID,
				"chat_jid":     normalizedChatJID,
				"message_text": messageText,
				"message_type": messageType,
				"timestamp":    msgTime,
				"is_from_me":   isFromMe,
			})
			
			s.logger.Infof("[%s] History message stored and broadcasted (IsFromMe: %v)", tenantID, isFromMe)
		}
	}
}

// SendMessage sends a text message to a WhatsApp recipient
func (s *ClientService) SendMessage(ctx context.Context, tenantID string, recipientJID string, message string) (string, error) {
	s.logger.Infof("[%s] SendMessage: starting, recipientJID=%s", tenantID, recipientJID)
	
	client, err := s.clientManager.GetClient(tenantID)
	if err != nil {
		s.logger.Errorf("[%s] SendMessage: client not found: %v", tenantID, err)
		return "", fmt.Errorf("WhatsApp client not found. Please connect first")
	}

	s.logger.Infof("[%s] SendMessage: client found, checking connection", tenantID)

	if !client.IsConnected() {
		s.logger.Errorf("[%s] SendMessage: client not connected", tenantID)
		return "", fmt.Errorf("WhatsApp not connected. Please reconnect")
	}

	s.logger.Infof("[%s] SendMessage: client connected, parsing JID", tenantID)

	// Parse recipient JID
	jid, err := types.ParseJID(recipientJID)
	if err != nil {
		s.logger.Errorf("[%s] SendMessage: invalid JID: %v", tenantID, err)
		return "", fmt.Errorf("invalid recipient JID: %w", err)
	}

	// Remove device part from JID (required by whatsmeow for sending messages)
	jid = jid.ToNonAD()
	s.logger.Infof("[%s] SendMessage: JID normalized: %s, creating message", tenantID, jid.String())

	// Create message
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	s.logger.Infof("[%s] SendMessage: sending message via WhatsApp...", tenantID)

	// Send message
	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		s.logger.Errorf("[%s] SendMessage: failed to send: %v", tenantID, err)
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	s.logger.Infof("[%s] SendMessage: message sent! ID=%s", tenantID, resp.ID)

	messageID := resp.ID
	timestamp := resp.Timestamp.Unix()

	// Store sent message in database
	query := `
		INSERT INTO whatsapp_messages (
			tenant_id, message_id, chat_jid, sender_jid, 
			message_type, message_text, is_from_me, is_group, 
			timestamp, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (tenant_id, message_id) DO NOTHING
	`

	// Get our own JID for sender
	senderJID := ""
	if client.Store != nil && client.Store.ID != nil {
		senderJID = client.Store.ID.String()
	}

	_, dbErr := s.db.ExecContext(ctx, query,
		tenantID,
		messageID,
		recipientJID,
		senderJID,
		"text",
		message,
		true, // is_from_me
		false, // is_group (assuming direct message)
		timestamp,
	)

	if dbErr != nil {
		s.logger.Errorf("Failed to store sent message: %v", dbErr)
		// Don't return error - message was sent successfully
	}

	// Broadcast sent message via WebSocket so frontend can update in real-time
	hub := websocket.GetHub()
	hub.BroadcastToTenant(tenantID, websocket.EventNewMessage, map[string]interface{}{
		"message_id":   messageID,
		"sender_jid":   senderJID,
		"chat_jid":     recipientJID,
		"message_text": message,
		"message_type": "text",
		"timestamp":    timestamp,
		"is_from_me":   true,
	})

	s.logger.Infof("[%s] Message sent to %s: %s", tenantID, recipientJID, messageID)
	return messageID, nil
}

// SendMessageResult holds the result of sending a message
type SendMessageResult struct {
	MessageID string    `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

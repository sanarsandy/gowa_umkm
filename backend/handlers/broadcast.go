package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gowa-backend/db"

	"github.com/labstack/echo/v4"
)

// Broadcast represents a broadcast message
type Broadcast struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	Name            string     `json:"name" db:"name"`
	MessageContent  string     `json:"message_content" db:"message_content"`
	TemplateID      *string    `json:"template_id" db:"template_id"`
	Status          string     `json:"status" db:"status"`
	ScheduledAt     *time.Time `json:"scheduled_at" db:"scheduled_at"`
	StartedAt       *time.Time `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
	TotalRecipients int        `json:"total_recipients" db:"total_recipients"`
	SentCount       int        `json:"sent_count" db:"sent_count"`
	DeliveredCount  int        `json:"delivered_count" db:"delivered_count"`
	FailedCount     int        `json:"failed_count" db:"failed_count"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// BroadcastRecipient represents a recipient in a broadcast
type BroadcastRecipient struct {
	ID           string     `json:"id" db:"id"`
	BroadcastID  string     `json:"broadcast_id" db:"broadcast_id"`
	CustomerID   string     `json:"customer_id" db:"customer_id"`
	CustomerJID  string     `json:"customer_jid" db:"customer_jid"`
	CustomerName string     `json:"customer_name" db:"customer_name"`
	Status       string     `json:"status" db:"status"`
	MessageID    *string    `json:"message_id" db:"message_id"`
	SentAt       *time.Time `json:"sent_at" db:"sent_at"`
	DeliveredAt  *time.Time `json:"delivered_at" db:"delivered_at"`
	ErrorMessage *string    `json:"error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// GetBroadcasts returns all broadcasts for a tenant
func GetBroadcasts(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty list instead of 401
		// User is authenticated but hasn't created a tenant yet
		return c.JSON(http.StatusOK, map[string]interface{}{
			"broadcasts": []Broadcast{},
			"total":      0,
		})
	}

	status := c.QueryParam("status")

	var broadcasts []Broadcast
	var query string
	var args []interface{}

	if status != "" && status != "all" {
		query = `
			SELECT id, tenant_id, name, message_content, template_id, status,
			       scheduled_at, started_at, completed_at, total_recipients,
			       sent_count, delivered_count, failed_count, created_at, updated_at
			FROM broadcasts
			WHERE tenant_id = $1 AND status = $2
			ORDER BY created_at DESC
		`
		args = []interface{}{tenantID, status}
	} else {
		query = `
			SELECT id, tenant_id, name, message_content, template_id, status,
			       scheduled_at, started_at, completed_at, total_recipients,
			       sent_count, delivered_count, failed_count, created_at, updated_at
			FROM broadcasts
			WHERE tenant_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{tenantID}
	}

	if err := db.DB.Select(&broadcasts, query, args...); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch broadcasts")
	}

	if broadcasts == nil {
		broadcasts = []Broadcast{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"broadcasts": broadcasts,
		"total":      len(broadcasts),
	})
}

// GetBroadcast returns a single broadcast with recipients
func GetBroadcast(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	broadcastID := c.Param("id")

	var broadcast Broadcast
	query := `
		SELECT id, tenant_id, name, message_content, template_id, status,
		       scheduled_at, started_at, completed_at, total_recipients,
		       sent_count, delivered_count, failed_count, created_at, updated_at
		FROM broadcasts
		WHERE id = $1 AND tenant_id = $2
	`

	if err := db.DB.Get(&broadcast, query, broadcastID, tenantID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Broadcast not found")
	}

	// Get recipients
	var recipients []BroadcastRecipient
	recipientQuery := `
		SELECT br.id, br.broadcast_id, br.customer_id, br.customer_jid, 
		       COALESCE(ci.customer_name, ci.customer_phone, br.customer_jid) as customer_name,
		       br.status, br.message_id, br.sent_at, br.delivered_at, br.error_message, br.created_at
		FROM broadcast_recipients br
		LEFT JOIN customer_insights ci ON ci.id = br.customer_id
		WHERE br.broadcast_id = $1
		ORDER BY br.created_at ASC
	`

	if err := db.DB.Select(&recipients, recipientQuery, broadcastID); err != nil {
		recipients = []BroadcastRecipient{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"broadcast":  broadcast,
		"recipients": recipients,
	})
}

// CreateBroadcast creates a new broadcast
func CreateBroadcast(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	var req struct {
		Name           string   `json:"name"`
		MessageContent string   `json:"message_content"`
		TemplateID     *string  `json:"template_id"`
		CustomerIDs    []string `json:"customer_ids"`
		ScheduledAt    *string  `json:"scheduled_at"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" || req.MessageContent == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name and message content are required")
	}

	if len(req.CustomerIDs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one customer is required")
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	// Create broadcast
	var broadcast Broadcast
	var scheduledAt *time.Time
	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err == nil {
			scheduledAt = &t
		}
	}

	status := "draft"
	if scheduledAt != nil {
		status = "scheduled"
	}

	insertQuery := `
		INSERT INTO broadcasts (tenant_id, name, message_content, template_id, status, scheduled_at, total_recipients)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, tenant_id, name, message_content, template_id, status, scheduled_at, started_at, completed_at, total_recipients, sent_count, delivered_count, failed_count, created_at, updated_at
	`

	if err := tx.Get(&broadcast, insertQuery, tenantID, req.Name, req.MessageContent, req.TemplateID, status, scheduledAt, len(req.CustomerIDs)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create broadcast")
	}

	// Add recipients
	for _, customerID := range req.CustomerIDs {
		// Get customer JID
		var customerJID string
		jidQuery := `SELECT customer_jid FROM customer_insights WHERE id = $1 AND tenant_id = $2`
		if err := tx.Get(&customerJID, jidQuery, customerID, tenantID); err != nil {
			continue // Skip invalid customers
		}

		recipientQuery := `
			INSERT INTO broadcast_recipients (broadcast_id, customer_id, customer_jid)
			VALUES ($1, $2, $3)
		`
		tx.Exec(recipientQuery, broadcast.ID, customerID, customerJID)
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save broadcast")
	}

	return c.JSON(http.StatusCreated, broadcast)
}

// SendBroadcast starts sending a broadcast
func SendBroadcast(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	broadcastID := c.Param("id")

	// Get broadcast
	var broadcast Broadcast
	query := `SELECT * FROM broadcasts WHERE id = $1 AND tenant_id = $2`
	if err := db.DB.Get(&broadcast, query, broadcastID, tenantID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Broadcast not found")
	}

	if broadcast.Status != "draft" && broadcast.Status != "scheduled" {
		return echo.NewHTTPError(http.StatusBadRequest, "Broadcast cannot be sent in current status")
	}

	// Update status to sending
	updateQuery := `UPDATE broadcasts SET status = 'sending', started_at = NOW(), updated_at = NOW() WHERE id = $1`
	db.DB.Exec(updateQuery, broadcastID)

	// Get recipients
	var recipients []BroadcastRecipient
	recipientQuery := `SELECT id, customer_id, customer_jid FROM broadcast_recipients WHERE broadcast_id = $1 AND status = 'pending'`
	if err := db.DB.Select(&recipients, recipientQuery, broadcastID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get recipients")
	}

	// Send messages in background
	go sendBroadcastMessages(tenantID, broadcastID, broadcast.MessageContent, recipients)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Broadcast started",
		"status":  "sending",
	})
}

// sendBroadcastMessages sends messages to all recipients (runs in background)
func sendBroadcastMessages(tenantID, broadcastID, messageTemplate string, recipients []BroadcastRecipient) {
	ctx := context.Background()
	sentCount := 0
	failedCount := 0

	for _, recipient := range recipients {
		// Get customer name for personalization
		var customerName string
		nameQuery := `SELECT COALESCE(customer_name, customer_phone, '') FROM customer_insights WHERE id = $1`
		db.DB.Get(&customerName, nameQuery, recipient.CustomerID)

		// If no name found, use phone number from JID
		if customerName == "" {
			parts := strings.Split(recipient.CustomerJID, "@")
			if len(parts) > 0 {
				customerName = parts[0]
			}
		}

		// Personalize message - replace placeholders
		personalizedMessage := personalizeMessage(messageTemplate, customerName)

		// Send message via WhatsApp service
		messageID, err := whatsappService.SendMessage(ctx, tenantID, recipient.CustomerJID, personalizedMessage)

		if err != nil {
			// Mark as failed
			failedCount++
			errMsg := err.Error()
			updateQuery := `UPDATE broadcast_recipients SET status = 'failed', error_message = $1 WHERE id = $2`
			db.DB.Exec(updateQuery, errMsg, recipient.ID)
		} else {
			// Mark as sent
			sentCount++
			updateQuery := `UPDATE broadcast_recipients SET status = 'sent', message_id = $1, sent_at = NOW() WHERE id = $2`
			db.DB.Exec(updateQuery, messageID, recipient.ID)
		}

		// Update broadcast counts
		updateBroadcastQuery := `UPDATE broadcasts SET sent_count = $1, failed_count = $2, updated_at = NOW() WHERE id = $3`
		db.DB.Exec(updateBroadcastQuery, sentCount, failedCount, broadcastID)

		// Add small delay between messages to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	// Mark broadcast as completed
	completeQuery := `UPDATE broadcasts SET status = 'completed', completed_at = NOW(), updated_at = NOW() WHERE id = $1`
	db.DB.Exec(completeQuery, broadcastID)

	fmt.Printf("[Broadcast] Completed: %s - Sent: %d, Failed: %d\n", broadcastID, sentCount, failedCount)
}

// personalizeMessage replaces placeholders with actual values
// Supported placeholders:
// - {{nama}} or {{name}} - Customer name
// - {{phone}} - Customer phone number
func personalizeMessage(template, customerName string) string {
	result := template

	// Replace name placeholders (case insensitive)
	result = strings.ReplaceAll(result, "{{nama}}", customerName)
	result = strings.ReplaceAll(result, "{{Nama}}", customerName)
	result = strings.ReplaceAll(result, "{{NAMA}}", customerName)
	result = strings.ReplaceAll(result, "{{name}}", customerName)
	result = strings.ReplaceAll(result, "{{Name}}", customerName)
	result = strings.ReplaceAll(result, "{{NAME}}", customerName)

	return result
}

// CancelBroadcast cancels a broadcast
func CancelBroadcast(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	broadcastID := c.Param("id")

	// Allow cancelling draft, scheduled, and active (recurring) broadcasts
	query := `UPDATE broadcasts SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND tenant_id = $2 AND status IN ('draft', 'scheduled', 'active')`
	result, err := db.DB.Exec(query, broadcastID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cancel broadcast")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Broadcast cannot be cancelled")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Broadcast cancelled"})
}

// DeleteBroadcast deletes a broadcast
func DeleteBroadcast(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	broadcastID := c.Param("id")

	// Only allow deleting draft or cancelled broadcasts
	query := `DELETE FROM broadcasts WHERE id = $1 AND tenant_id = $2 AND status IN ('draft', 'cancelled')`
	result, err := db.DB.Exec(query, broadcastID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete broadcast")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Broadcast cannot be deleted")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Broadcast deleted"})
}

// GetBroadcastStats returns broadcast statistics
func GetBroadcastStats(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty stats instead of 401
		return c.JSON(http.StatusOK, map[string]interface{}{
			"total_broadcasts":   0,
			"total_messages_sent": 0,
			"total_delivered":    0,
			"total_failed":       0,
		})
	}

	var stats struct {
		TotalBroadcasts   int `db:"total_broadcasts"`
		TotalMessagesSent int `db:"total_messages_sent"`
		TotalDelivered    int `db:"total_delivered"`
		TotalFailed       int `db:"total_failed"`
	}

	query := `
		SELECT 
			COUNT(*) as total_broadcasts,
			COALESCE(SUM(sent_count), 0) as total_messages_sent,
			COALESCE(SUM(delivered_count), 0) as total_delivered,
			COALESCE(SUM(failed_count), 0) as total_failed
		FROM broadcasts
		WHERE tenant_id = $1
	`

	if err := db.DB.Get(&stats, query, tenantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get stats")
	}

	return c.JSON(http.StatusOK, stats)
}

// Placeholder for SQL NULL handling
var _ = sql.NullString{}
var _ = strings.TrimSpace


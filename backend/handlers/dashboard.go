package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"gowa-backend/db"

	"github.com/labstack/echo/v4"
)

// DashboardStats represents the dashboard statistics
type DashboardStats struct {
	TotalMessages     int     `json:"total_messages"`
	TotalCustomers    int     `json:"total_customers"`
	HotLeads          int     `json:"hot_leads"`
	PositiveSentiment float64 `json:"positive_sentiment"`
	MessagesToday     int     `json:"messages_today"`
	NewCustomersToday int     `json:"new_customers_today"`
	IsConnected       bool    `json:"is_connected"`
	ConnectedNumber   string  `json:"connected_number,omitempty"`
}

// RecentMessage represents a recent message for dashboard
type RecentMessage struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"customer_name"`
	CustomerJID  string    `json:"customer_jid"`
	MessageText  string    `json:"message_text"`
	Timestamp    time.Time `json:"timestamp"`
	IsFromMe     bool      `json:"is_from_me"`
}

// RecentCustomer represents a recent customer for dashboard
type RecentCustomer struct {
	ID              string    `json:"id"`
	CustomerName    string    `json:"customer_name"`
	CustomerJID     string    `json:"customer_jid"`
	Status          string    `json:"status"`
	LastMessageAt   time.Time `json:"last_message_at"`
	MessageCount    int       `json:"message_count"`
	LastMessageText string    `json:"last_message_text"`
}

// GetDashboardStats returns dashboard statistics for the authenticated tenant
func GetDashboardStats(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty stats instead of error
		// User is authenticated but hasn't created a tenant yet
		return c.JSON(http.StatusOK, DashboardStats{})
	}

	stats := DashboardStats{}

	// Get total messages
	var totalMessages sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM whatsapp_messages WHERE tenant_id = $1
	`, tenantID).Scan(&totalMessages)
	if err == nil && totalMessages.Valid {
		stats.TotalMessages = int(totalMessages.Int64)
	}

	// Get total customers
	var totalCustomers sql.NullInt64
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM customer_insights WHERE tenant_id = $1
	`, tenantID).Scan(&totalCustomers)
	if err == nil && totalCustomers.Valid {
		stats.TotalCustomers = int(totalCustomers.Int64)
	}

	// Get hot leads count
	var hotLeads sql.NullInt64
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM customer_insights 
		WHERE tenant_id = $1 AND status = 'hot_lead'
	`, tenantID).Scan(&hotLeads)
	if err == nil && hotLeads.Valid {
		stats.HotLeads = int(hotLeads.Int64)
	}

	// Get positive sentiment percentage
	var positiveCount, totalSentiment sql.NullInt64
	err = db.DB.QueryRow(`
		SELECT 
			COUNT(*) FILTER (WHERE sentiment = 'positive'),
			COUNT(*) FILTER (WHERE sentiment IS NOT NULL)
		FROM customer_insights 
		WHERE tenant_id = $1
	`, tenantID).Scan(&positiveCount, &totalSentiment)
	if err == nil && totalSentiment.Valid && totalSentiment.Int64 > 0 {
		stats.PositiveSentiment = float64(positiveCount.Int64) / float64(totalSentiment.Int64) * 100
	} else {
		stats.PositiveSentiment = 0
	}

	// Get messages today
	var messagesToday sql.NullInt64
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM whatsapp_messages 
		WHERE tenant_id = $1 AND created_at >= CURRENT_DATE
	`, tenantID).Scan(&messagesToday)
	if err == nil && messagesToday.Valid {
		stats.MessagesToday = int(messagesToday.Int64)
	}

	// Get new customers today
	var newCustomersToday sql.NullInt64
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM customer_insights 
		WHERE tenant_id = $1 AND created_at >= CURRENT_DATE
	`, tenantID).Scan(&newCustomersToday)
	if err == nil && newCustomersToday.Valid {
		stats.NewCustomersToday = int(newCustomersToday.Int64)
	}

	// Get WhatsApp connection status
	var isConnected bool
	var jid sql.NullString
	err = db.DB.QueryRow(`
		SELECT is_connected, jid FROM whatsapp_devices 
		WHERE tenant_id = $1 LIMIT 1
	`, tenantID).Scan(&isConnected, &jid)
	if err == nil {
		stats.IsConnected = isConnected
		if jid.Valid {
			stats.ConnectedNumber = jid.String
		}
	}

	return c.JSON(http.StatusOK, stats)
}

// GetRecentMessages returns the most recent messages for the dashboard
func GetRecentMessages(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty list instead of error
		return c.JSON(http.StatusOK, []RecentMessage{})
	}

	// Get limit from query param, default to 10
	limit := 10

	rows, err := db.DB.Query(`
		SELECT 
			m.id,
			COALESCE(ci.customer_name, m.sender_jid) as customer_name,
			m.sender_jid,
			COALESCE(m.message_text, '') as message_text,
			to_timestamp(m.timestamp) as timestamp,
			m.is_from_me
		FROM whatsapp_messages m
		LEFT JOIN customer_insights ci ON m.tenant_id = ci.tenant_id AND m.sender_jid = ci.customer_jid
		WHERE m.tenant_id = $1 AND m.is_from_me = false
		ORDER BY m.timestamp DESC
		LIMIT $2
	`, tenantID, limit)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch messages",
		})
	}
	defer rows.Close()

	messages := []RecentMessage{}
	for rows.Next() {
		var msg RecentMessage
		err := rows.Scan(
			&msg.ID,
			&msg.CustomerName,
			&msg.CustomerJID,
			&msg.MessageText,
			&msg.Timestamp,
			&msg.IsFromMe,
		)
		if err != nil {
			continue
		}

		// Truncate message text for display
		if len(msg.MessageText) > 100 {
			msg.MessageText = msg.MessageText[:100] + "..."
		}

		messages = append(messages, msg)
	}

	return c.JSON(http.StatusOK, messages)
}

// GetRecentCustomers returns the most recent/active customers
func GetRecentCustomers(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty list instead of error
		return c.JSON(http.StatusOK, []RecentCustomer{})
	}

	// Get limit from query param, default to 10
	limit := 10

	rows, err := db.DB.Query(`
		SELECT 
			id,
			COALESCE(customer_name, customer_jid) as customer_name,
			customer_jid,
			COALESCE(status, 'new') as status,
			COALESCE(last_message_at, created_at) as last_message_at,
			message_count,
			COALESCE(last_message_summary, '') as last_message_text
		FROM customer_insights
		WHERE tenant_id = $1
		ORDER BY last_message_at DESC NULLS LAST
		LIMIT $2
	`, tenantID, limit)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch customers",
		})
	}
	defer rows.Close()

	customers := []RecentCustomer{}
	for rows.Next() {
		var cust RecentCustomer
		err := rows.Scan(
			&cust.ID,
			&cust.CustomerName,
			&cust.CustomerJID,
			&cust.Status,
			&cust.LastMessageAt,
			&cust.MessageCount,
			&cust.LastMessageText,
		)
		if err != nil {
			continue
		}

		// Truncate message text for display
		if len(cust.LastMessageText) > 100 {
			cust.LastMessageText = cust.LastMessageText[:100] + "..."
		}

		customers = append(customers, cust)
	}

	return c.JSON(http.StatusOK, customers)
}



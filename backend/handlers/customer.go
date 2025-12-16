package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"gowa-backend/db"

	"github.com/labstack/echo/v4"
)

// Customer represents a customer/contact from WhatsApp
type Customer struct {
	ID                 string     `json:"id" db:"id"`
	TenantID           string     `json:"tenant_id" db:"tenant_id"`
	CustomerJID        string     `json:"customer_jid" db:"customer_jid"`
	CustomerName       *string    `json:"customer_name" db:"customer_name"`
	CustomerPhone      *string    `json:"customer_phone" db:"customer_phone"`
	Status             string     `json:"status" db:"status"`
	Sentiment          *string    `json:"sentiment" db:"sentiment"`
	Intent             *string    `json:"intent" db:"intent"`
	ProductInterest    *string    `json:"product_interest" db:"product_interest"`
	LastMessageSummary *string    `json:"last_message_summary" db:"last_message_summary"`
	MessageCount       int        `json:"message_count" db:"message_count"`
	LastMessageAt      *time.Time `json:"last_message_at" db:"last_message_at"`
	FirstMessageAt     *time.Time `json:"first_message_at" db:"first_message_at"`
	NeedsFollowUp      bool       `json:"needs_follow_up" db:"needs_follow_up"`
	Tags               *string    `json:"tags" db:"tags"`
	LeadScore          int        `json:"lead_score" db:"lead_score"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// CustomerListResponse represents paginated customer list
type CustomerListResponse struct {
	Customers  []Customer `json:"customers"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

// CustomerMessage represents a message in customer's chat history
type CustomerMessage struct {
	ID          string    `json:"id"`
	MessageText string    `json:"message_text"`
	MessageType string    `json:"message_type"`
	MediaURL    string    `json:"media_url,omitempty"`
	IsFromMe    bool      `json:"is_from_me"`
	Timestamp   time.Time `json:"timestamp"`
}

// CustomerDetailResponse represents customer with chat history
type CustomerDetailResponse struct {
	Customer Customer          `json:"customer"`
	Messages []CustomerMessage `json:"messages"`
}

// GetCustomers returns paginated list of customers with optional search/filter
func GetCustomers(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty list instead of error
		// User is authenticated but hasn't created a tenant yet
		return c.JSON(http.StatusOK, CustomerListResponse{
			Customers:  []Customer{},
			Total:      0,
			Page:       1,
			Limit:      20,
			TotalPages: 0,
		})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit
	search := c.QueryParam("search")
	status := c.QueryParam("status")
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "last_message_at"
	}
	sortOrder := c.QueryParam("sort_order")
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	// Build query
	baseQuery := `
		FROM customer_insights 
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argCount := 1

	// Add search filter
	if search != "" {
		argCount++
		baseQuery += ` AND (
			customer_name ILIKE $` + strconv.Itoa(argCount) + ` OR 
			customer_jid ILIKE $` + strconv.Itoa(argCount) + ` OR
			customer_phone ILIKE $` + strconv.Itoa(argCount) + ` OR
			last_message_summary ILIKE $` + strconv.Itoa(argCount) + `
		)`
		args = append(args, "%"+search+"%")
	}

	// Add status filter
	if status != "" && status != "all" {
		argCount++
		baseQuery += ` AND status = $` + strconv.Itoa(argCount)
		args = append(args, status)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := db.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to count customers",
		})
	}

	// Get customers with pagination
	// Validate sort column to prevent SQL injection
	validSortColumns := map[string]bool{
		"last_message_at": true,
		"created_at":      true,
		"message_count":   true,
		"customer_name":   true,
		"status":          true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "last_message_at"
	}

	selectQuery := `
		SELECT 
			id, tenant_id, customer_jid, customer_name, customer_phone,
			COALESCE(status, 'new') as status, sentiment, intent,
			product_interest::text, last_message_summary,
			message_count, last_message_at, first_message_at,
			needs_follow_up, tags::text, COALESCE(lead_score, 0) as lead_score,
			created_at, updated_at
		` + baseQuery + `
		ORDER BY ` + sortBy + ` ` + sortOrder + ` NULLS LAST
		LIMIT $` + strconv.Itoa(argCount+1) + ` OFFSET $` + strconv.Itoa(argCount+2)
	
	args = append(args, limit, offset)

	rows, err := db.DB.Query(selectQuery, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch customers: " + err.Error(),
		})
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		var cust Customer
		err := rows.Scan(
			&cust.ID, &cust.TenantID, &cust.CustomerJID, &cust.CustomerName,
			&cust.CustomerPhone, &cust.Status, &cust.Sentiment, &cust.Intent,
			&cust.ProductInterest, &cust.LastMessageSummary, &cust.MessageCount,
			&cust.LastMessageAt, &cust.FirstMessageAt, &cust.NeedsFollowUp,
			&cust.Tags, &cust.LeadScore, &cust.CreatedAt, &cust.UpdatedAt,
		)
		if err != nil {
			continue
		}
		customers = append(customers, cust)
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, CustomerListResponse{
		Customers:  customers,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetCustomerDetail returns customer details with chat history
func GetCustomerDetail(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Customer ID is required",
		})
	}

	// Get customer
	var cust Customer
	query := `
		SELECT 
			id, tenant_id, customer_jid, customer_name, customer_phone,
			COALESCE(status, 'new') as status, sentiment, intent,
			product_interest::text, last_message_summary,
			message_count, last_message_at, first_message_at,
			needs_follow_up, tags::text, COALESCE(lead_score, 0) as lead_score,
			created_at, updated_at
		FROM customer_insights
		WHERE id = $1 AND tenant_id = $2
	`
	
	err := db.DB.QueryRow(query, customerID, tenantID).Scan(
		&cust.ID, &cust.TenantID, &cust.CustomerJID, &cust.CustomerName,
		&cust.CustomerPhone, &cust.Status, &cust.Sentiment, &cust.Intent,
		&cust.ProductInterest, &cust.LastMessageSummary, &cust.MessageCount,
		&cust.LastMessageAt, &cust.FirstMessageAt, &cust.NeedsFollowUp,
		&cust.Tags, &cust.LeadScore, &cust.CreatedAt, &cust.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Customer not found",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch customer",
		})
	}

	// Get chat history
	messagesQuery := `
		SELECT 
			id, COALESCE(message_text, '') as message_text, 
			message_type, COALESCE(media_url, '') as media_url,
			is_from_me, to_timestamp(timestamp) as timestamp
		FROM whatsapp_messages
		WHERE tenant_id = $1 AND (sender_jid = $2 OR chat_jid = $2)
		ORDER BY timestamp DESC
		LIMIT 50
	`
	
	rows, err := db.DB.Query(messagesQuery, tenantID, cust.CustomerJID)
	if err != nil {
		// Return customer without messages if query fails
		return c.JSON(http.StatusOK, CustomerDetailResponse{
			Customer: cust,
			Messages: []CustomerMessage{},
		})
	}
	defer rows.Close()

	messages := []CustomerMessage{}
	for rows.Next() {
		var msg CustomerMessage
		err := rows.Scan(&msg.ID, &msg.MessageText, &msg.MessageType, &msg.MediaURL, &msg.IsFromMe, &msg.Timestamp)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return c.JSON(http.StatusOK, CustomerDetailResponse{
		Customer: cust,
		Messages: messages,
	})
}

// UpdateCustomer updates customer information
func UpdateCustomer(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Customer ID is required",
		})
	}

	// Parse request body
	var req struct {
		CustomerName  *string `json:"customer_name"`
		Status        *string `json:"status"`
		NeedsFollowUp *bool   `json:"needs_follow_up"`
		Tags          *string `json:"tags"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argCount := 0

	if req.CustomerName != nil {
		argCount++
		updates = append(updates, "customer_name = $"+strconv.Itoa(argCount))
		args = append(args, *req.CustomerName)
	}

	if req.Status != nil {
		// Validate status
		validStatuses := map[string]bool{
			"new": true, "hot_lead": true, "warm_lead": true,
			"cold_lead": true, "customer": true, "complaint": true, "spam": true,
		}
		if !validStatuses[*req.Status] {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid status value",
			})
		}
		argCount++
		updates = append(updates, "status = $"+strconv.Itoa(argCount))
		args = append(args, *req.Status)
	}

	if req.NeedsFollowUp != nil {
		argCount++
		updates = append(updates, "needs_follow_up = $"+strconv.Itoa(argCount))
		args = append(args, *req.NeedsFollowUp)
	}

	if req.Tags != nil {
		argCount++
		updates = append(updates, "tags = $"+strconv.Itoa(argCount)+"::jsonb")
		args = append(args, *req.Tags)
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No fields to update",
		})
	}

	// Add updated_at
	argCount++
	updates = append(updates, "updated_at = NOW()")

	// Build final query
	query := "UPDATE customer_insights SET "
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = $" + strconv.Itoa(argCount) + " AND tenant_id = $" + strconv.Itoa(argCount+1)
	args = append(args, customerID, tenantID)

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update customer",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Customer not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer updated successfully",
	})
}

// GetCustomerStats returns customer statistics
func GetCustomerStats(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty stats instead of error
		return c.JSON(http.StatusOK, map[string]interface{}{
			"total":         0,
			"new":           0,
			"hot_leads":     0,
			"warm_leads":    0,
			"cold_leads":    0,
			"customers":     0,
			"complaints":    0,
			"need_follow_up": 0,
		})
	}

	stats := struct {
		Total       int `json:"total"`
		New         int `json:"new"`
		HotLeads    int `json:"hot_leads"`
		WarmLeads   int `json:"warm_leads"`
		ColdLeads   int `json:"cold_leads"`
		Customers   int `json:"customers"`
		Complaints  int `json:"complaints"`
		NeedFollowUp int `json:"need_follow_up"`
	}{}

	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'new' OR status IS NULL) as new,
			COUNT(*) FILTER (WHERE status = 'hot_lead') as hot_leads,
			COUNT(*) FILTER (WHERE status = 'warm_lead') as warm_leads,
			COUNT(*) FILTER (WHERE status = 'cold_lead') as cold_leads,
			COUNT(*) FILTER (WHERE status = 'customer') as customers,
			COUNT(*) FILTER (WHERE status = 'complaint') as complaints,
			COUNT(*) FILTER (WHERE needs_follow_up = true) as need_follow_up
		FROM customer_insights
		WHERE tenant_id = $1
	`

	err := db.DB.QueryRow(query, tenantID).Scan(
		&stats.Total, &stats.New, &stats.HotLeads, &stats.WarmLeads,
		&stats.ColdLeads, &stats.Customers, &stats.Complaints, &stats.NeedFollowUp,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch customer stats",
		})
	}

	return c.JSON(http.StatusOK, stats)
}



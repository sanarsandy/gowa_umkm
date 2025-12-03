package handlers

import (
	"net/http"
	"time"

	"gowa-backend/db"

	"github.com/labstack/echo/v4"
)

// Tag represents a customer tag
type Tag struct {
	ID            string    `db:"id" json:"id"`
	TenantID      string    `db:"tenant_id" json:"tenant_id"`
	Name          string    `db:"name" json:"name"`
	Color         string    `db:"color" json:"color"`
	Description   string    `db:"description" json:"description"`
	CustomerCount int       `db:"customer_count" json:"customer_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// CustomerNote represents a note about a customer
type CustomerNote struct {
	ID         string    `db:"id" json:"id"`
	CustomerID string    `db:"customer_id" json:"customer_id"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// GetTags returns all tags for a tenant
// GET /api/tags
func GetTags(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var tags []Tag
	query := `
		SELECT 
			t.id, t.tenant_id, t.name, t.color, 
			COALESCE(t.description, '') as description,
			COUNT(cta.customer_id) as customer_count,
			t.created_at
		FROM customer_tags t
		LEFT JOIN customer_tag_assignments cta ON t.id = cta.tag_id
		WHERE t.tenant_id = $1
		GROUP BY t.id, t.tenant_id, t.name, t.color, t.description, t.created_at
		ORDER BY t.created_at DESC
	`

	err := db.DB.Select(&tags, query, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get tags")
	}

	return c.JSON(http.StatusOK, tags)
}

// CreateTag creates a new tag
// POST /api/tags
func CreateTag(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var req struct {
		Name        string `json:"name" validate:"required"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tag name is required")
	}

	if req.Color == "" {
		req.Color = "#6366f1"
	}

	var tag Tag
	query := `
		INSERT INTO customer_tags (tenant_id, name, color, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tenant_id, name, color, description, customer_count, created_at
	`

	err := db.DB.QueryRow(query, tenantID, req.Name, req.Color, req.Description).Scan(
		&tag.ID, &tag.TenantID, &tag.Name, &tag.Color, &tag.Description, &tag.CustomerCount, &tag.CreatedAt,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create tag: "+err.Error())
	}

	return c.JSON(http.StatusCreated, tag)
}

// UpdateTag updates a tag
// PUT /api/tags/:id
func UpdateTag(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	tagID := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	query := `
		UPDATE customer_tags 
		SET name = COALESCE(NULLIF($1, ''), name),
			color = COALESCE(NULLIF($2, ''), color),
			description = $3,
			updated_at = NOW()
		WHERE id = $4 AND tenant_id = $5
	`

	result, err := db.DB.Exec(query, req.Name, req.Color, req.Description, tagID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update tag")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Tag not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Tag updated"})
}

// DeleteTag deletes a tag
// DELETE /api/tags/:id
func DeleteTag(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	tagID := c.Param("id")

	result, err := db.DB.Exec(`DELETE FROM customer_tags WHERE id = $1 AND tenant_id = $2`, tagID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete tag")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Tag not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Tag deleted"})
}

// AssignTagToCustomer assigns a tag to a customer
// POST /api/customers/:id/tags
func AssignTagToCustomer(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")

	var req struct {
		TagID string `json:"tag_id" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Verify customer belongs to tenant
	var customerExists bool
	db.DB.Get(&customerExists, `SELECT EXISTS(SELECT 1 FROM customer_insights WHERE id = $1 AND tenant_id = $2)`, customerID, tenantID)
	if !customerExists {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	// Verify tag belongs to tenant
	var tagExists bool
	db.DB.Get(&tagExists, `SELECT EXISTS(SELECT 1 FROM customer_tags WHERE id = $1 AND tenant_id = $2)`, req.TagID, tenantID)
	if !tagExists {
		return echo.NewHTTPError(http.StatusNotFound, "Tag not found")
	}

	// Assign tag
	_, err := db.DB.Exec(`
		INSERT INTO customer_tag_assignments (customer_id, tag_id, assigned_by)
		VALUES ($1, $2, 'manual')
		ON CONFLICT DO NOTHING
	`, customerID, req.TagID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to assign tag")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Tag assigned"})
}

// RemoveTagFromCustomer removes a tag from a customer
// DELETE /api/customers/:id/tags/:tagId
func RemoveTagFromCustomer(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")
	tagID := c.Param("tagId")

	// Verify customer belongs to tenant
	var exists bool
	db.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM customer_insights WHERE id = $1 AND tenant_id = $2)`, customerID, tenantID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	_, err := db.DB.Exec(`DELETE FROM customer_tag_assignments WHERE customer_id = $1 AND tag_id = $2`, customerID, tagID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove tag")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Tag removed"})
}

// GetCustomerTags returns tags for a customer
// GET /api/customers/:id/tags
func GetCustomerTags(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")

	// Verify customer belongs to tenant
	var exists bool
	db.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM customer_insights WHERE id = $1 AND tenant_id = $2)`, customerID, tenantID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	var tags []Tag
	query := `
		SELECT t.id, t.tenant_id, t.name, t.color, COALESCE(t.description, '') as description, 0 as customer_count, t.created_at
		FROM customer_tags t
		JOIN customer_tag_assignments cta ON t.id = cta.tag_id
		WHERE cta.customer_id = $1
	`

	err := db.DB.Select(&tags, query, customerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get customer tags")
	}

	return c.JSON(http.StatusOK, tags)
}

// GetCustomerNotes returns notes for a customer
// GET /api/customers/:id/notes
func GetCustomerNotes(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")

	// Verify customer belongs to tenant
	var exists bool
	db.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM customer_insights WHERE id = $1 AND tenant_id = $2)`, customerID, tenantID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	var notes []CustomerNote
	query := `
		SELECT id, customer_id, content, created_at
		FROM customer_notes
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	err := db.DB.Select(&notes, query, customerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get customer notes")
	}

	return c.JSON(http.StatusOK, notes)
}

// CreateCustomerNote creates a note for a customer
// POST /api/customers/:id/notes
func CreateCustomerNote(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")

	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Content is required")
	}

	// Verify customer belongs to tenant
	var exists bool
	db.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM customer_insights WHERE id = $1 AND tenant_id = $2)`, customerID, tenantID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	var note CustomerNote
	query := `
		INSERT INTO customer_notes (customer_id, tenant_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, customer_id, content, created_at
	`

	err := db.DB.QueryRow(query, customerID, tenantID, req.Content).Scan(
		&note.ID, &note.CustomerID, &note.Content, &note.CreatedAt,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create note")
	}

	return c.JSON(http.StatusCreated, note)
}

// DeleteCustomerNote deletes a note
// DELETE /api/customers/:id/notes/:noteId
func DeleteCustomerNote(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")
	noteID := c.Param("noteId")

	result, err := db.DB.Exec(`
		DELETE FROM customer_notes 
		WHERE id = $1 AND customer_id = $2 AND tenant_id = $3
	`, noteID, customerID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete note")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Note not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Note deleted"})
}

// UpdateCustomerLeadScore updates a customer's lead score
// PUT /api/customers/:id/lead-score
func UpdateCustomerLeadScore(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	customerID := c.Param("id")

	var req struct {
		LeadScore int    `json:"lead_score"`
		Status    string `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate lead score
	if req.LeadScore < 0 || req.LeadScore > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "Lead score must be between 0 and 100")
	}

	// Validate status
	validStatuses := map[string]bool{
		"new": true, "hot_lead": true, "warm_lead": true,
		"cold_lead": true, "customer": true, "complaint": true, "spam": true,
	}
	if req.Status != "" && !validStatuses[req.Status] {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid status")
	}

	query := `
		UPDATE customer_insights 
		SET lead_score = $1, status = COALESCE(NULLIF($2, ''), status), updated_at = NOW()
		WHERE id = $3 AND tenant_id = $4
	`

	result, err := db.DB.Exec(query, req.LeadScore, req.Status, customerID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update lead score")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Lead score updated"})
}

// GetCustomersByTag returns customers with a specific tag
// GET /api/tags/:id/customers
func GetCustomersByTag(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	tagID := c.Param("id")

	// Verify tag belongs to tenant
	var exists bool
	db.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM customer_tags WHERE id = $1 AND tenant_id = $2)`, tagID, tenantID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Tag not found")
	}

	type CustomerWithTags struct {
		ID          string    `db:"id" json:"id"`
		PhoneNumber string    `db:"phone_number" json:"phone_number"`
		Name        string    `db:"name" json:"name"`
		LeadScore   int       `db:"lead_score" json:"lead_score"`
		Status      string    `db:"status" json:"lead_status"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
	}

	var customers []CustomerWithTags
	query := `
		SELECT c.id, 
			   COALESCE(c.customer_phone, c.customer_jid) as phone_number, 
			   COALESCE(c.customer_name, c.customer_phone, c.customer_jid) as name, 
			   COALESCE(c.lead_score, 0) as lead_score, 
			   COALESCE(c.status, 'new') as status, 
			   c.created_at
		FROM customer_insights c
		JOIN customer_tag_assignments cta ON c.id = cta.customer_id
		WHERE cta.tag_id = $1 AND c.tenant_id = $2
		ORDER BY c.created_at DESC
	`

	err := db.DB.Select(&customers, query, tagID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get customers")
	}

	return c.JSON(http.StatusOK, customers)
}


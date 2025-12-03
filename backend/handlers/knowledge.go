package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"gowa-backend/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// Knowledge represents a knowledge base entry
type Knowledge struct {
	ID          string         `json:"id" db:"id"`
	TenantID    string         `json:"tenant_id" db:"tenant_id"`
	Title       string         `json:"title" db:"title"`
	Content     string         `json:"content" db:"content"`
	Category    string         `json:"category" db:"category"`
	Keywords    pq.StringArray `json:"keywords" db:"keywords"`
	Tags        pq.StringArray `json:"tags" db:"tags"`
	Priority    int            `json:"priority" db:"priority"`
	UsageCount  int            `json:"usage_count" db:"usage_count"`
	LastUsedAt  *time.Time     `json:"last_used_at" db:"last_used_at"`
	IsActive    bool           `json:"is_active" db:"is_active"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// GetKnowledgeBase retrieves all knowledge base entries for a tenant
// GET /api/knowledge
func GetKnowledgeBase(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	category := c.QueryParam("category")
	search := c.QueryParam("search")
	activeOnly := c.QueryParam("active") == "true"

	query := `
		SELECT id, tenant_id, title, content, category, keywords, tags, 
		       priority, usage_count, last_used_at, is_active, created_at, updated_at
		FROM knowledge_base
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argCount := 1

	if category != "" {
		argCount++
		query += ` AND category = $` + string(rune(argCount+'0'))
		args = append(args, category)
	}

	if search != "" {
		argCount++
		query += ` AND (title ILIKE $` + string(rune(argCount+'0')) + ` OR content ILIKE $` + string(rune(argCount+'0')) + `)`
		args = append(args, "%"+search+"%")
	}

	if activeOnly {
		query += ` AND is_active = true`
	}

	query += ` ORDER BY priority DESC, created_at DESC`

	var knowledge []Knowledge
	err := db.DB.Select(&knowledge, query, args...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch knowledge base")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"knowledge": knowledge,
		"total":     len(knowledge),
	})
}

// CreateKnowledge creates a new knowledge base entry
// POST /api/knowledge
func CreateKnowledge(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var req struct {
		Title    string   `json:"title" validate:"required"`
		Content  string   `json:"content" validate:"required"`
		Category string   `json:"category"`
		Keywords []string `json:"keywords"`
		Tags     []string `json:"tags"`
		Priority int      `json:"priority"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Title == "" || req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title and content are required")
	}

	if req.Priority == 0 {
		req.Priority = 5 // default priority
	}

	id := uuid.New().String()

	query := `
		INSERT INTO knowledge_base (id, tenant_id, title, content, category, keywords, tags, priority, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
		RETURNING id, tenant_id, title, content, category, keywords, tags, priority, usage_count, is_active, created_at, updated_at
	`

	var knowledge Knowledge
	err := db.DB.Get(&knowledge, query, id, tenantID, req.Title, req.Content, req.Category,
		pq.Array(req.Keywords), pq.Array(req.Tags), req.Priority)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create knowledge: "+err.Error())
	}

	return c.JSON(http.StatusCreated, knowledge)
}

// UpdateKnowledge updates an existing knowledge base entry
// PUT /api/knowledge/:id
func UpdateKnowledge(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	knowledgeID := c.Param("id")

	var req struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Category string   `json:"category"`
		Keywords []string `json:"keywords"`
		Tags     []string `json:"tags"`
		Priority int      `json:"priority"`
		IsActive *bool    `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	query := `
		UPDATE knowledge_base
		SET title = COALESCE(NULLIF($1, ''), title),
		    content = COALESCE(NULLIF($2, ''), content),
		    category = COALESCE(NULLIF($3, ''), category),
		    keywords = COALESCE($4, keywords),
		    tags = COALESCE($5, tags),
		    priority = COALESCE(NULLIF($6, 0), priority),
		    is_active = COALESCE($7, is_active),
		    updated_at = NOW()
		WHERE id = $8 AND tenant_id = $9
		RETURNING id, tenant_id, title, content, category, keywords, tags, priority, usage_count, is_active, created_at, updated_at
	`

	var knowledge Knowledge
	err := db.DB.Get(&knowledge, query, req.Title, req.Content, req.Category,
		pq.Array(req.Keywords), pq.Array(req.Tags), req.Priority, req.IsActive, knowledgeID, tenantID)

	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Knowledge not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update knowledge")
	}

	return c.JSON(http.StatusOK, knowledge)
}

// DeleteKnowledge deletes a knowledge base entry
// DELETE /api/knowledge/:id
func DeleteKnowledge(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	knowledgeID := c.Param("id")

	result, err := db.DB.Exec(`
		DELETE FROM knowledge_base
		WHERE id = $1 AND tenant_id = $2
	`, knowledgeID, tenantID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete knowledge")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Knowledge not found")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Knowledge deleted successfully",
	})
}

// GetKnowledgeStats returns statistics about knowledge base
// GET /api/knowledge/stats
func GetKnowledgeStats(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var stats struct {
		Total          int `db:"total"`
		Active         int `db:"active"`
		ByCategory     map[string]int
		MostUsed       []Knowledge
		RecentlyAdded  []Knowledge
	}

	// Get total and active count
	err := db.DB.Get(&stats, `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_active = true) as active
		FROM knowledge_base
		WHERE tenant_id = $1
	`, tenantID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stats")
	}

	// Get most used knowledge
	err = db.DB.Select(&stats.MostUsed, `
		SELECT id, title, category, usage_count, last_used_at
		FROM knowledge_base
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY usage_count DESC
		LIMIT 5
	`, tenantID)

	if err != nil {
		stats.MostUsed = []Knowledge{}
	}

	// Get recently added
	err = db.DB.Select(&stats.RecentlyAdded, `
		SELECT id, title, category, created_at
		FROM knowledge_base
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT 5
	`, tenantID)

	if err != nil {
		stats.RecentlyAdded = []Knowledge{}
	}

	return c.JSON(http.StatusOK, stats)
}

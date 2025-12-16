package handlers

import (
	"fmt"
	"net/http"
	"time"

	"gowa-backend/db"

	"github.com/labstack/echo/v4"
)

// Template represents a message template
type Template struct {
	ID         string    `json:"id" db:"id"`
	TenantID   string    `json:"tenant_id" db:"tenant_id"`
	Name       string    `json:"name" db:"name"`
	Category   string    `json:"category" db:"category"`
	Content    string    `json:"content" db:"content"`
	Variables  string    `json:"variables" db:"variables"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	UsageCount int       `json:"usage_count" db:"usage_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// GetTemplates returns all templates for a tenant
func GetTemplates(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		// Tenant not found - return empty list instead of 401
		// User is authenticated but hasn't created a tenant yet
		return c.JSON(http.StatusOK, map[string]interface{}{
			"templates": []Template{},
			"total":     0,
		})
	}

	category := c.QueryParam("category")

	var templates []Template
	var query string
	var args []interface{}

	if category != "" && category != "all" {
		query = `
			SELECT id, tenant_id, name, category, content, variables, is_active, usage_count, created_at, updated_at
			FROM message_templates
			WHERE tenant_id = $1 AND category = $2
			ORDER BY usage_count DESC, name ASC
		`
		args = []interface{}{tenantID, category}
	} else {
		query = `
			SELECT id, tenant_id, name, category, content, variables, is_active, usage_count, created_at, updated_at
			FROM message_templates
			WHERE tenant_id = $1
			ORDER BY usage_count DESC, name ASC
		`
		args = []interface{}{tenantID}
	}

	if err := db.DB.Select(&templates, query, args...); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch templates")
	}

	if templates == nil {
		templates = []Template{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
	})
}

// CreateTemplate creates a new message template
func CreateTemplate(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	var req struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Content   string   `json:"content"`
		Variables []string `json:"variables"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" || req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name and content are required")
	}

	if req.Category == "" {
		req.Category = "general"
	}

	// Convert variables to JSON string
	variablesJSON := "[]"
	if len(req.Variables) > 0 {
		variablesJSON = `["` + joinStrings(req.Variables, `","`) + `"]`
	}

	var template Template
	query := `
		INSERT INTO message_templates (tenant_id, name, category, content, variables)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, tenant_id, name, category, content, variables, is_active, usage_count, created_at, updated_at
	`

	if err := db.DB.Get(&template, query, tenantID, req.Name, req.Category, req.Content, variablesJSON); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create template")
	}

	return c.JSON(http.StatusCreated, template)
}

// UpdateTemplate updates an existing template
func UpdateTemplate(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	templateID := c.Param("id")

	var req struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Content   string   `json:"content"`
		Variables []string `json:"variables"`
		IsActive  *bool    `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Build update query dynamically
	query := `UPDATE message_templates SET updated_at = NOW()`
	args := []interface{}{}
	argNum := 1

	if req.Name != "" {
		query += `, name = $` + itoa(argNum)
		args = append(args, req.Name)
		argNum++
	}
	if req.Category != "" {
		query += `, category = $` + itoa(argNum)
		args = append(args, req.Category)
		argNum++
	}
	if req.Content != "" {
		query += `, content = $` + itoa(argNum)
		args = append(args, req.Content)
		argNum++
	}
	if len(req.Variables) > 0 {
		variablesJSON := `["` + joinStrings(req.Variables, `","`) + `"]`
		query += `, variables = $` + itoa(argNum)
		args = append(args, variablesJSON)
		argNum++
	}
	if req.IsActive != nil {
		query += `, is_active = $` + itoa(argNum)
		args = append(args, *req.IsActive)
		argNum++
	}

	query += ` WHERE id = $` + itoa(argNum) + ` AND tenant_id = $` + itoa(argNum+1)
	args = append(args, templateID, tenantID)

	query += ` RETURNING id, tenant_id, name, category, content, variables, is_active, usage_count, created_at, updated_at`

	var template Template
	if err := db.DB.Get(&template, query, args...); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Template not found")
	}

	return c.JSON(http.StatusOK, template)
}

// DeleteTemplate deletes a template
func DeleteTemplate(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	templateID := c.Param("id")

	query := `DELETE FROM message_templates WHERE id = $1 AND tenant_id = $2`
	result, err := db.DB.Exec(query, templateID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete template")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Template not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Template deleted"})
}

// IncrementTemplateUsage increments usage count when a template is used
func IncrementTemplateUsage(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant not found. Please create a tenant first.")
	}

	templateID := c.Param("id")

	query := `UPDATE message_templates SET usage_count = usage_count + 1, updated_at = NOW() WHERE id = $1 AND tenant_id = $2`
	if _, err := db.DB.Exec(query, templateID, tenantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update usage count")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Usage count updated"})
}

// Helper functions
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}


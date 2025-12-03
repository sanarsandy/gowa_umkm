package handlers

import (
	"database/sql"
	"net/http"

	"gowa-backend/db"
	"gowa-backend/models"

	"github.com/labstack/echo/v4"
)

// CreateTenant creates a new tenant for the authenticated user
func CreateTenant(c echo.Context) error {
	var req struct {
		BusinessName        string `json:"business_name"`
		BusinessType        string `json:"business_type"`
		BusinessDescription string `json:"business_description"`
		BusinessPhone       string `json:"business_phone"`
		BusinessAddress     string `json:"business_address"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Get user ID from JWT (for now, use a placeholder)
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Create tenant
	query := `
		INSERT INTO tenants (user_id, business_name, business_type, business_description, business_phone, business_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var tenant models.Tenant
	err := db.DB.QueryRow(query,
		userID,
		req.BusinessName,
		req.BusinessType,
		req.BusinessDescription,
		req.BusinessPhone,
		req.BusinessAddress,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create tenant",
		})
	}

	tenant.UserID = userID
	tenant.BusinessName = req.BusinessName
	tenant.BusinessType = req.BusinessType
	tenant.BusinessDescription = req.BusinessDescription
	tenant.BusinessPhone = req.BusinessPhone
	tenant.BusinessAddress = req.BusinessAddress
	tenant.IsActive = true

	return c.JSON(http.StatusCreated, tenant)
}

// GetMyTenant returns the tenant for the authenticated user
func GetMyTenant(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var tenant models.Tenant
	query := `SELECT * FROM tenants WHERE user_id = $1 LIMIT 1`
	
	err := db.DB.QueryRow(query, userID).Scan(
		&tenant.ID,
		&tenant.UserID,
		&tenant.BusinessName,
		&tenant.BusinessType,
		&tenant.BusinessDescription,
		&tenant.BusinessPhone,
		&tenant.BusinessAddress,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Tenant not found",
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get tenant",
		})
	}

	return c.JSON(http.StatusOK, tenant)
}

// UpdateTenant updates the tenant information for the authenticated user
func UpdateTenant(c echo.Context) error {
	var req struct {
		BusinessName        string `json:"business_name"`
		BusinessType        string `json:"business_type"`
		BusinessDescription string `json:"business_description"`
		BusinessPhone       string `json:"business_phone"`
		BusinessAddress     string `json:"business_address"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Request tidak valid",
		})
	}

	// Get user ID from JWT
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Check if tenant exists
	var existingTenantID string
	checkQuery := `SELECT id FROM tenants WHERE user_id = $1 LIMIT 1`
	err := db.DB.QueryRow(checkQuery, userID).Scan(&existingTenantID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Tenant tidak ditemukan",
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal memeriksa tenant",
		})
	}

	// Update tenant
	updateQuery := `
		UPDATE tenants 
		SET business_name = $1, 
		    business_type = $2, 
		    business_description = $3, 
		    business_phone = $4, 
		    business_address = $5,
		    updated_at = NOW()
		WHERE id = $6
		RETURNING id, user_id, business_name, business_type, business_description, business_phone, business_address, is_active, created_at, updated_at
	`

	var tenant models.Tenant
	err = db.DB.QueryRow(updateQuery,
		req.BusinessName,
		req.BusinessType,
		req.BusinessDescription,
		req.BusinessPhone,
		req.BusinessAddress,
		existingTenantID,
	).Scan(
		&tenant.ID,
		&tenant.UserID,
		&tenant.BusinessName,
		&tenant.BusinessType,
		&tenant.BusinessDescription,
		&tenant.BusinessPhone,
		&tenant.BusinessAddress,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengupdate tenant",
		})
	}

	return c.JSON(http.StatusOK, tenant)
}

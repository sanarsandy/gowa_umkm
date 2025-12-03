package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"gowa-backend/db"
	"gowa-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	req := new(models.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Request tidak valid"})
	}

	// Check if email already exists
	var existingID string
	checkQuery := `SELECT id FROM users WHERE email = $1 AND (auth_provider = 'email' OR auth_provider = 'both')`
	err := db.DB.QueryRow(checkQuery, req.Email).Scan(&existingID)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email sudah terdaftar"})
	} else if err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Terjadi kesalahan pada database"})
	}

	// Validate password length
	if len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password minimal 6 karakter"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengenkripsi password"})
	}

	// Insert into DB
	query := `INSERT INTO users (email, password_hash, full_name, role, auth_provider) VALUES ($1, $2, $3, 'parent', 'email') RETURNING id, created_at`
	var user models.User
	user.Email = req.Email
	user.FullName = req.FullName
	user.Role = "parent"
	user.AuthProvider = "email"

	err = db.DB.QueryRow(query, req.Email, string(hashedPassword), req.FullName).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat akun: " + err.Error()})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Konfigurasi server tidak valid"})
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat token"})
	}

	// Auto-create tenant dengan data default
	tenantQuery := `INSERT INTO tenants (user_id, business_name, business_type, business_description, business_phone, business_address, is_active)
	                VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	var tenantID string
	var tenantCreatedAt, tenantUpdatedAt time.Time
	businessName := req.FullName + "'s Business"
	err = db.DB.QueryRow(tenantQuery, user.ID, businessName, "UMKM", "", "", "", true).
		Scan(&tenantID, &tenantCreatedAt, &tenantUpdatedAt)
	if err != nil {
		// Log error but don't fail registration - tenant can be created later
		// In production, you might want to handle this differently
		c.Logger().Warnf("Failed to auto-create tenant for user %s: %v", user.ID, err)
	}

	// Set cookies for frontend
	// Token cookie is HttpOnly for security (prevents XSS attacks)
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour) // 3 days
	cookie.Path = "/"
	cookie.HttpOnly = true // SECURITY: Prevent JavaScript access to prevent XSS
	cookie.Secure = os.Getenv("ENV") == "production" // HTTPS only in production
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	// Set user cookie (as JSON)
	userJSON, _ := json.Marshal(user)
	userCookie := new(http.Cookie)
	userCookie.Name = "user"
	userCookie.Value = string(userJSON)
	userCookie.Expires = time.Now().Add(72 * time.Hour)
	userCookie.Path = "/"
	userCookie.HttpOnly = false
	userCookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(userCookie)

	return c.JSON(http.StatusCreated, models.AuthResponse{
		Token: t,
		User:  user,
	})
}

func Login(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Request tidak valid"})
	}

	// Find user - only look for email/password users
	var user models.User
	var googleID sql.NullString
	var authProvider sql.NullString
	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider 
	          FROM users 
	          WHERE email = $1 AND (auth_provider = 'email' OR auth_provider = 'both')`
	err := db.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &googleID, &authProvider)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User belum terdaftar"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Terjadi kesalahan pada database"})
	}

	// Set optional fields
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "email"
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
	}

	// Generate JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Konfigurasi server tidak valid"})
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat token"})
	}

	// Check if user has a tenant, create one if not
	var tenantID string
	tenantCheckQuery := `SELECT id FROM tenants WHERE user_id = $1 AND is_active = true LIMIT 1`
	err = db.DB.QueryRow(tenantCheckQuery, user.ID).Scan(&tenantID)
	if err == sql.ErrNoRows {
		// User doesn't have a tenant, create one
		tenantQuery := `INSERT INTO tenants (user_id, business_name, business_type, business_description, business_phone, business_address, is_active)
		                VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		businessName := user.FullName + "'s Business"
		err = db.DB.QueryRow(tenantQuery, user.ID, businessName, "UMKM", "", "", "", true).Scan(&tenantID)
		if err != nil {
			// Log error but don't fail login
			c.Logger().Warnf("Failed to auto-create tenant for user %s on login: %v", user.ID, err)
		} else {
			c.Logger().Infof("Auto-created tenant %s for user %s on login", tenantID, user.ID)
		}
	} else if err != nil {
		// Database error, log but don't fail login
		c.Logger().Warnf("Failed to check tenant for user %s: %v", user.ID, err)
	}

	// Set cookies for frontend
	// Token cookie is HttpOnly for security (prevents XSS attacks)
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour) // 3 days
	cookie.Path = "/"
	cookie.HttpOnly = true // SECURITY: Prevent JavaScript access to prevent XSS
	cookie.Secure = os.Getenv("ENV") == "production" // HTTPS only in production
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	// Set user cookie (as JSON)
	userJSON, _ := json.Marshal(user)
	userCookie := new(http.Cookie)
	userCookie.Name = "user"
	userCookie.Value = string(userJSON)
	userCookie.Expires = time.Now().Add(72 * time.Hour)
	userCookie.Path = "/"
	userCookie.HttpOnly = false
	userCookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(userCookie)

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: t,
		User:  user,
	})
}

// GetMe returns the current authenticated user information from JWT
func GetMe(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Get user from database
	var user models.User
	var googleID sql.NullString
	var authProvider sql.NullString
	query := `SELECT id, email, full_name, role, google_id, auth_provider, created_at 
	          FROM users WHERE id = $1`
	err := db.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Role, 
		&googleID, &authProvider, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user",
		})
	}

	// Set optional fields
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "email"
	}

	return c.JSON(http.StatusOK, user)
}

// getUserIDFromContext extracts user ID from JWT claims in the context
func getUserIDFromContext(c echo.Context) string {
	// Get token from context (set by JWT middleware)
	user := c.Get("user")
	if user == nil {
		fmt.Printf("[DEBUG] getUserIDFromContext: user is nil in context\n")
		return ""
	}

	// Type assert to *jwt.Token
	token, ok := user.(*jwt.Token)
	if !ok {
		fmt.Printf("[DEBUG] getUserIDFromContext: failed to assert user to *jwt.Token, type: %T\n", user)
		return ""
	}

	// Log the actual type of claims
	fmt.Printf("[DEBUG] getUserIDFromContext: token.Claims type: %T\n", token.Claims)
	fmt.Printf("[DEBUG] getUserIDFromContext: token.Claims value: %+v\n", token.Claims)

	// Try to get claims - jwt.MapClaims is the expected type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("[DEBUG] getUserIDFromContext: failed to assert claims to jwt.MapClaims, actual type: %T\n", token.Claims)
		// Try using JSON marshal/unmarshal as fallback
		claimsJSON, err := json.Marshal(token.Claims)
		if err != nil {
			fmt.Printf("[DEBUG] getUserIDFromContext: failed to marshal claims: %v\n", err)
			return ""
		}
		var claimsMap map[string]interface{}
		if err := json.Unmarshal(claimsJSON, &claimsMap); err != nil {
			fmt.Printf("[DEBUG] getUserIDFromContext: failed to unmarshal claims: %v\n", err)
			return ""
		}
		// Extract from map
		userIDVal, exists := claimsMap["user_id"]
		if !exists {
			fmt.Printf("[DEBUG] getUserIDFromContext: user_id not found in claims. Claims keys: %v\n", getMapKeys(claimsMap))
			return ""
		}
		userID, ok := userIDVal.(string)
		if !ok {
			fmt.Printf("[DEBUG] getUserIDFromContext: user_id is not a string, type: %T, value: %v\n", userIDVal, userIDVal)
			return ""
		}
		fmt.Printf("[DEBUG] getUserIDFromContext: extracted user_id = '%s' (via JSON fallback)\n", userID)
		return userID
	}

	// Extract user_id from jwt.MapClaims
	userIDVal, exists := claims["user_id"]
	if !exists {
		fmt.Printf("[DEBUG] getUserIDFromContext: user_id not found in claims. Claims keys: %v\n", getMapKeys(claims))
		return ""
	}

	userID, ok := userIDVal.(string)
	if !ok {
		fmt.Printf("[DEBUG] getUserIDFromContext: user_id is not a string, type: %T, value: %v\n", userIDVal, userIDVal)
		return ""
	}

	fmt.Printf("[DEBUG] getUserIDFromContext: extracted user_id = '%s'\n", userID)
	return userID
}

// Helper function to get map keys for logging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}


package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"gowa-backend/db"
	"gowa-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var googleOAuthConfig *oauth2.Config

func init() {
	// Initialize Google OAuth config
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	
	if clientID == "" || clientSecret == "" {
		return
	}
	
	if redirectURL == "" {
		redirectURL = "http://localhost:3000/auth/google/callback"
	}

	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

func GetGoogleAuthURL(c echo.Context) error {
	if googleOAuthConfig == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Google OAuth is not configured",
		})
	}

	state := generateStateToken()
	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	
	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": url,
		"state":    state,
	})
}

func GoogleAuthCallback(c echo.Context) error {
	if googleOAuthConfig == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Google OAuth is not configured",
		})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization code not provided",
		})
	}

	// Exchange code for token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to exchange token: " + err.Error(),
		})
	}

	// Get user info from Google
	userInfo, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info: " + err.Error(),
		})
	}

	// Find or create user
	user, err := findOrCreateGoogleUser(userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create/find user: " + err.Error(),
		})
	}

	// Auto-create tenant if it doesn't exist (similar to email registration)
	// This is mandatory - login will fail if tenant creation fails
	var existingTenantID string
	tenantCheckQuery := `SELECT id FROM tenants WHERE user_id = $1 AND is_active = true LIMIT 1`
	err = db.DB.QueryRow(tenantCheckQuery, user.ID).Scan(&existingTenantID)
	if err == sql.ErrNoRows {
		// Tenant doesn't exist, create one (mandatory)
		businessName := user.FullName + "'s Business"
		if businessName == "'s Business" {
			businessName = user.Email + "'s Business"
		}
		tenantQuery := `INSERT INTO tenants (user_id, business_name, business_type, business_description, business_phone, business_address, is_active)
		                VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		var tenantID string
		err = db.DB.QueryRow(tenantQuery, user.ID, businessName, "UMKM", "", "", "", true).Scan(&tenantID)
		if err != nil {
			// Fail login if tenant creation fails
			c.Logger().Errorf("Failed to auto-create tenant for Google user %s: %v", user.ID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Gagal membuat tenant. Silakan coba lagi nanti atau hubungi administrator.",
			})
		}
		c.Logger().Infof("Auto-created tenant %s for Google user %s", tenantID, user.ID)
	} else if err != nil {
		// Database error while checking tenant - fail login
		c.Logger().Errorf("Error checking tenant for Google user %s: %v", user.ID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Terjadi kesalahan pada database. Silakan coba lagi nanti.",
		})
	}

	// Generate JWT
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	t, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	// Redirect to frontend
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	
	userJSON, _ := json.Marshal(user)
	redirectHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><title>Redirecting...</title></head>
		<body>
			<script>
				const token = '%s';
				const userData = %s;
				window.location.href = '%s/auth/google/callback?token=' + encodeURIComponent(token) + '&user=' + encodeURIComponent(JSON.stringify(userData));
			</script>
			<p>Redirecting...</p>
		</body>
		</html>
	`, t, userJSON, frontendURL)

	return c.HTML(http.StatusOK, redirectHTML)
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func findOrCreateGoogleUser(googleUser *GoogleUserInfo) (*models.User, error) {
	var user models.User
	var passwordHash sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	
	// Try to find by Google ID
	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, created_at 
	          FROM users WHERE google_id = $1`
	err := db.DB.QueryRow(query, googleUser.ID).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, 
		&user.Role, &googleID, &authProvider, &user.CreatedAt)
	
	if err == nil {
		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if authProvider.Valid {
			user.AuthProvider = authProvider.String
		} else {
			user.AuthProvider = "google"
		}
		return &user, nil
	}
	
	if err != sql.ErrNoRows {
		return nil, err
	}
	
	// Try to find by email
	query = `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, created_at 
	         FROM users WHERE email = $1`
	err = db.DB.QueryRow(query, googleUser.Email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, 
		&user.Role, &googleID, &authProvider, &user.CreatedAt)
	
	if err == nil {
		if authProvider.Valid && authProvider.String != "google" {
			return nil, fmt.Errorf("Email sudah terdaftar dengan akun email/password")
		}
		return &user, nil
	}
	
	if err != sql.ErrNoRows {
		return nil, err
	}
	
	// Create new user
	fullName := googleUser.Name
	if fullName == "" {
		fullName = googleUser.Email
	}
	
	insertQuery := `INSERT INTO users (email, full_name, google_id, auth_provider, role) 
	                VALUES ($1, $2, $3, 'google', 'parent') 
	                RETURNING id, email, full_name, role, auth_provider, created_at`
	err = db.DB.QueryRow(insertQuery, googleUser.Email, fullName, googleUser.ID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Role, &authProvider, &user.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	user.GoogleID = &googleUser.ID
	user.AuthProvider = "google"
	
	return &user, nil
}

func generateStateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}


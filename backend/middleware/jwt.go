package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// CRITICAL: JWT_SECRET must be set for security
		// Application should fail to start if not configured
		panic("FATAL: JWT_SECRET environment variable is not set. This is required for security.")
	}

	// Validate minimum length for security
	if len(jwtSecret) < 32 {
		panic("FATAL: JWT_SECRET must be at least 32 characters long for security.")
	}

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
		SigningKey: []byte(jwtSecret),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized: " + err.Error(),
			})
		},
	}
	return echojwt.WithConfig(config)
}


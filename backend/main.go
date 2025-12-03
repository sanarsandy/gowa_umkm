package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"gowa-backend/db"
	"gowa-backend/handlers"
	customMiddleware "gowa-backend/middleware"
	"gowa-backend/services/ai"
	"gowa-backend/services/scheduler"
	"gowa-backend/workers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Global AI service instance
var globalAIService *ai.AIService

func main() {
	// Validate critical environment variables before starting
	validateEnvironment()

	// Initialize Database
	db.Init()

	// Run database migrations
	if err := db.RunMigrations(db.DB.DB, "./migrations"); err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	// Initialize WhatsApp Service (includes Redis)
	handlers.InitWhatsAppService()

	// Initialize AI Service
	var err error
	globalAIService, err = ai.NewAIService()
	if err != nil {
		// Log error but don't fail - AI is optional
		println("Warning: Failed to initialize AI service:", err.Error())
		println("AI auto-reply features will be disabled")
	}

	// Start message worker in background
	if handlers.GetRedisClient() != nil {
		worker := workers.NewMessageWorker(handlers.GetRedisClient(), db.DB, handlers.GetWhatsAppService())
		ctx := context.Background()
		go worker.Start(ctx)
		
		// Start broadcast scheduler
		broadcastScheduler := scheduler.NewBroadcastScheduler(db.DB, handlers.GetRedisClient())
		go broadcastScheduler.Start()
		log.Println("✅ Broadcast scheduler started")
	}

	e := EchoServer()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	log.Println("Server stopped gracefully")
}

func EchoServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// CORS Configuration
	// Use environment variable for allowed origins in production
	corsOrigins := []string{"http://localhost:3000"}
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		// Split by comma and replace default
		corsOrigins = strings.Split(corsEnv, ",")
		for i, origin := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(origin)
		}
	}
	
	// Log CORS origins for debugging
	e.Logger.Infof("CORS allowed origins: %v", corsOrigins)
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Gowa UMKM WhatsApp API",
			"status":  "healthy",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "up",
		})
	})

	// Auth Routes (with stricter rate limiting)
	auth := e.Group("/api/auth")
	auth.Use(customMiddleware.AuthRateLimiterMiddleware())
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.GET("/google", handlers.GetGoogleAuthURL)
	auth.GET("/google/callback", handlers.GoogleAuthCallback)

	// WebSocket Route (handles its own auth via query param token)
	e.GET("/api/ws", handlers.HandleWebSocket)

	// Protected Routes (with rate limiting)
	api := e.Group("/api")
	api.Use(customMiddleware.APIRateLimiterMiddleware())
	api.Use(customMiddleware.JWTMiddleware())
	api.GET("/me", handlers.GetMe)

	// Tenant Routes
	api.POST("/tenant", handlers.CreateTenant)
	api.GET("/tenant", handlers.GetMyTenant)
	api.PUT("/tenant", handlers.UpdateTenant)

	// WhatsApp Routes
	whatsapp := api.Group("/whatsapp")
	whatsapp.POST("/connect", handlers.ConnectWhatsApp)
	whatsapp.DELETE("/disconnect", handlers.DisconnectWhatsApp)
	whatsapp.GET("/status", handlers.GetWhatsAppStatus)
	whatsapp.GET("/qr/stream", handlers.StreamQRCode)
	whatsapp.POST("/send", handlers.SendWhatsAppMessage)
	whatsapp.DELETE("/messages/:jid", handlers.ClearChatMessages)

	// Dashboard Routes
	dashboard := api.Group("/dashboard")
	dashboard.GET("/stats", handlers.GetDashboardStats)
	dashboard.GET("/messages/recent", handlers.GetRecentMessages)
	dashboard.GET("/customers/recent", handlers.GetRecentCustomers)

	// Customer Routes
	customers := api.Group("/customers")
	customers.GET("", handlers.GetCustomers)
	customers.GET("/stats", handlers.GetCustomerStats)
	customers.GET("/:id", handlers.GetCustomerDetail)
	customers.PUT("/:id", handlers.UpdateCustomer)
	customers.GET("/:id/tags", handlers.GetCustomerTags)
	customers.POST("/:id/tags", handlers.AssignTagToCustomer)
	customers.DELETE("/:id/tags/:tagId", handlers.RemoveTagFromCustomer)
	customers.GET("/:id/notes", handlers.GetCustomerNotes)
	customers.POST("/:id/notes", handlers.CreateCustomerNote)
	customers.DELETE("/:id/notes/:noteId", handlers.DeleteCustomerNote)
	customers.PUT("/:id/lead-score", handlers.UpdateCustomerLeadScore)

	// Template Routes
	templates := api.Group("/templates")
	templates.GET("", handlers.GetTemplates)
	templates.POST("", handlers.CreateTemplate)
	templates.PUT("/:id", handlers.UpdateTemplate)
	templates.DELETE("/:id", handlers.DeleteTemplate)
	templates.POST("/:id/use", handlers.IncrementTemplateUsage)

	// Broadcast Routes
	broadcasts := api.Group("/broadcasts")
	broadcasts.GET("", handlers.GetBroadcasts)
	broadcasts.GET("/stats", handlers.GetBroadcastStats)
	broadcasts.POST("", handlers.CreateBroadcast)
	broadcasts.GET("/:id", handlers.GetBroadcast)
	broadcasts.POST("/:id/send", handlers.SendBroadcast)
	broadcasts.POST("/:id/cancel", handlers.CancelBroadcast)
	broadcasts.DELETE("/:id", handlers.DeleteBroadcast)

	// AI Routes - Config routes are always available
	aiHandler := handlers.NewAIHandler(globalAIService)
	aiRoutes := api.Group("/ai")
	aiRoutes.GET("/config", aiHandler.GetAIConfig)
	aiRoutes.PUT("/config", aiHandler.UpdateAIConfig)
	aiRoutes.GET("/stats", aiHandler.GetAIStats)
	aiRoutes.GET("/providers", aiHandler.GetProviders)
	aiRoutes.GET("/providers/:provider/models", aiHandler.GetProviderModels)
	aiRoutes.POST("/test-connection", aiHandler.TestConnection)
	aiRoutes.POST("/test", aiHandler.TestAIResponse)
	
	// Knowledge Base Routes - always available
	knowledge := api.Group("/knowledge")
	knowledge.GET("", handlers.GetKnowledgeBase)
	knowledge.POST("", handlers.CreateKnowledge)
	knowledge.PUT("/:id", handlers.UpdateKnowledge)
	knowledge.DELETE("/:id", handlers.DeleteKnowledge)
	knowledge.GET("/stats", handlers.GetKnowledgeStats)

	// Analytics Routes
	analytics := api.Group("/analytics")
	analytics.GET("/overview", handlers.GetAnalyticsOverview)
	analytics.GET("/messages", handlers.GetAnalyticsMessages)
	analytics.GET("/customers", handlers.GetAnalyticsCustomers)
	analytics.GET("/ai", handlers.GetAnalyticsAI)
	analytics.GET("/top-customers", handlers.GetAnalyticsTopCustomers)
	analytics.GET("/hourly", handlers.GetAnalyticsHourly)
	analytics.GET("/intents", handlers.GetAnalyticsIntents)

	// Tags Routes
	tags := api.Group("/tags")
	tags.GET("", handlers.GetTags)
	tags.POST("", handlers.CreateTag)
	tags.PUT("/:id", handlers.UpdateTag)
	tags.DELETE("/:id", handlers.DeleteTag)
	tags.GET("/:id/customers", handlers.GetCustomersByTag)

	return e

}

// validateEnvironment checks that all critical environment variables are set
func validateEnvironment() {
	required := map[string]string{
		"JWT_SECRET":   "JWT secret for authentication",
		"DB_HOST":      "Database host",
		"DB_USER":      "Database user",
		"DB_PASSWORD":  "Database password",
		"DB_NAME":      "Database name",
	}

	missing := []string{}
	for key, description := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, fmt.Sprintf("%s (%s)", key, description))
		}
	}

	if len(missing) > 0 {
		fmt.Println("❌ FATAL: Missing required environment variables:")
		for _, m := range missing {
			fmt.Printf("   - %s\n", m)
		}
		fmt.Println("\nPlease set these variables in your .env file or environment.")
		fmt.Println("See .env.example for reference.")
		os.Exit(1)
	}

	// Validate JWT_SECRET length
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		fmt.Println("❌ FATAL: JWT_SECRET must be at least 32 characters long for security.")
		fmt.Printf("   Current length: %d characters\n", len(jwtSecret))
		fmt.Println("\nGenerate a secure secret with: openssl rand -base64 32")
		os.Exit(1)
	}

	// Warn about placeholder values
	placeholders := map[string]string{
		"JWT_SECRET":  "change_me",
		"DB_PASSWORD": "change_me",
	}

	for key, placeholder := range placeholders {
		value := os.Getenv(key)
		if strings.Contains(strings.ToLower(value), placeholder) {
			fmt.Printf("⚠️  WARNING: %s appears to contain placeholder value. Please change it!\n", key)
		}
	}

	fmt.Println("✅ Environment validation passed")
}

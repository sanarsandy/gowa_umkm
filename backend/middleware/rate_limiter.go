package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	UseRedis          bool
	RedisClient       *redis.Client
}

// InMemoryRateLimiter uses in-memory rate limiting (for development)
type InMemoryRateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter(requestsPerSecond float64, burst int) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

// GetLimiter returns the rate limiter for a given IP
func (rl *InMemoryRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = limiter
	}

	return limiter
}

// CleanupVisitors removes old entries (call this periodically)
func (rl *InMemoryRateLimiter) CleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Simple cleanup: clear all (in production, use more sophisticated approach)
	rl.visitors = make(map[string]*rate.Limiter)
}

// RateLimiterMiddleware creates a rate limiting middleware
func RateLimiterMiddleware(config RateLimiterConfig) echo.MiddlewareFunc {
	limiter := NewInMemoryRateLimiter(config.RequestsPerSecond, config.BurstSize)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.CleanupVisitors()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get client IP
			ip := c.RealIP()
			
			// Get limiter for this IP
			l := limiter.GetLimiter(ip)
			
			// Check if request is allowed
			if !l.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Terlalu banyak permintaan. Silakan coba lagi nanti.",
				})
			}

			return next(c)
		}
	}
}

// AuthRateLimiterMiddleware creates stricter rate limiting for auth endpoints
func AuthRateLimiterMiddleware() echo.MiddlewareFunc {
	return RateLimiterMiddleware(RateLimiterConfig{
		RequestsPerSecond: 2,  // 2 requests per second
		BurstSize:         5,  // Allow burst of 5
	})
}

// APIRateLimiterMiddleware creates rate limiting for general API endpoints
func APIRateLimiterMiddleware() echo.MiddlewareFunc {
	return RateLimiterMiddleware(RateLimiterConfig{
		RequestsPerSecond: 10, // 10 requests per second
		BurstSize:         20, // Allow burst of 20
	})
}

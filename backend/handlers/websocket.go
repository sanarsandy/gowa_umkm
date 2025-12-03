package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gowa-backend/db"
	ws "gowa-backend/services/websocket"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict in production
		return true
	},
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// HandleWebSocket handles WebSocket connections
// Note: This handler bypasses the standard JWT middleware because WebSocket
// connections pass the token as a query parameter instead of Authorization header
func HandleWebSocket(c echo.Context) error {
	// Get token from query parameter (WebSocket can't use headers)
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		fmt.Println("[WebSocket] No token provided")
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	// Parse and validate token manually
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		fmt.Printf("[WebSocket] Invalid token: %v\n", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid token",
		})
	}

	// Extract user_id from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("[WebSocket] Failed to parse claims")
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid claims",
		})
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		fmt.Println("[WebSocket] No user_id in claims")
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No user_id in token",
		})
	}

	// Get tenant_id from database
	var tenantID string
	query := `SELECT id FROM tenants WHERE user_id = $1 AND is_active = true LIMIT 1`
	err = db.DB.QueryRow(query, userID).Scan(&tenantID)
	if err != nil {
		fmt.Printf("[WebSocket] Failed to get tenant: %v\n", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Tenant not found",
		})
	}

	fmt.Printf("[WebSocket] Authenticated: user=%s, tenant=%s\n", userID, tenantID)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("[WebSocket] Failed to upgrade connection: %v\n", err)
		return err
	}

	// Create client
	client := &ws.Client{
		ID:       uuid.New().String(),
		TenantID: tenantID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	// Get hub and register client
	hub := ws.GetHub()
	hub.Register <- client

	fmt.Printf("[WebSocket] New connection: client=%s, tenant=%s\n", client.ID, tenantID)

	// Start goroutines for reading and writing
	go writePump(client, hub)
	go readPump(client, hub)

	return nil
}

// readPump pumps messages from the websocket connection to the hub
func readPump(client *ws.Client, hub *ws.Hub) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("[WebSocket] Read error: %v\n", err)
			}
			break
		}
		// We don't process incoming messages for now
		// Just keep connection alive
	}
}

// writePump pumps messages from the hub to the websocket connection
func writePump(client *ws.Client, hub *ws.Hub) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

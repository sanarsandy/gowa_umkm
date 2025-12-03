package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a single websocket connection
type Client struct {
	ID       string
	TenantID string
	Conn     *websocket.Conn
	Send     chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients per tenant
	clients map[string]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to specific tenant
	Broadcast chan *BroadcastMessage

	mu sync.RWMutex
}

// BroadcastMessage contains the message and target tenant
type BroadcastMessage struct {
	TenantID string
	Event    string
	Data     interface{}
}

// Event types
const (
	EventNewMessage      = "new_message"
	EventNewCustomer     = "new_customer"
	EventCustomerUpdated = "customer_updated"
	EventMessageSent     = "message_sent"
	EventConnectionStatus = "connection_status"
)

// WSMessage is the message format sent to clients
type WSMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Global hub instance
var globalHub *Hub
var once sync.Once

// GetHub returns the global hub instance
func GetHub() *Hub {
	once.Do(func() {
		globalHub = NewHub()
		go globalHub.Run()
	})
	return globalHub
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	fmt.Println("[WebSocket Hub] Starting...")
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.TenantID] == nil {
				h.clients[client.TenantID] = make(map[*Client]bool)
			}
			h.clients[client.TenantID][client] = true
			clientCount := len(h.clients[client.TenantID])
			h.mu.Unlock()
			fmt.Printf("[WebSocket Hub] Client registered: %s (tenant: %s, total: %d)\n", 
				client.ID, client.TenantID, clientCount)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.TenantID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.TenantID)
					}
				}
			}
			h.mu.Unlock()
			fmt.Printf("[WebSocket Hub] Client unregistered: %s (tenant: %s)\n", 
				client.ID, client.TenantID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			clients := h.clients[message.TenantID]
			h.mu.RUnlock()

			if clients == nil {
				continue
			}

			wsMsg := WSMessage{
				Event: message.Event,
				Data:  message.Data,
			}
			
			data, err := json.Marshal(wsMsg)
			if err != nil {
				fmt.Printf("[WebSocket Hub] Failed to marshal message: %v\n", err)
				continue
			}

			fmt.Printf("[WebSocket Hub] Broadcasting %s to %d client(s) for tenant %s\n", 
				message.Event, len(clients), message.TenantID)

			for client := range clients {
				select {
				case client.Send <- data:
				default:
					// Client buffer full, remove client
					h.mu.Lock()
					close(client.Send)
					delete(h.clients[message.TenantID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

// BroadcastToTenant sends a message to all clients of a specific tenant
func (h *Hub) BroadcastToTenant(tenantID string, event string, data interface{}) {
	h.Broadcast <- &BroadcastMessage{
		TenantID: tenantID,
		Event:    event,
		Data:     data,
	}
}

// GetClientCount returns the number of connected clients for a tenant
func (h *Hub) GetClientCount(tenantID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.clients[tenantID]; ok {
		return len(clients)
	}
	return 0
}


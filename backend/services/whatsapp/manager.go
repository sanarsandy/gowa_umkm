package whatsapp

import (
	"context"
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow"
)

// ClientManager manages multiple WhatsApp client instances for different tenants
type ClientManager struct {
	clients map[string]*whatsmeow.Client
	mu      sync.RWMutex
}

// NewClientManager creates a new ClientManager instance
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*whatsmeow.Client),
	}
}

// GetClient retrieves a client for a specific tenant
func (cm *ClientManager) GetClient(tenantID string) (*whatsmeow.Client, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	client, exists := cm.clients[tenantID]
	if !exists {
		return nil, fmt.Errorf("no client found for tenant: %s", tenantID)
	}

	return client, nil
}

// AddClient adds a new client for a tenant
func (cm *ClientManager) AddClient(tenantID string, client *whatsmeow.Client) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.clients[tenantID]; exists {
		return fmt.Errorf("client already exists for tenant: %s", tenantID)
	}

	cm.clients[tenantID] = client
	return nil
}

// RemoveClient removes a client for a tenant
func (cm *ClientManager) RemoveClient(tenantID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.clients[tenantID]
	if !exists {
		return fmt.Errorf("no client found for tenant: %s", tenantID)
	}

	// Disconnect the client if connected
	if client.IsConnected() {
		client.Disconnect()
	}

	delete(cm.clients, tenantID)
	return nil
}

// GetAllClients returns all active clients (for monitoring/debugging)
func (cm *ClientManager) GetAllClients() map[string]*whatsmeow.Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Return a copy to prevent external modification
	clients := make(map[string]*whatsmeow.Client, len(cm.clients))
	for k, v := range cm.clients {
		clients[k] = v
	}

	return clients
}

// GetClientCount returns the number of active clients
func (cm *ClientManager) GetClientCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.clients)
}

// DisconnectAll disconnects all clients (useful for graceful shutdown)
func (cm *ClientManager) DisconnectAll(ctx context.Context) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for tenantID, client := range cm.clients {
		if client.IsConnected() {
			client.Disconnect()
			fmt.Printf("Disconnected client for tenant: %s\n", tenantID)
		}
	}
}

// Global client manager instance
var globalClientManager *ClientManager
var once sync.Once

// GetGlobalClientManager returns the singleton ClientManager instance
func GetGlobalClientManager() *ClientManager {
	once.Do(func() {
		globalClientManager = NewClientManager()
	})
	return globalClientManager
}

package models

import "time"

// Tenant represents a business tenant in the system
type Tenant struct {
	ID                  string    `json:"id" db:"id"`
	UserID              string    `json:"user_id" db:"user_id"`
	BusinessName        string    `json:"business_name" db:"business_name"`
	BusinessType        string    `json:"business_type" db:"business_type"`
	BusinessDescription string    `json:"business_description" db:"business_description"`
	BusinessPhone       string    `json:"business_phone" db:"business_phone"`
	BusinessAddress     string    `json:"business_address" db:"business_address"`
	IsActive            bool      `json:"is_active" db:"is_active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// WhatsAppDevice represents a WhatsApp device connection
type WhatsAppDevice struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	JID             string     `json:"jid" db:"jid"`
	IsConnected     bool       `json:"is_connected" db:"is_connected"`
	LastConnectedAt *time.Time `json:"last_connected_at" db:"last_connected_at"`
	Platform        string     `json:"platform" db:"platform"`
	BusinessName    string     `json:"business_name" db:"business_name"`
	PushName        string     `json:"push_name" db:"push_name"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// WhatsAppConnectionRequest represents a request to connect WhatsApp
type WhatsAppConnectionRequest struct {
	TenantID string `json:"tenant_id"`
}

// WhatsAppConnectionResponse represents the response after initiating connection
type WhatsAppConnectionResponse struct {
	QRCode    string `json:"qr_code,omitempty"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	StreamURL string `json:"stream_url,omitempty"`
}

// WhatsAppStatusResponse represents the current WhatsApp connection status
type WhatsAppStatusResponse struct {
	IsConnected bool       `json:"is_connected"`
	JID         string     `json:"jid,omitempty"`
	PushName    string     `json:"push_name,omitempty"`
	LastSeen    *time.Time `json:"last_seen,omitempty"`
	Status      string     `json:"status"`
}

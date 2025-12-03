package whatsapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"bytes"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

// QRCodeGenerator handles QR code generation for WhatsApp pairing
type QRCodeGenerator struct{}

// NewQRCodeGenerator creates a new QR code generator
func NewQRCodeGenerator() *QRCodeGenerator {
	return &QRCodeGenerator{}
}

// GenerateQRCode generates a QR code from the pairing code
func (g *QRCodeGenerator) GenerateQRCode(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("pairing code cannot be empty")
	}

	// Generate QR code
	qr, err := qrcode.New(code, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Set size
	qr.DisableBorder = false

	// Convert to PNG bytes
	var buf bytes.Buffer
	err = png.Encode(&buf, qr.Image(256))
	if err != nil {
		return "", fmt.Errorf("failed to encode QR code: %w", err)
	}

	// Encode to base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	
	// Return as data URL
	return fmt.Sprintf("data:image/png;base64,%s", base64Str), nil
}

// QRChannelResult represents a QR code event
type QRChannelResult struct {
	Event   string `json:"event"`
	Code    string `json:"code,omitempty"`
	QRImage string `json:"qr_image,omitempty"`
	Error   string `json:"error,omitempty"`
}

// GetQRChannel returns a channel that receives QR codes from whatsmeow
func GetQRChannel(ctx context.Context, client *whatsmeow.Client) (<-chan QRChannelResult, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	qrChan := make(chan QRChannelResult, 5)
	generator := NewQRCodeGenerator()

	// Get QR channel from whatsmeow
	qrCodeChan, err := client.GetQRChannel(ctx)
	if err != nil {
		close(qrChan)
		return nil, fmt.Errorf("failed to get QR channel: %w", err)
	}

	// Start goroutine to process QR codes
	go func() {
		defer close(qrChan)

		for evt := range qrCodeChan {
			if evt.Event == "code" {
				// Generate QR image
				qrImage, err := generator.GenerateQRCode(evt.Code)
				if err != nil {
					qrChan <- QRChannelResult{
						Event: "error",
						Error: fmt.Sprintf("Failed to generate QR image: %v", err),
					}
					continue
				}

				qrChan <- QRChannelResult{
					Event:   "code",
					Code:    evt.Code,
					QRImage: qrImage,
				}
			} else if evt.Event == "success" {
				qrChan <- QRChannelResult{
					Event: "success",
				}
				return
			} else if evt.Event == "timeout" {
				qrChan <- QRChannelResult{
					Event: "timeout",
					Error: "QR code expired",
				}
				return
			}
		}
	}()

	return qrChan, nil
}

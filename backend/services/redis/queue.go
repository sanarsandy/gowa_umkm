package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Queue names
	AIQueueKey        = "ai:messages:queue"
	BroadcastQueueKey = "broadcast:messages:queue"
)

// MessagePayload represents a message in the AI processing queue
type MessagePayload struct {
	TenantID    string `json:"tenant_id"`
	MessageID   string `json:"message_id"`
	SenderJID   string `json:"sender_jid"`
	MessageText string `json:"message_text"`
	Timestamp   int64  `json:"timestamp"`
}

// BroadcastMessagePayload represents a broadcast message in the queue
type BroadcastMessagePayload struct {
	TenantID      string `json:"tenant_id"`
	BroadcastID   string `json:"broadcast_id"`
	RecipientID   string `json:"recipient_id"`
	CustomerJID   string `json:"customer_jid"`
	Message       string `json:"message"`
	CustomerName  string `json:"customer_name,omitempty"`
}

// PushToAIQueue adds a message to the AI processing queue
func (c *Client) PushToAIQueue(ctx context.Context, payload *MessagePayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return c.rdb.RPush(ctx, AIQueueKey, data).Err()
}

// PopFromAIQueue retrieves and removes a message from the AI queue
func (c *Client) PopFromAIQueue(ctx context.Context, timeout time.Duration) (*MessagePayload, error) {
	result, err := c.rdb.BLPop(ctx, timeout, AIQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No message available
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("unexpected result format")
	}

	var payload MessagePayload
	if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &payload, nil
}

// PushToBroadcastQueue adds a broadcast message to the queue
func (c *Client) PushToBroadcastQueue(ctx context.Context, payload *BroadcastMessagePayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast payload: %w", err)
	}

	return c.rdb.RPush(ctx, BroadcastQueueKey, data).Err()
}

// PopFromBroadcastQueue retrieves and removes a broadcast message from the queue
func (c *Client) PopFromBroadcastQueue(ctx context.Context, timeout time.Duration) (*BroadcastMessagePayload, error) {
	result, err := c.rdb.BLPop(ctx, timeout, BroadcastQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No message available
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("unexpected result format")
	}

	var payload BroadcastMessagePayload
	if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal broadcast payload: %w", err)
	}

	return &payload, nil
}

// GetQueueLength returns the number of items in a queue
func (c *Client) GetQueueLength(ctx context.Context, queueKey string) (int64, error) {
	return c.rdb.LLen(ctx, queueKey).Result()
}

package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gowa-backend/services/redis"

	"github.com/jmoiron/sqlx"
)

// BroadcastScheduler handles scheduled and recurring broadcasts
type BroadcastScheduler struct {
	db          *sqlx.DB
	redisClient *redis.Client
	ticker      *time.Ticker
	done        chan bool
}

// Broadcast represents a broadcast with scheduling info
type Broadcast struct {
	ID                string         `db:"id"`
	TenantID          string         `db:"tenant_id"`
	Name              string         `db:"name"`
	MessageContent    string         `db:"message_content"`
	TemplateID        sql.NullString `db:"template_id"`
	Status            string         `db:"status"`
	ScheduledAt       sql.NullTime   `db:"scheduled_at"`
	IsRecurring       bool           `db:"is_recurring"`
	RecurrenceType    sql.NullString `db:"recurrence_type"`
	RecurrenceInterval sql.NullInt32 `db:"recurrence_interval"`
	RecurrenceDays    sql.NullString `db:"recurrence_days"` // JSONB as string
	RecurrenceTime    sql.NullString `db:"recurrence_time"` // TIME as string
	RecurrenceEndDate sql.NullTime   `db:"recurrence_end_date"`
	RecurrenceCount   sql.NullInt32  `db:"recurrence_count"`
	LastExecutedAt    sql.NullTime   `db:"last_executed_at"`
	ExecutionCount    int            `db:"execution_count"`
}

// NewBroadcastScheduler creates a new broadcast scheduler
func NewBroadcastScheduler(db *sqlx.DB, redisClient *redis.Client) *BroadcastScheduler {
	return &BroadcastScheduler{
		db:          db,
		redisClient: redisClient,
		done:        make(chan bool),
	}
}

// Start begins the scheduler (checks every minute)
func (s *BroadcastScheduler) Start() {
	log.Println("[Scheduler] Starting broadcast scheduler...")
	s.ticker = time.NewTicker(1 * time.Minute)

	// Run immediately on start
	s.checkScheduledBroadcasts()

	// Then run every minute
	for {
		select {
		case <-s.ticker.C:
			s.checkScheduledBroadcasts()
		case <-s.done:
			log.Println("[Scheduler] Stopping broadcast scheduler...")
			return
		}
	}
}

// Stop stops the scheduler
func (s *BroadcastScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
}

// checkScheduledBroadcasts checks for broadcasts that need to be executed
func (s *BroadcastScheduler) checkScheduledBroadcasts() {
	ctx := context.Background()
	now := time.Now()

	// Query broadcasts that are ready to execute
	query := `
		SELECT id, tenant_id, name, message_content, template_id, status,
		       scheduled_at, is_recurring, recurrence_type, recurrence_interval,
		       recurrence_days, recurrence_time, recurrence_end_date, recurrence_count,
		       last_executed_at, execution_count
		FROM broadcasts
		WHERE (status = 'scheduled' OR (status = 'active' AND is_recurring = true))
		  AND scheduled_at <= $1
		ORDER BY scheduled_at ASC
		LIMIT 100
	`

	var broadcasts []Broadcast
	err := s.db.SelectContext(ctx, &broadcasts, query, now)
	if err != nil {
		log.Printf("[Scheduler] Error querying scheduled broadcasts: %v", err)
		return
	}

	if len(broadcasts) > 0 {
		log.Printf("[Scheduler] Found %d broadcasts ready to execute", len(broadcasts))
	}

	for _, broadcast := range broadcasts {
		s.executeBroadcast(ctx, &broadcast)
	}
}

// executeBroadcast executes a broadcast
func (s *BroadcastScheduler) executeBroadcast(ctx context.Context, broadcast *Broadcast) {
	log.Printf("[Scheduler] Executing broadcast: %s (ID: %s)", broadcast.Name, broadcast.ID)

	// Start transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[Scheduler] Error starting transaction: %v", err)
		return
	}
	defer tx.Rollback()

	// Update status to sending
	updateQuery := `
		UPDATE broadcasts
		SET status = 'sending',
		    started_at = NOW(),
		    last_executed_at = NOW(),
		    execution_count = execution_count + 1,
		    updated_at = NOW()
		WHERE id = $1 AND (status = 'scheduled' OR status = 'active')
		RETURNING id
	`

	var updatedID string
	err = tx.GetContext(ctx, &updatedID, updateQuery, broadcast.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Already being processed
			return
		}
		log.Printf("[Scheduler] Error updating broadcast status: %v", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[Scheduler] Error committing transaction: %v", err)
		return
	}

	log.Printf("[Scheduler] Broadcast %s marked as sending", broadcast.ID)

	// Trigger actual message sending
	go s.sendBroadcastMessages(broadcast.TenantID, broadcast.ID, broadcast.MessageContent)

	// Handle recurring broadcasts
	if broadcast.IsRecurring {
		s.scheduleNextRecurrence(ctx, broadcast)
	}
}

// sendBroadcastMessages queues messages to Redis for all recipients
func (s *BroadcastScheduler) sendBroadcastMessages(tenantID, broadcastID, messageTemplate string) {
	ctx := context.Background()
	
	// Get recipients
	type Recipient struct {
		ID          string `db:"id"`
		CustomerID  string `db:"customer_id"`
		CustomerJID string `db:"customer_jid"`
	}
	
	var recipients []Recipient
	recipientQuery := `SELECT id, customer_id, customer_jid FROM broadcast_recipients WHERE broadcast_id = $1 AND status = 'pending'`
	if err := s.db.SelectContext(ctx, &recipients, recipientQuery, broadcastID); err != nil {
		log.Printf("[Scheduler] Error getting recipients for broadcast %s: %v", broadcastID, err)
		return
	}
	
	if len(recipients) == 0 {
		log.Printf("[Scheduler] No recipients found for broadcast %s", broadcastID)
		// Mark as completed anyway
		s.db.ExecContext(ctx, `UPDATE broadcasts SET status = 'completed', completed_at = NOW() WHERE id = $1`, broadcastID)
		return
	}

	log.Printf("[Scheduler] Queueing broadcast %s to %d recipients via Redis", broadcastID, len(recipients))

	queuedCount := 0
	failedCount := 0

	for _, recipient := range recipients {
		// Get customer name for personalization
		var customerName string
		nameQuery := `SELECT COALESCE(customer_name, customer_phone, '') FROM customer_insights WHERE id = $1`
		s.db.Get(&customerName, nameQuery, recipient.CustomerID)

		// Personalize message
		personalizedMessage := messageTemplate
		personalizedMessage = strings.Replace(personalizedMessage, "{{nama}}", customerName, -1)
		personalizedMessage = strings.Replace(personalizedMessage, "{{name}}", customerName, -1)

		// Create broadcast message payload
		payload := &redis.BroadcastMessagePayload{
			TenantID:     tenantID,
			BroadcastID:  broadcastID,
			RecipientID:  recipient.ID,
			CustomerJID:  recipient.CustomerJID,
			Message:      personalizedMessage,
			CustomerName: customerName,
		}

		// Push to Redis broadcast queue
		if err := s.redisClient.PushToBroadcastQueue(ctx, payload); err != nil {
			log.Printf("[Scheduler] Error queuing message for recipient %s: %v", recipient.ID, err)
			failedCount++
			// Mark as failed
			s.db.ExecContext(ctx, `UPDATE broadcast_recipients SET status = 'failed', error_message = $1 WHERE id = $2`, err.Error(), recipient.ID)
		} else {
			log.Printf("[Scheduler] Queued message to %s: %s", recipient.CustomerJID, personalizedMessage)
			queuedCount++
			// Mark as queued
			s.db.ExecContext(ctx, `UPDATE broadcast_recipients SET status = 'queued' WHERE id = $1`, recipient.ID)
		}

		// Small delay to avoid overwhelming Redis
		time.Sleep(50 * time.Millisecond)
	}

	// Update broadcast counts
	updateBroadcastQuery := `UPDATE broadcasts SET sent_count = $1, failed_count = $2, status = 'sending', updated_at = NOW() WHERE id = $3`
	s.db.ExecContext(ctx, updateBroadcastQuery, queuedCount, failedCount, broadcastID)

	log.Printf("[Scheduler] Broadcast %s - Queued: %d, Failed: %d messages", broadcastID, queuedCount, failedCount)
}

// scheduleNextRecurrence calculates and schedules the next occurrence
func (s *BroadcastScheduler) scheduleNextRecurrence(ctx context.Context, broadcast *Broadcast) {
	// Check if should continue recurring
	if !s.shouldContinueRecurring(broadcast) {
		// Mark as completed
		_, err := s.db.ExecContext(ctx, `
			UPDATE broadcasts
			SET status = 'completed',
			    completed_at = NOW(),
			    updated_at = NOW()
			WHERE id = $1
		`, broadcast.ID)
		if err != nil {
			log.Printf("[Scheduler] Error marking recurring broadcast as completed: %v", err)
		} else {
			log.Printf("[Scheduler] Recurring broadcast %s completed", broadcast.ID)
		}
		return
	}

	// Calculate next execution time
	nextExecution := s.calculateNextExecution(broadcast)
	if nextExecution.IsZero() {
		log.Printf("[Scheduler] Could not calculate next execution for broadcast %s", broadcast.ID)
		return
	}

	// Update scheduled_at for next execution
	_, err := s.db.ExecContext(ctx, `
		UPDATE broadcasts
		SET scheduled_at = $1,
		    status = 'active',
		    updated_at = NOW()
		WHERE id = $2
	`, nextExecution, broadcast.ID)

	if err != nil {
		log.Printf("[Scheduler] Error scheduling next recurrence: %v", err)
	} else {
		log.Printf("[Scheduler] Next execution for %s scheduled at %s", broadcast.ID, nextExecution.Format(time.RFC3339))
	}
}

// shouldContinueRecurring checks if recurring broadcast should continue
func (s *BroadcastScheduler) shouldContinueRecurring(broadcast *Broadcast) bool {
	now := time.Now()

	// Check end date
	if broadcast.RecurrenceEndDate.Valid {
		if now.After(broadcast.RecurrenceEndDate.Time) {
			return false
		}
	}

	// Check execution count
	if broadcast.RecurrenceCount.Valid {
		if broadcast.ExecutionCount >= int(broadcast.RecurrenceCount.Int32) {
			return false
		}
	}

	return true
}

// calculateNextExecution calculates the next execution time for recurring broadcast
func (s *BroadcastScheduler) calculateNextExecution(broadcast *Broadcast) time.Time {
	if !broadcast.IsRecurring || !broadcast.RecurrenceType.Valid {
		return time.Time{}
	}

	lastExec := time.Now()
	if broadcast.LastExecutedAt.Valid {
		lastExec = broadcast.LastExecutedAt.Time
	}

	interval := 1
	if broadcast.RecurrenceInterval.Valid {
		interval = int(broadcast.RecurrenceInterval.Int32)
	}

	switch broadcast.RecurrenceType.String {
	case "hourly":
		return lastExec.Add(time.Duration(interval) * time.Hour)

	case "daily":
		return s.calculateDailyNext(lastExec, interval, broadcast.RecurrenceTime.String)

	case "weekly":
		return s.calculateWeeklyNext(lastExec, interval, broadcast.RecurrenceDays.String, broadcast.RecurrenceTime.String)

	default:
		return time.Time{}
	}
}

// calculateDailyNext calculates next daily execution
func (s *BroadcastScheduler) calculateDailyNext(lastExec time.Time, interval int, timeStr string) time.Time {
	// Parse time (format: HH:MM:SS)
	var hour, minute int
	if timeStr != "" {
		fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	}

	// Add interval days
	nextDate := lastExec.AddDate(0, 0, interval)

	// Set to specific time
	return time.Date(
		nextDate.Year(), nextDate.Month(), nextDate.Day(),
		hour, minute, 0, 0,
		lastExec.Location(),
	)
}

// calculateWeeklyNext calculates next weekly execution
func (s *BroadcastScheduler) calculateWeeklyNext(lastExec time.Time, interval int, daysJSON string, timeStr string) time.Time {
	// Parse time
	var hour, minute int
	if timeStr != "" {
		fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	}

	// Parse days
	var days []string
	if daysJSON != "" {
		json.Unmarshal([]byte(daysJSON), &days)
	}

	if len(days) == 0 {
		return time.Time{}
	}

	// Convert day names to weekday numbers
	dayMap := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	var targetDays []time.Weekday
	for _, day := range days {
		if wd, ok := dayMap[day]; ok {
			targetDays = append(targetDays, wd)
		}
	}

	// Find next occurrence
	current := lastExec.AddDate(0, 0, 1) // Start from next day
	for i := 0; i < 14; i++ {             // Check up to 2 weeks
		for _, targetDay := range targetDays {
			if current.Weekday() == targetDay {
				return time.Date(
					current.Year(), current.Month(), current.Day(),
					hour, minute, 0, 0,
					current.Location(),
				)
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	return time.Time{}
}

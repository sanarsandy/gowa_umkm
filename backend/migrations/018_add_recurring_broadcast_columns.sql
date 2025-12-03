-- Migration 018: Add Recurring Broadcast Columns
-- Adds missing columns for scheduled and recurring broadcast functionality

-- Add recurring broadcast columns
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS is_recurring BOOLEAN DEFAULT FALSE;
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_type VARCHAR(20); -- hourly, daily, weekly, monthly
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_interval INTEGER DEFAULT 1;
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_days JSONB; -- For weekly: ["monday", "friday"]
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_time TIME; -- Time of day for daily/weekly
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_end_date TIMESTAMP;
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_count INTEGER; -- Max number of executions
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS last_executed_at TIMESTAMP;
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS execution_count INTEGER DEFAULT 0;

-- Create indexes for scheduler performance
CREATE INDEX IF NOT EXISTS idx_broadcasts_scheduled ON broadcasts(tenant_id, scheduled_at) WHERE status IN ('scheduled', 'active');
CREATE INDEX IF NOT EXISTS idx_broadcasts_recurring ON broadcasts(is_recurring) WHERE is_recurring = TRUE;

COMMENT ON COLUMN broadcasts.is_recurring IS 'Whether this broadcast repeats on a schedule';
COMMENT ON COLUMN broadcasts.recurrence_type IS 'Type of recurrence: hourly, daily, weekly, monthly';
COMMENT ON COLUMN broadcasts.recurrence_interval IS 'Interval for recurrence (e.g., every 2 days)';
COMMENT ON COLUMN broadcasts.recurrence_days IS 'Days of week for weekly recurrence (JSON array)';
COMMENT ON COLUMN broadcasts.recurrence_time IS 'Time of day for daily/weekly recurrence';
COMMENT ON COLUMN broadcasts.execution_count IS 'Number of times this broadcast has been executed';

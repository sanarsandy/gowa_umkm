-- Migration 008: Add Recurring Broadcast Support
-- This migration adds columns to support recurring broadcasts (hourly, daily, weekly)

-- Add recurring broadcast columns to broadcasts table
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS is_recurring BOOLEAN DEFAULT false;
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_type VARCHAR(20); -- 'hourly', 'daily', 'weekly'
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_interval INTEGER DEFAULT 1; -- e.g., every 2 days, every 3 hours
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_days JSONB; -- For weekly: ["monday", "wednesday", "friday"]
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_time TIME; -- Time of day for daily/weekly (e.g., 10:00:00)
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_end_date TIMESTAMPTZ; -- When to stop recurring
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS recurrence_count INTEGER; -- Max number of occurrences (alternative to end_date)
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS last_executed_at TIMESTAMPTZ; -- Track last execution time
ALTER TABLE broadcasts ADD COLUMN IF NOT EXISTS execution_count INTEGER DEFAULT 0; -- Track number of executions

-- Add index for efficient querying of recurring broadcasts
CREATE INDEX IF NOT EXISTS idx_broadcasts_recurring ON broadcasts(tenant_id, is_recurring, status) WHERE is_recurring = true;
CREATE INDEX IF NOT EXISTS idx_broadcasts_next_execution ON broadcasts(tenant_id, scheduled_at, is_recurring) WHERE status IN ('scheduled', 'active');

-- Add comments for documentation
COMMENT ON COLUMN broadcasts.is_recurring IS 'Whether this broadcast is recurring';
COMMENT ON COLUMN broadcasts.recurrence_type IS 'Type of recurrence: hourly, daily, weekly';
COMMENT ON COLUMN broadcasts.recurrence_interval IS 'Interval for recurrence (e.g., every 2 days)';
COMMENT ON COLUMN broadcasts.recurrence_days IS 'Days of week for weekly recurrence (JSON array)';
COMMENT ON COLUMN broadcasts.recurrence_time IS 'Time of day for daily/weekly recurrence';
COMMENT ON COLUMN broadcasts.recurrence_end_date IS 'Date when recurring broadcast should stop';
COMMENT ON COLUMN broadcasts.recurrence_count IS 'Maximum number of executions';
COMMENT ON COLUMN broadcasts.last_executed_at IS 'Timestamp of last execution';
COMMENT ON COLUMN broadcasts.execution_count IS 'Number of times this broadcast has been executed';

-- Update status enum to include 'active' for recurring broadcasts
-- Note: 'active' status is for recurring broadcasts that are running
-- 'scheduled' is for one-time or next execution of recurring
-- 'completed' is when recurring broadcast reaches end_date or count

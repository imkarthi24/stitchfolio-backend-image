-- Migration: 004_add_task_entity-- Generated: 2026-02-08T14:56:04+05:30

-- ====================================
-- UP Migration
-- ====================================

-- Create table: stich.Tasks
CREATE TABLE IF NOT EXISTS stich."Tasks" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  is_completed BOOLEAN DEFAULT false,
  priority INTEGER,
  due_date TIMESTAMPTZ,
  reminder_date TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,
  assigned_to_id INTEGER,
  PRIMARY KEY (id)
);


-- Add foreign key to stich.Tasks
ALTER TABLE stich."Tasks" ADD CONSTRAINT fk_Task_assigned_to_id FOREIGN KEY (assigned_to_id) REFERENCES stich."Users" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;


-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually

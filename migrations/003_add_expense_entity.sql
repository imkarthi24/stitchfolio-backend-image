-- Migration: 003_add_expense_entity
-- Generated: 2026-01-28T20:36:54+05:30

-- ====================================
-- UP Migration
-- ====================================

-- Create table: stich.Expenses
CREATE TABLE IF NOT EXISTS stich."Expenses" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  purchase_date TIMESTAMPTZ,
  bill_number TEXT,
  company_name TEXT,
  material TEXT,
  price DOUBLE PRECISION,
  location TEXT,
  notes TEXT,
  PRIMARY KEY (id)
);


-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually
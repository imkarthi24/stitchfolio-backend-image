-- Migration: 002_person_entity_update
-- Generated: 2026-01-17T21:25:35+05:30

-- ====================================
-- UP Migration
-- ====================================

-- Add column to stitch.Persons
ALTER TABLE stitch."Persons" ADD COLUMN gender TEXT;

-- Add column to stitch.Persons
ALTER TABLE stitch."Persons" ADD COLUMN age INTEGER;


-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually

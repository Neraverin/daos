-- Migration: Replace compose_content with role_path
-- This is a breaking change - roles now store folder paths instead of inline content

ALTER TABLE roles DROP COLUMN compose_content;
ALTER TABLE roles ADD COLUMN role_path TEXT NOT NULL;

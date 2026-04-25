-- Migration: 022_add_pricing_tagline_cta
-- Add tagline and cta_button columns to pricing table

ALTER TABLE pricing ADD COLUMN IF NOT EXISTS tagline TEXT DEFAULT '';
ALTER TABLE pricing ADD COLUMN IF NOT EXISTS cta_button TEXT DEFAULT 'Get Started';

-- Update existing rows with default values
UPDATE pricing SET tagline = '' WHERE tagline IS NULL;
UPDATE pricing SET cta_button = 'Get Started' WHERE cta_button IS NULL;

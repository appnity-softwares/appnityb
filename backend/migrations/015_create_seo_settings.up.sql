CREATE TABLE IF NOT EXISTS seo_settings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_title VARCHAR(255),
  site_description TEXT,
  keywords TEXT[],
  og_image TEXT,
  twitter_handle VARCHAR(100),
  canonical_url TEXT,
  robots_txt TEXT,
  google_analytics_id VARCHAR(100),
  google_tag_manager VARCHAR(100),
  json_ld_schema JSONB,
  favicon TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_seo_settings_updated_at
BEFORE UPDATE ON seo_settings
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS themes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  is_active BOOLEAN DEFAULT false,
  colors JSONB DEFAULT '{}',
  fonts JSONB DEFAULT '{}',
  typography JSONB DEFAULT '{}',
  border_radius JSONB DEFAULT '{}',
  spacing JSONB DEFAULT '{}',
  shadows JSONB DEFAULT '{}',
  breakpoints JSONB DEFAULT '{}',
  background_images JSONB DEFAULT '{}',
  animations JSONB DEFAULT '{}',
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

CREATE TRIGGER update_themes_updated_at
BEFORE UPDATE ON themes
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

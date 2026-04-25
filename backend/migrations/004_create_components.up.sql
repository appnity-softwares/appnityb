CREATE TABLE IF NOT EXISTS components (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  page_id UUID REFERENCES pages(id) ON DELETE CASCADE,
  component_type VARCHAR(100) NOT NULL,
  is_visible BOOLEAN DEFAULT true,
  content JSONB DEFAULT '{}',
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_components_page_id ON components(page_id);
CREATE INDEX idx_components_order_index ON components(order_index);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_components_updated_at
BEFORE UPDATE ON components
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

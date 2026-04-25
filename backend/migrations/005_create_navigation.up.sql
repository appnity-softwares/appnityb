CREATE TABLE IF NOT EXISTS navigation (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  label VARCHAR(255) NOT NULL,
  url TEXT NOT NULL,
  parent_id UUID REFERENCES navigation(id) ON DELETE CASCADE,
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  is_external BOOLEAN DEFAULT false,
  icon VARCHAR(100),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_navigation_parent_id ON navigation(parent_id);
CREATE INDEX idx_navigation_order_index ON navigation(order_index);

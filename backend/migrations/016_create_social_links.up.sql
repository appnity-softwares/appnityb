CREATE TABLE IF NOT EXISTS social_links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  platform VARCHAR(100) NOT NULL,
  url TEXT NOT NULL,
  is_visible BOOLEAN DEFAULT true,
  icon VARCHAR(100),
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_social_links_order_index ON social_links(order_index);
CREATE INDEX idx_social_links_is_visible ON social_links(is_visible);

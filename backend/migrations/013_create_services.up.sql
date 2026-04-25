CREATE TABLE IF NOT EXISTS services (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  icon VARCHAR(100),
  image TEXT,
  features JSONB DEFAULT '[]',
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_services_order_index ON services(order_index);
CREATE INDEX idx_services_is_visible ON services(is_visible);

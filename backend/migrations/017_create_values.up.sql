CREATE TABLE IF NOT EXISTS "values" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  icon VARCHAR(100),
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_values_order_index ON "values"(order_index);
CREATE INDEX idx_values_is_visible ON "values"(is_visible);

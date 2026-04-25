CREATE TABLE IF NOT EXISTS pricing (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_name VARCHAR(255) NOT NULL,
  price_inr DECIMAL(10,2),
  price_usd DECIMAL(10,2),
  billing_period VARCHAR(50) DEFAULT 'monthly',
  description TEXT,
  features JSONB DEFAULT '[]',
  is_highlighted BOOLEAN DEFAULT false,
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pricing_order_index ON pricing(order_index);
CREATE INDEX idx_pricing_is_visible ON pricing(is_visible);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_pricing_updated_at
BEFORE UPDATE ON pricing
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

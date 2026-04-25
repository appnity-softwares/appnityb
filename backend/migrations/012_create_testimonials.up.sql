CREATE TABLE IF NOT EXISTS testimonials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  role VARCHAR(255),
  company VARCHAR(255),
  avatar TEXT,
  quote TEXT NOT NULL,
  rating INT DEFAULT 5,
  video_url TEXT,
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_testimonials_order_index ON testimonials(order_index);
CREATE INDEX idx_testimonials_is_visible ON testimonials(is_visible);
CREATE INDEX idx_testimonials_rating ON testimonials(rating);

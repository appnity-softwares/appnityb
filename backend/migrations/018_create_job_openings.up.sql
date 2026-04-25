CREATE TABLE IF NOT EXISTS job_openings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  department VARCHAR(100),
  location VARCHAR(255),
  type VARCHAR(100),
  description TEXT,
  requirements JSONB DEFAULT '[]',
  is_active BOOLEAN DEFAULT true,
  apply_url TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_job_openings_is_active ON job_openings(is_active);
CREATE INDEX idx_job_openings_department ON job_openings(department);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_job_openings_updated_at
BEFORE UPDATE ON job_openings
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

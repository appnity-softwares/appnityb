CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  full_name VARCHAR(255),
  role VARCHAR(50) DEFAULT 'admin',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  last_login TIMESTAMPTZ
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
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
CREATE TABLE IF NOT EXISTS pages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  slug VARCHAR(255) UNIQUE NOT NULL,
  title VARCHAR(255) NOT NULL,
  meta_description TEXT,
  meta_keywords TEXT,
  og_image TEXT,
  canonical_url TEXT,
  background_color VARCHAR(7),
  custom_css TEXT,
  is_published BOOLEAN DEFAULT true,
  sections JSONB DEFAULT '[]',
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

CREATE TRIGGER update_pages_updated_at
BEFORE UPDATE ON pages
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
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
CREATE TABLE IF NOT EXISTS blogs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  excerpt TEXT,
  content TEXT NOT NULL,
  cover_image TEXT,
  author VARCHAR(255),
  category VARCHAR(100),
  tags TEXT[],
  read_time INT,
  is_published BOOLEAN DEFAULT false,
  published_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_blogs_slug ON blogs(slug);
CREATE INDEX idx_blogs_is_published ON blogs(is_published);
CREATE INDEX idx_blogs_published_at ON blogs(published_at);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_blogs_updated_at
BEFORE UPDATE ON blogs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TABLE IF NOT EXISTS portfolios (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  cover_image TEXT,
  gallery_images JSONB DEFAULT '[]',
  category VARCHAR(100),
  client_name VARCHAR(255),
  project_url TEXT,
  metrics JSONB DEFAULT '{}',
  is_featured BOOLEAN DEFAULT false,
  is_visible BOOLEAN DEFAULT true,
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_portfolios_slug ON portfolios(slug);
CREATE INDEX idx_portfolios_is_visible ON portfolios(is_visible);
CREATE INDEX idx_portfolios_order_index ON portfolios(order_index);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_portfolios_updated_at
BEFORE UPDATE ON portfolios
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TABLE IF NOT EXISTS team (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name VARCHAR(255) NOT NULL,
  role VARCHAR(255) NOT NULL,
  bio TEXT,
  photo TEXT,
  social_links JSONB DEFAULT '{}',
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_team_order_index ON team(order_index);
CREATE INDEX idx_team_is_visible ON team(is_visible);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_team_updated_at
BEFORE UPDATE ON team
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TABLE IF NOT EXISTS awards (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  image TEXT,
  year INT,
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_awards_order_index ON awards(order_index);
CREATE INDEX idx_awards_is_visible ON awards(is_visible);
CREATE TABLE IF NOT EXISTS faqs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  question TEXT NOT NULL,
  answer TEXT NOT NULL,
  order_index INT DEFAULT 0,
  is_visible BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_faqs_order_index ON faqs(order_index);
CREATE INDEX idx_faqs_is_visible ON faqs(is_visible);
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
  icon VARCHAR(50) DEFAULT 'star',
  gradient TEXT DEFAULT 'from-white to-gray-50',
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

-- Seed 4 default pricing plans
DELETE FROM pricing;
INSERT INTO pricing (id, plan_name, price_inr, price_usd, billing_period, description, features, is_highlighted, order_index, is_visible, icon, gradient) VALUES
(gen_random_uuid(), 'Starter', 9999, 119, 'monthly', 'Perfect for individuals and small projects needing clean, fast design delivery.', '["You provide the wireframe", "Visual design using Figma", "Focused on website or branding only", "Weekday turnaround (Mon–Fri)", "Email support"]'::jsonb, false, 0, true, 'star', 'from-white to-gray-50'),
(gen_random_uuid(), 'Professional', 19999, 239, 'monthly', 'A complete design experience — tailored strategy, polished visuals, and flexible collaboration.', '["Help shaping your wireframe or brief", "Custom design for website, brand, or logo", "High-fidelity mockups using Figma & Framer", "Dedicated weekday focus & deeper involvement", "Priority email support"]'::jsonb, true, 1, true, 'zap', 'from-orange-600 to-orange-900'),
(gen_random_uuid(), 'Business', 39999, 479, 'monthly', 'For growing teams needing comprehensive design and development solutions.', '["Everything in Professional", "Full-stack development included", "Unlimited revisions", "24/7 priority support", "Monthly strategy calls", "Source code delivery"]'::jsonb, false, 2, true, 'crown', 'from-gray-900 to-black'),
(gen_random_uuid(), 'Enterprise', 79999, 959, 'monthly', 'Complete digital transformation for large organizations with dedicated team.', '["Everything in Business", "Dedicated project manager", "White-label deliverables", "Custom integrations", "SLA guarantee", "Quarterly reviews"]'::jsonb, false, 3, true, 'sparkles', 'from-violet-700 to-purple-900');
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
CREATE TABLE IF NOT EXISTS contacts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  service VARCHAR(255),
  message TEXT,
  status VARCHAR(50) DEFAULT 'new',
  admin_notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_contacts_status ON contacts(status);
CREATE INDEX idx_contacts_email ON contacts(email);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_contacts_updated_at
BEFORE UPDATE ON contacts
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
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
CREATE TABLE IF NOT EXISTS media (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  filename VARCHAR(255) NOT NULL,
  original_name VARCHAR(255),
  mime_type VARCHAR(100),
  size BIGINT,
  url TEXT NOT NULL,
  alt_text TEXT,
  folder VARCHAR(100) DEFAULT 'general',
  uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_media_folder ON media(folder);
CREATE INDEX idx_media_uploaded_by ON media(uploaded_by);
CREATE INDEX idx_media_mime_type ON media(mime_type);

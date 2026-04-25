-- CLEAR EXISTING DATA
TRUNCATE TABLE portfolios, team, services, pricing, testimonials, faqs, "values" CASCADE;

-- SEED PORTFOLIOS
INSERT INTO portfolios (id, title, slug, description, category, client_name, project_url, is_featured, order_index) VALUES
(gen_random_uuid(), 'Chhattisgarh Shadi', 'chhattisgarh-shadi', 'A community-focused matrimonial platform built for Chhattisgarh. Features real-time matching, profile verification, and secure chat.', 'Mobile Application', 'Chhatisgarh Shadi', 'https://chhattisgarhshadi.com', true, 0),
(gen_random_uuid(), 'Mitaan Express', 'mitaan-express', 'High-performance news website for Shri Dhar Rao. Real-time content updates, intuitive categorization, and optimized for speed.', 'News Portal', 'Shridhar Rao', 'https://mitaanexpress.com', true, 1),
(gen_random_uuid(), 'Shri Dhar Rao Portfolio', 'shri-dhar-rao-portfolio', 'Official professional portfolio for Shri Dhar Rao. A minimal, authority-driven design showcasing personal and political achievements.', 'Personal Brand', 'Shridhar Rao', 'https://shridharrao.com', true, 2),
(gen_random_uuid(), 'Crova', 'crova', 'A high-end boutique showcase for handcrafted embroidery and apparel. Built with a focus on minimal design and heritage craftsmanship.', 'Premium Showcase', 'Crova', 'https://crova.in', true, 3);

-- SEED TEAM
INSERT INTO team (id, full_name, role, bio, social_links, order_index) VALUES
(gen_random_uuid(), 'Saurabh Jain', 'Co-founder & CEO', 'Leading the business strategy and growth. Committed to building Appnity into a premier global engineering house.', '{"github": "https://github.com/jsaurabh334"}'::jsonb, 0),
(gen_random_uuid(), 'Pushp Raj Sharma', 'Managing Director', 'Overseeing operational excellence and strategic partnerships. Ensuring absolute precision in project delivery.', '{"linkedin": "https://www.linkedin.com/in/pushp-raj-sharma/", "github": "https://github.com/pushp314"}'::jsonb, 1),
(gen_random_uuid(), 'Neha Mourya', 'Full Stack Web Developer', 'Specialist in React, Node.js and distributed systems. Crafting seamless user experiences with robust backend logic.', '{"github": "https://github.com/nehamoury", "linkedin": "https://www.linkedin.com/in/nehamourya/"}'::jsonb, 2),
(gen_random_uuid(), 'Kunal Daharwal', 'App Developer', 'Mobile engineering expert specializing in React Native CLI. Building native-grade experiences for iOS and Android.', '{"github": "https://github.com/kunal592"}'::jsonb, 3),
(gen_random_uuid(), 'Lelin Helina Tandon', 'Head of Human Resources', 'Building a culture of excellence and senior-grade engineering talent. Managing global operations and team growth.', '{"linkedin": "https://www.linkedin.com/in/helina-tandan-27b9b93b9"}'::jsonb, 4),
(gen_random_uuid(), 'Jatin Kurrey', 'SDE & R&D Intern', 'Focusing on core software development and research into emerging UI/UX architectures. Delivering high-performance digital systems.', '{"github": "https://github.com/jatin-kurrey", "linkedin": "https://www.linkedin.com/in/jatin-kurrey-07a046251"}'::jsonb, 5);

-- SEED SERVICES (SOLUTIONS)
INSERT INTO services (id, title, description, icon, features, order_index) VALUES
(gen_random_uuid(), 'Sales Growth Engine', 'Stop losing leads. Organize your sales process and close deals faster with automated follow-ups.', 'Users', '["Lead Organizer", "Auto Follow-ups", "Sales Reports"]'::jsonb, 0),
(gen_random_uuid(), 'Knowledge & Training Hub', 'Train your employees or students anywhere in the world with a professional learning platform.', 'BookOpen', '["Video Lessons", "Student Progress", "Mobile App"]'::jsonb, 1),
(gen_random_uuid(), 'Subscription Business', 'Launch your own ''Netflix'' or ''Spotify'' style business with built-in recurring payments.', 'Layers', '["Easy Payments", "User Accounts", "Growth Ready"]'::jsonb, 2),
(gen_random_uuid(), 'Modern School Manager', 'Manage admissions, fees, and results from a single dashboard. No more manual paperwork.', 'School', '["Fee Collection", "Smart Reports", "Parent App"]'::jsonb, 3),
(gen_random_uuid(), 'Property & Real Estate', 'Showcase your properties beautifully and manage inquiries effortlessly on one site.', 'Building2', '["Property Gallery", "Agent Portal", "Map Search"]'::jsonb, 4),
(gen_random_uuid(), 'Vendor Marketplace', 'Build your own Amazon or Flipkart. Let multiple sellers list products on your platform.', 'ShoppingCart', '["Seller Dashboard", "Admin Control", "Easy Payouts"]'::jsonb, 5);

-- SEED PRICING
INSERT INTO pricing (id, plan_name, price_inr, billing_period, description, features, is_highlighted, order_index) VALUES
(gen_random_uuid(), 'Digital Presence', 25000, 'starting at', 'Premium static websites and high-performance landing pages for elite brands.', '["Custom UI/UX Design", "Next.js Static Export", "Perfect Lighthouse Score", "SEO & Metadata Hardening", "CMS Integration (Optional)", "CDN Global Deployment"]'::jsonb, false, 0),
(gen_random_uuid(), 'Dynamic Systems', 85000, 'starting at', 'Full-stack digital products, SaaS MVPs, and complex web applications.', '["Custom Web Application", "Database Architecture", "Auth & Security Protocols", "API Development (REST/GraphQL)", "Cloud Infrastructure Setup", "3 Months Post-Launch Support"]'::jsonb, true, 1),
(gen_random_uuid(), 'Enterprise Solutions', 0, 'quote based', 'Custom engineering for large-scale distributed systems and AI pipelines.', '["Distributed Architecture", "Legacy System Migration", "AI & Data Engineering", "Security Audits & Hardening", "Dedicated Engineering Unit", "On-prem / Hybrid Deployment"]'::jsonb, false, 2);

-- SEED TESTIMONIALS
INSERT INTO testimonials (id, name, role, company, quote, order_index) VALUES
(gen_random_uuid(), 'Rashmeet Kaur', 'CEO', 'Crova', 'Building a custom wardrobe management system was a huge challenge. Appnity delivered a solution that is both technically robust and beautiful. Their attention to detail is exceptional.', 0),
(gen_random_uuid(), 'Jaideep Gupta', 'Founder', 'Chhatisgarh Shadi', 'For a platform as personal as matrimonial services, we needed a partner who understood both technology and trust. Appnity built a secure, high-performance system.', 1),
(gen_random_uuid(), 'Prateek Tatode', 'CTO', 'GrowthHub', 'Transparent, fast, and professional. The team at Appnity feels like an extension of our own engineering department. They truly architect excellence.', 2),
(gen_random_uuid(), 'Shridhar Rao', 'Editor', 'Mitaan Express', 'Delivering absolute precision was key for my personal brand. Appnity shipped a high-performance site that exceeds all my expectations in terms of speed and UI.', 3);

-- SEED FAQS
INSERT INTO faqs (id, question, answer, order_index) VALUES
(gen_random_uuid(), 'How long does a typical project take?', 'Most MVP systems ship within 6-8 weeks. Larger enterprise platforms can take 3-5 months. We work in 2-week sprints so you see progress constantly.', 0),
(gen_random_uuid(), 'Do you provide post-launch support?', 'Yes. We offer managed cloud maintenance and ''Fractional CTO'' services to ensure your systems remain secure and scalable as you grow.', 1),
(gen_random_uuid(), 'Will I own the code?', 'Absolutely. Once the final payment is made, you own 100% of the IP and the source code. No vendor lock-in, ever.', 2),
(gen_random_uuid(), 'How do we communicate during development?', 'We set up a dedicated Slack channel for your team and conduct weekly video demos. You''ll have direct access to the senior engineers working on your project.', 3),
(gen_random_uuid(), 'Can you help me migrate from a legacy system?', 'Yes, we specialize in legacy modernization—moving your old systems to modern, cloud-native architectures without losing data or causing downtime.', 4);

-- SEED VALUES
INSERT INTO "values" (id, title, description, icon, order_index) VALUES
(gen_random_uuid(), 'Extreme Ownership', 'We don''t just write code; we take responsibility for the business outcome. If something isn''t right, we fix it.', 'Target', 0),
(gen_random_uuid(), 'Speed with Precision', 'We move fast, but we never compromise on quality. Automated testing and senior oversight are non-negotiable.', 'Zap', 1),
(gen_random_uuid(), 'Radical Candor', 'We tell you what you need to hear, not what you want to hear. Honesty is the foundation of our partnership.', 'Heart', 2);

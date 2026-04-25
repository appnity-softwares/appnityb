package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://user:password@localhost:5432/appnity?sslmode=require"
	}
	fmt.Println("Using DATABASE_URL:", databaseURL)

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
		return
	}
	defer pool.Close()

	ctx := context.Background()
	if err := pool.Ping(ctx); err != nil {
		log.Printf("Could not ping database: %v", err)
		log.Println("Seed script ready. Set DATABASE_URL environment variable and run again.")
		return
	}

	log.Println("Seeding database...")

	seedUsers(ctx, pool)
	seedThemes(ctx, pool)
	seedNavigation(ctx, pool)
	seedBlogs(ctx, pool)
	seedPortfolios(ctx, pool)
	seedTeam(ctx, pool)
	seedServices(ctx, pool)
	seedFAQs(ctx, pool)
	seedPricing(ctx, pool)
	seedTestimonials(ctx, pool)
	seedAwards(ctx, pool)
	seedValues(ctx, pool)
	seedSocialLinks(ctx, pool)
	seedSEO(ctx, pool)

	log.Println("Database seeded successfully!")
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Admin@123456"), bcrypt.DefaultCost)
	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash, updated_at = NOW()
	`, uuid.New(), "admin@appnity.co.in", string(hashedPassword), "Appnity Admin", "admin", true)
	if err != nil {
		log.Printf("Error seeding users: %v", err)
	}
}

func seedThemes(ctx context.Context, pool *pgxpool.Pool) {
	colors, _ := json.Marshal(map[string]string{
		"primary":        "#f97316",
		"secondary":      "#111111",
		"accent":         "#f97316",
		"background":     "#e6e6e6",
		"text":           "#111111",
		"surface":        "#ffffff",
		"border":         "#dbdbdb",
		"gradient_start": "#f97316",
		"gradient_end":   "#ff4d00",
	})
	fonts, _ := json.Marshal(map[string]interface{}{
		"heading": map[string]interface{}{"family": "Plus Jakarta Sans", "weights": []int{400, 500, 700, 800}},
		"body":    map[string]interface{}{"family": "Plus Jakarta Sans", "weights": []int{300, 400, 500}},
		"cursive": map[string]interface{}{"family": "Pacifico"},
	})
	typography, _ := json.Marshal(map[string]interface{}{
		"h1": map[string]interface{}{"size": "4rem", "weight": 800, "lineHeight": 1.1},
		"h2": map[string]interface{}{"size": "3rem", "weight": 700, "lineHeight": 1.2},
		"h3": map[string]interface{}{"size": "2rem", "weight": 700, "lineHeight": 1.3},
		"h4": map[string]interface{}{"size": "1.5rem", "weight": 600, "lineHeight": 1.3},
		"p":  map[string]interface{}{"size": "1rem", "weight": 400, "lineHeight": 1.6},
	})
	borderRadius, _ := json.Marshal(map[string]string{
		"sm": "0.25rem", "md": "0.5rem", "lg": "1rem", "xl": "1.5rem", "full": "9999px",
	})
	spacing, _ := json.Marshal(map[string]string{
		"xs": "0.25rem", "sm": "0.5rem", "md": "1rem", "lg": "1.5rem", "xl": "2rem", "2xl": "3rem",
	})

	_, err := pool.Exec(ctx, `
		INSERT INTO themes (id, name, is_active, colors, fonts, typography, border_radius, spacing, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), "Appnity Default", true, colors, fonts, typography, borderRadius, spacing)
	if err != nil {
		log.Printf("Error seeding themes: %v", err)
	}
}

func seedNavigation(ctx context.Context, pool *pgxpool.Pool) {
	items := []struct {
		label, url string
		order      int
	}{
		{"Home", "/", 0},
		{"Works", "/works", 1},
		{"Services", "/services", 2},
		{"About", "/about", 3},
		{"Team", "/team", 4},
		{"Blog", "/blog", 5},
		{"Portfolio", "/portfolio", 6},
		{"Contact", "/contact", 7},
	}
	for _, item := range items {
		_, err := pool.Exec(ctx, `
			INSERT INTO navigation (id, label, url, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), item.label, item.url, item.order, true)
		if err != nil {
			log.Printf("Error seeding navigation: %v", err)
		}
	}
}

func seedBlogs(ctx context.Context, pool *pgxpool.Pool) {
	blogs := []struct {
		title, slug, excerpt, content, coverImage, author, category string
		readTime                                                    int
	}{
		{
			"The Future of UI/UX Design in 2025",
			"future-ui-ux-design-2025",
			"Explore the latest trends shaping digital experiences",
			"<p>The design landscape is evolving rapidly. From AI-powered interfaces to immersive 3D experiences, 2025 brings exciting opportunities for designers and developers alike.</p><p>Key trends include micro-interactions, glassmorphism evolution, and AI-assisted design tools that streamline the creative process.</p>",
			"https://framerusercontent.com/images/placeholder1.png",
			"Pushpa Raj",
			"Design",
			5,
		},
		{
			"Building Scalable Design Systems",
			"building-scalable-design-systems",
			"How to create design systems that grow with your startup",
			"<p>A well-architected design system is the backbone of any successful product team. Learn how to build one from scratch.</p><p>We cover component architecture, token systems, documentation practices, and team workflows.</p>",
			"https://framerusercontent.com/images/placeholder2.png",
			"Pushpa Raj",
			"Development",
			8,
		},
		{
			"Startup Growth Through Better Design",
			"startup-growth-better-design",
			"Why design investment directly impacts your bottom line",
			"<p>Startups that prioritize design from day one see 2x more user engagement and 3x faster growth rates.</p><p>This article breaks down the ROI of design investment with real case studies from our portfolio.</p>",
			"https://framerusercontent.com/images/placeholder3.png",
			"Pushpa Raj",
			"Business",
			6,
		},
	}
	for _, b := range blogs {
		_, err := pool.Exec(ctx, `
			INSERT INTO blogs (id, title, slug, excerpt, content, cover_image, author, category, read_time, is_published, published_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), NOW())
			ON CONFLICT (slug) DO NOTHING
		`, uuid.New(), b.title, b.slug, b.excerpt, b.content, b.coverImage, b.author, b.category, b.readTime, true)
		if err != nil {
			log.Printf("Error seeding blogs: %v", err)
		}
	}
}

func seedPortfolios(ctx context.Context, pool *pgxpool.Pool) {
	portfolios := []struct {
		title, slug, description, coverImage, category, clientName, projectURL string
		isFeatured                                                             bool
	}{
		{
			"Crova - E-Commerce Platform",
			"crova-ecommerce",
			"Complete redesign of an e-commerce platform resulting in 40% increase in conversions",
			"https://framerusercontent.com/images/crova-cover.png",
			"E-Commerce",
			"Crova",
			"https://crova.vercel.app/",
			true,
		},
		{
			"Growth Hub - Marketing Agency",
			"growth-hub-agency",
			"Brand identity and website for a full-service marketing agency",
			"https://framerusercontent.com/images/growthhub-cover.png",
			"Branding",
			"Growth Hub",
			"https://www.growth-hub.co.in",
			false,
		},
		{
			"BrainWave - AI Platform",
			"brainwave-ai",
			"UI/UX design for an AI-powered mock interview platform",
			"https://framerusercontent.com/images/brainwave-cover.png",
			"Product Design",
			"BrainWave",
			"https://ai-mock-interviews-five-alpha.vercel.app",
			false,
		},
		{
			"UI Zone - Design Tool",
			"ui-zone-tool",
			"Interactive design tool for rapid prototyping and collaboration",
			"https://framerusercontent.com/images/uizone-cover.png",
			"Product Design",
			"UI Zone",
			"https://uizone314.vercel.app",
			false,
		},
	}
	for i, p := range portfolios {
		metrics, _ := json.Marshal(map[string]interface{}{
			"conversion_rate": "40%",
			"satisfaction":    "98%",
			"timeline":        "6 weeks",
		})
		_, err := pool.Exec(ctx, `
			INSERT INTO portfolios (id, title, slug, description, cover_image, category, client_name, project_url, metrics, is_featured, is_visible, order_index, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
			ON CONFLICT (slug) DO NOTHING
		`, uuid.New(), p.title, p.slug, p.description, p.coverImage, p.category, p.clientName, p.projectURL, metrics, p.isFeatured, true, i)
		if err != nil {
			log.Printf("Error seeding portfolios: %v", err)
		}
	}
}

func seedTeam(ctx context.Context, pool *pgxpool.Pool) {
	members := []struct {
		name, role, bio, photo string
		social                 map[string]string
	}{
		{
			"Pushpa Raj",
			"Founder & Lead Designer",
			"10+ years of experience in UI/UX design, leading startups to launch successful products.",
			"/founder.jpeg",
			map[string]string{"linkedin": "#", "instagram": "#", "twitter": "#"},
		},
		{
			"Team Member 2",
			"Senior Developer",
			"Full-stack developer with expertise in React, Go, and cloud architecture.",
			"/p-main.avif",
			map[string]string{"linkedin": "#", "github": "#"},
		},
		{
			"Team Member 3",
			"UI Designer",
			"Passionate about creating beautiful, functional interfaces that users love.",
			"/team3.avif",
			map[string]string{"dribbble": "#", "behance": "#"},
		},
	}
	for i, m := range members {
		social, _ := json.Marshal(m.social)
		_, err := pool.Exec(ctx, `
			INSERT INTO team (id, full_name, role, bio, photo, social_links, order_index, is_visible, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), m.name, m.role, m.bio, m.photo, social, i, true)
		if err != nil {
			log.Printf("Error seeding team: %v", err)
		}
	}
}

func seedServices(ctx context.Context, pool *pgxpool.Pool) {
	services := []struct {
		title, description, icon string
		features                 []string
	}{
		{"Brand Identity", "Complete brand design from logo to guidelines", "palette", []string{"Logo Design", "Color Palette", "Typography", "Brand Guidelines"}},
		{"Product Design", "End-to-end product design for web and mobile", "smartphone", []string{"UI/UX Design", "Prototyping", "User Research", "Design Systems"}},
		{"Web Development", "Modern, performant web applications", "code", []string{"React", "Next.js", "Go Backend", "Database Design"}},
	}
	for i, s := range services {
		features, _ := json.Marshal(s.features)
		_, err := pool.Exec(ctx, `
			INSERT INTO services (id, title, description, icon, features, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), s.title, s.description, s.icon, features, i, true)
		if err != nil {
			log.Printf("Error seeding services: %v", err)
		}
	}
}

func seedFAQs(ctx context.Context, pool *pgxpool.Pool) {
	faqs := []struct{ question, answer string }{
		{"What services does Appnity offer?", "We offer brand identity design, product design (UI/UX), web development, and ongoing design support for startups."},
		{"How long does a typical project take?", "Most projects take 4-8 weeks depending on scope. We work in agile sprints to deliver value quickly."},
		{"What is your pricing model?", "We offer monthly retainers starting at Rs 9,999/month for our Standard plan and custom pricing for Premium engagements."},
		{"Do you work with early-stage startups?", "Absolutely! We specialize in helping early-stage startups establish their brand and product design."},
		{"Can I see examples of your work?", "Check out our Portfolio page for detailed case studies of our recent projects."},
		{"How do we get started?", "Simply reach out through our Contact page and we'll schedule a free consultation to discuss your needs."},
	}
	for i, f := range faqs {
		_, err := pool.Exec(ctx, `
			INSERT INTO faqs (id, question, answer, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), f.question, f.answer, i, true)
		if err != nil {
			log.Printf("Error seeding FAQs: %v", err)
		}
	}
}

func seedPricing(ctx context.Context, pool *pgxpool.Pool) {
	plans := []struct {
		name        string
		priceINR    float64
		priceUSD    float64
		description string
		features    []string
		highlighted bool
		icon        string
		gradient    string
		tagline     string
		ctaButton   string
	}{
		{
			"Starter",
			9999,
			119,
			"Perfect for individuals and small projects needing clean, fast design delivery.",
			[]string{"You provide the wireframe", "Visual design using Figma", "Focused on website or branding only", "Weekday turnaround (Mon–Fri)", "Email support"},
			false,
			"star",
			"from-white to-gray-50",
			"Best Value",
			"Choose Starter",
		},
		{
			"Professional",
			19999,
			239,
			"A complete design experience — tailored strategy, polished visuals, and flexible collaboration.",
			[]string{"Help shaping your wireframe or brief", "Custom design for website, brand, or logo", "High-fidelity mockups using Figma & Framer", "Dedicated weekday focus & deeper involvement", "Priority email support"},
			true,
			"zap",
			"from-orange-600 to-orange-900",
			"Most Popular",
			"Choose Professional",
		},
		{
			"Business",
			39999,
			479,
			"For growing teams needing comprehensive design and development solutions.",
			[]string{"Everything in Professional", "Full-stack development included", "Unlimited revisions", "24/7 priority support", "Monthly strategy calls", "Source code delivery"},
			false,
			"crown",
			"from-gray-900 to-black",
			"Best Value",
			"Choose Business",
		},
		{
			"Enterprise",
			79999,
			959,
			"Complete digital transformation for large organizations with dedicated team.",
			[]string{"Everything in Business", "Dedicated project manager", "White-label deliverables", "Custom integrations", "SLA guarantee", "Quarterly reviews", "On-site workshops"},
			false,
			"sparkles",
			"from-violet-700 to-purple-900",
			"Premium",
			"Contact Sales",
		},
	}
	for i, p := range plans {
		features, _ := json.Marshal(p.features)
		_, err := pool.Exec(ctx, `
			INSERT INTO pricing (id, plan_name, price_inr, price_usd, billing_period, description, features, is_highlighted, order_index, is_visible, icon, gradient, tagline, cta_button, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), p.name, p.priceINR, p.priceUSD, "monthly", p.description, features, p.highlighted, i, true, p.icon, p.gradient, p.tagline, p.ctaButton)
		if err != nil {
			log.Printf("Error seeding pricing: %v", err)
		}
	}
}

func seedTestimonials(ctx context.Context, pool *pgxpool.Pool) {
	testimonials := []struct {
		name, role, company, quote, avatar string
		rating                             int
	}{
		{"Client A", "CEO", "TechStartup", "Appnity transformed our vision into reality. The design quality exceeded all expectations.", "/avatar.png", 5},
		{"Client B", "Founder", "GrowthCo", "Working with Appnity was seamless. They understood our brand and delivered beyond expectations.", "/avatar1.png", 5},
		{"Client C", "CTO", "InnovateLab", "The attention to detail and user-centric approach made a huge impact on our product metrics.", "/p1.png", 5},
	}
	for i, t := range testimonials {
		_, err := pool.Exec(ctx, `
			INSERT INTO testimonials (id, name, role, company, quote, avatar, rating, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), t.name, t.role, t.company, t.quote, t.avatar, t.rating, i, true)
		if err != nil {
			log.Printf("Error seeding testimonials: %v", err)
		}
	}
}

func seedAwards(ctx context.Context, pool *pgxpool.Pool) {
	awards := []struct {
		title, description, image string
		year                      int
	}{
		{"Best UI/UX Design Agency", "Recognized for outstanding design work in the startup ecosystem", "/firs.jpeg", 2024},
		{"Top 10 Design Agencies India", "Featured among India's leading design agencies", "/second1.jpeg", 2024},
		{"Innovation Award", "For innovative use of AI in design workflows", "/third3.png", 2023},
	}
	for i, a := range awards {
		_, err := pool.Exec(ctx, `
			INSERT INTO awards (id, title, description, image, year, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), a.title, a.description, a.image, a.year, i, true)
		if err != nil {
			log.Printf("Error seeding awards: %v", err)
		}
	}
}

func seedValues(ctx context.Context, pool *pgxpool.Pool) {
	values := []struct{ title, description, icon string }{
		{"Design-First", "Every pixel, every interaction is crafted with intention and care.", "sparkles"},
		{"Client-Centric", "Your success is our success. We partner with you, not just work for you.", "heart"},
		{"Innovation", "We push boundaries and explore new possibilities in design and technology.", "lightbulb"},
		{"Transparency", "Clear communication, honest pricing, no hidden agendas.", "eye"},
	}
	for i, v := range values {
		_, err := pool.Exec(ctx, `
			INSERT INTO "values" (id, title, description, icon, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), v.title, v.description, v.icon, i, true)
		if err != nil {
			log.Printf("Error seeding values: %v", err)
		}
	}
}

func seedSocialLinks(ctx context.Context, pool *pgxpool.Pool) {
	links := []struct{ platform, url, icon string }{
		{"instagram", "https://instagram.com/appnity", "instagram"},
		{"linkedin", "https://linkedin.com/company/appnity", "linkedin"},
		{"twitter", "https://twitter.com/appnity", "twitter"},
		{"behance", "https://behance.net/appnity", "behance"},
		{"dribbble", "https://dribbble.com/appnity", "dribbble"},
	}
	for i, l := range links {
		_, err := pool.Exec(ctx, `
			INSERT INTO social_links (id, platform, url, icon, order_index, is_visible, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), l.platform, l.url, l.icon, i, true)
		if err != nil {
			log.Printf("Error seeding social links: %v", err)
		}
	}
}

func seedSEO(ctx context.Context, pool *pgxpool.Pool) {
	keywords := []string{"UI/UX Design", "Startup Design", "Brand Identity", "Product Design", "Web Development", "Bhilai", "India"}
	jsonLd, _ := json.Marshal(map[string]interface{}{
		"@context":    "https://schema.org",
		"@type":       "Organization",
		"name":        "Appnity",
		"url":         "https://www.appnity.co.in",
		"description": "Effortless Design for Startups based in Bhilai, (C.G)",
	})
	_, err := pool.Exec(ctx, `
		INSERT INTO seo_settings (id, site_title, site_description, keywords, og_image, twitter_handle, canonical_url, google_analytics_id, json_ld_schema, favicon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), "Appnity - Effortless Design for Startups", "We make it easy for startups to launch, grow, and scale with clean, conversion-focused designs", keywords, "https://www.appnity.co.in/ap2.png", "@appnity", "https://www.appnity.co.in", "", jsonLd, "/ap2.png")
	if err != nil {
		log.Printf("Error seeding SEO: %v", err)
	}
}

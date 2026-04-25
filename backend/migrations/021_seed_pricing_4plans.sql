-- 4 Default Pricing Plans for Appnity
-- Run this after migration 012_add_pricing_style

-- Delete existing pricing plans first (optional, remove if you want to keep existing)
DELETE FROM pricing;

-- Starter Plan
INSERT INTO pricing (id, plan_name, price_inr, price_usd, billing_period, description, features, is_highlighted, order_index, is_visible, icon, gradient, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Starter',
    9999,
    119,
    'monthly',
    'Perfect for individuals and small projects needing clean, fast design delivery.',
    '["You provide the wireframe", "Visual design using Figma", "Focused on website or branding only", "Weekday turnaround (Mon–Fri)", "Email support"]'::jsonb,
    false,
    0,
    true,
    'star',
    'from-white to-gray-50',
    NOW(),
    NOW()
);

-- Professional Plan (Highlighted)
INSERT INTO pricing (id, plan_name, price_inr, price_usd, billing_period, description, features, is_highlighted, order_index, is_visible, icon, gradient, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Professional',
    19999,
    239,
    'monthly',
    'A complete design experience — tailored strategy, polished visuals, and flexible collaboration.',
    '["Help shaping your wireframe or brief", "Custom design for website, brand, or logo", "High-fidelity mockups using Figma & Framer", "Dedicated weekday focus & deeper involvement", "Priority email support"]'::jsonb,
    true,
    1,
    true,
    'zap',
    'from-orange-600 to-orange-900',
    NOW(),
    NOW()
);

-- Business Plan
INSERT INTO pricing (id, plan_name, price_inr, price_usd, billing_period, description, features, is_highlighted, order_index, is_visible, icon, gradient, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Business',
    39999,
    479,
    'monthly',
    'For growing teams needing comprehensive design and development solutions.',
    '["Everything in Professional", "Full-stack development included", "Unlimited revisions", "24/7 priority support", "Monthly strategy calls", "Source code delivery"]'::jsonb,
    false,
    2,
    true,
    'crown',
    'from-gray-900 to-black',
    NOW(),
    NOW()
);

-- Enterprise Plan
INSERT INTO pricing (id, plan_name, price_inr, price_usd, billing_period, description, features, is_highlighted, order_index, is_visible, icon, gradient, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Enterprise',
    79999,
    959,
    'monthly',
    'Complete digital transformation for large organizations with dedicated team.',
    '["Everything in Business", "Dedicated project manager", "White-label deliverables", "Custom integrations", "SLA guarantee", "Quarterly reviews", "On-site workshops"]'::jsonb,
    false,
    3,
    true,
    'sparkles',
    'from-violet-700 to-purple-900',
    NOW(),
    NOW()
);

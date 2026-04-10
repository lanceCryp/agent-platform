-- AgentHub 数据库初始化脚本
-- 版本: 001
-- 创建时间: 2026-04-10

-- ============================================
-- 扩展
-- ============================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================
-- 1. 用户表 (users)
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'vip', 'admin')),
    balance DECIMAL(10, 2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);

-- ============================================
-- 2. Agent 分类表 (categories)
-- ============================================
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    icon VARCHAR(50),  -- emoji 或图标名
    description TEXT,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_sort ON categories(sort_order);

-- ============================================
-- 3. Agent 表 (agents)
-- ============================================
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id VARCHAR(100) UNIQUE NOT NULL,  -- 唯一标识，如 "engineering-code-reviewer"
    name VARCHAR(200) NOT NULL,
    name_en VARCHAR(200),
    description TEXT,
    category_id UUID REFERENCES categories(id),
    category VARCHAR(50) NOT NULL,  -- 分类 slug
    tags TEXT[] DEFAULT '{}',
    tier INT DEFAULT 3 CHECK (tier BETWEEN 1 AND 4),
    runtime_type VARCHAR(20) NOT NULL CHECK (runtime_type IN ('openclaw', 'claude', 'openai')),
    price_per_request DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    price_per_token DECIMAL(10, 6),
    avg_duration_seconds INT DEFAULT 60,
    success_rate DECIMAL(5, 2) DEFAULT 100.00,
    rating DECIMAL(2, 1) DEFAULT 5.0 CHECK (rating BETWEEN 0 AND 5),
    total_tasks INT DEFAULT 0,
    total_success INT DEFAULT 0,
    total_failed INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    config JSONB DEFAULT '{}',
    input_example TEXT,
    output_example TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agents_agent_id ON agents(agent_id);
CREATE INDEX idx_agents_category ON agents(category);
CREATE INDEX idx_agents_tier ON agents(tier);
CREATE INDEX idx_agents_runtime ON agents(runtime_type);
CREATE INDEX idx_agents_rating ON agents(rating DESC);
CREATE INDEX idx_agents_tasks ON agents(total_tasks DESC);
CREATE INDEX idx_agents_active ON agents(is_active) WHERE is_active = true;
-- 全文搜索索引
CREATE INDEX idx_agents_search ON agents USING gin(name gin_trgm_ops);
CREATE INDEX idx_agents_search_desc ON agents USING gin(description gin_trgm_ops);

-- ============================================
-- 4. 任务表 (tasks)
-- ============================================
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    agent_id UUID NOT NULL REFERENCES agents(id),
    prompt TEXT NOT NULL,
    result TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    priority INT DEFAULT 5 CHECK (priority BETWEEN 1 AND 10),
    cost DECIMAL(10, 2) DEFAULT 0.00,
    tokens_used INT DEFAULT 0,
    duration_seconds INT DEFAULT 0,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    parent_task_id UUID REFERENCES tasks(id),
    context JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tasks_user ON tasks(user_id);
CREATE INDEX idx_tasks_agent ON tasks(agent_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created ON tasks(created_at DESC);
CREATE INDEX idx_tasks_user_status ON tasks(user_id, status);
CREATE INDEX idx_tasks_priority ON tasks(priority DESC) WHERE status IN ('pending', 'processing');

-- ============================================
-- 5. 任务消息表 (task_messages)
-- ============================================
CREATE TABLE IF NOT EXISTS task_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_messages_task ON task_messages(task_id);
CREATE INDEX idx_task_messages_created ON task_messages(created_at);

-- ============================================
-- 6. 套餐表 (plans)
-- ============================================
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL UNIQUE CHECK (type IN ('basic', 'pro', 'enterprise')),
    description TEXT,
    price_monthly DECIMAL(10, 2) NOT NULL,
    price_yearly DECIMAL(10, 2) NOT NULL,
    discount_percentage DECIMAL(5, 2) DEFAULT 0,
    features JSONB NOT NULL DEFAULT '{}',
    task_limit INT,  -- NULL 表示无限
    agent_tier_limit INT DEFAULT 2,
    api_access BOOLEAN DEFAULT false,
    priority_support BOOLEAN DEFAULT false,
    custom_agents BOOLEAN DEFAULT false,
    max_concurrent_tasks INT DEFAULT 5,
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_plans_type ON plans(type);
CREATE INDEX idx_plans_active ON plans(is_active) WHERE is_active = true;

-- ============================================
-- 7. 订阅表 (subscriptions)
-- ============================================
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    plan_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'cancelled', 'expired', 'pending')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('monthly', 'yearly')),
    auto_renew BOOLEAN DEFAULT true,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    payment_method VARCHAR(50),
    payment_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, plan_type)  -- 每个用户每种类型只能有一个活跃订阅
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_end_date ON subscriptions(end_date);

-- ============================================
-- 8. 交易记录表 (transactions)
-- ============================================
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('recharge', 'consumption', 'refund', 'subscription', 'bonus')),
    amount DECIMAL(10, 2) NOT NULL,
    balance_before DECIMAL(10, 2) NOT NULL,
    balance_after DECIMAL(10, 2) NOT NULL,
    description VARCHAR(500),
    reference_id UUID,  -- 关联的订阅ID、任务ID等
    reference_type VARCHAR(50),  -- 'task', 'subscription', 'recharge'
    payment_method VARCHAR(50),
    payment_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'completed' CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_created ON transactions(created_at DESC);
CREATE INDEX idx_transactions_reference ON transactions(reference_id, reference_type);

-- ============================================
-- 9. 用户会话表 (sessions)
-- ============================================
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255),
    device_info VARCHAR(500),
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    refresh_expires_at TIMESTAMP WITH TIME ZONE,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- ============================================
-- 10. 操作日志表 (audit_logs)
-- ============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);

-- ============================================
-- 触发器: 自动更新 updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_agents_updated_at
    BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_plans_updated_at
    BEFORE UPDATE ON plans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 种子数据
-- ============================================

-- 插入分类
INSERT INTO categories (name, slug, icon, description, sort_order) VALUES
('工程开发', 'engineering', '💻', '代码开发、测试、部署等工程相关任务', 1),
('设计创意', 'design', '🎨', 'UI设计、品牌视觉、插画等创意工作', 2),
('市场营销', 'marketing', '📈', '营销策划、文案创作、SEO优化等', 3),
('数据分析', 'data', '📊', '数据挖掘、可视化、BI分析等', 4),
('游戏开发', 'gaming', '🎮', '游戏设计、关卡设计、剧情创作等', 5),
('财务管理', 'finance', '💰', '财务分析、投资评估、预算规划等', 6),
('法务合规', 'legal', '⚖️', '合同审查、法律咨询、合规检查等', 7),
('学术研究', 'academic', '📚', '论文写作、文献综述、研究方法等', 8)
ON CONFLICT (slug) DO NOTHING;

-- 插入套餐
INSERT INTO plans (name, type, description, price_monthly, price_yearly, features, task_limit, agent_tier_limit, api_access, priority_support, custom_agents, max_concurrent_tasks, is_featured, sort_order) VALUES
('基础版', 'basic', '适合个人用户，入门级AI Agent服务', 99.00, 950.00,
 '{"task_limit": 100, "agent_tier_limit": 2, "api_access": false, "priority_support": false, "custom_agents": false}',
 100, 2, false, false, false, 5, false, 1),
('专业版', 'pro', '适合小型团队，解锁全部Agent能力', 299.00, 2800.00,
 '{"task_limit": null, "agent_tier_limit": 3, "api_access": true, "priority_support": true, "custom_agents": false}',
 NULL, 3, true, true, false, 20, true, 2),
('企业版', 'enterprise', '适合中大型企业，完整解决方案', 599.00, 5500.00,
 '{"task_limit": null, "agent_tier_limit": 4, "api_access": true, "priority_support": true, "custom_agents": true}',
 NULL, 4, true, true, true, 100, false, 3)
ON CONFLICT (type) DO NOTHING;

-- 插入示例 Agent
INSERT INTO agents (agent_id, name, name_en, description, category, tags, tier, runtime_type, price_per_request, avg_duration_seconds, success_rate, rating, input_example, output_example) VALUES
('engineering-code-reviewer', '代码审查专家', 'Code Review Expert',
 '专业代码审查，发现潜在bug、安全漏洞和性能问题。提供详细的改进建议。',
 'engineering', ARRAY['代码审查', '质量', '安全', '性能'], 1, 'claude', 9.90, 45, 98.50, 4.8,
 '请帮我审查这段Python代码的潜在问题:\n\ndef calculate(a, b):\n    return a / b',
 '✅ 代码审查完成\n\n发现 1 个潜在问题:\n\n**除零风险** (Line 2)\n- 问题: 当 b 为 0 时会抛出 ZeroDivisionError\n- 建议: 添加除零检查\n\n```python\ndef calculate(a, b):\n    if b == 0:\n        return 0  # 或抛出自定义异常\n    return a / b\n```'),
 
('design-brand-guardian', '品牌视觉设计师', 'Brand Guardian',
 '确保品牌视觉一致性，提供设计建议和品牌规范指南。',
 'design', ARRAY['品牌设计', '视觉', '规范', 'LOGO'], 2, 'openclaw', 29.90, 120, 96.20, 4.9,
 '请检查这个LOGO设计是否符合现代科技公司的品牌定位，并给出改进建议。',
 '✅ 品牌审查完成\n\n**整体评估**: 良好\n\n**优点**:\n✓ 色彩搭配专业\n✓ 图形简洁有力\n\n**改进建议**:\n1. 优化字体间距\n2. 增加在不同背景下的可读性测试'),

('marketing-content-writer', '营销文案专家', 'Marketing Content Writer',
 '创作高质量营销文案，包括广告、社交媒体、产品描述等。',
 'marketing', ARRAY['文案创作', '营销', 'SEO', '广告'], 1, 'claude', 14.90, 60, 97.80, 4.7,
 '为一款新的智能手表撰写产品描述，突出健康监测功能。',
 '✅ 营销文案已完成\n\n**产品名称**: 智护手表 Pro\n\n**产品描述**:\n\n👆 守护每一刻心跳\n\n智护手表 Pro，您的24小时健康管家。搭载全新一代生物传感芯片，精准监测心率、血氧、睡眠质量，让健康数据触手可及。\n\n【核心功能】\n• 实时心率监测\n• 血氧饱和度检测\n• 深度睡眠分析\n• 7天超长续航')
ON CONFLICT (agent_id) DO NOTHING;

-- 更新 Agent 分类 ID
UPDATE agents a SET category_id = c.id FROM categories c WHERE a.category = c.slug;

-- 插入管理员账户 (密码: admin123)
INSERT INTO users (email, username, password_hash, role, balance, email_verified) VALUES
('admin@agenthub.com', 'admin', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4.XQoTcW1HqvQPRu', 'admin', 10000.00, true)
ON CONFLICT (email) DO NOTHING;

COMMIT;

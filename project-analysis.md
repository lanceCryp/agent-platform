# AI Agent 服务平台 - 项目分析文档

## 1. 项目概述

### 1.1 项目背景

**项目名称**: AgentHub (智匠云)  
**项目类型**: SaaS 平台 / AI 服务  
**核心功能**: 提供 191 个专业 AI Agent 的按需调用服务  
**目标用户**: 中小企业、初创公司、个体开发者

### 1.2 核心价值主张

```
传统方式: 招聘专业人才 ¥30,000/月
       或: 外包服务 ¥5,000-50,000/项目

我们的方式: 订阅 AI Agent 服务 ¥99-599/月
           节省 90%+ 成本
```

### 1.3 产品形态

| 形态 | 说明 |
|------|------|
| **Web 应用** | 主要入口，用户浏览、下单、管理 |
| **API** | 供开发者集成到自己的系统 |
| **SDK** | 多语言客户端库 |

---

## 2. 功能需求分析

### 2.1 用户端功能

#### 2.1.1 用户注册与认证

```
功能点:
├── 用户注册 (邮箱/手机/微信/Google)
├── 用户登录 (密码/验证码/第三方OAuth)
├── 密码找回
├── 个人信息管理
├── 账户安全设置 (2FA/登录历史)
└── 会员等级 (免费/付费VIP/企业)

数据模型:
├── User {
│   ├── id: UUID
│   ├── email: string (唯一)
│   ├── username: string
│   ├── password_hash: string
│   ├── phone: string?
│   ├── avatar_url: string?
│   ├── role: enum (user/vip/admin)
│   ├── balance: decimal
│   ├── created_at: timestamp
│   └── updated_at: timestamp
│   }
│
└── UserSession {
    ├── id: UUID
    ├── user_id: FK
    ├── token: string
    ├── device_info: string
    ├── ip_address: string
    ├── expires_at: timestamp
    └── created_at: timestamp
    }
```

#### 2.1.2 Agent 市场

```
功能点:
├── Agent 分类浏览 (engineering/marketing/design/...)
├── Agent 搜索 (名称/标签/描述)
├── Agent 详情页 (介绍/价格/评价/示例)
├── 热门 Agent 推荐
├── Agent 收藏/关注
└── Agent 使用统计

数据模型:
├── Agent {
│   ├── id: UUID
│   ├── agent_id: string (唯一标识, 如 "design-brand-guardian")
│   ├── name: string (中文名)
│   ├── description: text
│   ├── category: enum
│   ├── tags: string[]
│   ├── tier: int (1/2/3/4)
│   ├── runtime_type: enum (openclaw/claude/openai)
│   ├── price_per_request: decimal
│   ├── price_per_token: decimal?
│   ├── avg_duration_seconds: int
│   ├── success_rate: decimal
│   ├── rating: decimal (1-5)
│   ├── total_tasks: int
│   ├── is_active: boolean
│   ├── config: json (Agent配置)
│   └── metadata: json (扩展信息)
│   }
│
└── AgentCategory {
    ├── id: UUID
    ├── name: string
    ├── slug: string
    ├── icon: string
    ├── description: text
    ├── sort_order: int
    └── agent_count: int
    }
```

#### 2.1.3 任务执行

```
功能点:
├── 创建任务 (选择Agent + 输入Prompt)
├── 任务状态跟踪 (pending/processing/completed/failed)
├── 任务历史记录
├── 任务结果查看
├── 任务重新执行
├── 任务取消
└── 批量任务创建

数据模型:
├── Task {
│   ├── id: UUID
│   ├── user_id: FK
│   ├── agent_id: FK
│   ├── prompt: text
│   ├── result: text?
│   ├── status: enum
│   ├── priority: int (1-10)
│   ├── cost: decimal
│   ├── tokens_used: int
│   ├── duration_seconds: int
│   ├── error_message: text?
│   ├── retry_count: int
│   ├── max_retries: int
│   ├── context: json?
│   ├── created_at: timestamp
│   ├── started_at: timestamp?
│   └── completed_at: timestamp?
│   }
│
└── TaskMessage {
    ├── id: UUID
    ├── task_id: FK
    ├── role: enum (user/assistant/system)
    ├── content: text
    └── created_at: timestamp
    }
```

#### 2.1.4 订阅与计费

```
功能点:
├── 订阅套餐展示
├── 订阅购买 (月/年)
├── 余额充值
├── 消费记录查看
├── 发票申请
├── 优惠券使用
└── 退款处理

数据模型:
├── Subscription {
│   ├── id: UUID
│   ├── user_id: FK
│   ├── plan_type: enum (basic/pro/enterprise)
│   ├── status: enum (active/cancelled/expired)
│   ├── start_date: date
│   ├── end_date: date
│   ├── price: decimal
│   ├── billing_cycle: enum (monthly/yearly)
│   └── auto_renew: boolean
│   }
│
├── Transaction {
│   ├── id: UUID
│   ├── user_id: FK
│   ├── type: enum (recharge/consumption/refund/subscription)
│   ├── amount: decimal
│   ├── balance_before: decimal
│   ├── balance_after: decimal
│   ├── description: string
│   ├── reference_id: UUID? (关联任务/订阅ID)
│   └── created_at: timestamp
│   }
│
└── Plan {
    ├── id: UUID
    ├── name: string
    ├── type: enum
    ├── price_monthly: decimal
    ├── price_yearly: decimal
    ├── features: json
    ├── task_limit: int? (null=无限)
    ├── agent_tier_limit: int? (允许使用的最高Tier)
    └── is_active: boolean
    }
```

### 2.2 管理端功能

#### 2.2.1 Agent 管理

```
功能点:
├── Agent 列表与筛选
├── Agent 添加/编辑/下架
├── Agent 定价调整
├── Agent 质量监控
├── Agent 自动分级
└── Agent 使用统计

权限: admin
```

#### 2.2.2 用户管理

```
功能点:
├── 用户列表与筛选
├── 用户详情查看
├── 用户余额调整
├── 用户权限管理
├── 用户封禁/解封
└── 用户消费统计

权限: admin
```

#### 2.2.3 运营管理

```
功能点:
├── 订阅订单管理
├── 消费记录对账
├── 客服工单处理
├── 系统公告发布
├── 促销活动配置
└── 数据报表导出

权限: admin/operator
```

---

## 3. 非功能需求分析

### 3.1 性能需求

| 指标 | 目标值 | 说明 |
|------|--------|------|
| **API 响应时间** | P95 < 200ms | 95% 请求在 200ms 内响应 |
| **任务执行时间** | P95 < 60s | Agent 执行 95% 在 60s 内完成 |
| **并发用户数** | 1,000+ | 同时在线用户 |
| **任务并发数** | 100+ | 同时执行的任务数 |
| **系统可用性** | 99.9% | 每月停机时间 < 45 分钟 |

### 3.2 安全需求

| 需求 | 说明 |
|------|------|
| **认证** | JWT + refresh token，7天有效期 |
| **授权** | RBAC 角色权限控制 |
| **数据加密** | TLS 1.3 传输，AES-256 存储 |
| **输入验证** | 所有用户输入严格校验 |
| **SQL注入防护** | 参数化查询 |
| **XSS防护** | 输出编码 |
| **CSRF防护** | Token 验证 |

### 3.3 可用性需求

| 需求 | 说明 |
|------|------|
| **多设备支持** | 桌面/平板/手机 |
| **浏览器兼容** | Chrome/Firefox/Safari/Edge 最新两个版本 |
| **无障碍** | 符合 WCAG 2.1 AA 标准 |

---

## 4. 技术架构设计

### 4.1 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                           用户层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   Web App   │  │  Mobile H5  │  │   API/SDK   │          │
│  │  (Next.js)  │  │   (React)   │  │  (REST/gRPC) │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼────────────────┼────────────────┼─────────────────────┘
          │                │                │
          └────────────────┴────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         网关层 (Kong)                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   限流      │  │   认证      │  │   监控      │          │
│  │ 100 req/min │  │   JWT验证   │  │   链路追踪   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         服务层 (Go微服务)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  User Svc  │  │  Task Svc   │  │  Agent Svc  │          │
│  │  :8081     │  │  :8082     │  │  :8083     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Billing Svc │  │ Notif Svc  │  │  Review Svc │          │
│  │  :8084     │  │  :8085     │  │  :8086     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         消息队列 (Kafka)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ task.created │  │task.completed│  │ billing.event│          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Agent 运行时层                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              Agent Worker Pool                          │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐         │  │
│  │  │ OpenClaw  │  │  Claude   │  │  OpenAI   │         │  │
│  │  │  Worker   │  │  Worker   │  │  Worker   │         │  │
│  │  └────────────┘  └────────────┘  └────────────┘         │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         数据层                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ PostgreSQL   │  │    Redis    │  │    OSS     │          │
│  │  主数据库    │  │    缓存     │  │   文件存储  │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 核心流程

#### 4.2.1 任务执行流程

```
1. 用户发起任务请求
   POST /api/v1/tasks
   
2. API Gateway 验证 JWT，限流
   
3. Task Service 接收请求
   ├── 创建 Task 记录 (status=pending)
   ├── 计算预估成本
   ├── 检查用户余额
   └── 发布消息到 Kafka (topic: task.created)
   
4. Task Service 立即返回
   HTTP 202 { task_id: "xxx", status: "pending" }
   
5. Agent Worker 消费消息
   ├── 选择合适的运行时 (OpenClaw/Claude/OpenAI)
   ├── 创建沙箱环境
   ├── 执行 Agent
   ├── 捕获输出
   └── 清理环境
   
6. 结果处理
   ├── 更新 Task 记录 (status=completed, result=xxx)
   ├── 计算实际成本
   ├── 扣减用户余额
   └── 发送 WebSocket 通知用户
   
7. 用户通过 WebSocket/轮询获取结果
   GET /api/v1/tasks/{id}
```

#### 4.2.2 订阅购买流程

```
1. 用户选择套餐
   POST /api/v1/subscriptions
   
2. 创建订单
   ├── 生成订单号
   ├── 记录订阅信息 (status=pending)
   └── 返回支付链接
   
3. 第三方支付 (Stripe/支付宝)
   ├── 用户完成支付
   └── 支付平台回调通知
   
4. 订阅生效
   ├── 更新订单状态 (status=paid)
   ├── 更新用户角色 (role=vip)
   ├── 开通套餐权限
   └── 发送欢迎邮件
```

### 4.3 数据库设计

```sql
-- 用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    role VARCHAR(20) DEFAULT 'user',
    balance DECIMAL(10,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Agent 表
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    name_en VARCHAR(200),
    description TEXT,
    category VARCHAR(50) NOT NULL,
    tags TEXT[],
    tier INT DEFAULT 3,
    runtime_type VARCHAR(20) NOT NULL,
    price_per_request DECIMAL(10,2) NOT NULL,
    price_per_token DECIMAL(10,6),
    avg_duration_seconds INT DEFAULT 60,
    success_rate DECIMAL(5,2) DEFAULT 100.00,
    rating DECIMAL(2,1) DEFAULT 5.0,
    total_tasks INT DEFAULT 0,
    total_success INT DEFAULT 0,
    total_failed INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    config JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agents_category ON agents(category);
CREATE INDEX idx_agents_tier ON agents(tier);
CREATE INDEX idx_agents_runtime ON agents(runtime_type);

-- 任务表
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    agent_id UUID NOT NULL REFERENCES agents(id),
    prompt TEXT NOT NULL,
    result TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    priority INT DEFAULT 5,
    cost DECIMAL(10,2) DEFAULT 0.00,
    tokens_used INT DEFAULT 0,
    duration_seconds INT DEFAULT 0,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    context JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_tasks_user ON tasks(user_id);
CREATE INDEX idx_tasks_agent ON tasks(agent_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created ON tasks(created_at DESC);

-- 订阅表
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    plan_id UUID NOT NULL,
    plan_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    billing_cycle VARCHAR(20) NOT NULL,
    auto_renew BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 交易记录表
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    balance_before DECIMAL(10,2) NOT NULL,
    balance_after DECIMAL(10,2) NOT NULL,
    description VARCHAR(500),
    reference_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_trans_user ON transactions(user_id);
CREATE INDEX idx_trans_type ON transactions(type);
```

### 4.4 API 设计

#### 4.4.1 用户认证

```yaml
# POST /api/v1/auth/register - 用户注册
Request:
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword123"
}
Response:
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe"
    },
    "token": "jwt_token",
    "refresh_token": "refresh_token"
  }
}

# POST /api/v1/auth/login - 用户登录
Request:
{
  "email": "user@example.com",
  "password": "securepassword123"
}
Response:
{
  "success": true,
  "data": {
    "token": "jwt_token",
    "expires_in": 3600
  }
}
```

#### 4.4.2 Agent 管理

```yaml
# GET /api/v1/agents - 获取Agent列表
Query: ?category=engineering&tier=1&page=1&limit=20
Response:
{
  "success": true,
  "data": {
    "total": 30,
    "page": 1,
    "limit": 20,
    "agents": [
      {
        "id": "uuid",
        "agent_id": "engineering-code-reviewer",
        "name": "代码审查专家",
        "category": "engineering",
        "tier": 1,
        "price_per_request": 9.90,
        "rating": 4.8,
        "total_tasks": 1500,
        "is_active": true
      }
    ]
  }
}

# GET /api/v1/agents/{agent_id} - 获取Agent详情
Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "agent_id": "engineering-code-reviewer",
    "name": "代码审查专家",
    "description": "...",
    "category": "engineering",
    "tags": ["代码", "审查", "质量"],
    "tier": 1,
    "runtime_type": "claude",
    "price_per_request": 9.90,
    "avg_duration_seconds": 45,
    "success_rate": 98.5,
    "rating": 4.8,
    "input_example": "请帮我审查这段代码的潜在问题...",
    "output_example": "代码审查完成，发现3个问题..."
  }
}
```

#### 4.4.3 任务执行

```yaml
# POST /api/v1/tasks - 创建任务
Headers: Authorization: Bearer {token}
Request:
{
  "agent_id": "engineering-code-reviewer",
  "prompt": "请帮我审查这段Python代码的潜在问题:\n\ndef calculate(a, b):\n    return a / b",
  "priority": 5,
  "max_retries": 3
}
Response (202 Accepted):
{
  "success": true,
  "data": {
    "task_id": "uuid",
    "status": "pending",
    "estimated_cost": 0.50,
    "estimated_duration": 45
  }
}

# GET /api/v1/tasks/{task_id} - 获取任务状态
Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "agent_id": "engineering-code-reviewer",
    "status": "completed",
    "result": "代码审查完成，发现1个潜在问题...",
    "cost": 0.45,
    "tokens_used": 1200,
    "duration_seconds": 38,
    "created_at": "2026-04-10T09:00:00Z",
    "completed_at": "2026-04-10T09:00:38Z"
  }
}

# GET /api/v1/tasks - 获取任务列表
Query: ?status=completed&page=1&limit=20
Response:
{
  "success": true,
  "data": {
    "total": 50,
    "tasks": [...]
  }
}
```

#### 4.4.4 订阅计费

```yaml
# GET /api/v1/plans - 获取套餐列表
Response:
{
  "success": true,
  "data": {
    "plans": [
      {
        "id": "uuid",
        "name": "基础版",
        "type": "basic",
        "price_monthly": 99.00,
        "price_yearly": 950.00,
        "features": {
          "task_limit": 100,
          "agent_tier_limit": 2,
          "api_access": false
        }
      }
    ]
  }
}

# POST /api/v1/subscriptions - 购买订阅
Request:
{
  "plan_id": "uuid",
  "billing_cycle": "monthly"
}
Response:
{
  "success": true,
  "data": {
    "subscription_id": "uuid",
    "payment_url": "https://..."
  }
}

# GET /api/v1/users/me/balance - 获取余额
Response:
{
  "success": true,
  "data": {
    "balance": 150.00,
    "subscription": {
      "plan_type": "pro",
      "end_date": "2026-05-10"
    }
  }
}
```

---

## 5. 项目结构

### 5.1 整体目录结构

```
agenthub/
├── frontend/                    # Next.js 前端
│   ├── app/                  # App Router
│   ├── components/            # React 组件
│   ├── hooks/                 # 自定义 Hooks
│   ├── lib/                   # 工具库
│   ├── stores/               # Zustand stores
│   ├── types/                # TypeScript 类型
│   └── public/               # 静态资源
│
├── backend/                   # Go 后端
│   ├── cmd/                   # 入口文件
│   │   ├── api-gateway/
│   │   ├── user-service/
│   │   ├── task-service/
│   │   ├── agent-service/
│   │   └── billing-service/
│   ├── pkg/                   # 公共库
│   │   ├── config/
│   │   ├── database/
│   │   ├── cache/
│   │   ├── queue/
│   │   ├── auth/
│   │   └── middleware/
│   ├── proto/                 # gRPC protobuf
│   └── migrations/            # 数据库迁移
│
├── agents/                    # 191个Agent配置
│   ├── engineering/
│   ├── design/
│   ├── marketing/
│   └── ...
│
├── infra/                     # 基础设施
│   ├── kubernetes/
│   ├── docker/
│   └── monitoring/
│
├── docs/                      # 文档
├── scripts/                    # 脚本
├── Makefile
└── README.md
```

### 5.2 前端项目结构

```
frontend/
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   └── forgot-password/page.tsx
│   ├── (dashboard)/
│   │   ├── layout.tsx
│   │   ├── agents/
│   │   │   ├── page.tsx              # Agent 市场
│   │   │   └── [agentId]/page.tsx    # Agent 详情
│   │   ├── tasks/
│   │   │   ├── page.tsx              # 任务列表
│   │   │   └── [taskId]/page.tsx     # 任务详情
│   │   ├── billing/
│   │   │   ├── page.tsx             # 账单
│   │   │   └── plans/page.tsx       # 套餐
│   │   └── settings/
│   │       ├── page.tsx             # 设置
│   │       └── profile/page.tsx     # 个人信息
│   ├── page.tsx                      # 首页
│   ├── layout.tsx                    # 根布局
│   └── globals.css
│
├── components/
│   ├── ui/                           # shadcn/ui 组件
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── input.tsx
│   │   └── ...
│   ├── agents/
│   │   ├── agent-card.tsx
│   │   ├── agent-list.tsx
│   │   └── agent-detail.tsx
│   ├── tasks/
│   │   ├── task-card.tsx
│   │   ├── task-form.tsx
│   │   └── task-result.tsx
│   └── layout/
│       ├── header.tsx
│       ├── sidebar.tsx
│       └── footer.tsx
│
├── hooks/
│   ├── use-auth.ts
│   ├── use-agents.ts
│   ├── use-tasks.ts
│   └── use-subscription.ts
│
├── lib/
│   ├── api.ts                        # API 客户端
│   ├── auth.ts                       # 认证
│   ├── utils.ts                      # 工具函数
│   └── constants.ts                  # 常量
│
├── stores/
│   ├── auth-store.ts
│   ├── agent-store.ts
│   └── task-store.ts
│
└── types/
    ├── api.ts
    ├── agent.ts
    └── user.ts
```

### 5.3 后端项目结构

```
backend/
├── cmd/
│   ├── api-gateway/
│   │   ├── main.go
│   │   ├── handler/
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   └── health.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       ├── ratelimit.go
│   │       └── logging.go
│   │
│   ├── user-service/
│   │   ├── main.go
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   │
│   ├── task-service/
│   │   ├── main.go
│   │   ├── handler/
│   │   │   ├── task.go
│   │   │   └── result.go
│   │   ├── service/
│   │   │   ├── task.go
│   │   │   └── scheduler.go
│   │   ├── worker/
│   │   │   ├── pool.go
│   │   │   └── executor.go
│   │   └── repository/
│   │
│   └── agent-service/
│       ├── main.go
│       ├── handler/
│       ├── service/
│       └── runtime/
│           ├── openclaw.go
│           ├── claude.go
│           └── openai.go
│
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── postgres.go
│   │   └── migrations/
│   ├── cache/
│   │   └── redis.go
│   ├── queue/
│   │   ├── kafka.go
│   │   └── consumer.go
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   ├── errors/
│   │   └── errors.go
│   └── logger/
│       └── logger.go
│
└── proto/
    ├── user.proto
    ├── task.proto
    └── agent.proto
```

---

## 6. 开发计划

### 6.1 Sprint 规划

#### Sprint 1: 项目初始化 (2周)

```
目标: 搭建完整开发环境，实现基础框架

Week 1:
├── [ ] Git 仓库初始化
├── [ ] 前端: Next.js 16 项目创建
├── [ ] 后端: Go 项目结构搭建
├── [ ] 数据库: PostgreSQL Schema 创建
├── [ ] CI/CD: GitHub Actions 配置
└── [ ] 开发规范: ESLint/Prettier/Commit规范

Week 2:
├── [ ] 用户认证: 注册/登录/JWT
├── [ ] 用户模块: CRUD + 积分
├── [ ] 前端: 登录/注册页面
├── [ ] 前端: 用户设置页面
└── [ ] 联调: 前后端对接
```

#### Sprint 2: 核心功能 MVP (2周)

```
目标: 实现 Agent 市场和任务执行

Week 3:
├── [ ] Agent 模块: Schema + API
├── [ ] Agent 运行时: OpenClaw 集成
├── [ ] 任务模块: 创建/状态/结果
├── [ ] 前端: Agent 市场列表
└── [ ] 前端: Agent 详情页

Week 4:
├── [ ] 任务执行: 异步队列
├── [ ] WebSocket: 实时结果推送
├── [ ] 前端: 创建任务页面
├── [ ] 前端: 任务列表/详情
└── [ ] 联调 + Bug修复
```

#### Sprint 3: 订阅计费 (2周)

```
目标: 实现订阅和计费系统

Week 5:
├── [ ] 订阅模块: Schema + API
├── [ ] 支付集成: Stripe/支付宝沙箱
├── [ ] 余额系统: 充值/消费
├── [ ] 前端: 套餐展示/购买
└── [ ] 前端: 余额充值

Week 6:
├── [ ] 交易记录: 完整流水
├── [ ] 优惠码系统
├── [ ] 发票申请
├── [ ] 管理后台: 订单管理
└── [ ] 完整测试 + 文档
```

#### Sprint 4: 优化上线 (2周)

```
目标: 性能优化 + 部署上线

Week 7:
├── [ ] 性能优化: 缓存/索引
├── [ ] Agent 质量监控
├── [ ] 自动扩容配置
├── [ ] 监控大盘配置
└── [ ] 安全审计

Week 8:
├── [ ] 域名/SSL 配置
├── [ ] 正式环境部署
├── [ ] 数据迁移
├── [ ] 监控告警测试
└── [ ] 正式上线
```

### 6.2 里程碑

| 里程碑 | 时间 | 交付物 |
|--------|------|--------|
| M1: MVP | 第4周末 | 基础功能可用 |
| M2: Beta | 第8周末 | 付费功能完整 |
| M3: Launch | 第12周末 | 正式上线 |

---

## 7. 风险评估

### 7.1 技术风险

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|---------|
| Agent 输出质量不稳定 | 高 | 中 | 内置 Reviewer Agent 质检 |
| 高并发性能瓶颈 | 高 | 中 | K8s 自动扩容，预估容量 |
| 多运行时兼容问题 | 中 | 中 | 统一接口抽象，逐步集成 |
| API 成本超支 | 中 | 低 | 缓存优化，按需扩容 |

### 7.2 业务风险

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|---------|
| 用户留存低 | 高 | 中 | 持续优化用户体验 |
| 竞品模仿 | 中 | 高 | 快速迭代，建立壁垒 |
| 监管政策变化 | 低 | 低 | 关注政策，合规运营 |

---

## 8. 成功指标

### 8.1 产品指标

| 指标 | 第1月 | 第3月 | 第6月 |
|------|-------|-------|-------|
| 注册用户 | 100 | 1,000 | 10,000 |
| 付费用户 | 10 | 100 | 1,000 |
| 付费转化率 | 10% | 10% | 10% |
| 月活跃用户 | 50 | 500 | 5,000 |

### 8.2 技术指标

| 指标 | 目标值 |
|------|--------|
| API 可用性 | 99.9% |
| 平均响应时间 | < 200ms |
| 任务成功率 | > 95% |
| 崩溃率 | < 0.1% |

### 8.3 商业指标

| 指标 | 第1月 | 第3月 | 第6月 |
|------|-------|-------|-------|
| MRR | ¥5,000 | ¥50,000 | ¥500,000 |
| ARPU | ¥500 | ¥500 | ¥500 |
| CAC | ¥200 | ¥150 | ¥100 |
| LTV | ¥1,500 | ¥3,000 | ¥5,000 |

---

**文档版本**: v1.0  
**最后更新**: 2026-04-10  
**负责人**: AgentHub Team

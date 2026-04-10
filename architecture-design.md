# AI Agent 服务平台 - 软件架构设计

## 1. 架构总览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              接入层 (Access Layer)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Web App    │  │  Mobile App  │  │   API/SDK    │  │   Webhook    │   │
│  │  (React)     │  │  (小程序)     │  │  (REST/gRPC) │  │              │   │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │
└─────────┼────────────────┼────────────────┼────────────────┼─────────────┘
          │                │                │                │
          └────────────────┴────────────────┴────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           网关层 (Gateway Layer)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                        API Gateway (Kong/AWS API Gateway)             │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │  │
│  │  │  Rate Limit  │  │   Auth/JWT   │  │    HTTPS     │  │  Logging  │ │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └───────────┘ │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           服务层 (Service Layer)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌───────────┐  │
│  │  User Service  │  │  Task Service  │  │ Billing Svc    │  │  Notif    │  │
│  │  (用户管理)     │  │  (任务调度)     │  │ (计费服务)     │  │ (通知)    │  │
│  └───────┬────────┘  └───────┬────────┘  └───────┬────────┘  └─────┬─────┘  │
│          │                   │                   │               │        │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌───────────┐  │
│  │  Agent Service │  │  Knowledge Svc │  │  Review Svc    │  │  Analytics│  │
│  │  (Agent管理)   │  │  (知识库)       │  │  (质检审核)    │  │ (分析)    │  │
│  └───────┬────────┘  └───────┬────────┘  └───────┬────────┘  └─────┬─────┘  │
└──────────┼───────────────────┼───────────────────┼────────────────┼──────────┘
           │                   │                   │                │
           └───────────────────┴───────────────────┴────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        消息队列层 (Message Queue)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                    Apache Kafka / RabbitMQ / AWS SQS                  │  │
│  │                                                                       │  │
│  │   Topics:                                                             │  │
│  │   ├─ task.created      (任务创建)                                      │  │
│  │   ├─ task.assigned     (任务分配)                                      │  │
│  │   ├─ task.completed    (任务完成)                                      │  │
│  │   ├─ agent.available   (Agent可用)                                    │  │
│  │   └─ billing.events    (计费事件)                                      │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Agent 执行层 (Agent Runtime)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                     Agent Scheduler (调度器)                          │  │
│  │                                                                       │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────────┐  │  │
│  │  │  Worker 1  │  │  Worker 2  │  │  Worker 3  │  │   Worker N     │  │  │
│  │  │ (engineer) │  │  (design)  │  │ (marketing)│  │  (academic)    │  │  │
│  │  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘  └────────┬───────┘  │  │
│  │        │               │               │                  │         │  │
│  │        └───────────────┴───────────────┴──────────────────┘         │  │
│  │                              │                                       │  │
│  │                              ▼                                       │  │
│  │  ┌──────────────────────────────────────────────────────────────┐   │  │
│  │  │              Runtime Pool (执行运行时)                        │   │  │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐            │   │  │
│  │  │  │  OpenClaw  │  │ Claude Code│  │  OpenAI    │            │   │  │
│  │  │  │  Runtime   │  │  Runtime   │  │  API       │            │   │  │
│  │  │  └────────────┘  └────────────┘  └────────────┘            │   │  │
│  │  └──────────────────────────────────────────────────────────────┘   │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         数据层 (Data Layer)                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  PostgreSQL  │  │    Redis     │  │ Elasticsearch│  │    MinIO     │   │
│  │  (主数据库)   │  │   (缓存)     │  │   (搜索)     │  │  (对象存储)   │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      基础设施层 (Infrastructure)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  Kubernetes  │  │  Docker      │  │   Nginx      │  │  Prometheus  │   │
│  │  (容器编排)   │  │  (容器化)    │  │  (反向代理)   │  │  (监控)      │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Grafana    │  │    Jaeger    │  │    Vault     │  │    Consul    │   │
│  │  (可视化)    │  │  (链路追踪)   │  │  (密钥管理)   │  │  (服务发现)   │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. 核心组件详解

### 2.1 接入层 (Access Layer)

#### 2.1.1 Web 应用 (Next.js 14)

```
frontend/
├── app/                          # Next.js App Router
│   ├── page.tsx                  # 首页
│   ├── agents/                   # Agent 市场
│   │   ├── page.tsx              # Agent 列表
│   │   └── [id]/                 # Agent 详情
│   │       └── page.tsx
│   ├── tasks/                    # 任务管理
│   │   ├── page.tsx              # 任务列表
│   │   └── [id]/                 # 任务详情
│   │       └── page.tsx
│   ├── billing/                  # 计费管理
│   │   └── page.tsx
│   └── dashboard/                # 用户仪表盘
│       └── page.tsx
├── components/                   # React 组件
│   ├── ui/                       # 基础 UI 组件
│   ├── agents/                   # Agent 相关组件
│   ├── tasks/                    # 任务相关组件
│   └── layout/                   # 布局组件
├── lib/                          # 工具库
│   ├── api.ts                    # API 客户端
│   ├── auth.ts                   # 认证工具
│   └── utils.ts                  # 通用工具
└── styles/                       # 样式文件
    └── globals.css
```

**技术选型**:
- **框架**: Next.js 14 (App Router)
- **语言**: TypeScript
- **样式**: Tailwind CSS
- **UI 库**: shadcn/ui
- **状态管理**: Zustand
- **数据获取**: TanStack Query

#### 2.1.2 API 设计 (REST + gRPC)

```yaml
# OpenAPI 3.0 规范
openapi: 3.0.0
info:
  title: AgentHub API
  version: 1.0.0
paths:
  # Agent 管理
  /api/v1/agents:
    get:
      summary: 获取 Agent 列表
      parameters:
        - name: category
          in: query
          schema:
            type: string
            enum: [engineering, design, marketing, ...]
      responses:
        '200':
          description: Agent 列表
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Agent'

  /api/v1/agents/{id}:
    get:
      summary: 获取 Agent 详情
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Agent 详情

  # 任务管理
  /api/v1/tasks:
    post:
      summary: 创建任务
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                agent_id:
                  type: string
                prompt:
                  type: string
                priority:
                  type: string
                  enum: [low, medium, high]
      responses:
        '201':
          description: 任务创建成功

  /api/v1/tasks/{id}:
    get:
      summary: 获取任务状态
      responses:
        '200':
          description: 任务详情
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Task'

  # 用户管理
  /api/v1/users/me:
    get:
      summary: 获取当前用户信息
      security:
        - bearerAuth: []
      responses:
        '200':
          description: 用户信息

  /api/v1/users/me/balance:
    get:
      summary: 获取账户余额
      security:
        - bearerAuth: []
      responses:
        '200':
          description: 余额信息

components:
  schemas:
    Agent:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
        description:
          type: string
        category:
          type: string
        price_per_request:
          type: number
        avg_rating:
          type: number
        total_tasks:
          type: integer

    Task:
      type: object
      properties:
        id:
          type: string
        agent_id:
          type: string
        status:
          type: string
          enum: [pending, processing, completed, failed]
        prompt:
          type: string
        result:
          type: string
        cost:
          type: number
        created_at:
          type: string
          format: date-time
        completed_at:
          type: string
          format: date-time

  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

---

### 2.2 服务层 (Service Layer)

#### 2.2.1 微服务架构

```
services/
├── user-service/                 # 用户服务
│   ├── src/
│   │   ├── handlers/
│   │   │   ├── auth.go           # 认证处理
│   │   │   ├── user.go           # 用户管理
│   │   │   └── billing.go        # 计费管理
│   │   ├── models/
│   │   │   └── user.go           # 数据模型
│   │   ├── repository/
│   │   │   └── user_repo.go      # 数据访问
│   │   └── main.go               # 服务入口
│   ├── Dockerfile
│   └── go.mod
│
├── task-service/                 # 任务服务
│   ├── src/
│   │   ├── handlers/
│   │   │   ├── task.go           # 任务管理
│   │   │   └── queue.go          # 队列管理
│   │   ├── scheduler/
│   │   │   └── scheduler.go      # 任务调度器
│   │   └── main.go
│   ├── Dockerfile
│   └── go.mod
│
├── agent-service/                # Agent 服务
│   ├── src/
│   │   ├── handlers/
│   │   │   ├── agent.go          # Agent 管理
│   │   │   └── runtime.go        # 运行时管理
│   │   ├── runtime/
│   │   │   ├── openclaw.go       # OpenClaw 集成
│   │   │   ├── claude.go         # Claude Code 集成
│   │   │   └── openai.go         # OpenAI API 集成
│   │   └── main.go
│   ├── Dockerfile
│   └── go.mod
│
├── billing-service/              # 计费服务
│   ├── src/
│   │   ├── handlers/
│   │   │   ├── billing.go        # 计费处理
│   │   │   └── subscription.go   # 订阅管理
│   │   └── main.go
│   ├── Dockerfile
│   └── go.mod
│
└── notification-service/         # 通知服务
    ├── src/
    │   ├── handlers/
    │   │   ├── email.go          # 邮件通知
    │   │   ├── sms.go            # 短信通知
    │   │   └── webhook.go        # Webhook 推送
    │   └── main.go
    ├── Dockerfile
    └── go.mod
```

**技术选型**:
- **语言**: Go (高并发、低延迟)
- **框架**: Gin / Echo
- **通信**: gRPC (服务间) + REST (外部)
- **服务发现**: Consul / etcd

#### 2.2.2 任务调度器设计

```go
// task-service/scheduler/scheduler.go
package scheduler

import (
    "context"
    "log"
    "time"
    
    "github.com/robfig/cron/v3"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
    queue     TaskQueue          // 任务队列
    workers   map[string]*Worker // Agent 工作池
    cron      *cron.Cron         // 定时任务
    metrics   MetricsCollector   // 指标收集
}

// Worker Agent 工作节点
type Worker struct {
    ID         string
    AgentID    string
    Runtime    RuntimeType  // openclaw / claude / openai
    Status     WorkerStatus // idle / busy / offline
    Capacity   int          // 并发处理能力
    Queue      chan *Task   // 任务队列
}

// Task 任务定义
type Task struct {
    ID          string
    UserID      string
    AgentID     string
    Prompt      string
    Priority    int          // 1-10
    Status      TaskStatus   // pending / processing / completed / failed
    MaxRetries  int
    RetryCount  int
    CreatedAt   time.Time
    StartedAt   *time.Time
    CompletedAt *time.Time
    Cost        float64
    Result      string
    Error       string
}

// 调度算法
func (s *TaskScheduler) Schedule(ctx context.Context, task *Task) error {
    // 1. 获取可用 Worker
    worker := s.selectWorker(task.AgentID)
    if worker == nil {
        // 无可用 Worker，加入等待队列
        return s.queue.Push(ctx, task)
    }
    
    // 2. 发送到 Worker
    select {
    case worker.Queue <- task:
        task.Status = StatusProcessing
        task.StartedAt = time.Now()
        s.metrics.IncTaskStarted(task.AgentID)
        return nil
    case <-time.After(5 * time.Second):
        // Worker 繁忙，加入队列
        return s.queue.Push(ctx, task)
    }
}

// Worker 选择算法
func (s *TaskScheduler) selectWorker(agentID string) *Worker {
    var bestWorker *Worker
    minLoad := float64(999)
    
    for _, worker := range s.workers {
        // 筛选条件
        if worker.AgentID != agentID {
            continue
        }
        if worker.Status != StatusIdle {
            continue
        }
        
        // 负载计算
        load := float64(len(worker.Queue)) / float64(worker.Capacity)
        if load < minLoad {
            minLoad = load
            bestWorker = worker
        }
    }
    
    return bestWorker
}

// Worker 池管理
func (s *TaskScheduler) ScaleWorkers(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // 根据队列长度自动扩缩容
            queueLen := s.queue.Len()
            workerCount := len(s.workers)
            
            if queueLen > workerCount*10 {
                // 扩容
                s.addWorker(ctx)
            } else if queueLen < workerCount*2 && workerCount > 3 {
                // 缩容
                s.removeWorker(ctx)
            }
        }
    }
}
```

---

### 2.3 Agent 执行层 (Agent Runtime)

#### 2.3.1 OpenClaw 运行时集成

```go
// agent-service/runtime/openclaw.go
package runtime

import (
    "context"
    "encoding/json"
    "os/exec"
    "time"
)

// OpenClawRuntime OpenClaw 运行时
type OpenClawRuntime struct {
    workspace   string        // 工作目录
    timeout     time.Duration // 超时时间
    maxRetries  int           // 最大重试次数
}

// Execute 执行任务
func (r *OpenClawRuntime) Execute(ctx context.Context, agentID string, prompt string) (*ExecutionResult, error) {
    // 1. 准备执行环境
    workspace := r.prepareWorkspace(agentID)
    
    // 2. 构建执行命令
    cmd := exec.CommandContext(ctx, "openclaw", "agents", "run", agentID, "--prompt", prompt)
    cmd.Dir = workspace
    
    // 3. 设置超时
    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()
    
    // 4. 执行并捕获输出
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, &ExecutionError{
            AgentID: agentID,
            Prompt:  prompt,
            Output:  string(output),
            Error:   err.Error(),
        }
    }
    
    // 5. 解析结果
    return &ExecutionResult{
        AgentID:   agentID,
        Output:    string(output),
        Tokens:    r.extractTokenUsage(output),
        Cost:      r.calculateCost(output),
        Duration:  time.Since(start),
    }, nil
}

// 批量执行优化
func (r *OpenClawRuntime) ExecuteBatch(ctx context.Context, tasks []*Task) ([]*ExecutionResult, error) {
    results := make([]*ExecutionResult, len(tasks))
    errChan := make(chan error, len(tasks))
    
    // 并行执行
    for i, task := range tasks {
        go func(idx int, t *Task) {
            result, err := r.Execute(ctx, t.AgentID, t.Prompt)
            if err != nil {
                errChan <- err
                return
            }
            results[idx] = result
            errChan <- nil
        }(i, task)
    }
    
    // 等待所有任务完成
    for i := 0; i < len(tasks); i++ {
        if err := <-errChan; err != nil {
            return nil, err
        }
    }
    
    return results, nil
}
```

#### 2.3.2 Claude Code 运行时集成

```go
// agent-service/runtime/claude.go
package runtime

import (
    "context"
    "fmt"
    "os/exec"
    "path/filepath"
)

// ClaudeRuntime Claude Code 运行时
type ClaudeRuntime struct {
    agentsDir   string        // Agent 配置目录 ~/.claude/agents/
    timeout     time.Duration
}

// Execute 执行 Claude Agent
func (r *ClaudeRuntime) Execute(ctx context.Context, agentID string, prompt string) (*ExecutionResult, error) {
    // 1. 验证 Agent 存在
    agentFile := filepath.Join(r.agentsDir, agentID+".md")
    if _, err := os.Stat(agentFile); os.IsNotExist(err) {
        return nil, fmt.Errorf("agent %s not found", agentID)
    }
    
    // 2. 构建命令
    cmd := exec.CommandContext(ctx, "claude",
        "--agent", agentID,
        "--bare",  // 跳过交互式提示
        "-p", prompt,
    )
    
    // 3. 执行
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, err
    }
    
    return &ExecutionResult{
        AgentID: agentID,
        Output:  string(output),
        Runtime: "claude",
    }, nil
}
```

#### 2.3.3 运行时选择策略

```go
// agent-service/runtime/manager.go
package runtime

// RuntimeManager 运行时管理器
type RuntimeManager struct {
    runtimes map[RuntimeType]Runtime
    strategy SelectionStrategy
}

// Runtime 运行时接口
type Runtime interface {
    Execute(ctx context.Context, agentID string, prompt string) (*ExecutionResult, error)
    Health() error
    Cost() float64
}

// SelectionStrategy 运行时选择策略
type SelectionStrategy int

const (
    StrategyCost      SelectionStrategy = iota  // 成本优先
    StrategySpeed                               // 速度优先
    StrategyQuality                             // 质量优先
    StrategyBalanced                            // 平衡模式
)

// SelectRuntime 选择最优运行时
func (m *RuntimeManager) SelectRuntime(agentID string, strategy SelectionStrategy) (Runtime, error) {
    candidates := []Runtime{
        m.runtimes[RuntimeOpenClaw],
        m.runtimes[RuntimeClaude],
        m.runtimes[RuntimeOpenAI],
    }
    
    switch strategy {
    case StrategyCost:
        // 选择成本最低的
        return m.selectByCost(candidates)
    case StrategySpeed:
        // 选择响应最快的
        return m.selectByLatency(candidates)
    case StrategyQuality:
        // 选择质量最高的
        return m.selectByQuality(agentID, candidates)
    case StrategyBalanced:
        // 综合评分
        return m.selectByScore(agentID, candidates)
    default:
        return candidates[0], nil
    }
}
```

---

### 2.4 数据层 (Data Layer)

#### 2.4.1 数据库 Schema 设计

```sql
-- PostgreSQL Schema

-- 用户表
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    username        VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    avatar_url      VARCHAR(500),
    role            VARCHAR(20) DEFAULT 'user', -- admin / user / vip
    status          VARCHAR(20) DEFAULT 'active',
    balance         DECIMAL(10,2) DEFAULT 0.00,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Agent 表
CREATE TABLE agents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id        VARCHAR(100) UNIQUE NOT NULL, -- 如: design-brand-guardian
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    category        VARCHAR(50) NOT NULL, -- engineering / design / marketing
    price_per_request DECIMAL(10,2) NOT NULL,
    price_per_token   DECIMAL(10,6),
    runtime_type    VARCHAR(20) NOT NULL, -- openclaw / claude / openai
    config          JSONB, -- Agent 配置
    is_active       BOOLEAN DEFAULT true,
    avg_rating      DECIMAL(2,1) DEFAULT 5.0,
    total_tasks     INTEGER DEFAULT 0,
    success_rate    DECIMAL(5,2) DEFAULT 100.00,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 任务表
CREATE TABLE tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    agent_id        VARCHAR(100) NOT NULL REFERENCES agents(agent_id),
    prompt          TEXT NOT NULL,
    result          TEXT,
    status          VARCHAR(20) DEFAULT 'pending', -- pending / processing / completed / failed
    priority        INTEGER DEFAULT 5, -- 1-10
    cost            DECIMAL(10,2) DEFAULT 0.00,
    tokens_used     INTEGER,
    runtime_type    VARCHAR(20),
    error_message   TEXT,
    retry_count     INTEGER DEFAULT 0,
    max_retries     INTEGER DEFAULT 3,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- 订阅表
CREATE TABLE subscriptions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    plan_type       VARCHAR(20) NOT NULL, -- basic / pro / enterprise
    status          VARCHAR(20) DEFAULT 'active', -- active / cancelled / expired
    start_date      TIMESTAMP NOT NULL,
    end_date        TIMESTAMP NOT NULL,
    price           DECIMAL(10,2) NOT NULL,
    billing_cycle   VARCHAR(20) NOT NULL, -- monthly / yearly
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 交易记录表
CREATE TABLE transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    type            VARCHAR(20) NOT NULL, -- recharge / consumption / refund
    amount          DECIMAL(10,2) NOT NULL,
    balance_after   DECIMAL(10,2) NOT NULL,
    description     VARCHAR(500),
    reference_id    UUID, -- 关联任务ID或订阅ID
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 知识库表
CREATE TABLE knowledge_base (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    agent_id        VARCHAR(100) REFERENCES agents(agent_id),
    title           VARCHAR(500) NOT NULL,
    content         TEXT NOT NULL,
    embedding       vector(1536), -- pgvector 扩展
    is_public       BOOLEAN DEFAULT false,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建 pgvector 扩展用于相似度搜索
CREATE EXTENSION IF NOT EXISTS vector;

-- 创建相似度搜索索引
CREATE INDEX idx_knowledge_embedding ON knowledge_base 
USING ivfflat (embedding vector_cosine_ops);
```

#### 2.4.2 Redis 缓存策略

```yaml
# Redis 缓存策略

# 1. 用户会话缓存
key: "session:{user_id}"
value: {user_id, email, role, permissions}
ttl: 86400  # 24小时

# 2. Agent 元数据缓存
key: "agent:{agent_id}"
value: {id, name, description, price, rating}
ttl: 3600   # 1小时

# 3. 任务状态缓存
key: "task:{task_id}"
value: {id, status, result, cost}
ttl: 1800   # 30分钟

# 4. 热门 Agent 列表
key: "agents:popular"
value: [{agent_id, name, total_tasks, rating}]
ttl: 300    # 5分钟

# 5. 用户余额缓存
key: "user:{user_id}:balance"
value: 123.45
ttl: 60     # 1分钟 (短TTL保证准确性)

# 6. 限流计数器
key: "rate_limit:{user_id}:{api_endpoint}"
value: 45   # 请求次数
ttl: 60     # 1分钟窗口

# 7. 分布式锁
key: "lock:task:{task_id}"
value: "worker_001"
ttl: 30     # 30秒自动释放
```

---

### 2.5 基础设施层 (Infrastructure)

#### 2.5.1 Kubernetes 部署配置

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: agenthub
  labels:
    name: agenthub

---
# k8s/deployment-api.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: agenthub
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: agenthub/api-gateway:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres-service
        - name: REDIS_HOST
          value: redis-service
        - name: KAFKA_BROKERS
          value: kafka-service:9092
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
# k8s/deployment-task-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-service
  namespace: agenthub
spec:
  replicas: 5
  selector:
    matchLabels:
      app: task-service
  template:
    metadata:
      labels:
        app: task-service
    spec:
      containers:
      - name: task-service
        image: agenthub/task-service:v1.0.0
        ports:
        - containerPort: 8081
        env:
        - name: KAFKA_BROKERS
          value: kafka-service:9092
        - name: MAX_WORKERS
          value: "10"
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"

---
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: task-service-hpa
  namespace: agenthub
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: task-service
  minReplicas: 5
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60

---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: api-gateway-service
  namespace: agenthub
spec:
  selector:
    app: api-gateway
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP

---
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agenthub-ingress
  namespace: agenthub
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.agenthub.com
    secretName: agenthub-tls
  rules:
  - host: api.agenthub.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-gateway-service
            port:
              number: 80
```

#### 2.5.2 Docker Compose 开发环境

```yaml
# docker-compose.yml
version: '3.8'

services:
  # API Gateway
  api-gateway:
    build: ./services/api-gateway
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - postgres
      - redis
      - kafka
    volumes:
      - ./services/api-gateway:/app
    command: air  # 热重载

  # Task Service
  task-service:
    build: ./services/task-service
    environment:
      - KAFKA_BROKERS=kafka:9092
      - DB_HOST=postgres
    depends_on:
      - postgres
      - kafka
    volumes:
      - ./services/task-service:/app
    command: air

  # PostgreSQL
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: agenthub
      POSTGRES_PASSWORD: agenthub123
      POSTGRES_DB: agenthub
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

  # Redis
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # Kafka
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    depends_on:
      - zookeeper

  # Frontend
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
    command: npm run dev

volumes:
  postgres_data:
  redis_data:
```

---

## 3. 关键设计决策

### 3.1 为什么选择微服务架构？

| 优势 | 说明 |
|------|------|
| **独立部署** | Agent 服务更新不影响其他服务 |
| **独立扩展** | Task Service 可单独扩容处理高峰期 |
| **技术异构** | Agent Runtime 可用 Python，其他用 Go |
| **故障隔离** | 单个服务崩溃不影响全局 |
| **团队并行** | 不同团队负责不同服务 |

### 3.2 为什么选择消息队列？

```
同步调用 vs 异步消息队列

同步调用:
用户 ──► API ──► Task Svc ──► Agent Runtime
     ◄── 等待 ◄── 等待 ◄── 等待
     
问题: 耦合度高、响应慢、无法削峰

异步消息队列:
用户 ──► API ──► Kafka ──► Task Svc ──► Agent Runtime
     ◄── 立即返回 Task ID
     
     稍后通过 Webhook/SSE 通知结果
     
优势:
✅ 解耦服务
✅ 削峰填谷
✅ 支持重试
✅ 可扩展性好
```

### 3.3 为什么选择多运行时支持？

```
┌─────────────────────────────────────────┐
│        任务类型 → 最优运行时选择         │
├─────────────────────────────────────────┤
│ 简单问答        → OpenAI API (快/便宜)   │
│ 复杂推理        → Claude Code (强/贵)    │
│ 本地工具调用    → OpenClaw (安全/可控)   │
│ 代码生成        → GitHub Copilot (集成)  │
└─────────────────────────────────────────┘

策略: 根据任务特征动态选择，平衡成本/质量/速度
```

---

## 4. 安全设计

### 4.1 认证授权

```
JWT Token 结构:
{
  "sub": "user_uuid",
  "email": "user@example.com",
  "role": "user",
  "permissions": ["agent:read", "task:create"],
  "iat": 1699123456,
  "exp": 1699209856
}

权限控制 (RBAC):
├─ admin: 所有权限
├─ user:  创建任务、查看结果
└─ guest: 仅浏览
```

### 4.2 数据安全

| 层面 | 措施 |
|------|------|
| **传输加密** | TLS 1.3 全链路 |
| **存储加密** | AES-256 数据库加密 |
| **密钥管理** | HashiCorp Vault |
| **输入过滤** | SQL 注入/XSS 防护 |
| **审计日志** | 所有操作记录 |

### 4.3 运行时安全

```go
// Agent 执行沙箱化
func ExecuteInSandbox(agentID, prompt string) {
    // 1. 创建隔离环境
    container := createContainer(
        Image: "agent-sandbox:latest",
        Limits: ResourceLimits{
            CPU:    "1",
            Memory: "512Mi",
            Disk:   "1GB",
        },
        Network: "none", // 隔离网络
    )
    
    // 2. 超时控制
    ctx, cancel := context.WithTimeout(5 * time.Minute)
    defer cancel()
    
    // 3. 执行
    result := container.Run(ctx, agentID, prompt)
    
    // 4. 清理
    defer container.Destroy()
    
    return result
}
```

---

## 5. 监控运维

### 5.1 监控指标体系

```yaml
# 黄金指标

延迟 (Latency):
- api_request_duration_seconds
- agent_execution_duration_seconds
- database_query_duration_seconds

流量 (Traffic):
- api_requests_per_second
- tasks_created_per_minute
- active_users

错误 (Errors):
- api_error_rate
- agent_failure_rate
- http_5xx_errors

饱和度 (Saturation):
- cpu_utilization
- memory_utilization
- queue_length
- database_connections
```

### 5.2 告警规则

```yaml
# Prometheus Alert Rules

groups:
- name: agenthub
  rules:
  # API 延迟告警
  - alert: HighLatency
    expr: api_request_duration_seconds{quantile="0.99"} > 2
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: API latency is high
      
  # 错误率告警
  - alert: HighErrorRate
    expr: rate(api_errors_total[5m]) / rate(api_requests_total[5m]) > 0.05
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: Error rate is above 5%
      
  # Agent 执行失败告警
  - alert: AgentExecutionFailures
    expr: rate(agent_executions_failed_total[5m]) > 10
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Multiple agent execution failures
```

### 5.3 链路追踪

```go
// OpenTelemetry 集成
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    tracer := otel.Tracer("task-service")
    ctx, span := tracer.Start(ctx, "CreateTask")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("agent_id", req.AgentID),
        attribute.String("user_id", req.UserID),
    )
    
    // 调用数据库
    task, err := db.CreateTask(ctx, req)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    // 发送到队列
    span.AddEvent("publishing to queue")
    if err := queue.Publish(ctx, task); err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    return task, nil
}
```

---

## 6. 部署架构演进

### 6.1 MVP 阶段 (单机部署)

```
┌─────────────────────────────────────┐
│           单台服务器                 │
│  ┌───────────────────────────────┐  │
│  │    Docker Compose             │  │
│  │  ┌─────────┐ ┌─────────────┐  │  │
│  │  │  Nginx  │ │   Next.js   │  │  │
│  │  └────┬────┘ └─────────────┘  │  │
│  │       │    ┌───────────────┐   │  │
│  │       └───►│  Go Services  │   │  │
│  │            └───────┬───────┘   │  │
│  │            ┌───────┴───────┐   │  │
│  │            │  PostgreSQL   │   │  │
│  │            │     Redis     │   │  │
│  │            └───────────────┘   │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘

成本: ~¥500/月
支持: 1000 日活用户
```

### 6.2 成长阶段 (K8s 集群)

```
┌─────────────────────────────────────────────┐
│              K8s 集群                        │
│  ┌───────────────────────────────────────┐  │
│  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐    │  │
│  │  │ API │ │Task │ │Agent│ │Bill │    │  │
│  │  │ x3  │ │Svc  │ │Svc  │ │Svc  │    │  │
│  │  │     │ │ x5  │ │ x5  │ │ x2  │    │  │
│  │  └─────┘ └─────┘ └─────┘ └─────┘    │  │
│  │  ┌───────────────────────────────┐   │  │
│  │  │      PostgreSQL Cluster       │   │  │
│  │  │       Redis Cluster           │   │  │
│  │  │       Kafka Cluster           │   │  │
│  │  └───────────────────────────────┘   │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘

成本: ~¥5000/月
支持: 10万 日活用户
```

### 6.3 扩展阶段 (多区域)

```
┌─────────────────────────────────────────────────────────┐
│                    多区域部署                            │
├─────────────────────────────────────────────────────────┤
│  中国华东 ◄────────► 中国华南 ◄────────► 海外新加坡      │
│     │                   │                   │           │
│  ┌──┴──┐             ┌──┴──┐             ┌──┴──┐       │
│  │K8s  │             │K8s  │             │K8s  │       │
│  │华东 │             │华南 │             │海外 │       │
│  └─────┘             └─────┘             └─────┘       │
│                                                         │
│  ┌─────────────────────────────────────────────────┐    │
│  │           Global Load Balancer                   │    │
│  │     (DNS 轮询 + 就近访问 + 故障转移)              │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘

成本: ~¥5万/月
支持: 100万+ 日活用户
```

---

## 7. 技术栈总结

| 层级 | 技术选型 | 理由 |
|------|---------|------|
| **前端** | Next.js 14 + TypeScript + Tailwind | SSR、性能、开发效率 |
| **API 网关** | Kong / Nginx | 成熟、高性能、插件丰富 |
| **微服务** | Go + Gin / Echo | 高并发、低延迟、云原生 |
| **数据库** | PostgreSQL 16 | ACID、JSON支持、扩展性强 |
| **缓存** | Redis 7 | 高性能、数据结构丰富 |
| **消息队列** | Apache Kafka | 高吞吐、持久化、可扩展 |
| **搜索** | Elasticsearch | 全文搜索、聚合分析 |
| **对象存储** | MinIO / OSS | S3兼容、低成本 |
| **容器** | Docker + Kubernetes | 标准化、弹性伸缩 |
| **监控** | Prometheus + Grafana | 云原生标准 |
| **追踪** | Jaeger + OpenTelemetry | 分布式追踪 |
| **CI/CD** | GitHub Actions / ArgoCD | 自动化部署 |

---

## 8. 下一步行动

### 本周
- [ ] 搭建开发环境 (Docker Compose)
- [ ] 初始化项目代码结构
- [ ] 设计数据库 Schema
- [ ] 搭建 CI/CD 流水线

### 本月
- [ ] 实现 API Gateway
- [ ] 实现 User Service
- [ ] 实现 Task Service
- [ ] 集成 5 个核心 Agent

### 下月
- [ ] 实现前端界面
- [ ] 接入支付系统
- [ ] 部署到测试环境
- [ ] 压力测试

---

*架构设计文档 v1.0*  
*最后更新: 2026-04-10*

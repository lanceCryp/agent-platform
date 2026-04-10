# AI Agent 全量服务平台 - 技术实施方案

## 目标
将 agency-agents-zh 全部 191 个 Agent 一次性包装为可调用服务

## 1. 全量 Agent 分类体系

```
191 Agents 分类矩阵:

├── Tier 1: 核心生产力 (30个) - 立即开放
│   ├── 工程开发 (8)
│   │   ├── engineering-code-reviewer
│   │   ├── engineering-senior-developer
│   │   ├── engineering-devops-automator
│   │   ├── engineering-security-engineer
│   │   ├── engineering-frontend-developer
│   │   ├── engineering-backend-architect
│   │   ├── engineering-ai-engineer
│   │   └── engineering-database-optimizer
│   │
│   ├── 设计创意 (7)
│   │   ├── design-brand-guardian
│   │   ├── design-ui-designer
│   │   ├── design-ux-architect
│   │   ├── design-image-prompt-engineer
│   │   ├── design-visual-storyteller
│   │   ├── design-ux-researcher
│   │   └── design-inclusive-visuals-specialist
│   │
│   ├── 市场营销 (8)
│   │   ├── marketing-seo-specialist
│   │   ├── marketing-content-creator
│   │   ├── marketing-social-media-strategist
│   │   ├── marketing-baidu-seo-specialist
│   │   ├── marketing-growth-hacker
│   │   ├── marketing-douyin-strategist
│   │   ├── marketing-xiaohongshu-operator
│   │   └── marketing-wechat-official-account
│   │
│   ├── 产品管理 (4)
│   │   ├── product-manager
│   │   ├── product-sprint-prioritizer
│   │   ├── product-behavioral-nudge-engine
│   │   └── product-trend-researcher
│   │
│   └── 项目管理 (3)
│       ├── project-manager-senior
│       ├── project-management-project-shepherd
│       └── project-management-jira-workflow-steward
│
├── Tier 2: 专业领域 (50个) - 1个月后开放
│   ├── 游戏开发 (12)
│   │   ├── unity-shader-graph-artist
│   │   ├── unity-architect
│   │   ├── unity-multiplayer-engineer
│   │   ├── unreal-world-builder
│   │   ├── unreal-technical-artist
│   │   ├── godot-shader-developer
│   │   ├── godot-gameplay-scripter
│   │   ├── godot-multiplayer-engineer
│   │   ├── game-designer
│   │   ├── game-audio-engineer
│   │   ├── level-designer
│   │   └── narrative-designer
│   │
│   ├── 财务金融 (6)
│   │   ├── finance-financial-forecaster
│   │   ├── finance-fraud-detector
│   │   ├── finance-invoice-manager
│   │   ├── accounts-payable-agent
│   │   ├── sales-pipeline-analyst
│   │   └── sales-deal-strategist
│   │
│   ├── 人力资源 (4)
│   │   ├── hr-recruiter
│   │   ├── hr-performance-reviewer
│   │   ├── recruitment-specialist
│   │   └── corporate-training-designer
│   │
│   ├── 法务合规 (4)
│   │   ├── legal-contract-reviewer
│   │   ├── legal-policy-writer
│   │   ├── compliance-auditor
│   │   └── blockchain-security-auditor
│   │
│   └── 其他专业 (24个)
│       ├── engineering-sre
│       ├── engineering-threat-detection-engineer
│       ├── testing-performance-benchmarker
│       ├── testing-api-tester
│       ├── testing-accessibility-auditor
│       ├── specialized-mcp-builder
│       ├── specialized-meeting-assistant
│       ├── marketing-bilibili-strategist
│       ├── marketing-kuaishou-strategist
│       ├── marketing-livestream-commerce-coach
│       ├── marketing-podcast-strategist
│       ├── marketing-book-co-author
│       ├── sales-account-strategist
│       ├── sales-coach
│       ├── sales-engineer
│       ├── sales-outbound-strategist
│       ├── sales-proposal-strategist
│       ├── project-management-experiment-tracker
│       ├── project-management-studio-operations
│       ├── project-management-studio-producer
│       ├── support-analytics-reporter
│       ├── support-finance-tracker
│       ├── support-infrastructure-maintainer
│       └── support-legal-compliance-checker
│
├── Tier 3: 垂直细分 (71个) - 3个月后开放
│   ├── 工程技术深度 (20)
│   │   ├── engineering-cms-developer
│   │   ├── engineering-mobile-app-builder
│   │   ├── engineering-embedded-firmware-engineer
│   │   ├── engineering-embedded-linux-driver-engineer
│   │   ├── engineering-fpga-digital-design-engineer
│   │   ├── engineering-iot-solution-architect
│   │   ├── engineering-solidity-smart-contract-engineer
│   │   ├── engineering-software-architect
│   │   ├── engineering-rapid-prototyper
│   │   ├── engineering-git-workflow-master
│   │   ├── engineering-incident-response-commander
│   │   ├── engineering-technical-writer
│   │   ├── engineering-dingtalk-integration-developer
│   │   ├── engineering-feishu-integration-developer
│   │   ├── engineering-wechat-mini-program-developer
│   │   ├── engineering-filament-optimization-specialist
│   │   ├── engineering-autonomous-optimization-architect
│   │   ├── engineering-ai-data-remediation-engineer
│   │   ├── engineering-data-engineer
│   │   └── engineering-email-intelligence-engineer
│   │
│   ├── 营销细分领域 (25)
│   │   ├── marketing-app-store-optimizer
│   │   ├── marketing-instagram-curator
│   │   ├── marketing-linkedin-content-creator
│   │   ├── marketing-twitter-engager
│   │   ├── marketing-reddit-community-builder
│   │   ├── marketing-carousel-growth-engine
│   │   ├── marketing-cross-border-ecommerce
│   │   ├── marketing-ecommerce-operator
│   │   ├── marketing-private-domain-operator
│   │   ├── marketing-video-optimization-specialist
│   │   ├── marketing-short-video-editing-coach
│   │   ├── marketing-ai-citation-strategist
│   │   ├── marketing-knowledge-commerce-strategist
│   │   ├── marketing-china-ecommerce-operator
│   │   ├── marketing-china-market-localization-strategist
│   │   ├── marketing-wechat-operator
│   │   ├── marketing-weibo-strategist
│   │   ├── marketing-weixin-channels-strategist
│   │   ├── marketing-zhihu-strategist
│   │   ├── marketing-tiktok-strategist
│   │   ├── marketing-xiaohongshu-specialist
│   │   ├── marketing-seo-specialist (重复)
│   │   └── ...
│   │
│   └── 其他垂直 (26)
│       ├── academic-anthropologist
│       ├── academic-geographer
│       ├── academic-historian
│       ├── academic-narratologist
│       ├── academic-psychologist
│       ├── academic-study-planner
│       ├── gaokao-college-advisor
│       ├── study-abroad-advisor
│       ├── supply-chain-inventory-forecaster
│       ├── supply-chain-route-optimizer
│       ├── supply-chain-vendor-evaluator
│       ├── identity-graph-operator
│       ├── report-distribution-agent
│       ├── data-consolidation-agent
│       ├── specialized-cultural-intelligence-strategist
│       ├── specialized-civil-engineer
│       ├── specialized-developer-advocate
│       ├── specialized-document-generator
│       ├── specialized-french-consulting-market
│       ├── specialized-korean-business-navigator
│       ├── specialized-model-qa
│       ├── specialized-pricing-optimizer
│       ├── specialized-risk-assessor
│       ├── specialized-salesforce-architect
│       └── specialized-workflow-architect
│
└── Tier 4: 实验性/冷门 (40个) - 按需开放
    ├── 学术深度研究 (6)
    ├── 政府/特殊行业 (4)
    ├── 新兴技术实验 (15)
    └── 小众市场 (15)

总统计:
- Tier 1 (核心): 30个 (16%)
- Tier 2 (专业): 50个 (26%)
- Tier 3 (垂直): 71个 (37%)
- Tier 4 (实验): 40个 (21%)
- 总计: 191个
```

## 2. 全量服务化技术方案

### 2.1 统一 Agent 接口标准

```go
// pkg/agent/types.go
package agent

// Agent 统一接口
type Agent interface {
    // 基础信息
    ID() string
    Name() string
    Category() string
    Tier() int  // 1/2/3/4
    
    // 执行
    Execute(ctx context.Context, req *Request) (*Response, error)
    
    // 元数据
    GetMetadata() *Metadata
    
    // 健康检查
    Health() error
    
    // 成本估算
    EstimateCost(req *Request) (*CostEstimate, error)
}

// Request 统一请求结构
type Request struct {
    AgentID     string                 `json:"agent_id"`
    UserID      string                 `json:"user_id"`
    Prompt      string                 `json:"prompt"`
    Context     map[string]interface{} `json:"context,omitempty"`
    MaxTokens   int                    `json:"max_tokens,omitempty"`
    Temperature float64                `json:"temperature,omitempty"`
    Priority    int                    `json:"priority,omitempty"` // 1-10
}

// Response 统一响应结构
type Response struct {
    AgentID      string        `json:"agent_id"`
    TaskID       string        `json:"task_id"`
    Status       string        `json:"status"` // success / failed
    Result       string        `json:"result"`
    Cost         CostBreakdown `json:"cost"`
    Duration     time.Duration `json:"duration"`
    TokensUsed   int           `json:"tokens_used"`
    Error        *AgentError   `json:"error,omitempty"`
    Metadata     *ResponseMeta `json:"metadata"`
}

// Metadata Agent 元数据
type Metadata struct {
    ID              string   `json:"id"`
    Name            string   `json:"name"`
    Description     string   `json:"description"`
    Category        string   `json:"category"`
    Tier            int      `json:"tier"`
    Tags            []string `json:"tags"`
    InputExample    string   `json:"input_example"`
    OutputExample   string   `json:"output_example"`
    AvgCost         float64  `json:"avg_cost"`      // 平均每次调用成本
    AvgDuration     int      `json:"avg_duration"`  // 平均执行时间(秒)
    SuccessRate     float64  `json:"success_rate"`  // 成功率
    Rating          float64  `json:"rating"`        // 用户评分
    TotalCalls      int64    `json:"total_calls"`   // 总调用次数
    SupportRuntime  []string `json:"support_runtime"` // 支持的运行时
}
```

### 2.2 全量 Agent 注册中心

```go
// pkg/agent/registry.go
package agent

// Registry Agent 注册中心
type Registry struct {
    agents map[string]Agent
    mu     sync.RWMutex
    
    // 索引
    byCategory map[string][]Agent
    byTier     map[int][]Agent
    byTag      map[string][]Agent
}

// NewRegistry 创建注册中心
func NewRegistry() *Registry {
    return &Registry{
        agents:     make(map[string]Agent),
        byCategory: make(map[string][]Agent),
        byTier:     make(map[int][]Agent),
        byTag:      make(map[string][]Agent),
    }
}

// Register 注册 Agent
func (r *Registry) Register(agent Agent) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    id := agent.ID()
    if _, exists := r.agents[id]; exists {
        return fmt.Errorf("agent %s already registered", id)
    }
    
    r.agents[id] = agent
    
    // 更新索引
    meta := agent.GetMetadata()
    r.byCategory[meta.Category] = append(r.byCategory[meta.Category], agent)
    r.byTier[meta.Tier] = append(r.byTier[meta.Tier], agent)
    
    for _, tag := range meta.Tags {
        r.byTag[tag] = append(r.byTag[tag], agent)
    }
    
    log.Printf("Agent registered: %s (Tier %d)", id, meta.Tier)
    return nil
}

// Get 获取 Agent
func (r *Registry) Get(id string) (Agent, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    agent, ok := r.agents[id]
    if !ok {
        return nil, fmt.Errorf("agent %s not found", id)
    }
    
    return agent, nil
}

// List 列出所有 Agent（支持过滤）
func (r *Registry) List(opts ListOptions) []Agent {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var result []Agent
    
    // 按 Tier 过滤
    if opts.Tier > 0 {
        result = r.byTier[opts.Tier]
    } else {
        // 返回所有
        for _, agent := range r.agents {
            result = append(result, agent)
        }
    }
    
    // 按 Category 过滤
    if opts.Category != "" {
        var filtered []Agent
        for _, agent := range result {
            if agent.GetMetadata().Category == opts.Category {
                filtered = append(filtered, agent)
            }
        }
        result = filtered
    }
    
    // 按 Tag 过滤
    if len(opts.Tags) > 0 {
        var filtered []Agent
        tagSet := make(map[string]bool)
        for _, tag := range opts.Tags {
            tagSet[tag] = true
        }
        
        for _, agent := range result {
            for _, tag := range agent.GetMetadata().Tags {
                if tagSet[tag] {
                    filtered = append(filtered, agent)
                    break
                }
            }
        }
        result = filtered
    }
    
    return result
}

// LoadAllFromDirectory 从目录加载全部 191 个 Agent
func (r *Registry) LoadAllFromDirectory(dir string) error {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }
    
    var successCount, failCount int
    
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        if !strings.HasSuffix(entry.Name(), ".md") {
            continue
        }
        
        agentID := strings.TrimSuffix(entry.Name(), ".md")
        
        // 解析 Agent 配置
        agent, err := r.parseAgentFile(filepath.Join(dir, entry.Name()))
        if err != nil {
            log.Printf("Failed to parse agent %s: %v", agentID, err)
            failCount++
            continue
        }
        
        if err := r.Register(agent); err != nil {
            log.Printf("Failed to register agent %s: %v", agentID, err)
            failCount++
            continue
        }
        
        successCount++
    }
    
    log.Printf("Loaded %d agents successfully, %d failed", successCount, failCount)
    return nil
}

// 自动分类 Tier
func (r *Registry) autoAssignTier(agentID string) int {
    // 核心生产力 Agent
    tier1Agents := []string{
        "engineering-code-reviewer",
        "engineering-senior-developer",
        "design-brand-guardian",
        "design-ui-designer",
        "marketing-seo-specialist",
        "marketing-content-creator",
        "product-manager",
        "project-manager-senior",
        // ... 其他 22 个
    }
    
    // 检查是否在 Tier 1 列表
    for _, id := range tier1Agents {
        if agentID == id {
            return 1
        }
    }
    
    // 根据调用频率自动分级
    // Tier 2: 中等使用频率
    // Tier 3: 低频但稳定
    // Tier 4: 实验性/冷门
    
    // 默认 Tier 3
    return 3
}
```

### 2.3 全量 Agent 服务化代码

```go
// cmd/agent-service/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "agenthub/pkg/agent"
    "agenthub/pkg/runtime"
    "agenthub/pkg/server"
)

func main() {
    // 1. 创建 Agent 注册中心
    registry := agent.NewRegistry()
    
    // 2. 加载全部 191 个 Agent
    agentDir := os.Getenv("AGENT_DIR")
    if agentDir == "" {
        agentDir = "/app/agents"  // Docker 容器内路径
    }
    
    log.Println("Loading all 191 agents...")
    if err := registry.LoadAllFromDirectory(agentDir); err != nil {
        log.Fatalf("Failed to load agents: %v", err)
    }
    
    // 3. 创建运行时管理器
    runtimeManager := runtime.NewManager(runtime.Config{
        OpenClaw: runtime.OpenClawConfig{
            Enabled:  true,
            Endpoint: os.Getenv("OPENCLAW_ENDPOINT"),
        },
        Claude: runtime.ClaudeConfig{
            Enabled: true,
            Path:    "/usr/local/bin/claude",
        },
        OpenAI: runtime.OpenAIConfig{
            Enabled: true,
            APIKey:  os.Getenv("OPENAI_API_KEY"),
        },
    })
    
    // 4. 创建 HTTP 服务器
    srv := server.New(server.Config{
        Port:           8080,
        Registry:       registry,
        RuntimeManager: runtimeManager,
        MaxConcurrency: 100,  // 最大并发 100 个任务
    })
    
    // 5. 优雅启动
    go func() {
        log.Println("Starting Agent Service on :8080")
        if err := srv.Start(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()
    
    // 6. 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}
```

### 2.4 统一 API 端点

```go
// pkg/server/handlers.go
package server

import (
    "encoding/json"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "agenthub/pkg/agent"
)

// Handler HTTP 处理器
type Handler struct {
    registry       *agent.Registry
    runtimeManager *runtime.Manager
}

// ListAgents 列出所有 Agent
func (h *Handler) ListAgents(c *gin.Context) {
    opts := agent.ListOptions{
        Tier:     c.GetInt("tier"),
        Category: c.Query("category"),
    }
    
    // 解析 tags
    if tags := c.QueryArray("tag"); len(tags) > 0 {
        opts.Tags = tags
    }
    
    agents := h.registry.List(opts)
    
    // 转换为响应格式
    var response []AgentDTO
    for _, a := range agents {
        meta := a.GetMetadata()
        response = append(response, AgentDTO{
            ID:          meta.ID,
            Name:        meta.Name,
            Category:    meta.Category,
            Tier:        meta.Tier,
            Description: meta.Description,
            Tags:        meta.Tags,
            AvgCost:     meta.AvgCost,
            Rating:      meta.Rating,
            IsAvailable: a.Health() == nil,
        })
    }
    
    c.JSON(http.StatusOK, gin.H{
        "total": len(response),
        "agents": response,
    })
}

// GetAgent 获取单个 Agent 详情
func (h *Handler) GetAgent(c *gin.Context) {
    agentID := c.Param("id")
    
    a, err := h.registry.Get(agentID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    meta := a.GetMetadata()
    c.JSON(http.StatusOK, AgentDetailDTO{
        ID:             meta.ID,
        Name:           meta.Name,
        Description:    meta.Description,
        Category:       meta.Category,
        Tier:           meta.Tier,
        Tags:           meta.Tags,
        InputExample:   meta.InputExample,
        OutputExample:  meta.OutputExample,
        AvgCost:        meta.AvgCost,
        AvgDuration:    meta.AvgDuration,
        SuccessRate:    meta.SuccessRate,
        Rating:         meta.Rating,
        TotalCalls:     meta.TotalCalls,
        SupportRuntime: meta.SupportRuntime,
    })
}

// ExecuteAgent 执行 Agent
func (h *Handler) ExecuteAgent(c *gin.Context) {
    agentID := c.Param("id")
    
    var req agent.Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    req.AgentID = agentID
    req.UserID = c.GetString("user_id") // 从 JWT 中获取
    
    // 检查 Agent 是否存在
    a, err := h.registry.Get(agentID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    // 检查 Tier 权限
    meta := a.GetMetadata()
    if meta.Tier > 2 { // Tier 3/4 需要特殊权限
        if !c.GetBool("is_vip") {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "This agent requires VIP subscription",
                "tier": meta.Tier,
            })
            return
        }
    }
    
    // 检查余额
    userBalance := c.GetFloat64("user_balance")
    estimate, _ := a.EstimateCost(&req)
    if userBalance < estimate.Total {
        c.JSON(http.StatusPaymentRequired, gin.H{
            "error": "Insufficient balance",
            "required": estimate.Total,
            "balance": userBalance,
        })
        return
    }
    
    // 异步执行（立即返回任务ID）
    taskID, err := h.runtimeManager.SubmitTask(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusAccepted, gin.H{
        "task_id": taskID,
        "status": "pending",
        "message": "Task submitted successfully",
    })
}

// GetTaskResult 获取任务结果
func (h *Handler) GetTaskResult(c *gin.Context) {
    taskID := c.Param("task_id")
    
    result, err := h.runtimeManager.GetTaskResult(c.Request.Context(), taskID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, result)
}

// GetAgentStats 获取 Agent 统计
func (h *Handler) GetAgentStats(c *gin.Context) {
    stats := h.registry.GetStats()
    
    c.JSON(http.StatusOK, gin.H{
        "total_agents": stats.TotalAgents,
        "by_tier": map[string]int{
            "tier_1": stats.ByTier[1],
            "tier_2": stats.ByTier[2],
            "tier_3": stats.ByTier[3],
            "tier_4": stats.ByTier[4],
        },
        "by_category": stats.ByCategory,
        "total_calls_today": stats.TodayCalls,
        "avg_response_time": stats.AvgResponseTime,
    })
}
```

### 2.5 Docker 部署（全量 191 Agent）

```dockerfile
# Dockerfile.agent-service
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o agent-service ./cmd/agent-service

# 运行时镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# 复制可执行文件
COPY --from=builder /app/agent-service .

# 复制全部 191 个 Agent 配置
COPY --from=builder /app/agents ./agents

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./agent-service"]
```

```yaml
# docker-compose.full.yml
version: '3.8'

services:
  agent-service:
    build:
      context: .
      dockerfile: Dockerfile.agent-service
    ports:
      - "8080:8080"
    environment:
      - AGENT_DIR=/app/agents
      - OPENCLAW_ENDPOINT=http://openclaw:5000
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - DB_HOST=postgres
      - REDIS_HOST=redis
      - KAFKA_BROKERS=kafka:9092
    volumes:
      - ./agents:/app/agents:ro  # 挂载 191 个 Agent 配置
    depends_on:
      - postgres
      - redis
      - kafka
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 4G
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # 其余服务与之前相同...
  postgres:
    image: postgres:16-alpine
    # ...

  redis:
    image: redis:7-alpine
    # ...

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    # ...

# 自动扩容配置
  agent-service-autoscaler:
    image: k8s.gcr.io/cluster-proportional-autoscaler:1.8.1
    environment:
      - SCALE_TARGET=agent-service
      - MIN_REPLICAS=3
      - MAX_REPLICAS=50
```

## 3. 质量保障机制

### 3.1 Agent 质量分级系统

```go
// pkg/quality/grader.go
package quality

// AgentGrader Agent 质量评分器
type AgentGrader struct {
    db *sql.DB
}

// Grade 对 Agent 进行质量评级
func (g *AgentGrader) Grade(agentID string) (*GradeResult, error) {
    // 查询最近 30 天的数据
    query := `
        SELECT 
            COUNT(*) as total_calls,
            SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as success_count,
            AVG(duration) as avg_duration,
            AVG(cost) as avg_cost,
            AVG(user_rating) as avg_rating
        FROM task_executions
        WHERE agent_id = $1
        AND created_at > NOW() - INTERVAL '30 days'
    `
    
    var result GradeResult
    err := g.db.QueryRow(query, agentID).Scan(
        &result.TotalCalls,
        &result.SuccessCount,
        &result.AvgDuration,
        &result.AvgCost,
        &result.AvgRating,
    )
    if err != nil {
        return nil, err
    }
    
    // 计算成功率
    result.SuccessRate = float64(result.SuccessCount) / float64(result.TotalCalls)
    
    // 综合评分 (0-100)
    result.Score = g.calculateScore(result)
    
    // 自动调整 Tier
    result.RecommendedTier = g.recommendTier(result)
    
    return &result, nil
}

func (g *AgentGrader) calculateScore(r GradeResult) float64 {
    // 评分权重
    successWeight := 0.4
    ratingWeight := 0.3
    usageWeight := 0.2
    costWeight := 0.1
    
    // 各项指标归一化 (0-100)
    successScore := r.SuccessRate * 100
    ratingScore := (r.AvgRating / 5.0) * 100
    usageScore := math.Min(float64(r.TotalCalls)/100.0, 1.0) * 100  // 使用频率
    costScore := math.Max(0, 100-r.AvgCost*10)  // 成本越低分越高
    
    return successScore*successWeight + 
           ratingScore*ratingWeight + 
           usageScore*usageWeight + 
           costScore*costWeight
}

func (g *AgentGrader) recommendTier(r GradeResult) int {
    switch {
    case r.Score >= 85 && r.TotalCalls >= 100:
        return 1  // Tier 1: 优秀
    case r.Score >= 70 && r.TotalCalls >= 50:
        return 2  // Tier 2: 良好
    case r.Score >= 60:
        return 3  // Tier 3: 及格
    default:
        return 4  // Tier 4: 需优化
    }
}
```

### 3.2 降级机制

```go
// pkg/fallback/manager.go
package fallback

// FallbackManager 降级管理器
type FallbackManager struct {
    primary  agent.Agent
    fallback []agent.Agent
}

// ExecuteWithFallback 带降级的执行
func (fm *FallbackManager) ExecuteWithFallback(ctx context.Context, req *agent.Request) (*agent.Response, error) {
    // 1. 尝试主 Agent
    resp, err := fm.primary.Execute(ctx, req)
    if err == nil {
        return resp, nil
    }
    
    log.Printf("Primary agent %s failed: %v, trying fallback...", fm.primary.ID(), err)
    
    // 2. 依次尝试备用 Agent
    for _, fb := range fm.fallback {
        resp, err = fb.Execute(ctx, req)
        if err == nil {
            resp.Metadata.FallbackUsed = true
            resp.Metadata.FallbackAgent = fb.ID()
            return resp, nil
        }
    }
    
    return nil, fmt.Errorf("all agents failed")
}

// 配置示例
var FallbackChains = map[string][]string{
    "engineering-code-reviewer": {
        "engineering-senior-developer",
        "engineering-software-architect",
    },
    "design-brand-guardian": {
        "design-ui-designer",
        "design-ux-architect",
    },
    // ...
}
```

## 4. 性能优化策略

### 4.1 智能预加载

```go
// 根据使用模式预加载热门 Agent
type Preloader struct {
    cache *lru.Cache
    stats *UsageStats
}

func (p *Preloader) PreloadHotAgents() {
    // 获取当前小时热门 Agent
    hotAgents := p.stats.GetTopAgents(time.Hour, 20)
    
    for _, agentID := range hotAgents {
        // 预热运行时
        go p.warmupRuntime(agentID)
    }
}
```

### 4.2 连接池优化

```go
// 运行时连接池
type RuntimePool struct {
    openclaw *pool.Pool  // 10 个连接
    claude   *pool.Pool  // 5 个连接  
    openai   *pool.Pool  // 20 个连接
}
```

## 5. 监控大盘

```yaml
# 全量 Agent 监控指标

# Agent 健康度
agent_health_score{agent_id} 0-100

# 按 Tier 统计
tier_usage{tier="1"} 1500  # Tier 1 今日调用
tier_usage{tier="2"} 800
tier_usage{tier="3"} 300
tier_usage{tier="4"} 50

# 失败率 Top 10
agent_failure_rate{agent_id="xxx"} 0.15

# 成本分布
agent_cost_histogram{bucket="0-1"} 50
agent_cost_histogram{bucket="1-5"} 120
agent_cost_histogram{bucket="5-10"} 30
agent_cost_histogram{bucket="10+"} 5
```

## 6. 实施路线图

### Phase 1: 基础设施 (1周)
- [ ] 搭建 Agent Registry
- [ ] 实现统一接口标准
- [ ] 加载全部 191 个 Agent

### Phase 2: 核心服务 (1周)
- [ ] API Gateway
- [ ] 任务调度器
- [ ] 运行时管理器

### Phase 3: 质量保障 (1周)
- [ ] 自动分级系统
- [ ] 降级机制
- [ ] 监控告警

### Phase 4: 灰度发布 (1周)
- [ ] Tier 1 开放 (30个)
- [ ] 收集反馈
- [ ] 优化调整

### Phase 5: 全量开放 (2周后)
- [ ] Tier 2 开放 (50个)
- [ ] Tier 3/4 按需开放

## 7. 成本预估

### 运行全部 191 个 Agent 的成本

| 项目 | 计算 | 月成本 |
|------|------|--------|
| **服务器** | 10台 8核16G | ¥8,000 |
| **API 调用** | 100万次/天 × ¥0.01 | ¥30,000 |
| **存储** | 1TB SSD | ¥500 |
| **带宽** | 10TB | ¥2,000 |
| **监控** | Prometheus/Grafana | ¥500 |
| **总计** | | **¥41,000/月** |

### 收入预估 (覆盖成本)

- 需要日活用户: 1,000 人
- 人均消费: ¥1.5/天
- 日收入: ¥1,500
- 月收入: ¥45,000 ✅

---

**总结**: 技术上完全可以一次性包装全部 191 个 Agent，但建议分级开放以保证质量。全量开放的月运营成本约 ¥4万，需要 1000 日活用户即可盈亏平衡。

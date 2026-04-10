# AgentHub 开发任务清单

## 📋 任务总览

| # | 任务 | 模块 | 优先级 | 状态 |
|---|------|------|--------|------|
| 1 | 创建数据库迁移脚本 | 后端 | P0 | ⬜ |
| 2 | 用户认证 API (注册/登录) | 后端 | P0 | ⬜ |
| 3 | Agent API | 后端 | P0 | ⬜ |
| 4 | 任务 API | 后端 | P0 | ⬜ |
| 5 | 订阅计费 API | 后端 | P0 | ⬜ |
| 6 | 注册页面 | 前端 | P0 | ⬜ |
| 7 | Dashboard 布局 | 前端 | P0 | ⬜ |
| 8 | Agent 详情页 | 前端 | P1 | ⬜ |
| 9 | 任务创建页面 | 前端 | P1 | ⬜ |
| 10 | 任务列表页 | 前端 | P1 | ⬜ |
| 11 | 个人设置页 | 前端 | P2 | ⬜ |
| 12 | Toast 通知组件 | 前端 | P2 | ⬜ |
| 13 | 单元测试 | 测试 | P0 | ⬜ |
| 14 | API 集成测试 | 测试 | P1 | ⬜ |

---

## 任务详情

### Task 1: 数据库迁移脚本
**文件**: `backend/migrations/001_initial_schema.sql`
**内容**:
- [ ] users 表
- [ ] agents 表
- [ ] tasks 表
- [ ] subscriptions 表
- [ ] transactions 表
- [ ] plans 表
- [ ] categories 表
- [ ] 索引创建
- [ ] 种子数据

### Task 2: 用户认证 API
**文件**: `backend/cmd/api-gateway/handler/auth.go`
**内容**:
- [ ] POST /auth/register - 用户注册
- [ ] POST /auth/login - 用户登录
- [ ] POST /auth/refresh - 刷新Token
- [ ] POST /auth/logout - 登出
- [ ] GET /users/me - 获取用户信息
- [ ] PATCH /users/me - 更新用户信息
- [ ] GET /users/me/balance - 获取余额

### Task 3: Agent API
**文件**: `backend/cmd/api-gateway/handler/agent.go`
**内容**:
- [ ] GET /agents - Agent列表 (分页/筛选)
- [ ] GET /agents/:id - Agent详情
- [ ] GET /agents/categories - 分类列表
- [ ] 缓存实现

### Task 4: 任务 API
**文件**: `backend/cmd/api-gateway/handler/task.go`
**内容**:
- [ ] POST /tasks - 创建任务
- [ ] GET /tasks - 任务列表
- [ ] GET /tasks/:id - 任务详情
- [ ] POST /tasks/:id/cancel - 取消任务
- [ ] POST /tasks/:id/retry - 重试任务
- [ ] Kafka 消息发布

### Task 5: 订阅计费 API
**文件**: `backend/cmd/api-gateway/handler/billing.go`
**内容**:
- [ ] GET /plans - 套餐列表
- [ ] GET /plans/:id - 套餐详情
- [ ] POST /subscriptions - 购买订阅
- [ ] DELETE /subscriptions/:id - 取消订阅
- [ ] GET /transactions - 交易记录
- [ ] POST /billing/recharge - 余额充值

### Task 6: 注册页面
**文件**: `frontend/app/(auth)/register/page.tsx`
**内容**:
- [ ] 表单设计
- [ ] 邮箱/用户名/密码验证
- [ ] 第三方登录按钮
- [ ] 协议勾选
- [ ] 错误处理

### Task 7: Dashboard 布局
**文件**: `frontend/app/(dashboard)/layout.tsx`
**内容**:
- [ ] 侧边栏导航
- [ ] 顶部栏
- [ ] 用户菜单
- [ ] 响应式布局
- [ ] 面包屑导航

### Task 8: Agent 详情页
**文件**: `frontend/app/(dashboard)/agents/[agentId]/page.tsx`
**内容**:
- [ ] Agent 信息展示
- [ ] 使用示例
- [ ] 用户评价
- [ ] 立即使用按钮
- [ ] 相关推荐

### Task 9: 任务创建页面
**文件**: `frontend/app/(dashboard)/tasks/new/page.tsx`
**内容**:
- [ ] Agent 选择器
- [ ] Prompt 输入框
- [ ] 参数配置
- [ ] 预估费用显示
- [ ] 提交处理

### Task 10: 任务列表页
**文件**: `frontend/app/(dashboard)/tasks/page.tsx`
**内容**:
- [ ] 任务列表展示
- [ ] 状态筛选
- [ ] 分页
- [ ] 任务详情查看
- [ ] 重新执行

### Task 11: 个人设置页
**文件**: `frontend/app/(dashboard)/settings/page.tsx`
**内容**:
- [ ] 个人信息编辑
- [ ] 密码修改
- [ ] 通知设置
- [ ] API Keys

### Task 12: Toast 通知组件
**文件**: `frontend/components/ui/toast.tsx`
**内容**:
- [ ] Toast 组件
- [ ] useToast Hook
- [ ] 成功/错误/警告/信息 类型
- [ ] 自动关闭
- [ ] 手动关闭

### Task 13: 单元测试
**目录**: `backend/pkg/**/*_test.go`
**内容**:
- [ ] config 测试
- [ ] auth 测试
- [ ] cache 测试
- [ ] database 测试

### Task 14: API 集成测试
**目录**: `backend/cmd/api-gateway/handler/*_test.go`
**内容**:
- [ ] auth handler 测试
- [ ] agent handler 测试
- [ ] task handler 测试
- [ ] billing handler 测试

---

## ✅ 完成记录

| 日期 | 完成任务 | 备注 |
|------|----------|------|
| 2026-04-10 | - | - |

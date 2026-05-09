# 宿舍管理系统 Harness 环境体系设计方案

> 基于 Harness Creator Skill，结合 Go + Vue3 双技术栈宿舍管理系统项目

---

## 一、Harness Creator 核心定位

### 1.1 是什么

**一句话**：为代码库构建 AI Agent 的操作系统——让 AI 代理能可靠地工作，而不是每次都盲人摸象。

```
Intelligence without infrastructure is just a demo.
Harness = 操作系统，LLM = CPU，代码库 = 唯一事实来源
```

### 1.2 它做什么 / 不做什么

| ✅ 做 | ❌ 不做 |
|------|--------|
| `AGENTS.md` — Agent 导航地图 | 业务代码（handler、service、repo） |
| `docs/ARCHITECTURE.md` — 架构文档 | Vue 组件代码 |
| `scripts/lint-*` — **层级检查器** | 数据库 schema 设计 |
| `harness/` — 运行环境 + eval 任务 | API 路由逻辑 |
| `Makefile` — 自动化入口 | 业务类型定义 |

### 1.3 核心哲学

1. **Repository 是唯一事实来源** — Agent 看不到 Slack/飞书/微信里的信息，所有知识必须进代码库
2. **AGENTS.md 是地图，不是手册** — 控制在 80-120 行，详细内容链向 `docs/`
3. **层级依赖用机械方式强制** — linter 报错必须包含：WHAT（什么错）+ WHY（为什么错）+ HOW（怎么修）
4. **Build to Delete** — 每个组件都应可替换，不要过度工程化

---

## 二、项目现状分析

### 2.1 技术栈

| 层次 | 技术 |
|------|------|
| 后端 | Go 1.x + Gin + pgx/v5 + PostgreSQL |
| 前端 | Vue 3 + TypeScript + Vite |
| 基础设施 | Docker + Make |

### 2.2 现有 Harness 组件检测

| 组件 | 状态 | 详情 |
|------|------|------|
| `AGENTS.md` | ⚠️ 偏简 | 73行，缺前端导航 |
| `docs/ARCHITECTURE.md` | ✅ 良好 | 已覆盖 Go 后端层 |
| `docs/DEVELOPMENT.md` | ✅ 良好 | 已有前后端启动说明 |
| `docs/QUALITY.md` | ✅ 良好 | 质量标准 |
| `scripts/lint-deps` | ⚠️ 仅 Go | 缺 TypeScript 版本 |
| `scripts/lint-quality` | ⚠️ 仅 Go | 缺 TypeScript 版本 |
| `harness/config/environment.json` | ⚠️ 仅后端 | 缺前端环境配置 |
| `harness/scripts/` | ❌ 缺失 | 无启动脚本 |
| `harness/tasks/` | ❌ 缺失 | 无 eval 任务 |

**结论**：Partial Harness，Go 后端基础良好，前端完全缺失。

---

## 三、推荐 Harness 体系架构

### 3.1 整体目录结构

```
dormitory-management/
├── AGENTS.md                          # Agent 入口导航（更新到 ~100行）
├── docs/
│   ├── ARCHITECTURE.md               # 更新：Go层 + Vue3结构
│   ├── DEVELOPMENT.md                # 更新：前后端启动命令
│   ├── QUALITY.md                    # 更新：质量标准
│   ├── LAYER_MAP.md                  # 新增：Go+Vue3 完整层级矩阵
│   └── design-docs/
│       ├── BUILDING_MANAGEMENT.md    # 新增：楼栋管理设计
│       ├── ALLOCATION_FLOW.md        # 新增：入住分配流程
│       └── VIOLATION_SYSTEM.md      # 新增：违规记录设计
├── scripts/
│   ├── lint-deps.go                  # 更新：Go层依赖检查（更全）
│   ├── lint-deps.ts                  # 新增：TypeScript 前端层级检查
│   ├── lint-quality.go               # 已有：Go代码质量
│   └── lint-quality.ts              # 新增：Vue3 代码质量
├── harness/
│   ├── config/
│   │   └── environment.json          # 更新：v2.0 schema（前后端）
│   ├── scripts/
│   │   ├── setup-env.sh             # 启动 PostgreSQL + Redis
│   │   ├── start-be.sh              # 启动 Go 服务（健康检查）
│   │   ├── start-fe.sh              # 启动 Vue3 dev server
│   │   └── teardown-env.sh          # 清理
│   └── tasks/
│       ├── bootstrap-dormitory-system/  # 已有的引导任务
│       ├── task-boilerplate.yaml    # 新增：标准任务模板
│       ├── task-crud-entity.yaml    # 新增：CRUD 实体任务
│       └── task-api-integration.yaml # 新增：前后端集成任务
└── Makefile                          # 更新：前端 lint 目标
```

---

## 四、Go 后端 Layer Map（已验证）

### 4.1 层级定义

```
┌─────────────────────────────────────────────┐
│ Layer 4 — Entry Points（顶层，依赖所有层）    │
│   cmd/server/, cmd/migrate/, cmd/seed/       │
├─────────────────────────────────────────────┤
│ Layer 3 — HTTP 层（依赖 0, 1, 2）             │
│   internal/handler/, internal/router/        │
├─────────────────────────────────────────────┤
│ Layer 2 — 业务逻辑层（依赖 0, 1）             │
│   internal/service/                         │
├─────────────────────────────────────────────┤
│ Layer 1 — 数据访问层（依赖 0）                │
│   internal/repository/                      │
├─────────────────────────────────────────────┤
│ Layer 0 — 基础设施层（零依赖）                │
│   pkg/database/, pkg/logger/, pkg/config/  │
└─────────────────────────────────────────────┘
```

### 4.2 禁止的反向依赖（已验证通过）

| 被调用方（低层） | 调用方（高层） | 错误说明 |
|----------------|--------------|---------|
| `internal/handler/` | `internal/service/` | Service 不能调 Handler |
| `internal/service/` | `internal/repository/` | Repo 不能调 Service |
| `internal/repository/` | `pkg/` | Repository 不能依赖 config/logger |
| `pkg/config/` | `internal/` | Config 层不能依赖业务层 |

---

## 五、Vue3 前端 Layer Map（新增）

### 5.1 层级定义

```
┌──────────────────────────────────────────────┐
│ Layer 3 — 视图层（最高层）                     │
│   frontend/src/views/, frontend/src/components/ │
├──────────────────────────────────────────────┤
│ Layer 2 — 组合式逻辑层                         │
│   frontend/src/composables/                   │
│   (useBuilding, useAllocation, useStudent...) │
├──────────────────────────────────────────────┤
│ Layer 1 — API 数据层                           │
│   frontend/src/api/                           │
│   (building.ts, student.ts, room.ts...)       │
├──────────────────────────────────────────────┤
│ Layer 0 — 工具基础设施层                       │
│   frontend/src/utils/, frontend/src/types/   │
│   (axios 封装、工具函数、共享类型)             │
└──────────────────────────────────────────────┘
```

### 5.2 禁止的反向依赖

| 被调用方（低层） | 调用方（高层） | 错误说明 |
|----------------|--------------|---------|
| `frontend/src/api/` | `frontend/src/views/` | View 不能直接调 API，应用 composable |
| `frontend/src/composables/` | `frontend/src/views/` | Composable 不能依赖 View |
| `frontend/src/utils/` | `frontend/src/api/` | API 层可用 utils，utils 不能依赖 api |
| `frontend/src/types/` | `frontend/src/views/` | View 应通过 props/emit 与类型交互 |

### 5.3 跨禁止依赖（前后端同检）

| 被调用方 | 调用方 | 说明 |
|---------|--------|------|
| `cmd/server/` | `cmd/migrate/` | Server 不应依赖迁移工具 |
| `cmd/seed/` | `cmd/migrate/` | Seed 不应依赖迁移工具 |
| `internal/handler/` | `internal/repository/` | Handler 应通过 Service 间接访问 Repo |

---

## 六、层级检查器设计

### 6.1 `scripts/lint-deps.go`（Go，已存在，更新）

```go
// 层级依赖矩阵 — 每条记录表示"上层可以依赖下层"
var ALLOWED = [][2]string{
    // Layer 0 基础设施
    {"pkg/database", "pkg/logger"},   // database 可被 logger 使用
    {"pkg/config", "pkg/database"},   // config 可用 database 类型
    // Layer 1 数据访问
    {"internal/repository", "pkg/database"},
    {"internal/repository", "pkg/logger"},
    {"internal/repository", "pkg/config"},
    // Layer 2 业务逻辑
    {"internal/service", "internal/repository"},
    {"internal/service", "pkg/database"},
    {"internal/service", "pkg/logger"},
    // Layer 3 HTTP 层
    {"internal/handler", "internal/service"},
    {"internal/handler", "internal/repository"},
    {"internal/handler", "pkg/logger"},
    // Layer 4 入口点
    {"cmd/server", "internal/handler"},
    {"cmd/server", "internal/service"},
    {"cmd/server", "internal/repository"},
}

// 禁止的反向依赖（高→低）
var FORBIDDEN = [][2]string{
    {"internal/service", "internal/handler"},   // ✗ Service 调 Handler
    {"internal/repository", "internal/service"}, // ✗ Repo 调 Service
    {"pkg/config", "internal/"},               // ✗ Config 依赖业务层
}
```

**错误信息模板**（必须包含 HOW TO FIX）：
```
[LINT-DEPS] internal/service/allocation_svc.go:12
  imports internal/handler (layer 3 → layer 2, reverse direction)

  Layer rule: internal/service (layer 2) may ONLY import:
    - internal/repository (layer 1)
    - pkg/database, pkg/logger (layer 0)

  Fix options:
  1. Move handler logic to a service method
  2. Define an interface in service, implement in handler
  3. Use dependency injection to break the cycle
```

### 6.2 `scripts/lint-deps.ts`（新增，TypeScript）

```typescript
const LAYER_MAP: Record<string, number> = {
  'frontend/src/utils': 0,
  'frontend/src/types': 0,
  'frontend/src/api': 1,
  'frontend/src/composables': 2,
  'frontend/src/views': 3,
  'frontend/src/components': 3,
}

const FORBIDDEN: [string, string][] = [
  // View 不能直接调 API
  ['frontend/src/views', 'frontend/src/api'],
  // Composable 不能依赖 View
  ['frontend/src/composables', 'frontend/src/views'],
  // 跨语言：Go server handler 不能直接调 repo（要经 service）
  ['cmd/server', 'internal/repository'],
]

// 错误输出示例：
// [LINT-DEPS] frontend/src/views/BuildingList.vue:5
//   imports frontend/src/api/building.ts (layer 3 → layer 1, reverse)
//   Fix: Use useBuilding() composable instead of direct API call
```

### 6.3 `scripts/lint-quality.ts`（新增）

检查项：
- Vue SFC：`<script setup>` 必填，禁用 Options API
- 组件名：PascalCase，必有 `.vue` 扩展名
- API 文件：每个实体一个文件（building.ts, student.ts, room.ts）
- Composable：必以 `use` 开头，return 完整类型
- 禁 `any` 类型（`@typescript-eslint/no-explicit-any`）

---

## 七、harness/ 配置文件

### 7.1 `harness/config/environment.json`（v2.0）

```json
{
  "version": "2.0",
  "services": [
    {
      "name": "postgres",
      "image": "postgres:15-alpine",
      "port": 5432,
      "env": {
        "POSTGRES_DB": "dormitory",
        "POSTGRES_USER": "dormuser",
        "POSTGRES_PASSWORD": "${POSTGRES_PASSWORD}"
      },
      "healthcheck": "pg_isready -U dormuser -d dormitory"
    },
    {
      "name": "redis",
      "image": "redis:7-alpine",
      "port": 6379,
      "env": {
        "REDIS_PASSWORD": "${REDIS_PASSWORD}"
      },
      "healthcheck": "redis-cli ping"
    }
  ],
  "env_files": [".env.be", ".env.fe"],
  "startup": {
    "backend": {
      "command": "go run cmd/server/main.go",
      "port": 8080,
      "health_endpoint": "/health",
      "depends_on": ["postgres", "redis"]
    },
    "frontend": {
      "command": "cd frontend && npm run dev",
      "port": 5173,
      "depends_on": ["backend"]
    }
  },
  "secrets": [
    {"name": "DATABASE_URL", "prompt": "PostgreSQL 连接字符串"},
    {"name": "POSTGRES_PASSWORD", "prompt": "PostgreSQL 密码"},
    {"name": "REDIS_PASSWORD", "prompt": "Redis 密码"},
    {"name": "JWT_SECRET", "prompt": "JWT 签名密钥"}
  ]
}
```

### 7.2 `harness/scripts/setup-env.sh`

```bash
#!/bin/bash
set -e

# 1. Start PostgreSQL
docker run -d \
  --name dormitory-postgres \
  -e POSTGRES_DB=dormitory \
  -e POSTGRES_USER=dormuser \
  -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
  -p 5432:5432 \
  postgres:15-alpine

# 2. Wait for postgres
until pg_isready -h localhost -U dormuser -d dormitory; do
  echo "Waiting for postgres..."
  sleep 1
done

# 3. Start Redis
docker run -d \
  --name dormitory-redis \
  -e REDIS_PASSWORD="${REDIS_PASSWORD}" \
  -p 6379:6379 \
  redis:7-alpine

echo "✅ All services started"
```

### 7.3 `harness/scripts/start-be.sh`

```bash
#!/bin/bash
set -e

cd /opt/data/dormitory-management

# Run migrations
go run cmd/migrate/main.go up

# Start server
go run cmd/server/main.go
```

### 7.4 `harness/scripts/start-fe.sh`

```bash
#!/bin/bash
set -e

cd /opt/data/dormitory-management/frontend
npm install
npm run dev
```

---

## 八、Eval 任务设计

### 8.1 `harness/tasks/task-boilerplate.yaml`

```yaml
name: task-boilerplate
description: 标准 CRUD 实体任务模板
trigger: |
  当用户要求"为 {entity} 实现 CRUD"时触发

steps:
  - id: create-entity-types
    description: 在 internal/types/ 创建实体定义
    verify: grep -q "type {Entity} struct" internal/types/

  - id: create-repo
    description: 在 internal/repository/ 创建数据访问层
    verify: grep -q "func.*Create.*{Entity}" internal/repository/

  - id: create-service
    description: 在 internal/service/ 创建业务逻辑层
    verify: grep -q "func.*Create.*{Entity}" internal/service/

  - id: create-handler
    description: 在 internal/handler/ 创建 HTTP 处理函数
    verify: grep -q "func.*Create.*{Entity}" internal/handler/

  - id: register-route
    description: 在 router 中注册路由
    verify: grep -q "/{entity}" internal/router/

  - id: write-unit-tests
    description: 为 service 层写单元测试
    verify: test -f internal/service/{entity}_svc_test.go

constraints:
  - 层依赖必须符合 LAYER_MAP.md
  - 所有 public 方法必须有注释
  - 测试覆盖率 ≥ 70%
```

### 8.2 `harness/tasks/task-api-integration.yaml`

```yaml
name: task-api-integration
description: 前后端 API 集成任务
trigger: |
  当用户要求"对接 {entity} API"或"联调前后端"时触发

steps:
  - id: define-api-contract
    description: 在 docs/API_CONTRACT.md 中定义 API 契约
    verify: grep -q "## .*{Entity}" docs/API_CONTRACT.md

  - id: backend-implementation
    description: 实现后端 API（符合 RESTful 规范）
    verify: curl -s http://localhost:8080/api/v1/{entity}s | grep -q "items"

  - id: frontend-api-client
    description: 在 frontend/src/api/ 创建 API 客户端
    verify: test -f frontend/src/api/{entity}.ts

  - id: frontend-composable
    description: 创建 use{Entity} composable
    verify: test -f frontend/src/composables/use{Entity}.ts

  - id: frontend-view
    description: 实现列表/详情/新建/编辑视图
    verify: test -f frontend/src/views/{Entity}List.vue

constraints:
  - 前端 view 不得直接 import API 模块（须经 composable）
  - API 响应须符合 docs/API_CONTRACT.md 定义
  - 前后端联调须通过 `make verify-fe`
```

---

## 九、AGENTS.md 更新计划

### 9.1 扩容目标：73行 → ~100行

**新增章节：**

```markdown
## 7. 前端开发（Vue3）
- 目录结构：views/ > composables/ > api/ > utils/
- 组件规范：PascalCase、`<script setup>` 必填
- API 调用：永远经由 composable，不在 view 里直接 fetch
- Lint 前端：`make lint-fe`
- 前端层级检查：`make lint-deps-fe`

## 8. 前后端集成
- 同时启动：`make setup && make start-be & make start-fe`
- API 契约文档：`docs/API_CONTRACT.md`
- 联调验证：`make verify-fe`

## 9. Harness 环境
- 环境配置：`harness/config/environment.json`
- 层级规则：`docs/LAYER_MAP.md`
- 新增任务：参考 `harness/tasks/task-boilerplate.yaml`

## 10. Quality Gates（新增）
- `make all` — 全部检查通过才能提交
- 层级检查失败：查看 `scripts/lint-deps` 输出
- 覆盖率：`make test-be` 显示覆盖率报告
```

---

## 十、执行计划

### Phase 1：Delta 检测（已完成）

| 缺失项 | 优先级 | 工作量 |
|--------|--------|--------|
| `scripts/lint-deps.ts` | P0 | 中 |
| `scripts/lint-quality.ts` | P0 | 中 |
| `harness/scripts/setup-env.sh` | P0 | 小 |
| `harness/scripts/start-be.sh` | P0 | 小 |
| `harness/scripts/start-fe.sh` | P0 | 小 |
| `harness/config/environment.json`（更新） | P1 | 小 |
| `docs/LAYER_MAP.md` | P1 | 中 |
| `harness/tasks/task-boilerplate.yaml` | P1 | 中 |
| `harness/tasks/task-api-integration.yaml` | P2 | 中 |
| `AGENTS.md`（更新） | P1 | 小 |
| `docs/API_CONTRACT.md` | P2 | 中 |
| `docs/design-docs/*.md` | P2 | 大 |

### Phase 2-4：并行创建（预计用时 30 分钟）

启动 3 个并行 subagent：
- **Doc Agent**：创建 `docs/LAYER_MAP.md`、`docs/API_CONTRACT.md`、design docs
- **Linter Agent**：创建 `scripts/lint-deps.ts`、`scripts/lint-quality.ts`
- **Harness Agent**：创建 `harness/scripts/*.sh`、更新 `environment.json`、task yaml

### Phase 5：验证

```bash
make all           # Go 构建 + 测试 + lint
make lint-deps-fe  # TypeScript 层级检查（新增）
make verify-fe     # 前端类型检查 + lint + 构建
```

---

## 十一、关键设计决策

### 11.1 为什么前后端分开 lint？

Go 和 TypeScript 是不同生态系统：
- Go 的 import 图可以通过静态分析完全解析
- TypeScript/Vue 的依赖需要考虑 SFC 的 `<script setup>` 和自动导入（Auto Import）
- 分开写可以在各自语言生态内做到最精确的检查

### 11.2 为什么不用单个 `lint-all`？

| 考量 | 单个 lint-all | 分开 lint |
|------|--------------|----------|
| 增量速度 | 慢（每次全量） | 快（只跑改过的语言） |
| 错误定位 | 模糊 | 清晰（哪个语言的问题） |
| 并行执行 | 不能 | 能（Go 和 TS 并行 lint） |
| 调试成本 | 高 | 低 |

### 11.3 为什么 `verify-fe` 依赖 `verify-be`？

前端 API 调用后端，前端构建成功不代表功能正确。通过 `depends_on: backend` 确保后端服务健康后再验证前端集成。

---

## 十二、如何开始

**选项 A：立即执行**
告诉执行者"执行完整 Phase 2-5"，会自动完成所有缺失文件的创建。

**选项 B：分步执行**
先执行 P0 项（lint-deps.ts + 环境脚本），验证通过后再做 P1、P2。

**关键验收标准**：
- `make lint-deps` 输出无层级错误
- `make lint-deps-fe` 报告无跨层依赖问题
- `AGENTS.md` 行数在 90-110 之间
- 所有 `scripts/*` 有执行权限

# 宿舍管理系统 - 智能体指南

学生宿舍管理系统 - Go Web API + Vue3 Frontend

## 1 快速开始

- [架构概述](docs/ARCHITECTURE.md) — 系统设计、层级、数据流
- [开发环境搭建](docs/DEVELOPMENT.md) — 构建、测试、环境配置
- [层级地图](docs/LAYER_MAP.md) — **必读**：Go + Vue3 层级依赖规则

## 2 架构

| 章节 | 文档 | 描述 |
|---------|----------|-------------|
| 2.1 | [系统架构](docs/ARCHITECTURE.md) | 层级结构、依赖图、数据流 |
| 2.2 | [层级地图](docs/LAYER_MAP.md) | Go + Vue3 完整层级矩阵，禁止依赖清单 |
| 2.3 | [设计文档](docs/design-docs/index.md) | 组件深度解析 |

## 3 质量与标准

| 章节 | 文档 | 描述 |
|---------|----------|-------------|
| 3.1 | [代码质量](docs/QUALITY.md) | 黄金原则、linter 规则 |
| 3.2 | [测试标准](docs/TESTING.md) | 测试模式、覆盖率 |
| 3.3 | [安全策略](docs/SECURITY.md) | 安全注意事项 |
| 3.4 | [质量评分](docs/QUALITY_SCORE.md) | 各领域质量等级 |

## 4 开发

```bash
# 后端
make all           # 全部检查（构建 + 层级 + 质量 + 测试）
go run ./scripts/lint-deps   # Go 层级检查

# 前端
make verify-fe     # 类型 + ESLint + 层级 + 构建
node scripts/lint-deps.mjs   # Vue3 层级检查

# 环境
./harness/scripts/setup-env.sh    # PostgreSQL + Redis
./harness/scripts/start-be.sh     # migrate + Go 服务
./harness/scripts/start-fe.sh     # Vue3 dev server
```

详情请参阅 [开发环境搭建](docs/DEVELOPMENT.md) 和 [层级地图](docs/LAYER_MAP.md)。

## 5 执行计划

- [进行中的计划](docs/exec-plans/active/) — 当前工作进度
- [已完成计划](docs/exec-plans/completed/) — 历史记录
- [技术债务追踪](docs/exec-plans/tech-debt-tracker.md)

## 6 核心目录

### 6.1 Go 后端

| 目录 | 层级 | 用途 |
|-----------|-------|---------|
| `cmd/server/` | L4 | 应用入口 |
| `cmd/migrate/` | L4 | 数据库迁移 |
| `cmd/seed/` | L4 | 测试数据初始化 |
| `internal/handler/` | L3 | HTTP 处理函数 |
| `internal/router/` | L3 | 路由注册 |
| `internal/middleware/` | L3 | 中间件（认证、CORS、日志） |
| `internal/response/` | L3 | 响应构建器 |
| `internal/request/` | **L0** | 请求 DTO（业务数据结构，与 types 同级） |
| `internal/service/` | L2 | 业务逻辑（事务编排） |
| `internal/repository/` | L1 | 数据访问（SQL） |
| `internal/cache/` | L1 | Redis 缓存访问 |
| `internal/types/` | L0 | 业务类型定义 |
| `internal/model/` | L0 | 数据库实体（Entity → Type 转换） |
| `pkg/database/` | **L-1** | PostgreSQL 连接池（任意层可导入） |
| `pkg/logger/` | **L-1** | Zap 日志封装（任意层可导入） |
| `pkg/config/` | **L-1** | Viper 配置读取（任意层可导入） |

> **层级 -1 说明**：`pkg/*` 包属于特殊基础设施层，任何层都可以导入，但它们自己不能依赖任何 `internal/*` 包。

### 6.2 Vue3 前端

| 目录 | 层级 | 用途 |
|-----------|-------|---------|
| `frontend/src/types/` | L0 | TypeScript 类型定义 |
| `frontend/src/utils/` | L0 | Axios 封装、工具函数 |
| `frontend/src/api/` | L1 | REST API 客户端（每个实体一个文件） |
| `frontend/src/stores/` | L2 | Pinia 状态管理 |
| `frontend/src/views/` | L3 | 页面视图组件 |
| `frontend/src/components/` | L4 | 通用 UI 组件 |

**重要**：前端 `views/` 禁止直接 import `api/`（层级 3 → 层级 1 违规）。必须经由 stores/composables 间接访问。

## 7 设计文档

| 文档 | 描述 |
|----------|------|
| [BUILDING_MANAGEMENT.md](docs/design-docs/BUILDING_MANAGEMENT.md) | 楼栋管理：数据模型、业务规则 |
| [ROOM_MANAGEMENT.md](docs/design-docs/ROOM_MANAGEMENT.md) | 房间管理：房型/床位联动、状态机 |
| [ALLOCATION_FLOW.md](docs/design-docs/ALLOCATION_FLOW.md) | 入住分配：校验规则、床位联动 |
| [VIOLATION_SYSTEM.md](docs/design-docs/VIOLATION_SYSTEM.md) | 违规系统：违规类型、处理流程 |
| [REPAIR_SYSTEM.md](docs/design-docs/REPAIR_SYSTEM.md) | 维修系统：状态机、报修/派单/维修/评价 |

完整索引：[design-docs/index.md](docs/design-docs/index.md)

## 8 Harness 环境

| 文件 | 用途 |
|------|---------|
| `harness/config/environment.json` | 运行环境配置（服务、端口、环境变量） |
| `harness/scripts/setup-env.sh` | 启动 PostgreSQL + Redis |
| `harness/scripts/start-be.sh` | 运行 migrate + 启动 Go 服务 |
| `harness/scripts/start-fe.sh` | 启动 Vue3 dev server |
| `harness/scripts/teardown-env.sh` | 清理 Docker 容器 |
| `harness/tasks/task-boilerplate.yaml` | 新增 CRUD 实体的标准流程 |
| `harness/tasks/task-api-integration.yaml` | 前后端 API 集成流程 |
| `harness/tasks/task-batch-allocation.yaml` | 批量分配场景（开学入住） |
| `harness/tasks/task-batch-checkout.yaml` | 批量退宿场景（毕业） |
| `harness/tasks/task-room-transfer.yaml` | 宿舍调换场景（换房） |
| `docs/HARNESS_DESIGN.md` | Harness 体系完整设计方案 |

## 9 任务模板

| 场景 | 模板 | 触发关键词 |
|------|------|-----------|
| 新增实体 CRUD | [task-boilerplate.yaml](harness/tasks/task-boilerplate.yaml) | "实现学生CRUD"、"新增模块" |
| 前后端集成 | [task-api-integration.yaml](harness/tasks/task-api-integration.yaml) | "对接API"、"联调前后端" |
| 批量分配（开学） | [task-batch-allocation.yaml](harness/tasks/task-batch-allocation.yaml) | "批量分配"、"开学入住" |
| 批量退宿（毕业） | [task-batch-checkout.yaml](harness/tasks/task-batch-checkout.yaml) | "批量退宿"、"毕业退房" |
| 宿舍调换 | [task-room-transfer.yaml](harness/tasks/task-room-transfer.yaml) | "换房"、"换宿"、"宿舍调整" |

---

**注意**：本文档是一个导航索引（约 130 行）。详细内容位于链接的文档中。
# Layer Map — 宿舍管理系统层级依赖规范

> 所有 AI Agent 必须严格遵守本文件的层级依赖规则。违规的代码不会被合并。

---

## 一、概述

本项目采用**分层架构**，每层只能依赖自身及更低层（数字越小越底层）：

```
Layer 0（基础设施）  ←  不得依赖任何内部包
Layer 1（数据访问）  ←  只能依赖 Layer 0
Layer 2（业务逻辑）  ←  只能依赖 Layer 0, 1
Layer 3（HTTP 层）    ←  只能依赖 Layer 0, 1, 2
Layer 4（入口点）    ←  可以依赖所有层
```

**禁止反向依赖**：高层不能依赖低层（反向调用）。

---

## 二、Go 后端 Layer Map

### 2.1 层级定义

| Layer | 包路径 | 职责 | 依赖规则 |
|-------|--------|------|---------|
| **-1** | `pkg/database/`, `pkg/logger/`, `pkg/config/` | 数据库连接池、日志、配置读取 | **任意层均可导入**，无内部依赖 |
| **0** | `internal/types/`, `internal/model/`, `internal/request/` | 业务类型定义、数据库实体、请求 DTO | **零内部依赖**，可依赖 Layer -1 |
| **1** | `internal/repository/` | 数据访问（CRUD + SQL） | 可依赖 Layer -1, 0 |
| **2** | `internal/service/` | 业务逻辑、事务编排 | 可依赖 Layer -1, 0, 1 |
| **3** | `internal/handler/`, `internal/router/`, `internal/middleware/`, `internal/response/` | HTTP 处理、路由、中间件 | 可依赖 Layer -1, 0, 1, 2 |
| **4** | `cmd/server/`, `cmd/migrate/`, `cmd/seed/` | 应用入口、数据库迁移、数据初始化 | 可依赖所有层 |

> **注意**：`pkg/*` 包属于 Layer -1（特殊层），任何层都可以导入，但它们自己不能依赖任何 `internal/*` 包。

### 2.2 允许的依赖矩阵

```
FROM (row) → TO (column)
         pkg/db  pkg/logger  pkg/config  repository  service  handler  cmd/
pkg/db     ✓         ✗         ✗         ✗         ✗       ✗       ✗
pkg/logger ✓         ✓         ✗         ✗         ✗       ✗       ✗
pkg/config ✓         ✗         ✓         ✗         ✗       ✗       ✗
repository  ✓         ✓         ✓         ✓(self)   ✗       ✗       ✗
service     ✓         ✓         ✗         ✓         ✓(self)  ✗       ✗
handler    ✓         ✓         ✗         ✓         ✓       ✓(self)  ✗
cmd/server  ✓         ✓         ✓         ✓         ✓       ✓       ✓(self)
```

### 2.3 禁止的反向依赖（Go）

| 违规模式 | 说明 |
|---------|------|
| `internal/service` → `internal/handler` | Service 层禁止调用 Handler（循环依赖） |
| `internal/repository` → `internal/service` | Repo 层禁止调用 Service |
| `pkg/config` → `internal/*` | 配置包禁止依赖业务代码 |
| `cmd/server` → `cmd/migrate` | Server 不应依赖迁移工具 |
| `cmd/seed` → `cmd/migrate` | Seed 不应依赖迁移工具 |

### 2.4 类型（Types）归属

- **业务类型**（`internal/types/`）：属于 Layer 0（纯数据结构，无依赖）
- **数据库 Model**（`internal/model/`）：属于 Layer 0
- 业务类型可以引用数据库 Model，但禁止引入 Service 或 Handler

---

## 三、Vue3 前端 Layer Map

### 3.1 层级定义

| Layer | 目录 | 职责 | 依赖规则 |
|-------|------|------|---------|
| **0** | `frontend/src/types/` | TypeScript 类型定义 | **零内部依赖**，只依赖外部库 |
| **0** | `frontend/src/utils/` | 工具函数、axios 封装 | 可依赖 Layer 0 |
| **1** | `frontend/src/api/` | REST API 客户端（每个实体一个文件） | 可依赖 Layer 0（types, utils） |
| **2** | `frontend/src/stores/` | Pinia 状态管理 | 可依赖 Layer 0, 1 |
| **3** | `frontend/src/views/` | 页面视图组件 | 可依赖 Layer 0, 1, 2 |
| **4** | `frontend/src/components/` | 通用 UI 组件 | 可依赖 Layer 0, 1, 2, 3 |

### 3.2 允许的依赖矩阵

```
FROM (row) → TO (column)
         types/   utils/   api/   stores/  views/  components/
types/     ✓(self)  ✗      ✗      ✗        ✗         ✗
utils/     ✓        ✓(self) ✗     ✗        ✗         ✗
api/       ✓        ✓      ✓(self) ✗        ✗         ✗
stores/    ✓        ✓      ✓       ✓(self)  ✗         ✗
views/     ✓        ✓      ✓       ✓       ✓(self)    ✗
components  ✓        ✓      ✓       ✓        ✓        ✓(self)
```

### 3.3 禁止的反向依赖（Vue3/TypeScript）

| 违规模式 | 说明 | 正确做法 |
|---------|------|---------|
| `views/` → `api/` | View 不能直接调用 API | 通过 stores 或 composables 间接调用 |
| `views/` → `stores/` 间接违规 | View 应通过 composable 访问 store | 创建 `use*Store()` composable |
| `stores/` → `views/` | Store 不能依赖 View | Store 应是纯数据逻辑 |
| `api/` → `stores/` | API 层不能依赖状态管理 | 只依赖 types 和 utils |
| `utils/` → `api/` | 工具层不能依赖 API 层 | utils 应是纯函数 |

### 3.4 Vue SFC（.vue 文件）特殊规则

- `<script setup>` 中只能 import 层内允许的模块
- `<template>` 和 `<style>` 无层级限制
- 组件名必须使用 PascalCase
- 公共组件放在 `components/`，页面组件放在 `views/`

---

## 四、跨语言（前后端）集成规则

| 规则 | 说明 |
|------|------|
| 前端 API 地址 | `VITE_API_BASE_URL` 环境变量指定 |
| 后端 CORS | `CORS_ORIGINS` 环境变量控制 |
| API 契约 | 前后端必须先在 `docs/API_CONTRACT.md` 达成一致，再开始实现 |
| 联调验证 | 前端 `make verify-fe` 必须先通过后端健康检查 |

---

## 五、层级检查命令

### Go 后端
```bash
go run ./scripts/lint-deps
```

### Vue3 前端
```bash
node scripts/lint-deps.mjs
# 或
make lint-deps-fe
```

### 全部检查
```bash
make all          # 后端全部
make verify-fe    # 前端全部（含层级检查）
```

---

## 六、违规示例与修复

### 示例 1：Service 调 Handler（Go）

```go
// ❌ internal/service/allocation_svc.go
import "github.com/example/dormitory-management/internal/handler"

// Layer 2 (service) cannot import Layer 3 (handler)
```

**修复**：在 `internal/service/` 中定义接口，由 Handler 实现：
```go
// ✅ internal/service/interfaces.go
type AllocationHandlerInterface interface {
    Create(ctx context.Context, req *request.CreateAllocationRequest) (*types.Allocation, error)
}

// ✅ internal/service/allocation_svc.go
type AllocationService struct {
    handler AllocationHandlerInterface  // 通过接口注入，不直接依赖 handler 包
}
```

### 示例 2：View 直接调 API（Vue3）

```vue
<!-- ❌ frontend/src/views/BuildingList.vue -->
<script setup>
import { buildingApi } from '@/api/building'  // Layer 3 → Layer 1，直接违规
</script>
```

**修复**：通过 Pinia store 访问：
```vue
<!-- ✅ frontend/src/views/BuildingList.vue -->
<script setup>
import { useBuildingStore } from '@/stores/building'
const store = useBuildingStore()
</script>
```

### 示例 3：Store 调 View（Vue3）

```ts
// ❌ frontend/src/stores/building.ts
import type { BuildingListView } from '@/views/BuildingList.vue'  // Layer 2 → Layer 3
```

**修复**：Store 不应依赖任何 View 类型，类型应定义在 `types/` 中：
```ts
// ✅ frontend/src/types/building.ts
export interface BuildingSummary {
  id: string
  name: string
  gender: string
}
```

---

## 七、Layer Map 生效机制

层级依赖不是文档约定，而是**机械执行的规则**：

1. **CI 强制**：PR 必须通过 `make all` 和 `make verify-fe`，其中包含层级检查
2. **linter 拦截**：违反规则的代码在 `go run ./scripts/lint-deps` 或 `node scripts/lint-deps.mjs` 时立即报错
3. **错误信息可操作**：linter 输出包含 WHAT + WHICH LAYER + HOW TO FIX，AI Agent 可直接根据提示修复

---

## 八、变更流程

如果业务需要突破层级限制（例如新功能确实需要跨层调用）：

1. 在 `docs/LAYER_MAP.md` 中提出 RFC，说明为什么必须打破层级
2. 由架构 Owner 评审
3. 评审通过后更新本文件，然后才能合并相关代码

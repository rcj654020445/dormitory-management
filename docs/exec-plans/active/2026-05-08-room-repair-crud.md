# Room CRUD + Repair 维修系统 CRUD

**创建时间**: 2026-05-08
**任务ID**: room-repair-crud-20260508-1330-20260508-1330

## 目标

为宿舍管理系统实现两个核心模块的完整 CRUD：
1. **Room** — 房间管理（含床位联动、状态机）
2. **Repair** — 维修系统（含状态机、互斥规则）

每个模块包含：Go 后端 CRUD + 前端 Vue3 集成 + 单元测试。

## 范围

### 阶段1: Room CRUD

**后端文件（按 task-boilerplate.yaml 流程）**：
- `internal/types/room.go` — Room 类型定义
- `internal/model/entity.go` — 添加 RoomEntity + ToRoom()
- `internal/repository/room_repo.go` — CRUD + ListByBuilding
- `internal/service/room_svc.go` — 业务逻辑 + 床位校验
- `internal/handler/room_handler.go` — HTTP 接口
- `internal/router/` — 路由注册
- `cmd/migrate/main.go` — rooms 建表
- `internal/service/room_svc_test.go` — 单元测试

**前端文件（按 task-boilerplate.yaml 流程）**：
- `frontend/src/types/index.ts` — 添加 Room 类型
- `frontend/src/api/room.ts` — Layer 1 API 客户端
- `frontend/src/stores/room.ts` — Layer 2 Pinia store
- `frontend/src/views/RoomList.vue` — Layer 3 视图
- `frontend/src/router/index.ts` — 路由注册

**验证**：
- `go build ./cmd/server`
- `go run ./scripts/lint-deps`
- `make verify-fe`
- `go test ./internal/service/... -run Room -coverprofile=coverage.out`

---

### 阶段2: Repair CRUD

**后端文件（按 task-boilerplate.yaml 流程）**：
- `internal/types/repair.go` — Repair 类型定义
- `internal/model/entity.go` — 添加 RepairEntity + ToRepair()
- `internal/repository/repair_repo.go` — CRUD + 按状态列表
- `internal/service/repair_svc.go` — 状态机 + 互斥规则
- `internal/handler/repair_handler.go` — HTTP 接口
- `internal/router/` — 路由注册
- `cmd/migrate/main.go` — repairs 建表
- `internal/service/repair_svc_test.go` — 单元测试

**前端文件（按 task-api-integration.yaml 流程）**：
- `frontend/src/types/index.ts` — 添加 Repair 类型
- `frontend/src/api/repair.ts` — Layer 1 API 客户端
- `frontend/src/stores/repair.ts` — Layer 2 Pinia store
- `frontend/src/views/RepairList.vue` — Layer 3 视图
- `frontend/src/router/index.ts` — 路由注册

**验证**：
- `go build ./cmd/server`
- `go run ./scripts/lint-deps`
- `make verify-fe`
- `go test ./internal/service/... -run Repair -coverprofile=coverage.out`

## 参考文档

- [ROOM_MANAGEMENT.md](docs/design-docs/ROOM_MANAGEMENT.md)
- [REPAIR_SYSTEM.md](docs/design-docs/REPAIR_SYSTEM.md)
- [LAYER_MAP.md](docs/LAYER_MAP.md)
- [task-boilerplate.yaml](harness/tasks/task-boilerplate.yaml)
- [task-api-integration.yaml](harness/tasks/task-api-integration.yaml)

## 关键业务规则

### Room
- `beds_used <= beds_total` 不变量
- 状态为 `inactive`/`maintenance` 时不可分配
- 同一楼栋内房间号唯一

### Repair
- 同一房间同时只能有一个 `repairing` 状态的维修单
- 状态转换：`pending → assigned → repairing → completed`
- 取消仅限 `pending`/`assigned` 状态

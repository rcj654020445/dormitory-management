# Repair System — 维修管理

## 一、业务概述

维修系统（Repair）管理宿舍设施的报修、派单、维修完成整个流程。学生或宿管可以提交维修申请，维修工人接单处理，状态全程可追踪。

**核心业务流程**：
```
报修单创建 → 待派单 → 已派单 → 维修中 → 已完成
                  ↓
               已取消
```

**特殊约束**：
- 报修只能关联到具体的 Room，不能关联到 Building
- 一个 Room 同时只能有一个 `repairing` 状态的维修单（避免重复维修冲突）
- 维修状态变更需要记录时间戳

---

## 二、数据模型

### 2.1 Database Schema

```sql
CREATE TABLE repairs (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id      UUID        NOT NULL REFERENCES rooms(id),
    reporter_id  UUID        NOT NULL REFERENCES students(id),
    repairer_id  UUID        REFERENCES students(id),
    type         VARCHAR(30) NOT NULL DEFAULT 'facility',
    description  TEXT        NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority     VARCHAR(10) NOT NULL DEFAULT 'normal',
    scheduled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cost         DECIMAL(10,2),
    rating       INT         CHECK (rating BETWEEN 1 AND 5),
    remark       TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_repairs_room_id ON repairs(room_id);
CREATE INDEX idx_repairs_status ON repairs(status);
```

### 2.2 TypeScript 类型（`frontend/src/types/repair.ts`）

```typescript
export type RepairType =
  | 'facility'    // 设施维修（门/窗/床/柜）
  | 'plumbing'    // 水管/漏水
  | 'electrical'  // 电路/灯具/插座
  | 'network'     // 网络故障
  | 'cleaning'    // 保洁
  | 'other'       // 其他

export type RepairStatus =
  | 'pending'     // 待派单
  | 'assigned'    // 已派单
  | 'repairing'   // 维修中
  | 'completed'   // 已完成
  | 'cancelled'   // 已取消

export type RepairPriority = 'urgent' | 'normal' | 'low'

export interface Repair {
  id: string
  roomId: string
  reporterId: string
  repairerId: string | null
  type: RepairType
  description: string
  status: RepairStatus
  priority: RepairPriority
  scheduledAt: string | null
  completedAt: string | null
  cost: number | null
  rating: number | null
  remark: string | null
  createdAt: string
  updatedAt: string
}

export interface RepairWithRoom extends Repair {
  roomNumber: string
  buildingName: string
  reporterName: string
  repairerName: string | null
}
```

### 2.3 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | UUID | 主键 |
| `room_id` | UUID | 报修房间 |
| `reporter_id` | UUID | 报修人（学生） |
| `repairer_id` | UUID | 维修工人（可为空） |
| `type` | ENUM | 维修类型 |
| `description` | TEXT | 故障描述 |
| `status` | ENUM | 当前状态 |
| `priority` | ENUM | 优先级 |
| `scheduled_at` | TIMESTAMPTZ | 预约维修时间 |
| `completed_at` | TIMESTAMPTZ | 实际完成时间 |
| `cost` | DECIMAL | 维修费用（学生承担部分） |
| `rating` | INT | 学生评价（1-5星） |
| `remark` | TEXT | 备注 |

---

## 三、业务规则

### 3.1 状态机

```
pending → assigned → repairing → completed
    ↓          ↓
cancelled   cancelled
```

| 状态 | 说明 | 允许的下一个状态 |
|------|------|-----------------|
| `pending` | 待派单 | `assigned`, `cancelled` |
| `assigned` | 已派单 | `repairing`, `cancelled` |
| `repairing` | 维修中 | `completed` |
| `completed` | 已完成 | —（终态） |
| `cancelled` | 已取消 | —（终态） |

### 3.2 互斥规则

**同一房间同时只能有一个 `repairing` 状态的维修单**：

```sql
-- 创建 repair 时校验
SELECT COUNT(*) FROM repairs
WHERE room_id = $1 AND status = 'repairing'
-- 如果 count > 0，拒绝创建
```

### 3.3 报修人约束

- 只有宿舍分配了床位的学生才能提交维修申请（通过 Allocation 验证 `student_id` 有效）
- 维修工人由系统管理员（Admin）指定（repairer_id）

### 3.4 状态变更权限

| 操作 | 权限 |
|------|------|
| `pending` → `assigned` | Admin |
| `assigned` → `repairing` | Repairer |
| `repairing` → `completed` | Repairer |
| `pending/assigned` → `cancelled` | Admin |
| 评价（rating） | Reporter（仅学生可评价自己的维修单） |

### 3.5 错误码

| 错误码 | HTTP Status | 说明 |
|--------|-------------|------|
| `REPAIR_001` | 404 | 维修单不存在 |
| `REPAIR_002` | 409 | 该房间已有正在维修中的维修单 |
| `REPAIR_003` | 422 | 状态转换不合法 |
| `REPAIR_004` | 404 | 报修的房间不存在 |
| `REPAIR_005` | 422 | 学生未入住该房间（无法报修） |
| `REPAIR_006` | 422 | 维修工人不存在 |
| `REPAIR_007` | 403 | 无权限执行该状态变更 |

---

## 四、维修业务流程

### 4.1 报修流程（学生端）

1. 学生登录 → 进入"报修申请"
2. 选择房间（自动列出自己当前入住的房间）
3. 选择维修类型，填写故障描述
4. 提交 → 状态为 `pending`

### 4.2 派单流程（Admin 端）

1. Admin 查看所有 `pending` 维修单
2. 选择维修单 → 指定维修工人（repairer_id）
3. 可选填写预约时间（scheduled_at）
4. 确认派单 → 状态变为 `assigned`

### 4.3 维修流程（Repairer 端）

1. Repairer 登录 → 查看自己被分配的维修单（`assigned`）
2. 开始维修 → 状态变为 `repairing`
3. 维修完成 → 填写完成时间、费用，提交 → 状态变为 `completed`

### 4.4 评价流程（Reporter 端）

1. 维修单变为 `completed` 后，学生收到通知
2. 学生对维修单进行 1-5 星评价
3. 评价后 repair 表 rating 字段被填充

---

## 五、API 关联

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/repairs` | 维修单列表（支持 status、room_id、priority 过滤） |
| GET | `/api/repairs/:id` | 获取维修单详情 |
| POST | `/api/repairs` | 创建维修单（学生） |
| PUT | `/api/repairs/:id/status` | 更新维修单状态（Admin/Repairer） |
| PUT | `/api/repairs/:id/rating` | 学生评价（Reporter） |
| DELETE | `/api/repairs/:id` | 取消维修单（仅 pending/assigned） |

详细 API 契约见 [API_CONTRACT.md](../API_CONTRACT.md)

---

## 六、前端视图（Vue3）

| 视图 | 路径 | 功能 |
|------|------|------|
| 维修单列表 | `/repairs` | 展示维修单列表、状态筛选 |
| 报修申请 | `/repairs/new` | 学生提交新维修单 |
| 维修单详情 | `/repairs/:id` | 查看详情、执行状态变更 |

---

## 七、实现检查清单

- [ ] `internal/types/repair.go` — Repair 类型定义
- [ ] `internal/model/entity.go` — 添加 `RepairEntity` 和 `ToRepair()` 方法
- [ ] `internal/repository/repair_repo.go` — Repository（含互斥校验 SQL）
- [ ] `internal/service/repair_svc.go` — Service（含状态机、权限校验）
- [ ] `internal/handler/repair_handler.go` — HTTP handler
- [ ] `cmd/migrate/main.go` — 添加 repairs 建表语句
- [ ] `frontend/src/types/index.ts` — 添加 Repair 类型
- [ ] `frontend/src/api/repair.ts` — Layer 1 API client
- [ ] `frontend/src/stores/repair.ts` — Layer 2 Pinia store
- [ ] `frontend/src/views/RepairList.vue` — Layer 3 视图
- [ ] `frontend/src/router/index.ts` — 路由注册

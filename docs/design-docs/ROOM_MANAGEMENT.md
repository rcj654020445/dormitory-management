# Room Management — 房间管理

## 一、业务概述

房间（Room）是宿舍管理系统的核心资源单位。每个房间从属于一个楼栋（Building），具有固定的床位数、设施配置和当前状态。房间管理直接影响入住分配的业务能力。

**核心约束**：
- 房间的 `BedsUsed` 不得超过 `BedsTotal`
- 房间状态为 `inactive` / `maintenance` 时不得分配床位
- 房间变更楼栋时需校验楼栋性别属性（男/女楼）

---

## 二、数据模型

### 2.1 Database Schema

```sql
CREATE TABLE rooms (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    building_id UUID        NOT NULL REFERENCES buildings(id),
    number     VARCHAR(20) NOT NULL,
    floor      INT         NOT NULL CHECK (floor >= 1),
    type       VARCHAR(20) NOT NULL DEFAULT 'double',
    beds_total INT         NOT NULL CHECK (beds_total BETWEEN 1 AND 8),
    beds_used  INT         NOT NULL DEFAULT 0 CHECK (beds_used >= 0),
    has_bathroom BOOLEAN   NOT NULL DEFAULT false,
    has_ac     BOOLEAN     NOT NULL DEFAULT false,
    status     VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (building_id, number)
);
```

### 2.2 TypeScript 类型（`frontend/src/types/room.ts`）

```typescript
export type RoomType = 'single' | 'double' | 'quad' | 'hex' | 'oct'

export type RoomStatus = 'active' | 'inactive' | 'maintenance'

export interface Room {
  id: string
  buildingId: string
  number: string       // 如 "101", "1203"
  floor: number
  type: RoomType
  bedsTotal: number
  bedsUsed: number
  hasBathroom: boolean
  hasAC: boolean
  status: RoomStatus
  createdAt: string
  updatedAt: string
}

export interface RoomWithBuilding extends Room {
  buildingName: string
  buildingGender: Gender
}
```

### 2.3 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | UUID | 主键 |
| `building_id` | UUID | 所属楼栋 |
| `number` | VARCHAR(20) | 房间号，如 "101" 表示 1 楼 01 室 |
| `floor` | INT | 楼层 |
| `type` | ENUM | 房型：single(1), double(2), quad(4), hex(6), oct(8) |
| `beds_total` | INT | 床位总数 |
| `beds_used` | INT | 已占用床位数 |
| `has_bathroom` | BOOLEAN | 是否有独立卫生间 |
| `has_ac` | BOOLEAN | 是否有空调 |
| `status` | ENUM | active / inactive / maintenance |

---

## 三、业务规则

### 3.1 房间状态机

```
active ←→ inactive
  ↓
maintenance
```

| 状态 | 描述 | 可分配床位 | 可编辑 |
|------|------|-----------|--------|
| `active` | 正常可用 | ✅ | ✅ |
| `inactive` | 停用（装修/改造） | ❌ | ✅ |
| `maintenance` | 维修中 | ❌ | ✅ |

### 3.2 床位校验

**核心不变量**：任何时刻 `beds_used <= beds_total`

- 创建房间时：`beds_total >= 1`
- 更新房间时：禁止将 `beds_used` 设置为大于 `beds_total` 的值
- 删除房间时：必须 `beds_used == 0` 才能删除

### 3.3 房型与床位联动

| type | beds_total | 说明 |
|------|-----------|------|
| `single` | 1 | 单人间 |
| `double` | 2 | 双人间（默认） |
| `quad` | 4 | 四人间 |
| `hex` | 6 | 六人间 |
| `oct` | 8 | 八人间 |

修改 `type` 时，应同步更新 `beds_total`：
- 用户指定 `beds_total` 时优先使用用户值
- 用户仅改变 `type` 时，`beds_total` 自动调整为上表对应值

### 3.4 楼栋性别属性约束

楼栋有 `gender` 属性（`male` / `female`），房间本身不存储性别，但分配床位时通过楼栋间接约束：
- 男生楼栋的房间只能分配给 `gender=male` 的学生
- 女生楼栋的房间只能分配给 `gender=female` 的学生

### 3.5 错误码

| 错误码 | HTTP Status | 说明 |
|--------|-------------|------|
| `ROOM_001` | 404 | 房间不存在 |
| `ROOM_002` | 409 | 房间号在楼栋内已存在（UNIQUE 冲突） |
| `ROOM_003` | 422 | `beds_used > beds_total` |
| `ROOM_004` | 422 | 无法删除非空房间 |
| `ROOM_005` | 422 | 楼栋内楼层范围校验失败 |
| `ROOM_006` | 409 | 房间状态不允许分配（inactive/maintenance） |

---

## 四、API 关联

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/rooms` | 房间列表（支持 building_id、floor、status 过滤） |
| GET | `/api/rooms/:id` | 获取房间详情 |
| POST | `/api/rooms` | 创建房间 |
| PUT | `/api/rooms/:id` | 更新房间 |
| DELETE | `/api/rooms/:id` | 删除空房间 |
| GET | `/api/buildings/:id/rooms` | 获取指定楼栋的房间列表 |

详细 API 契约见 [API_CONTRACT.md](../API_CONTRACT.md)

---

## 五、前端视图（Vue3）

| 视图 | 路径 | 功能 |
|------|------|------|
| 房间列表 | `/rooms` | 列表、搜索（楼栋/楼层/状态）、新增 |
| 房间详情 | `/rooms/:id` | 查看房间详情、设施、床位占用情况 |

---

## 六、实现检查清单

为房间模块实现 CRUD 时需确认以下文件：

- [ ] `internal/types/room.go` — 类型定义
- [ ] `internal/model/entity.go` — 添加 `RoomEntity` 和 `ToRoom()` 方法
- [ ] `internal/repository/room_repo.go` — Repository 接口 + 实现
- [ ] `internal/service/room_svc.go` — Service（含床位校验逻辑）
- [ ] `internal/handler/room_handler.go` — HTTP handler
- [ ] `internal/handler/router.go` — 路由注册
- [ ] `cmd/migrate/main.go` — 添加 rooms 建表语句
- [ ] `frontend/src/types/index.ts` — 添加 Room 类型
- [ ] `frontend/src/api/room.ts` — Layer 1 API client
- [ ] `frontend/src/stores/room.ts` — Layer 2 Pinia store
- [ ] `frontend/src/views/RoomList.vue` — Layer 3 视图
- [ ] `frontend/src/router/index.ts` — 路由注册

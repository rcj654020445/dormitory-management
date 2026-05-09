# API Contract — 宿舍管理系统

> 前后端 API 契约文档。所有 API 必须先在本文件定义格式，再开始实现。
> 未在此定义的 API 不允许出现在代码中。

---

## 一、概述

### 1.1 基础信息

| 项目 | 值 |
|------|---|
| Base URL | `http://localhost:8080/api/v1` |
| Content-Type | `application/json` |
| 认证方式 | Bearer Token（在 `Authorization` header 中传递） |
| 分页 | Page-based，Query 参数 `page` + `page_size` |

### 1.2 通用响应格式

**成功响应**：
```json
// 单个资源
{
  "data": { ... }

  // 分页列表
  {
    "data": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

**错误响应**：
```json
{
  "error": {
    "code": "STUDENT_NOT_FOUND",
    "message": "学生不存在或已删除"
  }
}
```

### 1.3 通用错误码

| HTTP Status | Error Code | 说明 |
|-------------|-----------|------|
| 400 | `BAD_REQUEST` | 请求参数校验失败 |
| 401 | `UNAUTHORIZED` | 未认证（无 token） |
| 403 | `FORBIDDEN` | 无权限 |
| 404 | `NOT_FOUND` | 资源不存在 |
| 409 | `CONFLICT` | 业务冲突（如重复插入） |
| 500 | `INTERNAL_ERROR` | 服务器内部错误 |

---

## 二、楼栋管理（Buildings）

### 2.1 创建楼栋 — `POST /buildings`

**Request**：
```json
{
  "name": "男生宿舍楼1号楼",
  "gender": "male",
  "floor_count": 4,
  "room_per_floor": 10,
  "description": "2020年建成，共4层"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | ✅ | 楼栋名称 |
| `gender` | string | ✅ | `male` 或 `female` |
| `floor_count` | int | ✅ | 楼层数 |
| `room_per_floor` | int | ✅ | 每层房间数 |
| `description` | string | ❌ | 描述 |

**Response** `201 Created`：
```json
{
  "data": {
    "id": "uuid-string",
    "name": "男生宿舍楼1号楼",
    "gender": "male",
    "floor_count": 4,
    "room_per_floor": 10,
    "status": "active",
    "description": "2020年建成，共4层",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

### 2.2 获取楼栋 — `GET /buildings/:id`

**Response** `200 OK`：
```json
{
  "data": {
    "id": "uuid-string",
    "name": "男生宿舍楼1号楼",
    "gender": "male",
    "floor_count": 4,
    "room_per_floor": 10,
    "status": "active",
    "description": "2020年建成，共4层",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

### 2.3 楼栋列表 — `GET /buildings`

**Query Parameters**：
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | int | 1 | 页码 |
| `page_size` | int | 20 | 每页数量 |
| `gender` | string | — | 按性别过滤（male/female） |
| `status` | string | — | 按状态过滤（active/inactive） |

**Response** `200 OK`：
```json
{
  "data": [
    { "id": "...", "name": "...", "gender": "male", ... },
    { "id": "...", "name": "...", "gender": "female", ... }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

### 2.4 更新楼栋 — `PUT /buildings/:id`

**Request**：
```json
{
  "name": "男生宿舍楼1号楼（改造后）",
  "floor_count": 5,
  "status": "inactive"
}
```

**Response** `200 OK`：返回更新后的楼栋对象。

### 2.5 删除楼栋 — `DELETE /buildings/:id`

**Response** `204 No Content`（无 body）

---

## 三、房间管理（Rooms）

### 3.1 创建房间 — `POST /rooms`

**Request**：
```json
{
  "building_id": "uuid-string",
  "floor": 1,
  "number": "101",
  "beds_total": 4,
  "status": "available"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `building_id` | string | ✅ | 所属楼栋 ID |
| `floor` | int | ✅ | 楼层（1-indexed） |
| `number` | string | ✅ | 房间号（如 101、102） |
| `beds_total` | int | ✅ | 床位数（1-8） |
| `status` | string | ❌ | `available`/`occupied`/`maintenance`，默认 `available` |

**Response** `201 Created`：返回房间对象。

### 3.2 获取房间 — `GET /rooms/:id`

**Response** `200 OK`：返回房间对象（含楼栋信息）。

### 3.3 房间列表 — `GET /rooms`

**Query Parameters**：
| 参数 | 类型 | 说明 |
|------|------|------|
| `page`, `page_size` | int | 分页 |
| `building_id` | string | 按楼栋过滤 |
| `floor` | int | 按楼层过滤 |
| `status` | string | 按状态过滤 |
| `gender` | string | 按性别过滤（通过楼栋关联） |

**Response** `200 OK`：返回分页房间列表（含 beds_used 字段）。

### 3.4 更新房间 — `PUT /rooms/:id`

**Request**：
```json
{
  "beds_total": 6,
  "status": "maintenance"
}
```

**Response** `200 OK`：返回更新后的房间对象。

### 3.5 删除房间 — `DELETE /rooms/:id`

**Response** `204 No Content`。

---

## 四、学生管理（Students）

### 4.1 创建学生 — `POST /students`

**Request**：
```json
{
  "student_no": "2024001",
  "name": "张三",
  "gender": "male",
  "phone": "13800138000",
  "email": "zhangsan@example.com",
  "major": "计算机科学与技术",
  "grade": 2025
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `student_no` | string | ✅ | 学号（唯一） |
| `name` | string | ✅ | 姓名 |
| `gender` | string | ✅ | `male` 或 `female` |
| `phone` | string | ❌ | 联系电话 |
| `email` | string | ❌ | 邮箱 |
| `major` | string | ✅ | 专业 |
| `grade` | int | ✅ | 年级（2000-2100） |

### 4.2 学生详情 — `GET /students/:id`

**Response** `200 OK`：返回学生对象（含当前 allocation 信息如果有的话）。

### 4.3 学生列表 — `GET /students`

**Query Parameters**：
| 参数 | 类型 | 说明 |
|------|------|------|
| `page`, `page_size` | int | 分页 |
| `gender` | string | 按性别过滤 |
| `status` | string | 按状态过滤 |
| `building_id` | string | 按入住楼栋过滤（通过 allocation 关联） |
| `grade` | int | 按年级过滤 |

### 4.4 更新学生 — `PUT /students/:id`

**Request**：
```json
{
  "phone": "13900139000",
  "email": "new@example.com",
  "status": "graduated"
}
```

### 4.5 删除学生 — `DELETE /students/:id`

**Response** `204 No Content`。

---

## 五、入住分配（Allocations）

### 5.1 创建分配 — `POST /allocations`

**Request**：
```json
{
  "student_id": "uuid-string",
  "room_id": "uuid-string",
  "bed_number": 1
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `student_id` | string | ✅ | 学生 ID |
| `room_id` | string | ✅ | 房间 ID |
| `bed_number` | int | ✅ | 床位号（1 - room.capacity） |

**业务规则**：
- 学生不能有 active allocation（`CONFLICT: student already allocated`）
- 房间未满（`CONFLICT: room is full`）
- bed_number 不能超过房间 capacity
- 学生性别必须与楼栋性别一致

**Response** `201 Created`：返回 allocation 对象。

### 5.2 获取分配 — `GET /allocations/:id`

### 5.3 分配列表 — `GET /allocations`

**Query Parameters**：
| 参数 | 类型 | 说明 |
|------|------|------|
| `page`, `page_size` | int | 分页 |
| `student_id` | string | 按学生过滤 |
| `room_id` | string | 按房间过滤 |
| `status` | string | `active` 或 `checked_out` |
| `check_in_after` | date | 入住日期下限（ISO 8601） |

### 5.4 取消分配（退宿）— `DELETE /allocations/:id`

**Request**：
```json
{
  "reason": "学生毕业离校"
}
```

**Response** `200 OK`：返回更新后的 allocation（含 `check_out_at` 时间戳）。

---

## 六、违规记录（Violations）

### 6.1 创建违规 — `POST /violations`

**Request**：
```json
{
  "student_id": "uuid-string",
  "violation_type": "late_return",
  "description": "晚归超过22:00共计3次",
  "penalty": "口头警告"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `student_id` | string | ✅ | 学生 ID |
| `violation_type` | string | ✅ | 违规类型（见下表） |
| `description` | string | ✅ | 详细描述 |
| `penalty` | string | ❌ | 处罚措施 |

**Violation Types**：
| type | 说明 |
|------|------|
| `late_return` | 晚归 |
| `property_damage` | 财产损坏 |
| `noise_violation` | 噪音扰民 |
| `unauthorized_visit` | 违规来访 |
| `smoking` | 违规吸烟 |
| `other` | 其他 |

### 6.2 违规列表 — `GET /violations`

**Query Parameters**：
| 参数 | 类型 | 说明 |
|------|------|------|
| `student_id` | string | 按学生过滤 |
| `violation_type` | string | 按类型过滤 |
| `page`, `page_size` | int | 分页 |

### 6.3 删除违规 — `DELETE /violations/:id`

---

## 七、维修记录（Repairs）

### 7.1 创建维修 — `POST /repairs`

**Request**：
```json
{
  "room_id": "uuid-string",
  "reporter_name": "张三",
  "description": "热水器故障",
  "contact": "13800138000",
  "status": "pending"
}
```

### 7.2 维修列表 — `GET /repairs`

**Query Parameters**：`room_id`, `status`, `page`, `page_size`

### 7.3 更新维修 — `PUT /repairs/:id`

**Request**：
```json
{
  "status": "completed",
  "handler_name": "维修员李四",
  "result": "已更换加热元件"
}
```

---

## 八、健康检查

### `GET /health`

**Response** `200 OK`：
```json
{
  "status": "healthy",
  "timestamp": "2026-01-01T00:00:00Z",
  "services": {
    "database": "connected",
    "redis": "connected"
  }
}
```
# Student Management — Design Document

> 本文档定义学生管理的业务逻辑，供 AI Agent 实现功能时参考。

---

## 1. 业务概述

学生（Student）是宿舍管理系统的核心实体之一。每个学生有明确的性别、学号（唯一）、年级等信息，通过入住分配（Allocation）关联到具体房间。

---

## 2. 数据模型

### Student 结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键，由 service 层生成 |
| student_id | string | 学号，唯一约束 |
| name | string | 姓名 |
| gender | string | `male` 或 `female` |
| phone | string | 联系电话（可选） |
| email | string | 邮箱（可选） |
| major | string | 所修专业 |
| grade | int | 年级（2000-2100） |
| room_id | *string | 当前分配的宿舍ID（可选，Allocation 关联后填充） |
| check_in_at | *time.Time | 入住时间（可选） |
| status | string | `pending`（待入住）、`checked_in`（已入住）、`graduated`（已毕业）、`suspended`（休学） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

### 状态流转

```
pending → checked_in → graduated
              ↓
         suspended
```

| 当前状态 | 允许转入 | 说明 |
|---------|---------|------|
| `pending` | `checked_in`, `suspended` | 待入住学生可入住或休学 |
| `checked_in` | `graduated`, `suspended` | 已入住学生可毕业或休学退宿 |
| `graduated` | — | 终态 |
| `suspended` | `checked_in` | 休学后可复学重新入住 |

---

## 3. 业务规则

### 3.1 学号唯一性

`student_id`（学号）在全校范围内唯一。重复创建会触发数据库唯一约束错误，返回 `409 Conflict`。

### 3.2 性别与楼栋匹配

学生只能入住与其性别一致的楼栋（由 Room → Building 的 gender 字段保证）。业务层在 Allocation 分配时应校验。

### 3.3 状态与入住资格

| status | 可否分配房间 | 说明 |
|--------|------------|------|
| `pending` | ✅ | 待入住，可分配 |
| `checked_in` | ❌ | 已入住，不能重复分配（需先 Vacate） |
| `graduated` | ❌ | 已毕业，不能再入住 |
| `suspended` | ❌ | 休学状态，需复学后才能入住 |

### 3.4 删除学生（P3 — 待定）

禁止删除有 active allocation 的学生。返回 `409 Conflict`。

---

## 4. API 规范（与 API Contract 对齐）

### 4.1 创建学生 — `POST /api/v1/students`

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

**Errors**：
| 场景 | HTTP Status |
|------|------------|
| 学号已存在 | 409 Conflict |
| 必填字段缺失 | 400 Bad Request |
| `grade` 超范围 | 400 Bad Request |

### 4.2 获取学生 — `GET /api/v1/students/:id`

返回学生对象（含当前 allocation 信息）。

### 4.3 学生列表 — `GET /api/v1/students`

**Query Parameters**：

| 参数 | 类型 | 说明 |
|------|------|------|
| `page`, `page_size` | int | 分页 |
| `gender` | string | 按性别过滤 |
| `status` | string | 按状态过滤 |
| `grade` | int | 按年级过滤 |

### 4.4 更新学生 — `PUT /api/v1/students/:id`

**Request**：
```json
{
  "phone": "13900139000",
  "email": "new@example.com",
  "status": "graduated"
}
```

### 4.5 删除学生 — `DELETE /api/v1/students/:id`

返回 `204 No Content`。有 active allocation 时返回 `409 Conflict`。

### 4.6 分配房间 — `POST /api/v1/students/:id/allocate`

将学生分配到指定房间，触发 Allocation 创建 + Student.RoomID 更新。

### 4.7 退宿 — `POST /api/v1/students/:id/vacate`

学生退宿，释放 allocation + 清空 Student.RoomID。

---

## 5. 错误场景

| 场景 | 错误码 | HTTP Status |
|------|--------|------------|
| 学号已存在 | `CONFLICT` | 409 |
| 获取不存在的学生 | `NOT_FOUND` | 404 |
| `grade` 不在 2000-2100 范围 | `BAD_REQUEST` | 400 |
| 删除有 active allocation 的学生 | `CONFLICT` | 409 |
| 对 `checked_in` 状态学生再次分配 | `CONFLICT` | 409 |

---

## 6. 与其他模块的关系

- **Allocation**：学生通过 Allocation 与 Room 关联，一个 Allocation 记录一次入住周期
- **Room**：通过 `Student.RoomID` 指向当前入住房间
- **Violation**：学生可产生多条违规记录（由 Violation.student_id 关联）
- **Repair**：学生可以提交报修单（由 Repair.reporter_id 关联，但 reporter_id 实际上可以是任意人员 ID）

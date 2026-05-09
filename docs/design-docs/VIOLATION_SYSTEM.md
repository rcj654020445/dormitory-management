# 违规系统 — 设计文档

> 本文档定义学生违规记录的业务逻辑，供 AI Agent 实现功能时参考。

---

## 1. 业务概述

宿舍管理系统需要记录学生的违规行为，作为宿舍管理的重要依据。违规记录与学生（Student）绑定，支持多种违规类型、处罚措施，以及基本的统计分析。

---

## 2. 数据模型

### Violation 结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| student_id | UUID | 学生 ID（FK） |
| violation_type | enum | 违规类型（见 3.1） |
| description | string | 详细描述（必填） |
| penalty | string | 处罚措施 |
| recorded_at | timestamp | 记录时间（默认 now()） |
| recorded_by | string | 记录人（宿管姓名） |
| created_at | timestamp | 创建时间 |

---

## 3. 业务规则

### 3.1 违规类型

| type | 名称 | 说明 |
|------|------|------|
| `late_return` | 晚归 | 超过规定时间（22:00）返回宿舍 |
| `property_damage` | 财产损坏 | 损坏宿舍公共设施或他人财物 |
| `noise_violation` | 噪音扰民 | 在休息时间制造噪音影响他人 |
| `unauthorized_visit` | 违规来访 | 异性进入异性宿舍楼 |
| `smoking` | 违规吸烟 | 在宿舍楼内吸烟 |
| `unauthorized_stay` | 违规留宿 | 未经许可留宿外来人员 |
| `other` | 其他 | 其他违规行为 |

### 3.2 创建违规

**前置条件**：
- 学生必须存在且状态为 `active`
- `violation_type` 必须在上述类型列表中
- `description` 不能为空（最少 5 个字符）

**自动行为**：
- `recorded_at` 默认为当前时间
- `recorded_by` 从请求中获取（宿管系统登录用户）

### 3.3 违规次数统计

用于前端展示学生的违规记录次数：

```sql
SELECT student_id, COUNT(*) as violation_count
FROM violations
GROUP BY student_id
ORDER BY violation_count DESC
```

### 3.4 删除违规

- 违规记录一般**不允许删除**（保留历史）
- 如需删除，需管理员权限，理由必须填写
- 删除操作需要记录操作日志（`recorded_by` + `deleted_by`）

---

## 4. 业务场景

### 4.1 多次晚归预警

当学生 `late_return` 违规次数 ≥ 3 次时：
- 前端应在列表页显示警告标识
- 可触发通知（如有通知系统）

### 4.2 入住资格审查

在创建 Allocation 前，可以查询学生违规记录：
- 如果存在未处理的 `property_damage` 或其他重大违规
- 系统应在分配界面给出警告提示（不影响分配，仅提示）

---

## 5. 错误场景

| 场景 | 错误码 | HTTP Status |
|------|--------|------------|
| 学生不存在 | `STUDENT_NOT_FOUND` | 404 |
| 违规类型无效 | `INVALID_VIOLATION_TYPE` | 400 |
| description 长度不足 | `BAD_REQUEST` | 400 |
| 重复创建（同一学生同一时间同一类型） | `CONFLICT` | 409 |
# Allocation Flow — 设计文档

> 本文档定义学生入住分配的业务流程，供 AI Agent 实现功能时参考。

---

## 1. 业务流程概述

入住分配（Allocation）是学生入住宿舍的核心业务，将学生（Student）和房间（Room）绑定在一起，分配床位（Bed）。

```
学生申请入住 → 校验资格 → 选择房间 → 分配床位 → 登记入住
```

---

## 2. 数据模型

### Allocation 结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| student_id | UUID | 学生 ID（FK） |
| room_id | UUID | 房间 ID（FK） |
| bed_number | int | 床位号（1 - room.capacity） |
| status | enum | `active`（在住）、`checked_out`（已退宿） |
| check_in_at | timestamp | 入住时间 |
| check_out_at | timestamp | 退宿时间（nullable） |
| reason | string | 退宿原因（nullable） |
| created_at | timestamp | 创建时间 |

---

## 3. 入住流程（Create Allocation）

### 3.1 前置校验

```
Step 1: 校验学生是否存在，且状态为 active
Step 2: 校验学生是否已有 active allocation
        → 有：返回 409 Conflict（学生已有在住房间）
Step 3: 校验房间是否存在，状态为 available 或 occupied
        → 房间处于 maintenance：返回 409 Conflict
Step 4: 校验房间床位是否已满（beds_used >= capacity）
        → 已满：返回 409 Conflict（房间已满）
Step 5: 校验床位号是否有效（1 ≤ bed_number ≤ capacity）
Step 6: 校验学生性别与楼栋性别是否一致
        → 不一致：返回 409 Conflict（性别不匹配）
Step 7: 校验 bed_number 是否已被占用
        → 已占用：返回 409 Conflict（床位已被分配）
```

### 3.2 分配操作

所有校验通过后，执行以下操作（原子性事务）：

1. 在 `allocations` 表插入记录（status = `active`, check_in_at = now()）
2. 更新 `students` 表：将 `room_id` 设为分配的房间（用于快速查询当前住所）
3. 更新 `rooms` 表：`beds_used = beds_used + 1`
   - 如果 `beds_used + 1 == capacity`，自动将 `status` 设为 `occupied`

### 3.3 异常处理

如果步骤 2 或 3 失败，需要回滚已完成的操作。

---

## 4. 退宿流程（Cancel Allocation）

### 4.1 前置校验

```
Step 1: 校验 allocation 是否存在
Step 2: 校验 allocation.status == 'active'
        → 非 active：返回 409 Conflict（该分配已退宿）
```

### 4.2 退宿操作

1. 更新 `allocations` 表：
   - `status` → `checked_out`
   - `check_out_at` → now()
   - `reason` → 请求中的退宿原因
2. 更新 `students` 表：将 `room_id` 设为 null
3. 更新 `rooms` 表：
   - `beds_used = beds_used - 1`
   - 如果 `beds_used < capacity`，自动将 `status` 设为 `available`

---

## 5. 床位号设计

每个房间有 1-6 个床位（由 `room.capacity` 决定）。

**床位编号规则**：
- 上床位号 > 下床位号（如果房间有上下铺设计）
- 无上下铺时，按从门到窗的顺序编号

**前端展示**：
- 房间详情页展示床位图（1-6 个床位格）
- 已分配床位显示学生姓名
- 未分配床位显示为空闲状态

---

## 6. 房间状态与床位联动

```
Room.beds_used 计算：
  SELECT COUNT(*) FROM allocations
  WHERE room_id = ? AND status = 'active'

Room.status 规则：
  if beds_used == 0       → 'available'
  if 0 < beds_used < capacity → 'occupied'（部分入住）
  if beds_used == capacity    → 'full'
```

---

## 7. 错误场景

| 场景 | 错误码 | 说明 |
|------|--------|------|
| 学生已有 active allocation | `STUDENT_ALREADY_ALLOCATED` | 同一学生不能同时入住多个房间 |
| 房间已满 | `ROOM_FULL` | 找不到可用床位 |
| 床位已被占用 | `BED_OCCUPIED` | 指定床位号已被分配 |
| 性别不匹配 | `GENDER_MISMATCH` | 学生性别与楼栋性别不一致 |
| 房间在维护中 | `ROOM_UNDER_MAINTENANCE` | 房间不可分配 |
| 重复退宿 | `ALLOCATION_ALREADY_CLOSED` | allocation 已退宿 |
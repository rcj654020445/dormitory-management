# 楼栋管理 - 设计文档

> 本文档定义楼栋管理的业务逻辑，供 AI Agent 实现功能时参考。

---

## 1. 业务概述

楼栋（Building）是宿舍管理系统的最顶层物理资源。每个楼栋有明确的性别属性（male/female），学生必须入住与其性别一致的楼栋。

---

## 2. 数据模型

### Building 结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| name | string | 楼栋名称，如"男生宿舍楼1号楼" |
| gender | enum | `male` 或 `female`，楼栋性别属性 |
| floor_count | int | 楼层数 |
| room_per_floor | int | 每层房间数 |
| status | enum | `active`（在用）、`inactive`（停用） |
| description | string | 备注描述 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

### 衍生计算

```
总房间数 = floor_count × room_per_floor
总床位数 = 总房间数 × 平均容量（由各房间 capacity 决定）
```

---

## 3. 业务规则

### 3.1 创建楼栋

**前置条件**：
- `floor_count` ≥ 1
- `room_per_floor` ≥ 1
- `gender` 必须为 `male` 或 `female`

**自动行为**：
- 系统**不会**自动创建房间（Rooms 表为空），需要手动创建或通过批量工具生成
- 建议创建完楼栋后立即创建房间：`POST /rooms` 批量创建

**示例**：创建 4 层、每层 10 间房的楼栋，需要手动创建 40 个 Room 记录。

### 3.2 楼栋与房间的关系

- **一对多**：一个楼栋对应多个房间
- 楼栋的 `gender` 决定了只能接收该性别的学生
- 房间创建时必须指定 `building_id`，由系统保证一致性

### 3.3 楼栋状态

| status | 含义 | 对学生入住的影响 |
|--------|------|----------------|
| `active` | 正常使用中 | 可以分配学生入住 |
| `inactive` | 停用/装修中 | 不允许新增分配，既住户不受影响 |

### 3.4 删除楼栋

- 禁止删除有房间的楼栋（有 active allocation 的房间也不能删）
- 删除楼栋前必须先删除其所有房间
- 返回 `409 Conflict` 如果条件不满足

---

## 4. 房间生成逻辑（辅助）

批量创建房间的标准模式（可在 seed 或管理工具中实现）：

```python
for floor in range(1, building.floor_count + 1):
    for room_num in range(1, building.room_per_floor + 1):
        room_number = f"{floor:02d}{room_num:02d}"  # 101, 102, ... 410
        capacity = 4 if floor < building.floor_count else 6  # 顶楼可设6人间
        create_room(building_id, floor, room_number, capacity)
```

---

## 5. 错误场景

| 场景 | 错误码 | HTTP Status |
|------|--------|------------|
| 创建楼栋 name 为空 | `BAD_REQUEST` | 400 |
| gender 不是 male/female | `BAD_REQUEST` | 400 |
| 删除有房间的楼栋 | `CONFLICT` | 409 |
| 获取不存在的楼栋 | `NOT_FOUND` | 404 |

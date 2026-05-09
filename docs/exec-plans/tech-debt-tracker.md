# Tech Debt Tracker

## Active Debt

| ID | Description | Priority | Owner | Created |
|----|-------------|----------|-------|---------|
| TD-001 | DB schema 与代码不一致：DDL 用 `enrollment_year`，代码用 `grade`；DDL 缺 `email`、`grade`、`room_id` 列；多了 `emergency_contact` | P0 | — | 2026-05-08 |
| TD-002 | DB schema `students` 表 `student_id` 列名应为 `student_no` | P0 | — | 2026-05-08 |
| TD-003 | 所有 Service Create 方法未生成 UUID，导致 INSERT 的 id 为空字符串 | P0 | — | 2026-05-08 |
| TD-004 | Building/Room/Student Entity 的 `CreatedAt`/`UpdatedAt` 定义为 `string`，但 PostgreSQL 存 `timestamptz`，scan 时类型不匹配 | P0 | — | 2026-05-08 |
| TD-005 | `Student.Grade` DTO validation `max=10` 应为 `max=2100`（年份范围） | P1 | — | 2026-05-08 |
| TD-006 | `CreateRoomRequest.Capacity` 应为 `beds_total`（与业务字段名一致） | P1 | — | 2026-05-08 |
| TD-007 | `API_CONTRACT.md` Student 章节仍使用 `enrollment_year`/`emergency_contact`，与实际代码不一致 | P1 | — | 2026-05-08 |
| TD-008 | `API_CONTRACT.md` Room 章节使用 `capacity`，与实际 `beds_total` 不一致 | P2 | — | 2026-05-08 |
| TD-009 | `cmd/migrate/main.go` 使用硬编码内联 DDL，未使用 `migrations/*.sql` 文件规范管理 | P2 | — | 2026-05-08 |
| TD-010 | Building/Room/Student/Repair 的 Repository INSERT 语句字段顺序与 DB 建表 DDL 顺序不一致（潜在可维护性问题） | P2 | — | 2026-05-08 |
| TD-011 | `student_repo.go` INSERT/SELECT 仍使用旧的 `student_id` 列名（已在内存中修复，待重新 migrate 验证） | P0 | — | 2026-05-08 |
| TD-012 | `violation_repo.go`、`allocation_repo.go`、`repair_repo.go` 是否存在同样的列名/字段不匹配问题（未验证） | P1 | — | 2026-05-08 |

## Resolved Debt

| ID | Description | Resolution | Date |
|----|-------------|-----------|------|
| RD-001 | DDL `students` 表列名从 `student_id` 改为 `student_no` | 修改 `cmd/migrate/main.go` DDL + 修改 `student_repo.go` 所有 SQL | 2026-05-08 |
| RD-002 | DDL `students` 表增加 `email`、`grade`、`room_id` 列；移除 `enrollment_year`、`emergency_contact` | 修改 `cmd/migrate/main.go` DDL | 2026-05-08 |
| RD-003 | 所有 Service Create 方法添加 `uuid.New().String()` ID 生成 | 在 `building_svc.go`、`room_svc.go`、`student_svc.go`、`violation_svc.go`、`allocation_svc.go` 中添加 `entity.ID = uuid.New().String()` | 2026-05-08 |
| RD-004 | `building_repo.go`/`room_repo.go` scan 方法增加 `time.Time` 中间变量接收 DB `timestamptz`，再 Format 为 string 赋值 | 在 `scanBuilding`/`scanRoom` 中添加 `createdAt, updatedAt time.Time` 中间变量 | 2026-05-08 |
| RD-005 | `Student.Grade` validation `max=10` → `max=2100` | 修改 `internal/request/student.go` | 2026-05-08 |
| RD-006 | `CreateRoomRequest.Capacity` 改为 `beds_total` + `UpdateRoomRequest.Capacity` → `BedsTotal` | 修改 `internal/request/room.go` + 修复 `room_svc.go` 引用 | 2026-05-08 |
| RD-007 | `repair_svc.go` 已生成 UUID 但未赋值给 `entity.ID`（已修复） | 在 `Create` 方法中添加 `repair.ID = generateUUID()` | 2026-05-08 |
| RD-008 | `building_svc.go` 多余 `time` import 已移除 | 移除 `time` import | 2026-05-08 |
| RD-009 | `room_svc.go` `GetRoomTypeCapacity` 定义后从未调用（死代码），留待后续清理 | P3，不阻塞功能 | 2026-05-08 |
| RD-010 | `API_CONTRACT.md` 中 Repair 章节 `reporter_id` 校验逻辑未说明（需学生 ID，但文档未注明） | 待补充文档 | — |

## How to Track Debt

When you notice tech debt during development:
1. Add a row to the Active Debt table
2. Assign a priority (P0-P3)
3. Reference it in the relevant execution plan if applicable

### Priority Levels

- **P0**: Critical — blocks major features or causes data corruption
- **P1**: High — significant impact, should address soon
- **P2**: Medium — worth fixing but not urgent
- **P3**: Low — nice to have, can defer indefinitely

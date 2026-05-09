# 设计文档索引

| 文档 | 状态 | 最后更新 | 描述 |
|----------|--------|---------------|-------------|
| [STUDENT_MANAGEMENT.md](STUDENT_MANAGEMENT.md) | ✅ Complete | 2026-05-08 | 学生管理数据模型、状态机、CRUD API、学号唯一性约束 |
| [BUILDING_MANAGEMENT.md](BUILDING_MANAGEMENT.md) | ✅ Complete | 2026-05-08 | 楼栋管理业务逻辑、数据模型、业务规则 |
| [ALLOCATION_FLOW.md](ALLOCATION_FLOW.md) | ✅ Complete | 2026-05-08 | 入住分配流程、校验规则、房间床位联动 |
| [VIOLATION_SYSTEM.md](VIOLATION_SYSTEM.md) | ✅ Complete | 2026-05-08 | 违规记录类型、业务规则、错误场景 |
| [ROOM_MANAGEMENT.md](ROOM_MANAGEMENT.md) | ✅ Complete | 2026-05-08 | 房间管理、房型/床位联动、业务规则 |
| [REPAIR_SYSTEM.md](REPAIR_SYSTEM.md) | ✅ Complete | 2026-05-08 | 维修系统状态机、报修/派单/维修/评价流程 |

## 如何添加设计文档

1. 创建新文件: `docs/design-docs/{component-name}.md`
2. 使用以下结构:
   - 业务概述
   - 数据模型（字段表）
   - 业务规则（含错误码）
   - API 关联（如有）
3. 将条目添加到本索引
4. 如果是重要组件，请从 AGENTS.md 添加链接
# Harness Trace Format

> Trace 是 harness-executor 执行任务时的完整操作记录，用于复盘、调试和验证。

## Overview

```
harness/
├── trace/                          # 执行痕迹根目录
│   ├── [task-name]-[timestamp]/    # 每次任务运行的独立目录
│   │   ├── state/                  # 状态快照
│   │   │   ├── 00_initial.json     # 任务开始时
│   │   │   ├── 01_after_phase.json # 各阶段完成后
│   │   │   └── ...
│   │   ├── checkpoints/             # 检查点（关键决策点）
│   │   │   └── [n]-[description].json
│   │   ├── memory/                  # Agent 记忆状态
│   │   │   └── memory.json
│   │   ├── verification-report.json # 验证报告
│   │   └── summary.json            # 执行摘要
│   └── current -> [latest]/        # 符号链接指向最新任务
└── config/
    └── environment.json             # 运行时环境配置
```

## Trace Entry Types

### 1. State Snapshot (`state/`)

任务执行过程中关键节点的状态快照。

```json
{
  "_meta": {
    "type": "state_snapshot",
    "version": "1.0",
    "timestamp": "2026-05-08T13:30:00Z"
  },
  "phase": "phase_2_service_creation",
  "step": 3,
  "files_created": [
    "internal/service/room_service.go",
    "internal/handler/room_handler.go"
  ],
  "files_modified": [
    "internal/router/router.go"
  ],
  "lint_results": {
    "lint-deps": "passed",
    "lint-quality": "passed"
  },
  "build_status": "passed"
}
```

**字段说明：**

| 字段 | 类型 | 描述 |
|------|------|------|
| `_meta.type` | string | 固定为 `state_snapshot` |
| `_meta.version` | string | 格式版本号 |
| `_meta.timestamp` | ISO8601 | 快照时间 |
| `phase` | string | 当前阶段名 |
| `step` | int | 阶段内步骤编号 |
| `files_created` | string[] | 本次创建的文件路径 |
| `files_modified` | string[] | 本次修改的文件路径 |
| `lint_results` | object | 各 linter 的执行结果 |
| `build_status` | string | 构建状态：`passed` \| `failed` \| `not_run` |

---

### 2. Checkpoint (`checkpoints/`)

关键决策点或需要确认的节点。Agent 在此处暂停等待验证。

```json
{
  "_meta": {
    "type": "checkpoint",
    "version": "1.0",
    "timestamp": "2026-05-08T13:35:00Z"
  },
  "id": "01",
  "description": "API endpoint structure verified",
  "trigger": "after_handler_creation",
  "verification": {
    "method": "http_test",
    "endpoint": "/api/v1/rooms",
    "status_code": 200,
    "response_fields": ["id", "building_id", "room_number", "beds_total", "beds_used"]
  },
  "status": "passed",
  "notes": "All CRUD endpoints functional"
}
```

**字段说明：**

| 字段 | 类型 | 描述 |
|------|------|------|
| `_meta.type` | string | 固定为 `checkpoint` |
| `id` | string | 检查点序号（补零，如 `01`, `02`） |
| `description` | string | 检查点描述 |
| `trigger` | string | 触发时机 |
| `verification.method` | string | 验证方式：`http_test` \| `file_check` \| `command` \| `manual` |
| `verification.endpoint` | string | HTTP 端点（http_test 时） |
| `verification.status_code` | int | 期望状态码 |
| `verification.response_fields` | string[] | 期望存在的响应字段 |
| `status` | string | `passed` \| `failed` \| `skipped` |
| `notes` | string | 额外说明 |

---

### 3. Memory State (`memory/`)

Agent 记忆状态快照，记录当前已完成的 context 和决策。

```json
{
  "_meta": {
    "type": "memory_snapshot",
    "version": "1.0",
    "timestamp": "2026-05-08T13:40:00Z"
  },
  "context": {
    "task_goal": "Implement Room CRUD with repair status management",
    "current_phase": "phase_3_testing",
    "phase_index": 2
  },
  "learned": {
    "student_dto_uses_student_no": true,
    "allocation_uses_checkin_checkout": true,
    "repair_status_flow": "pending→assigned→repairing→completed"
  },
  "decisions": [
    {
      "id": "D-01",
      "at": "2026-05-08T13:32:00Z",
      "decision": "Use CheckInAt/CheckOutAt for Allocation instead of StartDate/EndDate",
      "reason": "Consistent with existing types and API design"
    }
  ],
  "pending_tasks": [
    "Write unit tests for room service",
    "Update router with new endpoints"
  ]
}
```

---

### 4. Verification Report (`verification-report.json`)

任务完成后的最终验证报告。

```json
{
  "_meta": {
    "type": "verification_report",
    "version": "1.0",
    "generated_at": "2026-05-08T14:00:00Z"
  },
  "overall_status": "passed",
  "verification_sections": {
    "build": {
      "status": "passed",
      "command": "go build ./cmd/server",
      "duration_ms": 4521
    },
    "layer_check": {
      "status": "passed",
      "command": "go run ./scripts/lint-deps",
      "violations": 0
    },
    "quality_check": {
      "status": "passed",
      "command": "go run ./scripts/lint-quality"
    },
    "tests": {
      "status": "passed",
      "command": "go test ./internal/service/... -run Room",
      "passed": 12,
      "failed": 0
    }
  },
  "local_verification_required": {
    "description": "Commands requiring host machine with Go toolchain:",
    "commands": [
      "go build ./cmd/server",
      "go run ./scripts/lint-deps --json",
      "go test ./internal/service/... -run Room -coverprofile=coverage.out"
    ]
  }
}
```

**注意**：Docker 容器内无 Go toolchain（`go` 命令不可用），代码验证必须在宿主机执行。

---

### 5. Summary (`summary.json`)

执行摘要，人可读。

```json
{
  "_meta": {
    "type": "execution_summary",
    "version": "1.0",
    "generated_at": "2026-05-08T14:05:00Z"
  },
  "task_id": "room-repair-crud-20260508-1330",
  "task_type": "boilerplate",
  "status": "completed",
  "duration_seconds": 1847,
  "phases_completed": [
    "phase_1_planning",
    "phase_2_code_generation",
    "phase_3_testing",
    "phase_4_integration"
  ],
  "files_created": 8,
  "files_modified": 3,
  "overall_score": {
    "correctness": 0.95,
    "style": 0.90,
    "completeness": 1.0,
    "overall": 0.95
  }
}
```

---

## Directory Naming Convention

每次任务运行创建独立目录：

```
[task-name]-[YYYYMMDD]-[HHMM]-[YYYYMMDD]-[HHMM]
```

例如：`room-repair-crud-20260508-1330-20260508-1330`

- 前缀：task name（来自 task yaml 的 filename 或 task_id）
- 中间：任务开始时间
- 后缀：任务结束时间

`harness/trace/current` 符号链接始终指向最新任务目录。

---

## Versioning

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0 | 2026-05-08 | 初始版本 |

---

## Implementation Notes

### Why JSON?

- **机器可读**：便于 CI/CD 集成和自动化验证
- **人类可读**：JSON 格式化后可直接查看
- **可扩展**：`_meta` 封装版本信息，方便格式演进

### Trace is Append-Only

Trace 文件一旦写入**不修改**，只追加新条目。这确保了完整的历史可追溯性。

### Performance Consideration

对于大型任务（100+ 文件），state snapshot 可能会很大。使用 `jq` 过滤需要的字段：

```bash
# 只看整体状态
cat harness/trace/current/summary.json | jq '{status, duration_seconds, files_created}'

# 只看违规
go run ./scripts/lint-deps --json | jq '.violations[] | {file, package}'
```

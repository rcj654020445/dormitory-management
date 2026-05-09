# 质量评分

最后更新：2026-05-07

## 1 领域等级

| 领域 | 等级 | 问题数 | 趋势 | 备注 |
|------|------|--------|------|------|
| Backend Architecture | N/A | 0 | — | 新项目 |
| Frontend Architecture | N/A | 0 | — | 新项目 |
| Data Access | N/A | 0 | — | 新项目 |
| API Layer | N/A | 0 | — | 新项目 |
| Business Logic | N/A | 0 | — | 新项目 |

## 2 架构层级

| Layer | Coverage | Lint Pass | Doc Fresh |
|-------|----------|-----------|-----------|
| Types (L0) | N/A | ✓ | ✓ |
| Model (L0) | N/A | ✓ | ✓ |
| Repository (L1) | N/A | ✓ | ✓ |
| Cache (L1) | N/A | ✓ | ✓ |
| Service (L2) | N/A | ✓ | ✓ |
| Handler (L3) | N/A | ✓ | ✓ |
| Middleware (L3) | N/A | ✓ | ✓ |
| Request (L3) | N/A | ✓ | ✓ |
| Response (L3) | N/A | ✓ | ✓ |
| Entry Points (L4) | N/A | ✓ | ✓ |
| Infrastructure (L-1) | N/A | ✓ | ✓ |

## 3 黄金原则

1. 结构化日志优于原始 print/log 语句
2. 带上下文的类型化错误
3. 文件不超过 500 行 (Go) / 300 行 (Vue)
4. Linter 强制执行层级层次结构
5. 所有公共 API 必须有文档
6. 版本控制中不得包含 secrets

> 由以下脚本强制执行：[`scripts/lint-quality.go`](scripts/lint-quality.go)

## 4 质量趋势

| Date | Overall | Notes |
|------|---------|-------|
| 2026-05-07 | Baseline | 项目使用 harness 基础设施初始化 |

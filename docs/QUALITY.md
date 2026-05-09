# 质量标准

由 `scripts/lint-deps.go` 和 `scripts/lint-quality.go` 强制执行的黄金原则。

## 1 结构化日志

> 强制执行：[`scripts/lint-quality.go`](scripts/lint-quality.go)

```go
// ✓ 正确 — 结构化、可解析、可查询
zapLogger.Info("operation completed", zap.String("key", value))

// ✗ 错误 — 非结构化，难以解析
log.Printf("operation completed: %v", value)
log.Println("operation completed")
```

## 2 层级结构

> 强制执行：[`scripts/lint-deps.go`](scripts/lint-deps.go)

每个包都属于某个层级。上层可以导入下层，但绝不能反向导入。

```
pkg/* (L-1) ← 任意层均可导入
cmd/* (L4) ← 顶层
internal/handler/*, middleware/*, request/*, response/* (L3)
internal/service/* (L2)
internal/repository/*, cache/* (L1)
internal/types/*, model/* (L0) ← 底层
```

## 3 错误处理

### 3.1 使用类型化错误

```go
// ✓ 正确 — 类型化、机器可读
return nil, types.NewNotFoundError("student")

// ✗ 错误 — 字符串类型
return nil, fmt.Errorf("student not found")
```

### 3.2 包装时保留上下文

```go
// ✓ 正确 — 保留错误链
return fmt.Errorf("reading config: %w", err)

// ✗ 错误 — 丢失上下文
return err
```

## 4 文件大小限制

> 强制执行：[`scripts/lint-quality.go`](scripts/lint-quality.go)

- **Go 文件**：每个文件最多 **500 行**
- **Vue 文件**：每个文件最多 **300 行**
- 将大文件拆分为专注于特定功能的模块
- 当文件接近 400 行时，规划拆分方案

## 5 命名规范

| Category | Convention | Example |
|----------|-----------|---------|
| Go packages | 小写 | `repository`, `middleware` |
| Go types | PascalCase | `StudentService` |
| Go functions | PascalCase (导出), camelCase (未导出) | `GetStudent`, `createStudent` |
| TypeScript interfaces | PascalCase | `Student`, `ApiResponse` |
| Vue components | PascalCase | `StudentList.vue` |
| Database tables | snake_case | `students`, `room_allocations` |

## 6 执行

```bash
# 运行架构检查
make lint-arch

# 运行所有检查器
make lint

# 前端检查器
cd frontend && npm run lint
```

## 7 前端质量规则

> 强制执行：[`scripts/lint-quality.ts`](scripts/lint-quality.ts)

- 非测试文件中不允许使用 `console.log` / `console.error`
- 禁止使用裸 `any` 类型，除非有说明性注释
- `.ts` 文件最多 500 行，`.vue` 文件最多 300 行
- 使用 `@/` 路径别名进行内部导入（不要使用相对路径）
# 安全策略

## 1 原则

1. **绝不提交 secrets** — API keys、密码、tokens 永不进入版本控制
2. **验证所有外部输入** — 用户数据在处理前进行验证
3. **使用参数化查询** — 通过 pgx prepared statements 防止 SQL 注入
4. **记录安全相关事件** — 认证失败、权限违规
5. **保持依赖更新** — 定期执行 `go mod tidy` 和 `npm audit`

## 2 敏感文件

不要提交以下文件：

```
.env                    # 实际 secrets
*.pem, *.key            # 私钥
credentials.json        # API 凭证
docker-compose.yml      # 可能包含 env vars 中的 secrets
```

使用 `.env.example` 并填写占位值作为参考。

## 3 认证与授权

- JWT tokens 用于 API 认证
- Tokens 24 小时后过期
- Refresh token 流程（待实现）
- 基于角色的访问控制：`admin`、`staff`、`student`

### 3.1 密码要求

- 最少 8 个字符
- Bcrypt 哈希，成本因子为 12
- 不存储明文密码

## 4 输入验证

所有用户输入使用 `go-playground/validator` 进行验证：

```go
type CreateStudentRequest struct {
    StudentID string `json:"student_id" binding:"required"`
    Name      string `json:"name" binding:"required"`
    Gender    string `json:"gender" binding:"required,oneof=male female"`
    Phone     string `json:"phone" binding:"required"`
    Email     string `json:"email" binding:"required,email"`
    Major     string `json:"major" binding:"required"`
    Grade     int    `json:"grade" binding:"required,min=1,max=10"`
}
```

## 5 CORS 配置

CORS origins 通过环境变量配置：

```bash
CORS_ORIGINS=http://localhost:3000,https://admin.example.com
```

生产环境中仅允许已配置的 origins。

## 6 依赖管理

```bash
# Go: 检查漏洞
go mod verify
govulncheck ./...

# Node: 审计依赖
cd frontend && npm audit
```

## 7 报告安全问题

安全问题请联系：security@example.com
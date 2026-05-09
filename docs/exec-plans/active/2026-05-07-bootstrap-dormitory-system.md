# 引导：宿舍管理系统

|**创建日期**: 2026-05-07
|**任务ID**: bootstrap-dormitory-system-20260507-1505

## 目标
将 Go + Vue3 宿舍管理系统从零启动到正常运行：
- 下载 Go 依赖
- 运行数据库迁移
- 初始化种子数据
- 验证服务器启动并正常响应

## 范围
- **待修改文件**: go.mod, go.sum, .env
- **待创建文件**: go.sum (由 go mod tidy 生成)
- **服务**: PostgreSQL (docker), Redis (docker)

## 阶段

### 阶段 1：下载依赖
- [ ] 1.1 修复损坏的导入路径 (github.comgin 拼写错误)
- [ ] 1.2 使用 goproxy.cn 运行 `go mod tidy`
- [ ] 1.3 安装前端 npm 依赖
- **验证方式**: `go build ./...` (必须成功)

### 阶段 2：数据库设置
- [ ] 2.1 运行 `docker-compose up -d postgres redis`
- [ ] 2.2 运行 `go run ./cmd/migrate up`
- [ ] 2.3 运行 `go run ./cmd/seed`
- **验证方式**: `psql $DATABASE_URL -c "SELECT 1"` (必须成功)

### 阶段 3：服务器验证
- [ ] 3.1 启动服务器：`go run ./cmd/server`
- [ ] 3.2 健康检查：`curl http://localhost:8080/health`
- [ ] 3.3 正常停止服务器
- **验证方式**: /health 端点返回 HTTP 200

## 经验总结
- goproxy.cn 解决 Go 模块的中国网络问题
- 导入路径拼写错误 "github.comgin-gonic" 必须在 go mod tidy 之前修复

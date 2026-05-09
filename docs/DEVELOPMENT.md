# 开发环境搭建

## 1 前提条件

- **Go**: 1.22+
- **Node.js**: 20+
- **PostgreSQL**: 15+（通过 Docker 或本地安装）
- **Redis**: 7+（可选，用于缓存）
- **Make**: 任意近期版本

## 2 快速开始

```bash
# 1. 克隆代码并设置环境
make setup

# 2. 启动开发服务器（后端 + 前端）
make run

# 3. 运行测试
make test

# 4. 运行 linters
make lint
```

## 3 构建命令

### 后端（Go）

| Command | Description | Duration |
|---------|-------------|----------|
| `make build` | Build the server binary | ~5s |
| `make test` | Run all tests | ~10s |
| `make lint-arch` | Run architecture linters | ~3s |
| `make lint` | Run all linters (arch + quality) | ~5s |
| `make clean` | Remove build artifacts | ~1s |

### 前端（Vue3）

```bash
cd frontend
npm install          # 安装依赖
npm run dev          # 启动开发服务器（端口 3000）
npm run build        # 生产环境构建
npm run lint         # 运行 ESLint
npm run test         # 运行 Vitest
```

## 4 项目结构

```
.
├── cmd/                    # 入口点（Layer 4）
│   ├── server/main.go      # HTTP 服务器
│   ├── migrate/main.go     # 数据库迁移
│   └── seed/main.go        # 测试数据初始化
├── internal/               # 私有包
│   ├── types/             # 核心类型（Layer 0）
│   ├── model/             # 数据库实体（Layer 0）
│   ├── repository/        # 数据访问（Layer 1）
│   ├── cache/             # Redis 缓存（Layer 1）
│   ├── service/           # 业务逻辑（Layer 2）
│   ├── handler/           # HTTP 处理器（Layer 3）
│   ├── middleware/        # HTTP 中间件（Layer 3）
│   ├── request/           # 请求 DTO（Layer 3）
│   └── response/          # 响应构建器（Layer 3）
├── pkg/                    # 共享基础设施（Layer -1）
│   ├── logger/            # 结构化日志
│   ├── config/            # 配置加载
│   └── database/          # 数据库连接
├── frontend/               # Vue3 前端
│   └── src/
│       ├── types/         # TypeScript 接口（Layer 0）
│       ├── api/           # API 客户端（Layer 1）
│       ├── stores/        # Pinia 状态管理（Layer 2）
│       ├── views/         # 页面组件（Layer 3）
│       └── components/    # UI 组件（Layer 4）
├── harness/               # Harness 基础设施
│   ├── config/            # 环境契约
│   └── scripts/           # 设置/拆除脚本
├── scripts/               # Linters
│   ├── lint-deps.go       # Go 架构 linter
│   ├── lint-quality.go    # Go 质量 linter
│   ├── lint-deps.ts       # Vue3 架构 linter
│   └── lint-quality.ts    # Vue3 质量 linter
└── docs/                  # 文档
```

## 5 环境变量

将 `.env.example` 复制到 `.env` 并进行配置：

```bash
cp .env.example .env
# 使用实际值编辑 .env
```

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `DATABASE_URL` | — | **Yes** | PostgreSQL 连接字符串 |
| `DB_PASSWORD` | — | **Yes** | 数据库密码 |
| `JWT_SECRET` | — | **Yes** | JWT 签名密钥 |
| `PORT` | `8080` | No | 服务器端口 |
| `ENV` | `development` | No | 环境 |
| `LOG_LEVEL` | `debug` | No | 日志详细程度 |
| `REDIS_URL` | `redis://localhost:6379` | No | Redis 连接 |
| `CORS_ORIGINS` | `http://localhost:3000` | No | 允许的 CORS 源 |

## 6 数据库设置

```bash
# 运行迁移
go run ./cmd/migrate up

# 初始化开发数据
go run ./cmd/seed

# 回滚迁移
go run ./cmd/migrate down
```

## 7 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 启动带有热重载的开发服务器
npm run dev

# 运行类型检查
npm run type-check

# 生产环境构建
npm run build
```

API 代理已配置为将 `/api/*` 请求转发到 `http://localhost:8080`。

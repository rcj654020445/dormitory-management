# Dormitory Management System / 学生宿舍管理系统

一个基于 **Go** (REST API) + **Vue3** (前端) 构建的学生宿舍管理系统。

## Features / 功能特性

- 学生管理（入住登记、入住/退宿）
- 楼栋与房间管理
- 自动房间分配
- 费用管理
- 查寝记录

## Tech Stack / 技术栈

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22 + Gin |
| Frontend | Vue3 + TypeScript + Vite |
| Database | PostgreSQL 15 |
| Cache | Redis 7 |
| Auth | JWT |

## Quick Start / 快速开始

```bash
# 环境配置
make setup

# 运行开发服务器
make run

# 运行测试
make test

# 运行代码检查
make lint
```

完整的配置说明请参阅 [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)。

## Architecture / 系统架构

```
dormitory-management/
├── cmd/                 # 入口点 (server, migrate, seed)
├── internal/            # 私有包 (分层 L0-L4)
│   ├── types/          # 核心类型定义
│   ├── model/          # 数据库实体
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑
│   └── handler/        # HTTP 处理器
├── pkg/                 # 共享基础设施包
├── frontend/            # Vue3 单页应用
├── harness/             # 测试框架基础设施
├── scripts/             # 代码检查脚本
└── docs/                # 文档
```

更多代理相关文档请参阅 [AGENTS.md](AGENTS.md)。

## License / 许可证

MIT

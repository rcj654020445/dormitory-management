# =============================================
# 学生宿舍管理系统 - Harness Makefile
# =============================================
#
# 使用方法：
#   make help          # 查看所有可用命令
#   make all           # 运行所有检查（后端+前端）
#   make verify-be        # 仅后端验证（构建+架构+质量+测试）
#   make verify-fe     # 仅前端验证（类型检查+lint+构建）
#   make quick         # 快速自检（推荐提交前运行）
#
# =============================================

.PHONY: help all verify-be verify-fe deps-be deps-fe \
        build-be build-fe lint-be lint-quality-be lint-fe \
        test-be type-check-fe clean setup quick

.DEFAULT_GOAL := help

# ============ 帮助信息 ============
help:
	@echo ""
	@echo "🏠 学生宿舍管理系统 - Harness Makefile"
	@echo ""
	@echo "📦 完整验证（后端+前端）："
	@echo "   make all          运行全部检查"
	@echo ""
	@echo "🔧 分模块验证："
	@echo "   make verify-be       后端：构建 + 架构 + 质量 + 测试"
	@echo "   make verify-fe    前端：类型检查 + lint + 构建"
	@echo ""
	@echo "📥 安装依赖："
	@echo "   make deps-be      安装 Go 依赖"
	@echo "   make deps-fe      安装前端依赖"
	@echo "   make setup        安装所有依赖（后端+前端）"
	@echo ""
	@echo "🔨 单独步骤："
	@echo "   make build-be          Go 构建"
	@echo "   make lint-be           Go 架构层级检查"
	@echo "   make lint-quality-be   Go 代码质量检查"
	@echo "   make test-be           Go 单元测试"
	@echo "   make type-check-fe     前端 TypeScript 类型检查"
	@echo "   make lint-fe           前端 ESLint"
	@echo "   make build-fe          前端构建"
	@echo ""
	@echo "⚡ 快速自检："
	@echo "   make quick        推荐提交前运行（构建+架构+测试）"
	@echo ""
	@echo "🧹 清理："
	@echo "   make clean        清理构建产物"
	@echo ""

# ============ 安装依赖 ============
setup: deps-be deps-fe
	@echo "✅ 依赖安装完成"

deps-be:
	@echo "📥 安装 Go 依赖..."
	@go mod tidy
	@echo "✅ Go 依赖安装完成"

deps-fe:
	@echo "📥 安装前端依赖..."
	@cd frontend && npm install
	@echo "✅ 前端依赖安装完成"

# ============ 完整验证 ============
all: verify-be verify-fe
	@echo ""
	@echo "🎉 所有检查通过！"

# ============ 后端验证 ============
verify-be: build-be lint-be lint-quality-be test-be
	@echo "✅ 后端验证全部通过"

build-be:
	@echo "🔨 [Go] 构建..."
	@go build ./...
	@echo "✅ Go 构建成功"

lint-be:
	@echo "🔍 [Go] 架构层级检查..."
	@go run ./scripts/lint-deps
	@echo "✅ Go 架构检查通过"

lint-quality-be:
	@echo "🔍 [Go] 代码质量检查..."
	@go run ./scripts/lint-quality
	@echo "✅ Go 质量检查通过"

test-be:
	@echo "🧪 [Go] 运行单元测试..."
	@go test ./...
	@echo "✅ Go 单元测试通过"

# ============ 前端验证 ============
verify-fe: type-check-fe lint-fe lint-deps-fe build-fe
	@echo "✅ 前端验证全部通过"

type-check-fe:
	@echo "🔍 [Vue3] TypeScript 类型检查..."
	@cd frontend && npx vue-tsc --noEmit
	@echo "✅ 前端类型检查通过"

lint-fe:
	@echo "🔍 [Vue3] ESLint 检查..."
	@cd frontend && npm run lint
	@echo "✅ 前端 ESLint 通过"

lint-deps-fe:
	@echo "🔍 [Vue3] 前端层级架构检查..."
	@node scripts/lint-deps.mjs
	@echo "✅ 前端层级检查通过"

build-fe:
	@echo "🔨 [Vue3] 构建生产包..."
	@cd frontend && npm run build
	@echo "✅ 前端构建成功"

# ============ 快速自检 ============
quick: build-be lint-be test-be
	@echo "✅ 快速检查完成（建议提交前运行）"

# ============ 清理 ============
clean:
	@echo "🧹 清理构建产物..."
	@go clean
	@rm -rf frontend/dist
	@echo "✅ 清理完成"

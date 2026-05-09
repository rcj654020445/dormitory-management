#!/bin/bash
# harness/scripts/setup-env.sh
# Starts PostgreSQL and Redis for local development
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.. && pwd")"

echo "🔧 Setting up dormitory management environment..."

# Load .env files if they exist
if [ -f "$PROJECT_DIR/.env.be" ]; then
    export "$(grep -v '^#' "$PROJECT_DIR/.env.be" | xargs)"
fi

# ── PostgreSQL ───────────────────────────────────────────────────────────────
if docker ps --format '{{.Names}}' | grep -q "^dormitory-postgres$"; then
    echo "✅ PostgreSQL already running"
else
    echo "🚀 Starting PostgreSQL..."
    docker run -d \
        --name dormitory-postgres \
        -e POSTGRES_DB="${POSTGRES_DB:-dormitory}" \
        -e POSTGRES_USER="${POSTGRES_USER:-dormuser}" \
        -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-dormsecret}" \
        -p "${POSTGRES_PORT:-5432}:5432" \
        postgres:15-alpine

    echo "   Waiting for PostgreSQL to be ready..."
    local attempt=0
    until docker exec dormitory-postgres pg_isready -U "${POSTGRES_USER:-dormuser}" -d "${POSTGRES_DB:-dormitory}" > /dev/null 2>&1; do
        attempt=$((attempt + 1))
        if [ $attempt -gt 30 ]; then
            echo "❌ PostgreSQL failed to start within 30 seconds"
            exit 1
        fi
        sleep 1
    done
    echo "✅ PostgreSQL ready"
fi

# ── Redis ───────────────────────────────────────────────────────────────────
if docker ps --format '{{.Names}}' | grep -q "^dormitory-redis$"; then
    echo "✅ Redis already running"
else
    echo "🚀 Starting Redis..."
    docker run -d \
        --name dormitory-redis \
        -e REDIS_PASSWORD="${REDIS_PASSWORD:-redissecret}" \
        -p "${REDIS_PORT:-6379}:6379" \
        redis:7-alpine \
        redis-server --requirepass "${REDIS_PASSWORD:-redissecret}"

    echo "   Waiting for Redis to be ready..."
    local attempt=0
    until docker exec dormitory-redis redis-cli -a "${REDIS_PASSWORD:-redissecret}" ping 2>/dev/null | grep -q PONG; do
        attempt=$((attempt + 1))
        if [ $attempt -gt 15 ]; then
            echo "❌ Redis failed to start within 15 seconds"
            exit 1
        fi
        sleep 1
    done
    echo "✅ Redis ready"
fi

# ── Export derived env vars ───────────────────────────────────────────────────
export DATABASE_URL="${DATABASE_URL:-postgres://${POSTGRES_USER:-dormuser}:${POSTGRES_PASSWORD:-dormsecret}@localhost:${POSTGRES_PORT:-5432}/${POSTGRES_DB:-dormitory}?sslmode=disable}"
export REDIS_URL="${REDIS_URL:-redis://:${REDIS_PASSWORD:-redissecret}@localhost:${REDIS_PORT:-6379}/0}"

echo ""
echo "✅ All services started"
echo "   DATABASE_URL=$DATABASE_URL"
echo "   (environment vars exported for current shell)"

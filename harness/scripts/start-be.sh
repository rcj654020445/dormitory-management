#!/bin/bash
# harness/scripts/start-be.sh
# Runs database migrations then starts the Go backend server
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.. && pwd")"

cd "$PROJECT_DIR"

# Ensure DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    source .env.be 2>/dev/null || true
fi

if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL is not set. Run 'source harness/scripts/setup-env.sh' first."
    exit 1
fi

echo "🔄 Running migrations..."
go run cmd/migrate/main.go up

echo "🚀 Starting Go server..."
go run cmd/server/main.go

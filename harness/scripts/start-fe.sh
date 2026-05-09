#!/bin/bash
# harness/scripts/start-fe.sh
# Starts the Vue3 frontend dev server
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.. && pwd")"

cd "$PROJECT_DIR/frontend"

if ! command -v npm > /dev/null 2>&1; then
    echo "❌ npm not found. Please install Node.js."
    exit 1
fi

# Install deps if node_modules missing
if [ ! -d "node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    npm install
fi

echo "🚀 Starting Vue3 dev server (http://localhost:5173)..."
npm run dev

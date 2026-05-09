#!/bin/bash
# harness/scripts/teardown-env.sh
# Stops and removes Docker containers for local development
set -e

echo "🧹 Tearing down dormitory management environment..."

CONTAINERS="dormitory-postgres dormitory-redis"

for name in $CONTAINERS; do
    if docker ps --format '{{.Names}}' | grep -q "^${name}$"; then
        echo "   Stopping $name..."
        docker stop "$name" > /dev/null 2>&1
        docker rm "$name" > /dev/null 2>&1
        echo "   ✅ $name removed"
    else
        echo "   ⏭️  $name not running, skipping"
    fi
done

echo "✅ Environment torn down"

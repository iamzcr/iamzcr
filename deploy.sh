#!/bin/bash
set -e

echo "=== Pulling latest code ==="
git pull

echo "=== Building Go backends ==="
make build

echo "=== Restarting services ==="
supervisorctl restart admin
supervisorctl restart frontend

echo "=== Building Admin SPA ==="
cd web/admin && npm run build

echo "=== Building Blog SPA ==="
cd ../frontend && npm run build

echo "=== Deploy complete ==="

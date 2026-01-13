#!/bin/bash

# Script para ver eventos fallidos en la DLQ (PostgreSQL)
# Uso: ./scripts/dlq-events.sh [limit]

set -e

LIMIT="${1:-10}"

echo "=== Dead Letter Queue Events ==="
echo "Showing last $LIMIT failed events"
echo ""

docker-compose exec -T postgres psql -U userapi -d userapi -c "
SELECT 
    id,
    event_type,
    user_id,
    error,
    attempts,
    created_at
FROM failed_events 
ORDER BY created_at DESC 
LIMIT $LIMIT;
"

echo ""
echo "Total failed events:"
docker-compose exec -T postgres psql -U userapi -d userapi -t -c "SELECT COUNT(*) FROM failed_events;"

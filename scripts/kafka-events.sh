#!/bin/bash

# Script para ver eventos en Kafka
# Uso: ./scripts/kafka-events.sh [topic] [from-beginning]

set -e

TOPIC="${1:-user-events}"
FROM_BEGINNING="${2:-}"

echo "=== Kafka Events Viewer ==="
echo "Topic: $TOPIC"
echo "Press Ctrl+C to stop"
echo ""

if [ "$FROM_BEGINNING" = "from-beginning" ]; then
    echo "Reading from beginning..."
    docker-compose exec -T kafka kafka-console-consumer \
        --bootstrap-server localhost:9092 \
        --topic "$TOPIC" \
        --from-beginning \
        --property print.key=true \
        --property key.separator=" | "
else
    echo "Reading new messages only..."
    docker-compose exec -T kafka kafka-console-consumer \
        --bootstrap-server localhost:9092 \
        --topic "$TOPIC" \
        --property print.key=true \
        --property key.separator=" | "
fi

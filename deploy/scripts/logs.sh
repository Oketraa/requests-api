#!/bin/bash

SERVICE=$1

cd "$(dirname "$0")/.." || exit 1

if [ -n "$SERVICE" ]; then
  echo "Логи сервиса: $SERVICE"
  docker compose logs -f "$SERVICE"
else
  echo "Логи всех сервисов (Ctrl+C для выхода):"
  docker compose logs -f
fi

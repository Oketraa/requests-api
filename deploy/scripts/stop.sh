#!/bin/bash

cd "$(dirname "$0")/.." || exit 1

echo "Остановка проекта..."
docker compose down

echo "Контейнеры остановлены."

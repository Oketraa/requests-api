#!/bin/bash

cd "$(dirname "$0")/.." || exit 1

echo "Запуск проекта..."
docker compose up -d

echo "Статус контейнеров:"
docker compose ps

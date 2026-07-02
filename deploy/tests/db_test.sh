#!/bin/bash
set -e

echo "=== Запуск Теста Базы Данных ==="

# Нам нужно находиться в папке deploy, чтобы docker compose сработал
cd "$(dirname "$0")/.."

# Запускаем утилиту проверки готовности pg_isready внутри контейнера postgres
if docker compose exec -T postgres pg_isready > /dev/null 2>&1; then
    echo "✅ Тест БД: PostgreSQL запущена и готова принимать подключения"
else
    echo "❌ Тест БД провален! PostgreSQL не отвечает."
    exit 1
fi

echo "🎉 Тест Базы Данных успешно пройден!"
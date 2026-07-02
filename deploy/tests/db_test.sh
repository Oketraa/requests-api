#!/bin/bash
set -e

echo "=== Запуск Теста Базы Данных ==="

# Нам нужно находиться в папке deploy, чтобы docker compose сработал
cd "$(dirname "$0")/.."

# Загружаем переменные из .env, чтобы скрипт знал имя пользователя БД
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

echo "Текущий статус контейнеров:"
docker compose ps

echo "Ожидание готовности PostgreSQL..."
RETRIES=30

while [ $RETRIES -gt 0 ]; do
    # Пытаемся подключиться через pg_isready, используя пользователя из .env
    if docker compose exec -T postgres pg_isready -U "${DB_USER:-postgres}" > /dev/null 2>&1; then
        echo "✅ Тест БД: PostgreSQL успешно запущена и готова к работе!"
        exit 0
    else
        echo "PostgreSQL еще настраивается... Осталось попыток: $RETRIES"
        RETRIES=$((RETRIES - 1))
        sleep 1
    fi
done

echo "❌ Тест БД провален! PostgreSQL не ответила в течение 30 секунд."
echo "=== ЛОГИ КОНТЕЙНЕРА POSTGRES ДЛЯ ОТЛАДКИ ==="
docker compose logs postgres
exit 1
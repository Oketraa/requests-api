#!/bin/bash
set -e

echo "=== Запуск Теста Базы Данных ==="

echo "Ожидание готовности PostgreSQL..."
RETRIES=30

# Запускаем цикл повторных проверок (до 30 попыток)
while [ $RETRIES -gt 0 ]; do
    # Проверяем утилитой pg_isready внутри контейнера
    if docker compose exec -T postgres pg_isready > /dev/null 2>&1; then
        echo "✅ Тест БД: PostgreSQL успешно запущена и готова к работе!"
        exit 0
    else
        echo "PostgreSQL еще настраивается... Осталось попыток: $RETRIES"
        RETRIES=$((RETRIES - 1))
        sleep 1
    fi
done

echo "❌ Тест БД провален! PostgreSQL не ответила в течение 30 секунд."
exit 1

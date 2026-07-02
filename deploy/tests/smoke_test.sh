#!/bin/bash
# Завершать работу при любой ошибке
set -e

echo "=== Запуск Smoke-теста ==="

# 1. Проверяем редирект HTTP -> HTTPS (должен вернуть HTTP-код 301)
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/health)
if [ "$HTTP_STATUS" -eq 301 ]; then
    echo "✅ Тест 1: HTTP-редирект работает исправно (HTTP 301)"
else
    echo "❌ Тест 1: Ошибка редиректа! Получен код: $HTTP_STATUS (ожидался 301)"
    exit 1
fi

# 2. Проверяем доступность HTTPS /health (должен вернуть HTTP-код 200)
HTTPS_STATUS=$(curl -k -s -o /dev/null -w "%{http_code}" https://localhost/health)
if [ "$HTTPS_STATUS" -eq 200 ]; then
    echo "✅ Тест 2: Эндпоинт HTTPS /health доступен (HTTP 200)"
else
    echo "❌ Тест 2: Сервис недоступен! Получен код: $HTTPS_STATUS (ожидался 200)"
    exit 1
fi

echo "🎉 Smoke-тест успешно пройден!"
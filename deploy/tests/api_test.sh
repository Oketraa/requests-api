#!/bin/bash
set -e

echo "=== Запуск API-теста ==="

# Проверяем эндпоинт получения заявок через HTTPS
API_STATUS=$(curl -k -s -o /dev/null -w "%{http_code}" https://localhost/api/requests)

if [ "$API_STATUS" -eq 200 ]; then
    echo "✅ Тест API: Эндпоинт https://localhost/api/requests доступен (HTTP 200)"
    
    # Получаем ответ и выводим его в консоль
    RESPONSE=$(curl -k -s https://localhost/api/requests)
    echo "Ответ от API: $RESPONSE"
else
    echo "❌ Тест API провален! Получен код: $API_STATUS (ожидался 200)"
    exit 1
fi

echo "🎉 API-тест успешно пройден!"
#!/bin/sh
set -e

# Проверка входных данных
input_check() {
    if [ -z "$ENV_NAME" ]; then
        echo "Ошибка: Переменная ENV_NAME не задана!" >&2
        echo "Использование: ENV_NAME=имя_ветки $0" >&2
        exit 1
    fi
}

# Определение базовых переменных проекта
variable_definition() {
    # Очистка имени ветки для безопасности
    ENV_NAME=$(echo "$ENV_NAME" | tr -cd 'a-zA-Z0-9_-' | tr 'A-Z' 'a-z')
    export ENV_NAME

    SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
    COMPOSE_FILE="$SCRIPT_DIR/../docker-compose.test.yml"
    
    # Исправлено: файл конфигурации лежит в корне deploy/
    ENV_FILE="$SCRIPT_DIR/../.env.$ENV_NAME"
    
    # Путь к директории мониторинга
    MONITORING_DIR=$(cd "$SCRIPT_DIR/../../monitoring" && pwd 2>/dev/null || echo "$SCRIPT_DIR/../monitoring")
}

# Удаление контейнеров, сетей, волумов и файлов конфигурации
remove_environment() {
    echo "Остановка контейнеров и удаление изолированных волумов данных..."
    
    # Проверяем наличие .env файла перед остановкой
    if [ -f "$ENV_FILE" ]; then
        # Передаем --env-file для корректного сворачивания всех ресурсов
        docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" -p "$ENV_NAME" down -v
        
        echo "Удаление файла конфигурации: $(basename "$ENV_FILE")"
        rm -f "$ENV_FILE"
    else
        # Если файла уже нет, чистим по имени проекта
        docker compose -f "$COMPOSE_FILE" -p "$ENV_NAME" down -v
    fi

    # Удаляем сгенерированный тестовый конфиг Prometheus
    local prom_config="$MONITORING_DIR/prometheus/prometheus.test.${ENV_NAME}.yml"
    if [ -f "$prom_config" ]; then
        echo "Удаление конфигурации Prometheus: $(basename "$prom_config")"
        rm -f "$prom_config"
    fi
}

# Главная функция управления логикой
main() {
    input_check
    variable_definition
    
    echo "Уничтожение окружения для ветки: $ENV_NAME"
    remove_environment
    
    echo "Окружение $ENV_NAME успешно уничтожено, ресурсы освобождены!"
}

# Точка входа в скрипт
main
#!/bin/sh

input_check()
{
    if [ -z "$ENV_NAME" ]; then
        echo "Ошибка: Переменная ENV_NAME не задана!"
        echo "Использование: ENV_NAME=имя_ветки $0"
        exit 1
    fi
}

variable_definition()
{
    ENV_NAME=$(echo "$ENV_NAME" | tr -cd 'a-zA-Z0-9_-' | tr 'A-Z' 'a-z')
    export ENV_NAME

    SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
    COMPOSE_FILE="$SCRIPT_DIR/../docker-compose.test.yml"
    ENV_FILE="$SCRIPT_DIR/.env.$ENV_NAME"
    MONITORING_DIR=$(cd "$SCRIPT_DIR/../../monitoring" && pwd 2>/dev/null || echo "$SCRIPT_DIR/../monitoring")
}

remove_enviroment()
{
    echo "Остановка контейнеров и удаление изолированных волумов данных..."
    # Флаг -v удаляет связанные с проектом волумы (очищает базу и метрики)
    docker compose -f "$COMPOSE_FILE" -p "$ENV_NAME" down -v

    if [ -f "$ENV_FILE" ]; then
        echo "Удаление файла конфигурации: $ENV_FILE"
        rm -f "$ENV_FILE"
    fi
    rm -f "$MONITORING_DIR/prometheus/prometheus.test.${ENV_NAME}.yml"
}

destroy_compose()
{
    input_check
    variable_definition
    
    echo "Уничтожение окружения для ветки: $ENV_NAME"
    remove_enviroment
    
    echo "Окружение $ENV_NAME успешно уничтожено, ресурсы освобождены!"
}

set -e
destroy_compose
#!/bin/sh
set -e

input_check() {
    if [ -z "$ENV_NAME" ]; then
        echo "Ошибка: Переменная ENV_NAME не задана!" >&2
        echo "Использование: ENV_NAME=имя_ветки $0" >&2
        exit 1
    fi
}

variable_definition() {
    ENV_NAME=$(echo "$ENV_NAME" | tr -cd 'a-zA-Z0-9_-' | tr 'A-Z' 'a-z')
    export ENV_NAME

    SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
    COMPOSE_FILE="$SCRIPT_DIR/../docker-compose.test.yml"
    ENV_FILE="$SCRIPT_DIR/../.env.$ENV_NAME"
    MONITORING_DIR=$(cd "$SCRIPT_DIR/../../monitoring" && pwd 2>/dev/null || echo "$SCRIPT_DIR/../monitoring")
}

find_free_port() {
    local port=$1
    while true; do
        if ! nc -z 127.0.0.1 "$port" >/dev/null 2>&1; then
            echo "$port"
            break
        fi
        port=$((port + 1))
    done
}

create_environment() {
    if [ -f "$ENV_FILE" ]; then
        echo "Найден готовый файл конфигурации: $ENV_FILE"
        while IFS= read -r line || [ -n "$line" ]; do
            case "$line" in
                ""|\#*) continue ;;
                *) export "$line" ;;
            esac
        done < "$ENV_FILE"
    else
        echo "Конфигурация не найдена. Запускаем автоматическую генерацию портов и секретов..."
        
        # Динамический подбор свободных портов
        export APP_PORT=8080 # Внутренний порт приложения
        export NGINX_PORT=$(find_free_port 8080)
        export PROMETHEUS_PORT=$(find_free_port 9090)
        export GRAFANA_PORT=$(find_free_port 3000)

        # Генерация учетных данных бд
        export DB_USER="user_${ENV_NAME}"
        export DB_NAME="db_${ENV_NAME}"
        export DB_PASSWORD=$(head -c 32 /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 16)
        
        # Генерация данных мониторинга
        export GRAFANA_USER="admin"
        export GRAFANA_PASSWORD=$(head -c 32 /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 12)

        # Сохраняем в файл строго по структуре нового шаблона (без хоста и порта базы)
        cat << EOF > "$ENV_FILE"
APP_PORT=$APP_PORT
NGINX_PORT=$NGINX_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
GRAFANA_USER=$GRAFANA_USER
GRAFANA_PASSWORD=$GRAFANA_PASSWORD
PROMETHEUS_PORT=$PROMETHEUS_PORT
GRAFANA_PORT=$GRAFANA_PORT
EOF
        echo "Конфигурация успешно сохранена в $ENV_FILE"
    fi
}

create_prometheus_environment() {
    local prom_template="$MONITORING_DIR/prometheus/prometheus.yml"
    local prom_target="$MONITORING_DIR/prometheus/prometheus.test.$ENV_NAME.yml"

    if [ -f "$prom_template" ]; then
        cat "$prom_template" > "$prom_target"
        cat << EOF >> "$prom_target"
    metric_relabel_configs:
      - source_labels: ['container_label_com_docker_compose_project']
        regex: "^${ENV_NAME}$"
        action: keep
EOF
    else
        echo "Предупреждение: Шаблон Prometheus не найден по пути $prom_template"
    fi
}

main() {
    input_check
    variable_definition
    
    echo "Запуск окружения для ветки: $ENV_NAME"
    create_environment
    create_prometheus_environment

    echo "Параметры окружения:"
    echo "  - Nginx (Приложение): http://localhost:$NGINX_PORT"
    echo "  - Prometheus:         http://localhost:$PROMETHEUS_PORT"
    echo "  - Grafana:            http://localhost:$GRAFANA_PORT (Логин: $GRAFANA_USER / Пароль: $GRAFANA_PASSWORD)"

    echo "Запускаем контейнеры через Docker Compose..."
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" -p "$ENV_NAME" up -d --build

    echo "Окружение $ENV_NAME успешно развернуто!"
}

main
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

find_free_port() {
    port=$1
    while true; do
        if ! nc -z 127.0.0.1 "$port" >/dev/null 2>&1; then
            echo "$port"
            break
        fi
        port=$((port + 1))
    done
}

create_enviroment()
{
    if [ -f "$ENV_FILE" ]; then
        echo "Найден готовый файл конфигурации: $ENV_FILE"
        while read -r line || [ -n "$line" ]; do
            case "$line" in
                ""|\#*) continue ;;
                *) export "$line" ;;
            esac
        done < "$ENV_FILE"
    else
        echo "Конфигурация не найдена. Запускаем автонастройку..."
        
        export NGINX_PORT=$(find_free_port 8080)
        export PROMETHEUS_PORT=$(find_free_port 9090)
        export GRAFANA_PORT=$(find_free_port 3000)

        export DB_USER="user_${ENV_NAME}"
        export DB_NAME="db_${ENV_NAME}"
        export DB_HOST="postgres"
        export DB_PORT="5432"
        
        export DB_PASSWORD=$(head -c 32 /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 16)
        export GRAFANA_USER="admin"
        export GRAFANA_PASSWORD=$(head -c 32 /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 12)

        # Сохраняем переменные
        cat << EOF > "$ENV_FILE"
ENV_NAME=$ENV_NAME
NGINX_PORT=$NGINX_PORT
PROMETHEUS_PORT=$PROMETHEUS_PORT
GRAFANA_PORT=$GRAFANA_PORT
DB_USER=$DB_USER
DB_NAME=$DB_NAME
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_PASSWORD=$DB_PASSWORD
GRAFANA_USER=$GRAFANA_USER
GRAFANA_PASSWORD=$GRAFANA_PASSWORD
EOF
        echo "Конфигурация сохранена в $ENV_FILE"
    fi
}
create_prometheus_enviroment()
{
    cat "$MONITORING_DIR/prometheus/prometheus.yml" > "$MONITORING_DIR/prometheus/prometheus.test.$ENV_NAME.yml"

# Дописываем в конец тестовую фильтрацию (главное соблюсти отступы YAML)
    cat <<EOF >> "$MONITORING_DIR/prometheus/prometheus.test.$ENV_NAME.yml"
    metric_relabel_configs:
      - source_labels: ['container_label_com_docker_compose_project']
        regex: "^${ENV_NAME}$"
        action: keep
EOF
}
start_compose()
{
    input_check
    variable_definition
    echo "Запуск окружения для ветки: $ENV_NAME"
    create_enviroment
    create_prometheus_enviroment

    echo "Параметры окружения:"
    echo "  - Nginx (Приложение): http://localhost:$NGINX_PORT"
    echo "  - Prometheus:         http://localhost:$PROMETHEUS_PORT"
    echo "  - Grafana:            http://localhost:$GRAFANA_PORT (Логин: $GRAFANA_USER / Пароль: $GRAFANA_PASSWORD)"

    echo "Запускаем контейнеры..."
    docker compose -f "$COMPOSE_FILE" -p "$ENV_NAME" up -d --build

    echo "Окружение $ENV_NAME успешно развернуто!"
}

set -e
start_compose
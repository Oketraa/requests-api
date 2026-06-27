#!/bin/sh
set -e

# Определение базовых переменных проекта
variable_definition() {
    SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
    DEPLOY_DIR=$(cd "$SCRIPT_DIR/.." && pwd)

    # По умолчанию проверяем .env, но можно передать имя конкретного файла аргументом
    TARGET_NAME="${1:-.env}"
    TARGET="$DEPLOY_DIR/$TARGET_NAME"
}

# Валидация содержимого файла конфигурации
validate_config() {
    if [ ! -f "$TARGET" ]; then
        echo "Ошибка: Файл $TARGET не найден. Сначала запустите generate_env.sh" >&2
        exit 1
    fi

    # Список строго обязательных секретов и параметров, которые не должны быть пустыми
    REQUIRED_VARS="NGINX_PORT DB_USER DB_PASSWORD DB_NAME GRAFANA_PASSWORD"
    local missing_count=0

    echo "Проверка обязательных переменных в $TARGET_NAME..."

    for var in $REQUIRED_VARS; do
        # Извлекаем значение переменной из файла
        local val=$(grep "^${var}=" "$TARGET" | cut -d'=' -f2- || true)
        
        if [ -z "$val" ]; then
            echo "Ошибка: Переменная $var отсутствует или пуста в файле $TARGET_NAME" >&2
            missing_count=$((missing_count + 1))
        fi
    done

    if [ "$missing_count" -ne 0 ]; then
        echo "Валидация завершилась ошибкой. Пожалуйста, заполните пустые секреты." >&2
        exit 1
    fi
}

# Главная функция управления логикой
main() {
    # Передаем все аргументы скрипта в функцию определения переменных
    variable_definition "$@"
    validate_config
    
    echo "Все обязательные переменные в $TARGET_NAME успешно заданы!"
}

# Точка входа
main "$@"
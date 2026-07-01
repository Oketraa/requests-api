#!/bin/sh
set -e

# Установка базовых путей
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
DEPLOY_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
TEMPLATE="$DEPLOY_DIR/.env.template"

# Проверка наличия шаблона
check_template() {
    if [ ! -f "$TEMPLATE" ]; then
        echo "Ошибка: Файл шаблона $TEMPLATE не найден!" >&2
        exit 1
    fi
}

# Определение имени итогового файла (с поддержкой переменной ENV_NAME)
set_target_file() {
    # Если ENV_NAME передана извне, используем её без лишних вопросов
    if [ -n "$ENV_NAME" ]; then
        if [ "$ENV_NAME" = "production" ]; then
            TARGET_FILE="$DEPLOY_DIR/.env"
        else
            # Очищаем имя окружения от греха подальше
            local clean_name=$(echo "$ENV_NAME" | tr -cd 'a-zA-Z0-9_-' | tr 'A-Z' 'a-z')
            TARGET_FILE="$DEPLOY_DIR/.env.${clean_name}"
        fi
        return 0
    fi

    # Если ENV_NAME пустая — включаем интерактивный режим
    echo "Выберите тип окружения для генерации конфигурации:"
    echo "1) Основное (будет создан .env)"
    echo "2) Тестовое (будет создан .env.имя_окружения)"
    printf "Введите номер варианта (1 или 2): "
    read -r ENV_TYPE < /dev/tty

    if [ "$ENV_TYPE" = "1" ]; then
        TARGET_FILE="$DEPLOY_DIR/.env"
    elif [ "$ENV_TYPE" = "2" ]; then
        printf "Введите имя тестового окружения (например, yota): "
        read -r ENV_NAME < /dev/tty
        if [ -z "$ENV_NAME" ]; then
            echo "Ошибка: Имя окружения не может быть пустым." >&2
            exit 1
        fi
        local clean_name=$(echo "$ENV_NAME" | tr -cd 'a-zA-Z0-9_-' | tr 'A-Z' 'a-z')
        TARGET_FILE="$DEPLOY_DIR/.env.${clean_name}"
    else
        echo "Ошибка: Неверный выбор." >&2
        exit 1
    fi
}

# Получение значения переменной с «умным» пропуском read
process_value() {
    local key="$1"
    local default="$2"
    local env_val=""

    # 1. Проверяем, пришло ли значение из внешнего окружения Shell
    eval "env_val=\$${key}"

    # Если значение передано извне, МЫ ПРОПУСКАЕМ read и сразу берем его
    if [ -n "$env_val" ]; then
        RET_VAL="$env_val"
        return 0
    fi

    # 2. Если в текущем Shell пусто, проверяем, вдруг файл .env уже существовал
    local current_val=""
    if [ -f "$TARGET_FILE" ]; then
        current_val=$(grep "^${key}=" "$TARGET_FILE" | cut -d'=' -f2- || true)
    fi

    # 3. Если в файле пусто, берем дефолт из шаблона
    if [ -z "$current_val" ]; then
        current_val="$default"
    fi

    # 4. Запрашиваем ввод только для того, чего не было во внешнем окружении
    if [ -n "$current_val" ]; then
        printf "Введите %s [%s]: " "$key" "$current_val"
        read -r user_input < /dev/tty
        if [ -z "$user_input" ]; then
            user_input="$current_val"
        fi
    else
        printf "Введите %s (обязательно): " "$key"
        read -r user_input < /dev/tty
    fi

    RET_VAL="$user_input"
}

# Генерация итогового файла на основе шаблона
generate_config() {
    local temp_file="${TARGET_FILE}.tmp"
    
    echo "--- Заполнение переменных конфигурации ---"
    touch "$TARGET_FILE"
    rm -f "$temp_file"

    while IFS= read -r line <&3 || [ -n "$line" ]; do
        if echo "$line" | grep -q '^[[:space:]]*#' || [ -z "$line" ]; then
            echo "$line" >> "$temp_file"
            continue
        fi

        local key=$(echo "$line" | cut -d'=' -f1)
        local template_default=$(echo "$line" | cut -d'=' -f2-)

        # Обрабатываем переменную
        process_value "$key" "$template_default"

        echo "${key}=${RET_VAL}" >> "$temp_file"
    done 3< "$TEMPLATE"

    mv "$temp_file" "$TARGET_FILE"
}

# Главная функция
main() {
    check_template
    
    TARGET_FILE=""
    RET_VAL=""
    
    set_target_file
    generate_config
    
    local target_name=$(basename "$TARGET_FILE")
    echo "----------------------------------------"
    echo "Успешно: Файл $target_name сгенерирован в папке deploy/."
}

main
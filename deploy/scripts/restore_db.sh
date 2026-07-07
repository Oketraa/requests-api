#!/bin/bash

# Сохраняем путь до cd
if [[ "$1" = /* ]]; then
  BACKUP_FILE="$1"
else
  BACKUP_FILE="$(pwd)/$1"
fi

cd "$(dirname "$0")/.." || exit 1

export $(grep -v '^#' .env | xargs)

# Проверяем что передан аргумент
if [ -z "$BACKUP_FILE" ]; then
  echo "Ошибка: укажи файл бэкапа!"
  echo "Пример: ./restore_db.sh ./backups/backup_2026-01-01_12-00-00.sql"
  exit 1
fi

# Проверяем что файл существует
if [ ! -f "$BACKUP_FILE" ]; then
  echo "Ошибка: файл $BACKUP_FILE не найден!"
  exit 1
fi

# Проверяем что файл не пустой
if [ ! -s "$BACKUP_FILE" ]; then
  echo "Ошибка: файл бэкапа пустой!"
  exit 1
fi

echo "Восстановление базы данных из $BACKUP_FILE..."

cat "$BACKUP_FILE" | docker compose exec -T postgres psql \
  -U "$DB_USER" \
  "$DB_NAME"

if [ $? -eq 0 ]; then
  echo "База данных успешно восстановлена!"
else
  echo "Ошибка при восстановлении!"
  exit 1
fi

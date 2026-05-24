#!/bin/bash

cd "$(dirname "$0")/.." || exit 1

export $(grep -v '^#' .env | xargs)

BACKUP_DIR="./backups"
mkdir -p "$BACKUP_DIR"

TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_FILE="$BACKUP_DIR/backup_${TIMESTAMP}.sql"

echo "Создание резервной копии базы данных..."

docker compose exec -T postgres pg_dump \
  -U "$DB_USER" \
  "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
  echo "Бэкап успешно сохранён: $BACKUP_FILE"
else
  echo "Ошибка при создании бэкапа!"
  exit 1
fi

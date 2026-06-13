#!/bin/bash

cd "$(dirname "$0")/.." || exit 1

export $(grep -v '^#' .env | xargs)

BACKUP_DIR="./backups"
mkdir -p "$BACKUP_DIR"

TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_FILE="$BACKUP_DIR/backup_${TIMESTAMP}.sql"

# Количество бэкапов которые хранить (остальные удаляются)
MAX_BACKUPS=5

echo "Создание резервной копии базы данных..."

docker compose exec -T postgres pg_dump \
  -U "$DB_USER" \
  "$DB_NAME" > "$BACKUP_FILE"

if [ $? -ne 0 ]; then
  echo "Ошибка при создании бэкапа!"
  rm -f "$BACKUP_FILE"
  exit 1
fi

# Проверяем что файл не пустой
if [ ! -s "$BACKUP_FILE" ]; then
  echo "Ошибка: бэкап пустой!"
  rm -f "$BACKUP_FILE"
  exit 1
fi

echo "Бэкап успешно сохранён: $BACKUP_FILE"

# Ротация — удаляем старые бэкапы
BACKUP_COUNT=$(ls -1 "$BACKUP_DIR"/*.sql 2>/dev/null | wc -l)
if [ "$BACKUP_COUNT" -gt "$MAX_BACKUPS" ]; then
  echo "Ротация бэкапов — удаляем старые..."
  ls -1t "$BACKUP_DIR"/*.sql | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -f
  echo "Оставлено последних $MAX_BACKUPS бэкапов"
fi

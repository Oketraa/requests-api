# Requests API

REST API сервис для учета заявок на Go.

---

# Технологии

- Go
- PostgreSQL
- Docker
- Nginx (Reverse Proxy)
- REST API
- Git

---

# Возможности

- Создание заявки
- Получение всех заявок
- Получение заявки по ID
- Обновление заявки
- Удаление заявки
- Health check endpoint

---

# Запуск проекта 

## 1. Установить Docker Desktop

Скачать:

https://www.docker.com/products/docker-desktop/

После установки запустить Docker Desktop.

---

## 2. Клонировать проект

```bash
git clone <repository_url>
```

Перейти в папку проекта:

```bash
cd requests-api
```

---

## 3. Создать .env файл

Создать файл `.env` в корне проекта.

Пример содержимого:

```env
APP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=requests_db
```

---

## 4. Запустить проект

```bash
docker compose up --build
```

После запуска сервер будет доступен:

```text
http://localhost
```

---

# Проверка работы API

---

# Health check

## GET /health

```bash
curl http://localhost/health
```

---

# Получение всех заявок

## GET /api/requests

```bash
curl http://localhost/api/requests
```

---

# Получение заявки по ID

## GET /api/requests/{id}

Пример:

```bash
curl http://localhost/api/requests/1
```

---

# Создание заявки

## POST /api/requests

```bash
curl -X POST http://localhost/api/requests/create ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Проверить сервер\",\"description\":\"Проверить nginx\"}"
```

---

# Обновление заявки

## PUT /api/requests/{id}

```bash
curl -X PUT http://localhost/api/requests/1 ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Проверить сервер\",\"description\":\"Проверить nginx\",\"status\":\"in_progress\"}"
```

---

# Удаление заявки

## DELETE /api/requests/{id}

```bash
curl -X DELETE http://localhost/api/requests/1
```

---
# Скрипты эксплуатации (DevOps 4)

## Запуск проекта

```bash
./deploy/scripts/start.sh
```

## Остановка проекта

```bash
./deploy/scripts/stop.sh
```

## Просмотр логов

```bash
# Все сервисы
./deploy/scripts/logs.sh

# Конкретный сервис
./deploy/scripts/logs.sh backend
```

## Резервное копирование базы данных

```bash
./deploy/scripts/backup_db.sh
```

Бэкап сохраняется в папку `deploy/backups/` с именем вида `backup_2026-01-01_12-00-00.sql`

## Восстановление базы из бэкапа

### Стратегия резервного копирования

- Бэкапы создаются вручную командой `./deploy/scripts/backup_db.sh`
- Каждый бэкап сохраняется в отдельный файл с датой и временем
- Файлы хранятся в папке `deploy/backups/`
- Формат имени файла: `backup_ГГГГ-ММ-ДД_ЧЧ-ММ-СС.sql`
- Рекомендуется делать бэкап перед каждым обновлением проекта

### Порядок восстановления

1. Убедись что проект запущен:

```bash
./deploy/scripts/start.sh
```

2. Посмотри список доступных бэкапов:

```bash
ls -la deploy/backups/
```

3. Восстанови нужный бэкап:

```bash
cat deploy/backups/backup_2026-01-01_12-00-00.sql | \
  docker compose exec -T postgres psql -U app_user requests_db
```

4. Проверь что данные восстановились:

```bash
curl http://localhost/api/requests
```
### Восстановление через скрипт

```bash
./deploy/scripts/restore_db.sh ./deploy/backups/backup_2026-01-01_12-00-00.sql
```

Скрипт автоматически проверяет что файл существует и не пустой.

### Восстановление на чистую базу

Если база данных была полностью удалена:

1. Запусти проект — база создастся автоматически:
```bash
./deploy/scripts/start.sh
```

2. Восстанови данные из бэкапа:
```bash
./deploy/scripts/restore_db.sh ./deploy/backups/backup_2026-01-01_12-00-00.sql
```

3. Проверь что данные восстановились:
```bash
curl http://localhost/api/requests
```
### Ротация бэкапов

Автоматически хранятся последние 5 бэкапов — старые удаляются при создании нового.

---

## Безопасность

Nginx настроен со следующими мерами защиты:

- Security headers — защита от XSS и clickjacking
- Rate limiting — не более 10 запросов в секунду с одного IP
- Ограничение размера запроса — максимум 1MB

Проверить headers:
```bash
curl -I http://localhost/health
```

---

## Healthcheck

Контейнеры имеют автоматические проверки состояния:

- `postgres` — проверяется через `pg_isready`
- `app` — проверяется через `GET /health`
- `nginx` — запускается только после того как `app` стал healthy

Проверить статус:
```bash
docker compose -f deploy/docker-compose.yml ps
```
---
# Структура проекта

```text
requests_api/
├── backend
│   ├── database
│   │   └── db.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── request.go
│   ├── main.go
│   ├── models
│   │   └── request.go
│   └── utils
│       └── response.go
├── deploy
│   ├── docker-compose.yml
│   ├── init.sql
│   └── nginx.conf
│   └── .env.example
└── README.md
```

---

# Автор

Backend project на Go.
Oketraa

# Requests API

REST API сервис для учета заявок на Go.

---

# Технологии

- Go
- PostgreSQL
- Postgres Exporter
- Docker
- Nginx (Reverse Proxy)
- REST API
- Git
- Grafana
- cAdvisor
- Prometheus
- Grafana Loki
- Grafana Alloy
- Prometheus Node Exporter
- GNU Make

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

```bash
make env
```

---

## 4. Запустить проект

```bash
make up
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
---

### Тестовые окружения

Для проверки новых функций в изолированных контурах (например, для ветки `yota`) на одной машине с рабочем окружением используются тестовые окружения. Изоляция обеспечивается на уровне метаданных Docker Compose. Управление тестовым контуром полностью автоматизировано локальными скриптами.

Если вы хотите определить свои переменные окружения. Скопируйте файл .env.example, переименуйте копию в формат .env.${ENV_NAME} (например, .env.yota). В противном случае, данный файл будет создан автоматически.

При запуске скрипт сообщит логин и пароль от Grafana, в случае их генерации, а также порты на которые будут проброшены контейнеры.

1. Запуск тестового окружения:
```bash
cd deploy/scripts
ENV_NAME=yota ./create_test_env.sh
```

2. Удаление тестового окружения
```bash
cd deploy/scripts
ENV_NAME=yota ./destroy_test_env.sh
```

---
# Структура проекта

```text
/requests-api
├── backend
│   ├── database
│   │   └── db.go
│   ├── Dockerfile
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   ├── health.go
│   │   └── request.go
│   ├── main.go
│   ├── metrics
│   │   ├── logic.go
│   │   └── model.go
│   ├── middleware
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── recovery.go
│   │   ├── request_id.go
│   │   └── user_agent.go
│   ├── models
│   │   └── request.go
│   └── utils
│       ├── response.go
│       └── validation.go
├── deploy
|   ├── .env.template
│   ├── backups
│   ├── docker-compose.test.yml
│   ├── docker-compose.yml
│   ├── init.sql
│   ├── init.test.sql
│   ├── nginx.conf
│   └── scripts
│       ├── backup_db.sh
│       ├── check_env.sh
│       ├── create_test_env.sh
│       ├── destroy_test_env.sh
│       ├── generate_env.sh
│       ├── logs.sh
│       ├── start.sh
│       └── stop.sh
├── documentation
│   ├── grafana.md
│   └── logs.md
├── Makefile
├── monitoring
│   ├── alloy
│   │   └── config.alloy
│   ├── grafana
│   │   ├── dashboards
│   │   │   └── Мониторинг системы технической поддержки.json
│   │   └── provisioning
│   │       ├── dashboards
│   │       │   └── dashboards.yaml
│   │       └── datasources
│   │           └── datasources.yaml
│   ├── loki
│   │   └── local-config.yaml
│   └── prometheus
│       ├── prometheus.test.production.yml
│       └── prometheus.yml
└── README.md
```

---

# Автор

Backend project на Go.
Oketraa

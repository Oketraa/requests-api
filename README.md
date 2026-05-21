# Requests API

REST API сервис для учета заявок на Go.

---

# Технологии

- Go
- PostgreSQL
- Docker
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
http://localhost:8080
```

---

# Проверка работы API

---

# Health check

## GET /health

```bash
curl http://localhost:8080/health
```

---

# Получение всех заявок

## GET /api/requests

```bash
curl http://localhost:8080/api/requests
```

---

# Получение заявки по ID

## GET /api/requests/{id}

Пример:

```bash
curl http://localhost:8080/api/requests/1
```

---

# Создание заявки

## POST /api/requests

```bash
curl -X POST http://localhost:8080/api/requests ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Проверить сервер\",\"description\":\"Проверить nginx\"}"
```

---

# Обновление заявки

## PUT /api/requests/{id}

```bash
curl -X PUT http://localhost:8080/api/requests/1 ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Проверить сервер\",\"description\":\"Проверить nginx\",\"status\":\"in_progress\"}"
```

---

# Удаление заявки

## DELETE /api/requests/{id}

```bash
curl -X DELETE http://localhost:8080/api/requests/1
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
│   └── .env.example
└── README.md
```

---

# Автор

Backend project на Go.
Oketraa

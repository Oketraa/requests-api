# Переменные по умолчанию
ENV_NAME ?= production
SHELL := /bin/sh

.PHONY: help env check up down restart logs ps test test-down backup clean

# Команда по умолчанию
help:
	@echo "Доступные команды для управления проектом:"
	@echo " - make env             - Генерация файла конфигурации (.env или .env.имя_ветки)"
	@echo " - make check           - Проверка заполнения обязательных переменных"
	@echo " - make up              - Запуск основного окружения"
	@echo " - make down            - Остановка основного окружения"
	@echo " - make restart         - Перезапуск основного окружения"
	@echo " - make logs            - Просмотр логов"
	@echo " - make ps              - Статус запущенных контейнеров основного окружения"
	@echo " - make test            - Развернуть тестовое окружение (Пример: make test ENV_NAME=yota)"
	@echo " - make test-down       - Уничтожить тестовое окружение (Пример: make test-down ENV_NAME=yota)"
	@echo " - make backup          - Создать резервную копию базы данных"
	@echo " - make clean           - Полная очистка временных файлов конфигурации"

env:
	@chmod +x deploy/scripts/generate_env.sh
	@ENV_NAME=$(ENV_NAME) ./deploy/scripts/generate_env.sh

check:
	@chmod +x deploy/scripts/check_env.sh
	@./deploy/scripts/check_env.sh

up: check
	@echo "Запуск основного окружения..."
	@chmod +x deploy/scripts/start.sh
	@./deploy/scripts/start.sh

down:
	@echo "Остановка основного окружения..."
	@chmod +x deploy/scripts/stop.sh
	@./deploy/scripts/stop.sh

restart: down up

logs:
	@chmod +x deploy/scripts/logs.sh
	@./deploy/scripts/logs.sh

ps:
	docker compose -f deploy/docker-compose.yml ps

test:
	@chmod +x deploy/scripts/create_test_env.sh
	@ENV_NAME=$(ENV_NAME) ./deploy/scripts/create_test_env.sh

test-down:
	@chmod +x deploy/scripts/destroy_test_env.sh
	@ENV_NAME=$(ENV_NAME) ./deploy/scripts/destroy_test_env.sh

backup:
	@echo "Запуск резервного копирования базы данных..."
	@chmod +x deploy/scripts/backup_db.sh
	@./deploy/scripts/backup_db.sh

clean:
	@echo "Удаление временных файлов генерации..."
	find deploy/ -maxdepth 1 -name "*.tmp" -exec rm -f {} +
	@echo "Очистка завершена."
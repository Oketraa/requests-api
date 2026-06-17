package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"study/database"
	"study/handlers"
	"study/metrics"
	"study/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const JSON_Log = true

func main() {
	if JSON_Log {
		// Включаем глобальный JSON-режим для slog (Пункт 13.6)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	} else {
		// Настройки дефолтного текстового логгера для 13.4 (если нужны кастомные префиксы)
		log.SetFlags(log.LstdFlags)
	}
	// 1. Инициализируем базу данных
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// 2. Создаем один изолированный роутер (mux) для всего приложения
	mux := http.NewServeMux()

	// 3. Регистрируем служебные ручки бэкендера (Healthchecks)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/live", handlers.Live)
	mux.HandleFunc("/ready", handlers.Ready)

	// 4. Регистрируем ручки бизнес-логики (API заявок)
	mux.HandleFunc("/api/requests", handlers.GetRequests)
	mux.HandleFunc("/api/requests/", handlers.RequestByID)
	mux.HandleFunc("/api/requests/create", handlers.CreateRequest)

	// 5. Регистрируем твою ручку сбора метрик Prometheus
	mux.Handle("/metrics", promhttp.Handler())

	// 6. Накатываем инфраструктурные мидлвари бэкендера на роутер
	// Цепочка идет снизу вверх: Cors -> Recovery -> Logger -> RequestID -> Роутер
	handler := middleware.Logger(JSON_Log)(mux)
	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(handler)
	handler = middleware.Cors(handler)

	// 7. Оборачиваем получившийся пирог в твою метрическую мидлварь,
	// чтобы Prometheus считал длительность и статусы ответов с учетом всех правил
	finalHandler := metrics.MetricsMiddleware(handler)

	// 8. Запускаем ОДИН единственный сервер на порту 8080
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatal(err)
	}
}

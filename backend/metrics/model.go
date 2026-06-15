package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество входящих HTTP-запросов",
		},
		[]string{"path", "status"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Количество запросов, обрабатываемых прямо сейчас",
		},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время обработки HTTP-запросов в секундах",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0}, // Интервалы для гистограммы
		},
		[]string{"path"},
	)

	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Общее количество ошибок 4xx и 5xx",
		},
		[]string{"path", "error_type"}, // error_type: "4xx" или "5xx"
	)

	requestsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_created_total",
			Help: "Общее количество созданных заявок",
		},
	)

	requestsUpdatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_updated_total",
			Help: "Общее количество обновленных заявок",
		},
	)

	requestsDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_deleted_total",
			Help: "Общее количество удаленных заявок",
		},
	)

	DatabaseQueryDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Время выполнения запросов к базе данных в секундах",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5}, // Сетки для быстрых запросов БД
		},
	)
)

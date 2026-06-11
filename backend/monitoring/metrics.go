package monitoring

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 1. Определяем метрики
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

// 2. Своя структура-обёртка для ResponseWriter, чтобы перехватить статус-код ответа
type responseWriterDelegator struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterDelegator) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// 3. Middleware для автоматического сбора метрик
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем сам эндпоинт метрик, чтобы не спамить в статистику
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Запрос вошел: увеличиваем инфлайт
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec() // Уменьшится автоматически при выходе из функции

		// Оборачиваем оригинальный ResponseWriter
		dwd := &responseWriterDelegator{ResponseWriter: w, statusCode: http.StatusOK}

		// Передаем управление дальше по цепочке
		next.ServeHTTP(dwd, r)

		// Запрос выполнился: считаем время
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.URL.Path).Observe(duration)

		// Метрика общего количества запросов (группируем по пути и статусу)
		statusStr := strconv.Itoa(dwd.statusCode)
		httpRequestsTotal.WithLabelValues(r.URL.Path, statusStr).Inc()

		// Отдельно считаем ошибки 4xx и 5xx
		if dwd.statusCode >= 400 && dwd.statusCode < 500 {
			httpErrorsTotal.WithLabelValues(r.URL.Path, "4xx").Inc()
		} else if dwd.statusCode >= 500 {
			httpErrorsTotal.WithLabelValues(r.URL.Path, "5xx").Inc()
		} else if r.URL.Path == "/api/requests/create" && dwd.statusCode == http.StatusCreated {
			requestsCreatedTotal.Inc()
		} else if strings.HasPrefix(r.URL.Path, "/api/requests/") && dwd.statusCode == http.StatusOK && r.Method == http.MethodPut {
			requestsUpdatedTotal.Inc()
		} else if strings.HasPrefix(r.URL.Path, "/api/requests/") && dwd.statusCode == http.StatusNoContent && r.Method == http.MethodDelete {
			requestsDeletedTotal.Inc()
		}
	})
}

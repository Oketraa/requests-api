package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Своя структура-обёртка для ResponseWriter, чтобы перехватить статус-код ответа
type responseWriterDelegator struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterDelegator) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware для автоматического сбора метрик
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		dwd := &responseWriterDelegator{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(dwd, r)

		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.URL.Path).Observe(duration)

		statusStr := strconv.Itoa(dwd.statusCode)
		httpRequestsTotal.WithLabelValues(r.URL.Path, statusStr).Inc()

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

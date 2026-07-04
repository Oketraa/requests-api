package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Передаем режим работы прямо в конструктор мидлвари при инициализации в main.go
func Logger(isJSON_Mode bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			reqID, _ := r.Context().Value(KeyRequestID).(string)

			if isJSON_Mode {
				// Извлекаем IP за прокси
				clientIP := r.Header.Get("X-Real-IP")
				if clientIP == "" {
					clientIP = r.RemoteAddr
				}

				slog.Info("http_request",
					slog.String("request_id", reqID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", rw.status),
					slog.Int64("latency_ms", time.Since(start).Milliseconds()),
					slog.String("client_ip", clientIP),
					slog.String("user_agent", r.UserAgent()),
				)
			} else {
				// Обычный текстовый вывод лога, как требовал первый пункт
				log.Printf("ID: %s | %s %s | Status: %d | Latency: %v",
					reqID, r.Method, r.URL.Path, rw.status, time.Since(start),
				)
			}
		})
	}
}

package middleware

import (
	"log"
	"net/http"
	"time"
)

// Правильная обертка для перехвата статуса
type responseWriter struct {
	http.ResponseWriter
	status int
}

// Переопределяем именно WriteHeader!
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK, // Дефолтный статус, если WriteHeader не вызовут явно
		}

		next.ServeHTTP(rw, r)

		// Теперь тут будет выводиться реальный статус (например, 404 или 500), а не всегда 200
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

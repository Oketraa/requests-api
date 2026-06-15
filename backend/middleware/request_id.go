package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const keyRequestID key = "request_id"

type key string

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()

		ctx := context.WithValue(r.Context(), keyRequestID, id)
		r = r.WithContext(ctx)

		w.Header().Set("X-request-ID", id)
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"
	"study/database"
)

func UserAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.UserAgent()

		query := `
			INSERT INTO client_user_agents (raw_user_agent)
			VALUES ($1)
		`

		database.DB.Exec(query, userAgent)
		next.ServeHTTP(w, r)
	})
}

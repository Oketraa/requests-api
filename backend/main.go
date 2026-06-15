package main

import (
	"fmt"
	"log"
	"net/http"
	"study/database"
	"study/handlers"
	"study/metrics"
	"study/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	err := database.ConnectDB()

	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/requests", handlers.GetRequests)
	http.HandleFunc("/api/requests/", handlers.RequestByID)
	http.HandleFunc("/api/requests/create", handlers.CreateRequest)
	http.Handle("/metrics", promhttp.Handler())

	wrappedServer := metrics.MetricsMiddleware(http.DefaultServeMux)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", wrappedServer); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/live", handlers.Live)
	mux.HandleFunc("/ready", handlers.Ready)
	mux.HandleFunc("/api/requests", handlers.GetRequests)
	mux.HandleFunc("/api/requests/", handlers.RequestByID)
	mux.HandleFunc("/api/requests/create", handlers.CreateRequest)

	handler := middleware.RequestID(mux)
	handler = middleware.Logger(handler)
	handler = middleware.Recovery(handler)
	handler = middleware.Cors(handler)

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", handler)
}

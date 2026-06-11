package main

import (
	"fmt"
	"log"
	"net/http"
	"study/database"
	"study/handlers"
	"study/monitoring"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}

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

	wrappedServer := monitoring.MetricsMiddleware(http.DefaultServeMux)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", wrappedServer); err != nil {
		log.Fatal(err)
	}
}

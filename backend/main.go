package main

import (
	"fmt"
	"net/http"
	"study/database"
	"study/handlers"
	"study/middleware"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading env")
		return
	}

	err = database.ConnectDB()

	if err != nil {
		fmt.Println("Database connection error:", err)
		return
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

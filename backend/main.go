package main

import (
	"fmt"
	"net/http"
	"study/database"
	"study/handlers"
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

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

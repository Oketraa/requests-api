package main

import (
	"fmt"
	"net/http"
	"study/database"
	"study/handlers"

	"github.com/joho/godotenv"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}

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

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/requests", handlers.GetRequests)
	http.HandleFunc("/api/requests/", handlers.RequestByID)
	http.HandleFunc("/api/requests/create", handlers.CreateRequest)

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

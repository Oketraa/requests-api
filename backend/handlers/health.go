package handlers

import (
	"net/http"
	"study/database"
	"study/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func Live(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "alive",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func Ready(w http.ResponseWriter, r *http.Request) {
	err := database.DB.Ping()

	if err != nil {
		response := map[string]string{
			"status":   "error",
			"database": "unavailable",
		}
		utils.WriteJSON(w, http.StatusServiceUnavailable, response)
		return
	}

	response := map[string]string{
		"status":   "ready",
		"database": "connected",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"study/database"
	"study/metrics"
	"study/models"
	"study/utils"
	"time"
)

func GetRequests(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM requests
		ORDER BY id ASC
	`

	start := time.Now()
	rows, err := database.DB.QueryContext(r.Context(), query)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.Observe(duration)

	defer rows.Close()
	var requests []models.Request

	for rows.Next() {
		var request models.Request

		err := rows.Scan(
			&request.ID,
			&request.Title,
			&request.Description,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		)

		if err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}

		requests = append(requests, request)
	}

	utils.WriteJSON(w, http.StatusOK, requests)
}

func GetRequestByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idString := strings.TrimPrefix(path, "/api/requests/")
	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM requests
		WHERE id = $1
	`

	var request models.Request

	start := time.Now()
	err = database.DB.QueryRowContext(r.Context(), query, id).Scan(
		&request.ID,
		&request.Title,
		&request.Description,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}

		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.Observe(duration)
	utils.WriteJSON(w, http.StatusOK, request)
}

func UpdateRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idString := strings.TrimPrefix(path, "/api/requests/")
	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	var input models.UpdateRequestInput

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !models.AllowedStatus[input.Status] {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE requests
		SET
			title = $1,
			description = $2,
			status = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, created_at, updated_at
	`

	var request models.Request

	request.Title = input.Title
	request.Description = input.Description
	request.Status = input.Status

	start := time.Now()
	err = database.DB.QueryRowContext(r.Context(), query, input.Title, input.Description, input.Status, id).Scan(
		&request.ID,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.Observe(duration)

	if err != nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, request)
}

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	var input models.CreateRequestInput
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO requests (title, description)
		VALUES ($1, $2)
		RETURNING id, status, created_at, updated_at
	`

	var request models.Request

	request.Title = input.Title
	request.Description = input.Description
	request.Status = "new"

	start := time.Now()
	err = database.DB.QueryRowContext(r.Context(), query, input.Title, input.Description).Scan(
		&request.ID,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.Observe(duration)

	utils.WriteJSON(w, http.StatusCreated, request)
}

func RequestByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetRequestByID(w, r)
		return
	}

	if r.Method == http.MethodPut {
		UpdateRequest(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		DeleteRequest(w, r)
		return
	}

	http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
}

func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idString := strings.TrimPrefix(path, "/api/requests/")
	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "ivalid request id", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM requests
		WHERE id = $1
	`
	start := time.Now()
	result, err := database.DB.ExecContext(r.Context(), query, id)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	metrics.DatabaseQueryDuration.Observe(duration)

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"message": "request deleted",
	}

	utils.WriteJSON(w, http.StatusNoContent, response)
}

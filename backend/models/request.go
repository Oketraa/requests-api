package models

import "time"

type Request struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRequestInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateRequestInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var AllowedStatus = map[string]bool{
	"new":         true,
	"in_progress": true,
	"done":        true,
	"canceled":    true,
}

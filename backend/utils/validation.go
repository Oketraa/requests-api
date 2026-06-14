package utils

import (
	"strings"
	"study/models"
)

func CreateValidator(input models.CreateRequestInput) string {
	title := strings.TrimSpace(input.Title)

	if title == "" {
		return "title is required"
	}

	if len(title) > 255 {
		return "title is too long"
	}

	return ""
}

func UpdateValidator(input models.UpdateRequestInput) string {
	title := strings.TrimSpace(input.Title)

	if title == "" {
		return "title is required"
	}

	if len(title) > 255 {
		return "title is too long"
	}

	if !models.AllowedStatus[input.Status] {
		return "invalid status"
	}

	return ""
}

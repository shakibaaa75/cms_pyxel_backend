package utils

import (
	"net/http"

	"cms-backend/database"
)

func RequireMongo(w http.ResponseWriter) bool {
	if database.MongoDB == nil {
		JSONErr(w, "database unavailable", http.StatusServiceUnavailable)
		return false
	}
	return true
}

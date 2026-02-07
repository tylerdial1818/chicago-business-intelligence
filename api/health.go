package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

// HealthHandler returns the health status of the API
func HealthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:   "ok",
			Database: "disconnected",
			Version:  "1.0.0",
		}

		// Check database connection
		if err := db.Ping(); err == nil {
			response.Database = "connected"
		}

		json.NewEncoder(w).Encode(response)
	}
}

package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type PipelineRun struct {
	ID              int       `json:"id"`
	DatasetName     string    `json:"dataset_name"`
	RecordsFetched  int       `json:"records_fetched"`
	RecordsLoaded   int       `json:"records_loaded"`
	RecordsRejected int       `json:"records_rejected"`
	StartedAt       time.Time `json:"started_at"`
	CompletedAt     time.Time `json:"completed_at"`
	Status          string    `json:"status"`
	ErrorMessage    string    `json:"error_message,omitempty"`
}

type PipelineStatusResponse struct {
	Runs       []PipelineRun `json:"runs"`
	TotalRuns  int           `json:"total_runs"`
	LastUpdate time.Time     `json:"last_update"`
}

// PipelineStatusHandler returns the status of all pipeline runs
func PipelineStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT id, dataset_name, records_fetched, records_loaded, records_rejected,
			       started_at, completed_at, status, COALESCE(error_message, '')
			FROM pipeline_runs
			ORDER BY started_at DESC
			LIMIT 50
		`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var runs []PipelineRun
		for rows.Next() {
			var run PipelineRun
			err := rows.Scan(
				&run.ID,
				&run.DatasetName,
				&run.RecordsFetched,
				&run.RecordsLoaded,
				&run.RecordsRejected,
				&run.StartedAt,
				&run.CompletedAt,
				&run.Status,
				&run.ErrorMessage,
			)
			if err != nil {
				continue
			}
			runs = append(runs, run)
		}

		// Get the last update time
		var lastUpdate time.Time
		if len(runs) > 0 {
			lastUpdate = runs[0].CompletedAt
		}

		response := PipelineStatusResponse{
			Runs:       runs,
			TotalRuns:  len(runs),
			LastUpdate: lastUpdate,
		}

		json.NewEncoder(w).Encode(response)
	}
}

package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type TrafficPattern struct {
	Date      time.Time `json:"date"`
	TripCount int       `json:"trip_count"`
	AvgMiles  float64   `json:"avg_miles"`
	AvgFare   float64   `json:"avg_fare"`
}

// TrafficPatternsHandler returns historical traffic patterns for a zip code
func TrafficPatternsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zip")
		if zipCode == "" {
			http.Error(w, "zip parameter is required", http.StatusBadRequest)
			return
		}

		query := `
			SELECT
				DATE(t.trip_start_timestamp) as date,
				COUNT(*) as trip_count,
				AVG(t.trip_miles) as avg_miles,
				AVG(t.trip_total) as avg_fare
			FROM taxi_trips t
			LEFT JOIN zip_community_map zcm ON (
				t.pickup_community_area = zcm.community_area
				OR t.dropoff_community_area = zcm.community_area
			)
			WHERE zcm.zip_code = $1
			AND t.trip_start_timestamp IS NOT NULL
			GROUP BY DATE(t.trip_start_timestamp)
			ORDER BY date DESC
			LIMIT 365
		`

		rows, err := db.Query(query, zipCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var patterns []TrafficPattern
		for rows.Next() {
			var p TrafficPattern
			var avgMiles, avgFare sql.NullFloat64

			err := rows.Scan(&p.Date, &p.TripCount, &avgMiles, &avgFare)
			if err != nil {
				continue
			}

			if avgMiles.Valid {
				p.AvgMiles = avgMiles.Float64
			}
			if avgFare.Valid {
				p.AvgFare = avgFare.Float64
			}

			patterns = append(patterns, p)
		}

		json.NewEncoder(w).Encode(patterns)
	}
}

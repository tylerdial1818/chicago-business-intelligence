package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AirportTraffic struct {
	Airport               string  `json:"airport"`
	DestinationZip        string  `json:"destination_zip"`
	DestinationNeighborhood string  `json:"destination_neighborhood"`
	TripCount             int     `json:"trip_count"`
	AvgMiles              float64 `json:"avg_miles"`
	AvgFare               float64 `json:"avg_fare"`
}

// AirportTrafficHandler returns trips from O'Hare and Midway airports
func AirportTrafficHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT
				CASE
					WHEN t.pickup_community_area = 76 THEN 'O''Hare'
					WHEN t.pickup_community_area = 56 THEN 'Midway'
				END as airport,
				zcm.zip_code as destination_zip,
				zcm.neighborhood_name as destination_neighborhood,
				COUNT(*) as trip_count,
				AVG(t.trip_miles) as avg_miles,
				AVG(t.trip_total) as avg_fare
			FROM taxi_trips t
			JOIN zip_community_map zcm ON t.dropoff_community_area = zcm.community_area
			WHERE t.pickup_community_area IN (76, 56)
			GROUP BY airport, zcm.zip_code, zcm.neighborhood_name
			ORDER BY trip_count DESC
			LIMIT 50
		`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var traffic []AirportTraffic
		for rows.Next() {
			var t AirportTraffic
			var avgMiles, avgFare sql.NullFloat64

			err := rows.Scan(
				&t.Airport,
				&t.DestinationZip,
				&t.DestinationNeighborhood,
				&t.TripCount,
				&avgMiles,
				&avgFare,
			)
			if err != nil {
				continue
			}

			if avgMiles.Valid {
				t.AvgMiles = avgMiles.Float64
			}
			if avgFare.Valid {
				t.AvgFare = avgFare.Float64
			}

			traffic = append(traffic, t)
		}

		json.NewEncoder(w).Encode(traffic)
	}
}

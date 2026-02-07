package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type CCVITrip struct {
	CommunityAreaName string  `json:"community_area_name"`
	CCVIScore         float64 `json:"ccvi_score"`
	CCVICategory      string  `json:"ccvi_category"`
	ZipCode           string  `json:"zip_code"`
	TripsFrom         int     `json:"trips_from"`
	TripsTo           int     `json:"trips_to"`
	TotalTrips        int     `json:"total_trips"`
}

// HighCCVIHandler returns trip data for high CCVI (vulnerability) neighborhoods
func HighCCVIHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT
				cv.community_area_name,
				cv.ccvi_score,
				cv.ccvi_category,
				zcm.zip_code,
				COUNT(CASE WHEN t.pickup_community_area = zcm.community_area THEN 1 END) as trips_from,
				COUNT(CASE WHEN t.dropoff_community_area = zcm.community_area THEN 1 END) as trips_to,
				COUNT(*) as total_trips
			FROM ccvi cv
			JOIN zip_community_map zcm ON cv.community_area_or_zip::integer = zcm.community_area
			LEFT JOIN taxi_trips t ON (
				t.pickup_community_area = zcm.community_area
				OR t.dropoff_community_area = zcm.community_area
			)
			WHERE cv.ccvi_category = 'HIGH'
			AND cv.geography_type = 'CA'
			GROUP BY cv.community_area_name, cv.ccvi_score, cv.ccvi_category, zcm.zip_code
			ORDER BY total_trips DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var trips []CCVITrip
		for rows.Next() {
			var t CCVITrip
			var ccviScore sql.NullFloat64

			err := rows.Scan(
				&t.CommunityAreaName,
				&ccviScore,
				&t.CCVICategory,
				&t.ZipCode,
				&t.TripsFrom,
				&t.TripsTo,
				&t.TotalTrips,
			)
			if err != nil {
				continue
			}

			if ccviScore.Valid {
				t.CCVIScore = ccviScore.Float64
			}

			trips = append(trips, t)
		}

		json.NewEncoder(w).Encode(trips)
	}
}

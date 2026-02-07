package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type CovidAlert struct {
	ZipCode                 string    `json:"zip_code"`
	WeekStart               time.Time `json:"week_start"`
	WeekEnd                 time.Time `json:"week_end"`
	CasesWeekly             int       `json:"cases_weekly"`
	CaseRateWeekly          float64   `json:"case_rate_weekly"`
	TestsWeekly             int       `json:"tests_weekly"`
	PercentTestedPositive   float64   `json:"percent_tested_positive"`
	TaxiTrips               int       `json:"taxi_trips"`
	AlertLevel              string    `json:"alert_level"`
}

// CovidAlertsHandler returns COVID-19 alerts for a specific zip code
func CovidAlertsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zip")
		if zipCode == "" {
			http.Error(w, "zip parameter is required", http.StatusBadRequest)
			return
		}

		query := `
			SELECT
				c.zip_code,
				c.week_start,
				c.week_end,
				c.cases_weekly,
				c.case_rate_weekly,
				c.tests_weekly,
				c.percent_tested_positive,
				COUNT(t.trip_id) as taxi_trips,
				CASE
					WHEN c.case_rate_weekly > 300 OR COUNT(t.trip_id) > 1500 THEN 'HIGH'
					WHEN c.case_rate_weekly > 100 OR COUNT(t.trip_id) > 500 THEN 'MEDIUM'
					ELSE 'LOW'
				END as alert_level
			FROM covid_cases c
			LEFT JOIN zip_community_map zcm ON c.zip_code = zcm.zip_code
			LEFT JOIN taxi_trips t ON (
				zcm.community_area = t.pickup_community_area
				OR zcm.community_area = t.dropoff_community_area
			)
			AND t.trip_start_timestamp BETWEEN c.week_start AND c.week_end
			WHERE c.zip_code = $1
			GROUP BY c.zip_code, c.week_start, c.week_end, c.cases_weekly,
					 c.case_rate_weekly, c.tests_weekly, c.percent_tested_positive
			ORDER BY c.week_start DESC
			LIMIT 50
		`

		rows, err := db.Query(query, zipCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var alerts []CovidAlert
		for rows.Next() {
			var alert CovidAlert
			err := rows.Scan(
				&alert.ZipCode,
				&alert.WeekStart,
				&alert.WeekEnd,
				&alert.CasesWeekly,
				&alert.CaseRateWeekly,
				&alert.TestsWeekly,
				&alert.PercentTestedPositive,
				&alert.TaxiTrips,
				&alert.AlertLevel,
			)
			if err != nil {
				continue
			}
			alerts = append(alerts, alert)
		}

		json.NewEncoder(w).Encode(alerts)
	}
}

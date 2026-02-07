package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tylerdial1818/chicago-business-intelligence/forecast"
)

type ForecastResponse struct {
	ZipCode    string                  `json:"zip_code"`
	Period     string                  `json:"period"`
	Historical []forecast.ForecastResult `json:"historical"`
	Forecast   []forecast.ForecastResult `json:"forecast"`
}

// ForecastHandler generates traffic forecasts for a zip code
func ForecastHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zip")
		if zipCode == "" {
			http.Error(w, "zip parameter is required", http.StatusBadRequest)
			return
		}

		period := r.URL.Query().Get("period")
		if period == "" {
			period = "d" // default to daily
		}

		// Determine aggregation query based on period
		var groupByClause string
		var periodsAhead int
		var windowSize int

		switch period {
		case "d", "daily":
			groupByClause = "DATE(t.trip_start_timestamp)"
			periodsAhead = 30 // forecast 30 days ahead
			windowSize = 14   // use last 14 days for trend
		case "w", "weekly":
			groupByClause = "DATE_TRUNC('week', t.trip_start_timestamp)"
			periodsAhead = 12 // forecast 12 weeks ahead
			windowSize = 8    // use last 8 weeks for trend
		case "m", "monthly":
			groupByClause = "DATE_TRUNC('month', t.trip_start_timestamp)"
			periodsAhead = 6 // forecast 6 months ahead
			windowSize = 6   // use last 6 months for trend
		default:
			http.Error(w, "period must be 'd', 'w', or 'm'", http.StatusBadRequest)
			return
		}

		query := fmt.Sprintf(`
			SELECT
				%s as period,
				COUNT(*) as trip_count
			FROM taxi_trips t
			LEFT JOIN zip_community_map zcm ON (
				t.pickup_community_area = zcm.community_area
				OR t.dropoff_community_area = zcm.community_area
			)
			WHERE zcm.zip_code = $1
			AND t.trip_start_timestamp IS NOT NULL
			GROUP BY period
			ORDER BY period ASC
		`, groupByClause)

		rows, err := db.Query(query, zipCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var historicalData []float64
		var periods []time.Time

		for rows.Next() {
			var periodTime time.Time
			var tripCount int

			err := rows.Scan(&periodTime, &tripCount)
			if err != nil {
				continue
			}

			periods = append(periods, periodTime)
			historicalData = append(historicalData, float64(tripCount))
		}

		if len(historicalData) == 0 {
			http.Error(w, "No data available for this zip code", http.StatusNotFound)
			return
		}

		// Generate forecast
		results := forecast.MovingAverageForecast(historicalData, periodsAhead, windowSize)

		// Split into historical and forecast
		var historical []forecast.ForecastResult
		var forecastData []forecast.ForecastResult

		for i, result := range results {
			if result.Historical {
				if i < len(periods) {
					result.Period = periods[i].Format("2006-01-02")
				}
				historical = append(historical, result)
			} else {
				// Calculate future period dates
				if len(periods) > 0 {
					lastPeriod := periods[len(periods)-1]
					futureIndex := i - len(historicalData) + 1

					var futurePeriod time.Time
					switch period {
					case "d", "daily":
						futurePeriod = lastPeriod.AddDate(0, 0, futureIndex)
					case "w", "weekly":
						futurePeriod = lastPeriod.AddDate(0, 0, futureIndex*7)
					case "m", "monthly":
						futurePeriod = lastPeriod.AddDate(0, futureIndex, 0)
					}
					result.Period = futurePeriod.Format("2006-01-02")
				}
				forecastData = append(forecastData, result)
			}
		}

		response := ForecastResponse{
			ZipCode:    zipCode,
			Period:     period,
			Historical: historical,
			Forecast:   forecastData,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// parseInt safely parses a string to int with a default value
func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

package queries

import (
	"database/sql"
	"fmt"
)

// Requirement 1: COVID-19 alerts by zip code with taxi trip correlation
func GetCOVIDAlertsByZip(db *sql.DB, zipCode string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			c.zip_code,
			c.week_start,
			c.week_end,
			c.cases_weekly,
			c.case_rate_weekly,
			COUNT(DISTINCT t.trip_id) as taxi_trips,
			CASE 
				WHEN CAST(c.case_rate_weekly AS FLOAT) < 50 THEN 'LOW'
				WHEN CAST(c.case_rate_weekly AS FLOAT) BETWEEN 50 AND 150 THEN 'MEDIUM'
				ELSE 'HIGH'
			END as alert_level
		FROM covid_cases c
		LEFT JOIN taxi_trips t ON (c.zip_code = t.pickup_zip_code OR c.zip_code = t.dropoff_zip_code)
			AND t.trip_start_timestamp::date BETWEEN c.week_start::date AND c.week_end::date
		WHERE c.zip_code = $1
		GROUP BY c.zip_code, c.week_start, c.week_end, c.cases_weekly, c.case_rate_weekly
		ORDER BY c.week_start DESC
		LIMIT 52
	`

	rows, err := db.Query(query, zipCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var zip, weekStart, weekEnd, casesWeekly, caseRate, alertLevel string
		var taxiTrips int

		err := rows.Scan(&zip, &weekStart, &weekEnd, &casesWeekly, &caseRate, &taxiTrips, &alertLevel)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"zip_code":      zip,
			"week_start":    weekStart,
			"week_end":      weekEnd,
			"cases_weekly":  casesWeekly,
			"case_rate":     caseRate,
			"taxi_trips":    taxiTrips,
			"alert_level":   alertLevel,
		})
	}

	return results, nil
}

// Requirement 2: Airport traffic to zip codes with COVID correlation
func GetAirportTrafficByZip(db *sql.DB) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			t.dropoff_zip_code as zip_code,
			COUNT(DISTINCT t.trip_id) as trip_count,
			AVG(CAST(c.case_rate_weekly AS FLOAT)) as avg_case_rate
		FROM taxi_trips t
		LEFT JOIN covid_cases c ON t.dropoff_zip_code = c.zip_code
		WHERE t.pickup_zip_code IN ('60666', '60638')  -- O'Hare (60666) and Midway (60638) approximate
		GROUP BY t.dropoff_zip_code
		HAVING COUNT(DISTINCT t.trip_id) > 5
		ORDER BY trip_count DESC
		LIMIT 50
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var zipCode string
		var tripCount int
		var avgCaseRate sql.NullFloat64

		err := rows.Scan(&zipCode, &tripCount, &avgCaseRate)
		if err != nil {
			continue
		}

		caseRate := 0.0
		if avgCaseRate.Valid {
			caseRate = avgCaseRate.Float64
		}

		results = append(results, map[string]interface{}{
			"zip_code":      zipCode,
			"trip_count":    tripCount,
			"avg_case_rate": caseRate,
		})
	}

	return results, nil
}

// Requirement 3: High CCVI neighborhoods with taxi trip volume
func GetHighCCVINeighborhoods(db *sql.DB) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			ccvi.community_area_name,
			ccvi.ccvi_category,
			ccvi.ccvi_score,
			COUNT(DISTINCT t.trip_id) as total_trips
		FROM ccvi
		LEFT JOIN taxi_trips t ON (
			ccvi.community_area_or_zip = t.pickup_zip_code OR 
			ccvi.community_area_or_zip = t.dropoff_zip_code
		)
		WHERE UPPER(ccvi.ccvi_category) = 'HIGH'
		GROUP BY ccvi.community_area_name, ccvi.ccvi_category, ccvi.ccvi_score
		ORDER BY total_trips DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var name, category, score string
		var trips int

		err := rows.Scan(&name, &category, &score, &trips)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"community_name": name,
			"ccvi_category":  category,
			"ccvi_score":     score,
			"total_trips":    trips,
		})
	}

	return results, nil
}

// Requirement 5: Top 5 neighborhoods by unemployment/poverty for investment
func GetInvestmentTargetNeighborhoods(db *sql.DB) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			community_area_name,
			unemployment,
			below_poverty_level,
			per_capita_income,
			COUNT(DISTINCT bp.permit_id) as permit_count
		FROM community_area_unemployment cau
		LEFT JOIN building_permits bp ON cau.community_area = bp.community_area
		WHERE 
			CAST(unemployment AS FLOAT) > 0 
			AND CAST(below_poverty_level AS FLOAT) > 0
		GROUP BY 
			cau.community_area_name, 
			cau.unemployment, 
			cau.below_poverty_level, 
			cau.per_capita_income
		ORDER BY 
			CAST(unemployment AS FLOAT) DESC, 
			CAST(below_poverty_level AS FLOAT) DESC
		LIMIT 5
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var name, unemployment, poverty, income string
		var permitCount int

		err := rows.Scan(&name, &unemployment, &poverty, &income, &permitCount)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"community_name":    name,
			"unemployment":      unemployment,
			"poverty_level":     poverty,
			"per_capita_income": income,
			"permit_count":      permitCount,
		})
	}

	return results, nil
}

// Requirement 6: Small business loan eligibility (low construction permits + low income)
func GetSmallBusinessLoanEligibility(db *sql.DB) ([]map[string]interface{}, error) {
	query := `
		WITH zip_permits AS (
			SELECT 
				contact_1_zipcode as zip_code,
				COUNT(*) as permit_count
			FROM building_permits
			WHERE permit_type = 'PERMIT - NEW CONSTRUCTION'
			GROUP BY contact_1_zipcode
		),
		zip_income AS (
			SELECT 
				community_area,
				per_capita_income
			FROM community_area_unemployment
			WHERE CAST(per_capita_income AS FLOAT) < 30000
		)
		SELECT 
			zp.zip_code,
			zp.permit_count,
			zi.per_capita_income,
			COUNT(DISTINCT bp.permit_id) as total_permits
		FROM zip_permits zp
		LEFT JOIN building_permits bp ON zp.zip_code = bp.contact_1_zipcode
		LEFT JOIN zip_income zi ON bp.community_area = zi.community_area
		WHERE zi.per_capita_income IS NOT NULL
		GROUP BY zp.zip_code, zp.permit_count, zi.per_capita_income
		ORDER BY zp.permit_count ASC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var zipCode, income string
		var permitCount, totalPermits int

		err := rows.Scan(&zipCode, &permitCount, &income, &totalPermits)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"zip_code":          zipCode,
			"new_construction":  permitCount,
			"per_capita_income": income,
			"total_permits":     totalPermits,
		})
	}

	return results, nil
}

// Get traffic patterns by zip code for forecasting (Requirement 4 & 9)
func GetTrafficPatternsByZip(db *sql.DB, zipCode string, days int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(`
		SELECT 
			DATE(trip_start_timestamp) as trip_date,
			COUNT(DISTINCT trip_id) as trip_count
		FROM taxi_trips
		WHERE (pickup_zip_code = $1 OR dropoff_zip_code = $1)
			AND trip_start_timestamp >= NOW() - INTERVAL '%d days'
		GROUP BY DATE(trip_start_timestamp)
		ORDER BY trip_date ASC
	`, days)

	rows, err := db.Query(query, zipCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var tripDate string
		var tripCount int

		err := rows.Scan(&tripDate, &tripCount)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"date":       tripDate,
			"trip_count": tripCount,
		})
	}

	return results, nil
}

// Get all zip codes with data
func GetActiveZipCodes(db *sql.DB) ([]string, error) {
	query := `
		SELECT DISTINCT zip_code 
		FROM covid_cases 
		WHERE zip_code IS NOT NULL 
		ORDER BY zip_code
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	zipCodes := make([]string, 0)
	for rows.Next() {
		var zipCode string
		err := rows.Scan(&zipCode)
		if err != nil {
			continue
		}
		zipCodes = append(zipCodes, zipCode)
	}

	return zipCodes, nil
}

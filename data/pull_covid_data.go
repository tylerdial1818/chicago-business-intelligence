package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// PullCovidCases fetches COVID-19 cases by zip code from the SODA API
func PullCovidCases(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "covid_cases",
		StartedAt: time.Now(),
	}

	log.Println("Pulling COVID-19 cases by zip code...")

	// Fetch from SODA API with pagination
	url := "https://data.cityofchicago.org/resource/yhhz-zm2v.json?$limit=50000"
	resp, err := http.Get(url)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("fetching data: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("reading response: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(body, &records); err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("parsing JSON: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}

	result.RecordsFetched = len(records)
	log.Printf("Fetched %d COVID case records", result.RecordsFetched)

	// Prepare insert statement with ON CONFLICT to handle duplicates
	stmt, err := db.Prepare(`
		INSERT INTO covid_cases (
			zip_code, week_number, week_start, week_end, cases_weekly,
			case_rate_weekly, tests_weekly, test_rate_weekly,
			percent_tested_positive, population
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("preparing statement: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}
	defer stmt.Close()

	// Insert records
	for _, rec := range records {
		zipCode, _ := rec["zip_code"].(string)
		if zipCode == "" {
			result.RecordsRejected++
			continue
		}

		weekNumber, _ := parseFloat(rec["week_number"])
		weekStart, _ := rec["week_start"].(string)
		weekEnd, _ := rec["week_end"].(string)
		casesWeekly, _ := parseFloat(rec["cases_weekly"])
		caseRateWeekly, _ := parseFloat(rec["case_rate_weekly"])
		testsWeekly, _ := parseFloat(rec["tests_weekly"])
		testRateWeekly, _ := parseFloat(rec["test_rate_weekly"])
		percentTested, _ := parseFloat(rec["percent_tested_positive_weekly"])
		population, _ := parseFloat(rec["population"])

		_, err := stmt.Exec(
			zipCode,
			int(weekNumber),
			parseDate(weekStart),
			parseDate(weekEnd),
			int(casesWeekly),
			caseRateWeekly,
			int(testsWeekly),
			testRateWeekly,
			percentTested,
			int(population),
		)
		if err != nil {
			log.Printf("Warning: failed to insert COVID case for zip %s: %v", zipCode, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("COVID cases loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

// PullCovidDaily fetches daily COVID-19 totals from the SODA API
func PullCovidDaily(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "covid_daily",
		StartedAt: time.Now(),
	}

	log.Println("Pulling COVID-19 daily totals...")

	url := "https://data.cityofchicago.org/resource/naz8-j4nc.json?$limit=50000"
	resp, err := http.Get(url)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("fetching data: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("reading response: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(body, &records); err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("parsing JSON: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}

	result.RecordsFetched = len(records)
	log.Printf("Fetched %d COVID daily records", result.RecordsFetched)

	stmt, err := db.Prepare(`
		INSERT INTO covid_daily (date, cases_total, deaths_total, hospitalizations_total)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (date) DO UPDATE SET
			cases_total = EXCLUDED.cases_total,
			deaths_total = EXCLUDED.deaths_total,
			hospitalizations_total = EXCLUDED.hospitalizations_total
	`)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("preparing statement: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}
	defer stmt.Close()

	for _, rec := range records {
		date, _ := rec["date"].(string)
		if date == "" {
			result.RecordsRejected++
			continue
		}

		casesTotal, _ := parseFloat(rec["cases_total"])
		deathsTotal, _ := parseFloat(rec["deaths_total"])
		hospitalizationsTotal, _ := parseFloat(rec["hospitalizations_total"])

		_, err := stmt.Exec(
			parseDate(date),
			int(casesTotal),
			int(deathsTotal),
			int(hospitalizationsTotal),
		)
		if err != nil {
			log.Printf("Warning: failed to insert COVID daily for date %s: %v", date, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("COVID daily loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

// Helper function to parse float values from interface{}
func parseFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		return f, err == nil
	default:
		return 0, false
	}
}

// Helper function to parse date strings
func parseDate(dateStr string) interface{} {
	if dateStr == "" {
		return nil
	}

	// Try various date formats
	formats := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return nil
}

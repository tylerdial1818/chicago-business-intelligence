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

// PullUnemploymentData fetches unemployment and poverty data from the SODA API
func PullUnemploymentData(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "unemployment",
		StartedAt: time.Now(),
	}

	log.Println("Pulling unemployment and poverty data...")

	// Fetch from SODA API
	url := "https://data.cityofchicago.org/resource/iqnk-2tcu.json?$limit=50000"
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
	log.Printf("Fetched %d unemployment records", result.RecordsFetched)

	// Prepare insert statement
	stmt, err := db.Prepare(`
		INSERT INTO unemployment (
			community_area, community_area_name, percent_housing_crowded,
			percent_below_poverty, percent_aged_16_unemployed,
			percent_no_high_school_diploma, percent_aged_under_18_over_64,
			per_capita_income, hardship_index
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (community_area) DO UPDATE SET
			community_area_name = EXCLUDED.community_area_name,
			percent_housing_crowded = EXCLUDED.percent_housing_crowded,
			percent_below_poverty = EXCLUDED.percent_below_poverty,
			percent_aged_16_unemployed = EXCLUDED.percent_aged_16_unemployed,
			percent_no_high_school_diploma = EXCLUDED.percent_no_high_school_diploma,
			percent_aged_under_18_over_64 = EXCLUDED.percent_aged_under_18_over_64,
			per_capita_income = EXCLUDED.per_capita_income,
			hardship_index = EXCLUDED.hardship_index
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
		communityArea, _ := parseFloat(rec["community_area"])
		if communityArea == 0 {
			result.RecordsRejected++
			continue
		}

		communityAreaName, _ := rec["community_area_name"].(string)
		percentHousingCrowded, _ := parseFloat(rec["crowded_housing"])
		percentBelowPoverty, _ := parseFloat(rec["below_poverty_level"])
		percentUnemployed, _ := parseFloat(rec["unemployment"])
		percentNoHSDiploma, _ := parseFloat(rec["no_high_school_diploma"])
		percentUnder18Over64, _ := parseFloat(rec["dependency"])
		perCapitaIncome, _ := parseFloat(rec["per_capita_income"])
		hardshipIndex, _ := parseFloat(rec["hardship_index"])

		_, err := stmt.Exec(
			int(communityArea),
			communityAreaName,
			floatOrNil(percentHousingCrowded),
			floatOrNil(percentBelowPoverty),
			floatOrNil(percentUnemployed),
			floatOrNil(percentNoHSDiploma),
			floatOrNil(percentUnder18Over64),
			intOrNil(perCapitaIncome),
			floatOrNil(hardshipIndex),
		)
		if err != nil {
			log.Printf("Warning: failed to insert unemployment data for community %d: %v", int(communityArea), err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("Unemployment data loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

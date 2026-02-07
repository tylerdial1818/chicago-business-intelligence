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

// PullCCVIData fetches CCVI (Community Vulnerability Index) data from the SODA API
func PullCCVIData(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "ccvi",
		StartedAt: time.Now(),
	}

	log.Println("Pulling CCVI data...")

	// Fetch from SODA API
	url := "https://data.cityofchicago.org/resource/xhc6-88s9.json?$limit=50000"
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
	log.Printf("Fetched %d CCVI records", result.RecordsFetched)

	// Clear existing data before inserting new data
	_, err = db.Exec("DELETE FROM ccvi")
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("clearing existing data: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}

	// Prepare insert statement
	stmt, err := db.Prepare(`
		INSERT INTO ccvi (
			geography_type, community_area_or_zip, community_area_name,
			ccvi_score, ccvi_category
		) VALUES ($1, $2, $3, $4, $5)
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
		geographyType, _ := rec["geography_type"].(string)
		communityAreaOrZip, _ := rec["community_area_or_zip"].(string)
		communityAreaName, _ := rec["community_area_name"].(string)
		ccviScore, _ := parseFloat(rec["ccvi_score"])
		ccviCategory, _ := rec["ccvi_category"].(string)

		// Skip records with missing critical fields
		if geographyType == "" || communityAreaOrZip == "" {
			result.RecordsRejected++
			continue
		}

		_, err := stmt.Exec(
			geographyType,
			communityAreaOrZip,
			communityAreaName,
			floatOrNil(ccviScore),
			ccviCategory,
		)
		if err != nil {
			log.Printf("Warning: failed to insert CCVI data for %s %s: %v", geographyType, communityAreaOrZip, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("CCVI data loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

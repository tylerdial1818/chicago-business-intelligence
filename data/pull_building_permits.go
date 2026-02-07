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

// PullBuildingPermits fetches building permit data from the SODA API
func PullBuildingPermits(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "building_permits",
		StartedAt: time.Now(),
	}

	log.Println("Pulling building permits data...")

	// Fetch from SODA API with pagination
	url := "https://data.cityofchicago.org/resource/ydr8-5enu.json?$limit=50000&$offset=0"
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
	log.Printf("Fetched %d building permit records", result.RecordsFetched)

	// Prepare insert statement
	stmt, err := db.Prepare(`
		INSERT INTO building_permits (
			id, permit_number, permit_type, review_type,
			application_start_date, issue_date, street_number,
			street_direction, street_name, suffix, work_description,
			building_fee, zoning_fee, other_fee, subtotal_paid,
			latitude, longitude, community_area, zip_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (id) DO NOTHING
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
		id, _ := rec["id"].(string)
		if id == "" {
			result.RecordsRejected++
			continue
		}

		permitNumber, _ := rec["permit_"].(string)
		permitType, _ := rec["permit_type"].(string)
		reviewType, _ := rec["review_type"].(string)
		applicationStartDate, _ := rec["application_start_date"].(string)
		issueDate, _ := rec["issue_date"].(string)
		streetNumber, _ := rec["street_number"].(string)
		streetDirection, _ := rec["street_direction"].(string)
		streetName, _ := rec["street_name"].(string)
		suffix, _ := rec["suffix"].(string)
		workDescription, _ := rec["work_description"].(string)

		buildingFee, _ := parseFloat(rec["building_fee_paid"])
		zoningFee, _ := parseFloat(rec["zoning_fee_paid"])
		otherFee, _ := parseFloat(rec["other_fee_paid"])
		subtotalPaid, _ := parseFloat(rec["subtotal_paid"])

		latitude, _ := parseFloat(rec["latitude"])
		longitude, _ := parseFloat(rec["longitude"])
		communityArea, _ := parseFloat(rec["community_area"])
		zipCode, _ := rec["zip_code"].(string)

		_, err := stmt.Exec(
			id,
			permitNumber,
			permitType,
			reviewType,
			parseDate(applicationStartDate),
			parseDate(issueDate),
			streetNumber,
			streetDirection,
			streetName,
			suffix,
			workDescription,
			floatOrNil(buildingFee),
			floatOrNil(zoningFee),
			floatOrNil(otherFee),
			floatOrNil(subtotalPaid),
			floatOrNil(latitude),
			floatOrNil(longitude),
			intOrNil(communityArea),
			zipCode,
		)
		if err != nil {
			log.Printf("Warning: failed to insert building permit %s: %v", id, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("Building permits loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

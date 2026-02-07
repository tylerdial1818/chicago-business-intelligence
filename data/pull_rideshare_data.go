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

// PullRideshareData fetches TNP/rideshare trip data from the SODA API
// This data is loaded into the same taxi_trips table with source='tnp'
func PullRideshareData(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "rideshare_trips",
		StartedAt: time.Now(),
	}

	log.Println("Pulling TNP/rideshare trips data...")

	// Fetch from SODA API with pagination
	url := "https://data.cityofchicago.org/resource/m6dm-c72p.json?$limit=50000&$offset=0"
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
	log.Printf("Fetched %d rideshare trip records", result.RecordsFetched)

	// Prepare insert statement
	// Note: TNP data may have different fields, but we map to the same table structure
	stmt, err := db.Prepare(`
		INSERT INTO taxi_trips (
			trip_id, trip_start_timestamp, trip_end_timestamp, trip_seconds,
			trip_miles, pickup_community_area, dropoff_community_area,
			pickup_centroid_latitude, pickup_centroid_longitude,
			dropoff_centroid_latitude, dropoff_centroid_longitude,
			fare, tips, trip_total, source
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (trip_id) DO NOTHING
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
		tripID, _ := rec["trip_id"].(string)
		if tripID == "" {
			result.RecordsRejected++
			continue
		}

		tripStartTime, _ := rec["trip_start_timestamp"].(string)
		tripEndTime, _ := rec["trip_end_timestamp"].(string)
		tripSeconds, _ := parseFloat(rec["trip_seconds"])
		tripMiles, _ := parseFloat(rec["trip_miles"])
		pickupCommunity, _ := parseFloat(rec["pickup_community_area"])
		dropoffCommunity, _ := parseFloat(rec["dropoff_community_area"])

		// TNP data may have different field names for lat/lon
		pickupLat, _ := parseFloat(rec["pickup_centroid_latitude"])
		pickupLon, _ := parseFloat(rec["pickup_centroid_longitude"])
		dropoffLat, _ := parseFloat(rec["dropoff_centroid_latitude"])
		dropoffLon, _ := parseFloat(rec["dropoff_centroid_longitude"])

		// TNP data may have fare info, but format may differ
		fare, _ := parseFloat(rec["fare"])
		tips, _ := parseFloat(rec["tips"])
		tripTotal, _ := parseFloat(rec["trip_total"])

		_, err := stmt.Exec(
			tripID,
			parseTimestamp(tripStartTime),
			parseTimestamp(tripEndTime),
			intOrNil(tripSeconds),
			floatOrNil(tripMiles),
			intOrNil(pickupCommunity),
			intOrNil(dropoffCommunity),
			floatOrNil(pickupLat),
			floatOrNil(pickupLon),
			floatOrNil(dropoffLat),
			floatOrNil(dropoffLon),
			floatOrNil(fare),
			floatOrNil(tips),
			floatOrNil(tripTotal),
			"tnp", // Mark as TNP/rideshare source
		)
		if err != nil {
			log.Printf("Warning: failed to insert rideshare trip %s: %v", tripID, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("Rideshare trips loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

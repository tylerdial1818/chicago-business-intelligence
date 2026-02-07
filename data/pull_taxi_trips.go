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

// PullTaxiTrips fetches taxi trip data from the SODA API
func PullTaxiTrips(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "taxi_trips",
		StartedAt: time.Now(),
	}

	log.Println("Pulling taxi trips data...")

	// Fetch from SODA API with pagination
	// Note: Using a smaller limit to avoid timeout; in production, implement pagination
	url := "https://data.cityofchicago.org/resource/wrvz-psew.json?$limit=50000&$offset=0"
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
	log.Printf("Fetched %d taxi trip records", result.RecordsFetched)

	// Prepare insert statement
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
		pickupLat, _ := parseFloat(rec["pickup_centroid_latitude"])
		pickupLon, _ := parseFloat(rec["pickup_centroid_longitude"])
		dropoffLat, _ := parseFloat(rec["dropoff_centroid_latitude"])
		dropoffLon, _ := parseFloat(rec["dropoff_centroid_longitude"])
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
			"taxi",
		)
		if err != nil {
			log.Printf("Warning: failed to insert taxi trip %s: %v", tripID, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("Taxi trips loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

// Helper function to parse timestamp strings
func parseTimestamp(timestampStr string) interface{} {
	if timestampStr == "" {
		return nil
	}

	formats := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t
		}
	}

	return nil
}

// Helper functions to return nil for zero values
func intOrNil(val float64) interface{} {
	if val == 0 {
		return nil
	}
	return int(val)
}

func floatOrNil(val float64) interface{} {
	if val == 0 {
		return nil
	}
	return val
}

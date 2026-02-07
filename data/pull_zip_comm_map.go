package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ZipCommunityMapping represents a mapping between zip code and community area
type ZipCommunityMapping struct {
	ZipCode           string
	CommunityArea     int
	CommunityAreaName string
	Neighborhood      string
}

// Chicago ZIP Code to Community Area mapping
// Source: Chicago Tribune, City of Chicago Boundaries dataset
// Note: ZIP codes can span multiple community areas
var zipCommunityMap = []ZipCommunityMapping{
	{"60601", 32, "Loop", "Loop"},
	{"60602", 32, "Loop", "Loop"},
	{"60603", 32, "Loop", "Loop"},
	{"60604", 32, "Loop", "Loop"},
	{"60605", 33, "Near South Side", "South Loop"},
	{"60606", 32, "Loop", "Loop"},
	{"60607", 28, "Near West Side", "West Loop"},
	{"60608", 31, "Lower West Side", "Pilsen"},
	{"60609", 61, "New City", "Back of the Yards"},
	{"60610", 8, "Near North Side", "Old Town"},
	{"60611", 8, "Near North Side", "Streeterville"},
	{"60612", 28, "Near West Side", "Near West Side"},
	{"60613", 6, "Lake View", "Lakeview"},
	{"60614", 7, "Lincoln Park", "Lincoln Park"},
	{"60615", 41, "Hyde Park", "Hyde Park"},
	{"60616", 33, "Near South Side", "Chinatown"},
	{"60617", 43, "South Shore", "South Shore"},
	{"60618", 5, "North Center", "North Center"},
	{"60619", 44, "Chatham", "Chatham"},
	{"60620", 67, "West Englewood", "West Englewood"},
	{"60621", 68, "Englewood", "Englewood"},
	{"60622", 24, "West Town", "Wicker Park"},
	{"60623", 30, "South Lawndale", "Little Village"},
	{"60624", 29, "North Lawndale", "North Lawndale"},
	{"60625", 4, "Lincoln Square", "Lincoln Square"},
	{"60626", 1, "Rogers Park", "Rogers Park"},
	{"60628", 49, "Roseland", "Roseland"},
	{"60629", 65, "West Lawn", "Chicago Lawn"},
	{"60630", 12, "Forest Glen", "Forest Glen"},
	{"60631", 10, "Norwood Park", "Norwood Park"},
	{"60632", 56, "Garfield Ridge", "Clearing"},
	{"60633", 55, "Hegewisch", "Hegewisch"},
	{"60634", 17, "Dunning", "Dunning"},
	{"60636", 66, "Chicago Lawn", "Chicago Lawn"},
	{"60637", 42, "Woodlawn", "Woodlawn"},
	{"60638", 56, "Garfield Ridge", "Garfield Ridge"},
	{"60639", 19, "Belmont Cragin", "Belmont Cragin"},
	{"60640", 3, "Uptown", "Uptown"},
	{"60641", 16, "Irving Park", "Irving Park"},
	{"60642", 24, "West Town", "West Town"},
	{"60643", 72, "Beverly", "Beverly"},
	{"60644", 25, "Austin", "Austin"},
	{"60645", 2, "West Ridge", "West Ridge"},
	{"60646", 13, "North Park", "North Park"},
	{"60647", 22, "Logan Square", "Logan Square"},
	{"60649", 43, "South Shore", "South Shore"},
	{"60651", 23, "Humboldt Park", "Humboldt Park"},
	{"60652", 65, "West Lawn", "West Lawn"},
	{"60653", 35, "Douglas", "Bronzeville"},
	{"60654", 8, "Near North Side", "River North"},
	{"60655", 72, "Beverly", "Morgan Park"},
	{"60656", 10, "Norwood Park", "Norwood Park"},
	{"60657", 6, "Lake View", "Lakeview"},
	{"60659", 2, "West Ridge", "West Ridge"},
	{"60660", 3, "Uptown", "Edgewater"},
	{"60661", 28, "Near West Side", "West Loop"},
	{"60666", 76, "O'Hare", "O'Hare"},
	{"60707", 25, "Austin", "Montclare"},
	{"60827", 54, "Riverdale", "Riverdale"},
}

// PullResult tracks the outcome of a data pull operation
type PullResult struct {
	Dataset         string
	RecordsFetched  int
	RecordsLoaded   int
	RecordsRejected int
	StartedAt       time.Time
	CompletedAt     time.Time
	Status          string
	Error           string
}

// LoadZipCommunityMap loads the hardcoded zip-to-community mapping into the database
func LoadZipCommunityMap(db *sql.DB) PullResult {
	result := PullResult{
		Dataset:   "zip_community_map",
		StartedAt: time.Now(),
	}

	log.Println("Loading ZIP code to community area mapping...")

	// Clear existing data
	_, err := db.Exec("DELETE FROM zip_community_map")
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("clearing existing data: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}

	result.RecordsFetched = len(zipCommunityMap)

	// Insert each mapping
	stmt, err := db.Prepare(`
		INSERT INTO zip_community_map (zip_code, community_area, community_area_name, neighborhood_name)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("preparing statement: %v", err)
		result.CompletedAt = time.Now()
		logPipelineRun(db, result)
		return result
	}
	defer stmt.Close()

	for _, mapping := range zipCommunityMap {
		_, err := stmt.Exec(
			mapping.ZipCode,
			mapping.CommunityArea,
			mapping.CommunityAreaName,
			mapping.Neighborhood,
		)
		if err != nil {
			log.Printf("Warning: failed to insert mapping for zip %s: %v", mapping.ZipCode, err)
			result.RecordsRejected++
			continue
		}
		result.RecordsLoaded++
	}

	result.Status = "SUCCESS"
	result.CompletedAt = time.Now()
	log.Printf("ZIP mapping loaded: %d records inserted, %d rejected", result.RecordsLoaded, result.RecordsRejected)

	logPipelineRun(db, result)
	return result
}

// logPipelineRun logs the result of a data pull to the pipeline_runs table
func logPipelineRun(db *sql.DB, result PullResult) {
	_, err := db.Exec(`
		INSERT INTO pipeline_runs (
			dataset_name, records_fetched, records_loaded, records_rejected,
			started_at, completed_at, status, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		result.Dataset,
		result.RecordsFetched,
		result.RecordsLoaded,
		result.RecordsRejected,
		result.StartedAt,
		result.CompletedAt,
		result.Status,
		result.Error,
	)
	if err != nil {
		log.Printf("Warning: failed to log pipeline run: %v", err)
	}
}

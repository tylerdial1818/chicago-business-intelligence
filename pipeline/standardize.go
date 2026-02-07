package pipeline

import (
	"database/sql"
	"log"
	"strings"
)

// StandardizeData performs data cleaning and normalization
func StandardizeData(db *sql.DB) error {
	log.Println("Starting data standardization...")

	// Standardize ZIP codes to 5 digits
	if err := standardizeZipCodes(db); err != nil {
		log.Printf("Warning: ZIP code standardization failed: %v", err)
	}

	// Assign ZIP codes to taxi trips using community area mapping
	if err := assignZipCodesToTrips(db); err != nil {
		log.Printf("Warning: ZIP code assignment failed: %v", err)
	}

	log.Println("Data standardization completed!")
	return nil
}

// standardizeZipCodes normalizes ZIP codes to 5-digit format
func standardizeZipCodes(db *sql.DB) error {
	log.Println("Standardizing ZIP codes...")

	tables := []struct {
		table  string
		column string
	}{
		{"covid_cases", "zip_code"},
		{"building_permits", "zip_code"},
		{"zip_community_map", "zip_code"},
	}

	for _, t := range tables {
		// Remove non-numeric characters and trim to 5 digits
		query := `
			UPDATE ` + t.table + `
			SET ` + t.column + ` = LEFT(REGEXP_REPLACE(` + t.column + `, '[^0-9]', '', 'g'), 5)
			WHERE ` + t.column + ` IS NOT NULL AND ` + t.column + ` != ''
		`
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
		log.Printf("Standardized ZIP codes in %s.%s", t.table, t.column)
	}

	return nil
}

// assignZipCodesToTrips assigns ZIP codes to taxi trips based on community area
func assignZipCodesToTrips(db *sql.DB) error {
	log.Println("Assigning ZIP codes to taxi trips...")

	// Update pickup_zip based on pickup_community_area
	query := `
		UPDATE taxi_trips t
		SET pickup_zip = zcm.zip_code
		FROM (
			SELECT DISTINCT ON (community_area) community_area, zip_code
			FROM zip_community_map
		) zcm
		WHERE t.pickup_community_area = zcm.community_area
		AND (t.pickup_zip IS NULL OR t.pickup_zip = '')
	`
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Assigned pickup ZIP codes to %d trips", rowsAffected)

	// Update dropoff_zip based on dropoff_community_area
	query = `
		UPDATE taxi_trips t
		SET dropoff_zip = zcm.zip_code
		FROM (
			SELECT DISTINCT ON (community_area) community_area, zip_code
			FROM zip_community_map
		) zcm
		WHERE t.dropoff_community_area = zcm.community_area
		AND (t.dropoff_zip IS NULL OR t.dropoff_zip = '')
	`
	result, err = db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, _ = result.RowsAffected()
	log.Printf("Assigned dropoff ZIP codes to %d trips", rowsAffected)

	return nil
}

// NormalizeString trims whitespace and converts to lowercase
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

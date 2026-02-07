package db

import (
	"database/sql"
	"log"

	"github.com/tylerdial1818/chicago-business-intelligence/data"
)

// SeedDatabase orchestrates all data pulls in the correct order
func SeedDatabase(db *sql.DB) error {
	log.Println("Starting database seeding...")

	// 1. Load ZIP code to community area mapping first (required for joins)
	log.Println("=== Step 1: Loading ZIP code to community area mapping ===")
	result := data.LoadZipCommunityMap(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: ZIP mapping load failed: %s", result.Error)
	}

	// 2. Load unemployment/poverty data (community area based)
	log.Println("=== Step 2: Loading unemployment and poverty data ===")
	result = data.PullUnemploymentData(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: Unemployment data load failed: %s", result.Error)
	}

	// 3. Load CCVI data (community area based)
	log.Println("=== Step 3: Loading CCVI data ===")
	result = data.PullCCVIData(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: CCVI data load failed: %s", result.Error)
	}

	// 4. Load COVID cases (zip code based)
	log.Println("=== Step 4: Loading COVID-19 cases by zip code ===")
	result = data.PullCovidCases(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: COVID cases load failed: %s", result.Error)
	}

	// 5. Load COVID daily data
	log.Println("=== Step 5: Loading COVID-19 daily totals ===")
	result = data.PullCovidDaily(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: COVID daily data load failed: %s", result.Error)
	}

	// 6. Load building permits
	log.Println("=== Step 6: Loading building permits data ===")
	result = data.PullBuildingPermits(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: Building permits load failed: %s", result.Error)
	}

	// 7. Load taxi trips
	log.Println("=== Step 7: Loading taxi trips data ===")
	result = data.PullTaxiTrips(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: Taxi trips load failed: %s", result.Error)
	}

	// 8. Load TNP/rideshare trips
	log.Println("=== Step 8: Loading rideshare/TNP trips data ===")
	result = data.PullRideshareData(db)
	if result.Status != "SUCCESS" {
		log.Printf("Warning: Rideshare trips load failed: %s", result.Error)
	}

	log.Println("Database seeding completed!")
	return nil
}

// RefreshData re-runs all data pulls (useful for the refresh endpoint)
func RefreshData(db *sql.DB) error {
	return SeedDatabase(db)
}

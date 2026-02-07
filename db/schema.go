package db

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateTables creates all required database tables
func CreateTables(db *sql.DB) error {
	log.Println("Creating database tables...")

	tables := []string{
		// Core data tables
		`CREATE TABLE IF NOT EXISTS taxi_trips (
			trip_id VARCHAR(50) PRIMARY KEY,
			trip_start_timestamp TIMESTAMP,
			trip_end_timestamp TIMESTAMP,
			trip_seconds INTEGER,
			trip_miles NUMERIC(10,2),
			pickup_community_area INTEGER,
			dropoff_community_area INTEGER,
			pickup_zip VARCHAR(10),
			dropoff_zip VARCHAR(10),
			pickup_centroid_latitude NUMERIC(12,8),
			pickup_centroid_longitude NUMERIC(12,8),
			dropoff_centroid_latitude NUMERIC(12,8),
			dropoff_centroid_longitude NUMERIC(12,8),
			fare NUMERIC(10,2),
			tips NUMERIC(10,2),
			trip_total NUMERIC(10,2),
			source VARCHAR(10) DEFAULT 'taxi'
		)`,

		`CREATE TABLE IF NOT EXISTS covid_cases (
			id SERIAL PRIMARY KEY,
			zip_code VARCHAR(10),
			week_number INTEGER,
			week_start DATE,
			week_end DATE,
			cases_weekly INTEGER,
			case_rate_weekly NUMERIC(10,2),
			tests_weekly INTEGER,
			test_rate_weekly NUMERIC(10,2),
			percent_tested_positive NUMERIC(5,2),
			population INTEGER
		)`,

		`CREATE TABLE IF NOT EXISTS covid_daily (
			date DATE PRIMARY KEY,
			cases_total INTEGER,
			deaths_total INTEGER,
			hospitalizations_total INTEGER
		)`,

		`CREATE TABLE IF NOT EXISTS building_permits (
			id VARCHAR(50) PRIMARY KEY,
			permit_number VARCHAR(50),
			permit_type VARCHAR(100),
			review_type VARCHAR(100),
			application_start_date DATE,
			issue_date DATE,
			street_number VARCHAR(20),
			street_direction VARCHAR(5),
			street_name VARCHAR(100),
			suffix VARCHAR(10),
			work_description TEXT,
			building_fee NUMERIC(12,2),
			zoning_fee NUMERIC(12,2),
			other_fee NUMERIC(12,2),
			subtotal_paid NUMERIC(12,2),
			latitude NUMERIC(12,8),
			longitude NUMERIC(12,8),
			community_area INTEGER,
			zip_code VARCHAR(10)
		)`,

		`CREATE TABLE IF NOT EXISTS ccvi (
			id SERIAL PRIMARY KEY,
			geography_type VARCHAR(20),
			community_area_or_zip VARCHAR(20),
			community_area_name VARCHAR(100),
			ccvi_score NUMERIC(8,4),
			ccvi_category VARCHAR(20)
		)`,

		`CREATE TABLE IF NOT EXISTS unemployment (
			community_area INTEGER PRIMARY KEY,
			community_area_name VARCHAR(100),
			percent_housing_crowded NUMERIC(5,2),
			percent_below_poverty NUMERIC(5,2),
			percent_aged_16_unemployed NUMERIC(5,2),
			percent_no_high_school_diploma NUMERIC(5,2),
			percent_aged_under_18_over_64 NUMERIC(5,2),
			per_capita_income INTEGER,
			hardship_index NUMERIC(5,2)
		)`,

		`CREATE TABLE IF NOT EXISTS zip_community_map (
			id SERIAL PRIMARY KEY,
			zip_code VARCHAR(10),
			community_area INTEGER,
			community_area_name VARCHAR(100),
			neighborhood_name VARCHAR(100)
		)`,

		// Pipeline tracking table
		`CREATE TABLE IF NOT EXISTS pipeline_runs (
			id SERIAL PRIMARY KEY,
			dataset_name VARCHAR(50),
			records_fetched INTEGER,
			records_loaded INTEGER,
			records_rejected INTEGER,
			started_at TIMESTAMP,
			completed_at TIMESTAMP,
			status VARCHAR(20),
			error_message TEXT
		)`,
	}

	for i, tableSQL := range tables {
		_, err := db.Exec(tableSQL)
		if err != nil {
			return fmt.Errorf("creating table %d: %w", i+1, err)
		}
	}

	log.Println("All database tables created successfully")
	return nil
}

// CreateIndexes creates indexes for improved query performance
func CreateIndexes(db *sql.DB) error {
	log.Println("Creating database indexes...")

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_taxi_trips_pickup_community ON taxi_trips(pickup_community_area)`,
		`CREATE INDEX IF NOT EXISTS idx_taxi_trips_dropoff_community ON taxi_trips(dropoff_community_area)`,
		`CREATE INDEX IF NOT EXISTS idx_taxi_trips_pickup_zip ON taxi_trips(pickup_zip)`,
		`CREATE INDEX IF NOT EXISTS idx_taxi_trips_dropoff_zip ON taxi_trips(dropoff_zip)`,
		`CREATE INDEX IF NOT EXISTS idx_taxi_trips_timestamp ON taxi_trips(trip_start_timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_taxi_trips_source ON taxi_trips(source)`,
		`CREATE INDEX IF NOT EXISTS idx_covid_cases_zip ON covid_cases(zip_code)`,
		`CREATE INDEX IF NOT EXISTS idx_covid_cases_week ON covid_cases(week_start)`,
		`CREATE INDEX IF NOT EXISTS idx_building_permits_community ON building_permits(community_area)`,
		`CREATE INDEX IF NOT EXISTS idx_building_permits_zip ON building_permits(zip_code)`,
		`CREATE INDEX IF NOT EXISTS idx_building_permits_type ON building_permits(permit_type)`,
		`CREATE INDEX IF NOT EXISTS idx_zip_community_zip ON zip_community_map(zip_code)`,
		`CREATE INDEX IF NOT EXISTS idx_zip_community_area ON zip_community_map(community_area)`,
		`CREATE INDEX IF NOT EXISTS idx_ccvi_category ON ccvi(ccvi_category)`,
		`CREATE INDEX IF NOT EXISTS idx_pipeline_runs_dataset ON pipeline_runs(dataset_name)`,
		`CREATE INDEX IF NOT EXISTS idx_pipeline_runs_status ON pipeline_runs(status)`,
	}

	for i, indexSQL := range indexes {
		_, err := db.Exec(indexSQL)
		if err != nil {
			return fmt.Errorf("creating index %d: %w", i+1, err)
		}
	}

	log.Println("All database indexes created successfully")
	return nil
}

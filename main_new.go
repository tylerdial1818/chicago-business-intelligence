package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tylerdial1818/chicago-business-intelligence/api"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	fmt.Println("Initializing database connection...")

	// Database connection - use environment variable for flexibility
	dbConnection := os.Getenv("DATABASE_URL")
	if dbConnection == "" {
		// Default for local development
		dbConnection = "user=postgres dbname=chicago_business_intelligence password=root host=localhost sslmode=disable port=5432"
	}

	db, err = sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal("Couldn't open database connection:", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Println("Warning: Couldn't immediately connect to database (it may not be ready yet)")
	} else {
		log.Println("✓ Database connection successful")
	}
}

func main() {
	// Register API routes
	api.RegisterRoutes(db)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Chicago Business Intelligence API starting on port %s", port)
	log.Printf("📊 API Endpoints:")
	log.Printf("   GET  /health - Health check")
	log.Printf("   GET  /api/covid-alerts?zip={zip} - COVID alerts by zip code")
	log.Printf("   GET  /api/airport-traffic - Airport traffic patterns")
	log.Printf("   GET  /api/high-ccvi - High vulnerability neighborhoods")
	log.Printf("   GET  /api/investment-targets - Top neighborhoods for investment")
	log.Printf("   GET  /api/small-business-loans - Small business loan eligibility")
	log.Printf("   GET  /api/traffic-patterns?zip={zip} - Traffic patterns by zip")
	log.Printf("   GET  /api/zip-codes - List of active zip codes")

	// Optional: Start background data collection
	// Uncomment if you want to collect data automatically
	// go startDataCollection()

	// Start HTTP server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Background data collection (optional - can be triggered separately)
func startDataCollection() {
	log.Println("Starting background data collection...")

	// Initial collection
	collectData()

	// Then collect once per day
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		collectData()
	}
}

func collectData() {
	log.Println("Collecting data from Chicago Data Portal...")

	// These functions are in the original main.go
	// For now, they can be triggered manually or via a separate endpoint
	// go GetCommunityAreaUnemployment(db)
	// go GetBuildingPermits(db)
	// go GetTaxiTrips(db)
	// go GetCovidDetails(db)
	// go GetCCVIDetails(db)

	log.Println("Data collection initiated")
}

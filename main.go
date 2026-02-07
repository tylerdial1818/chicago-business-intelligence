package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tylerdial1818/chicago-business-intelligence/api"
	"github.com/tylerdial1818/chicago-business-intelligence/db"
	"github.com/tylerdial1818/chicago-business-intelligence/pipeline"
)

func main() {
	log.Println("Starting Chicago Business Intelligence Platform...")

	// Connect to database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tables
	log.Println("Creating database tables...")
	if err := db.CreateTables(database); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Create indexes
	log.Println("Creating database indexes...")
	if err := db.CreateIndexes(database); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Check if we should seed the database
	// Only seed if SEED_DATABASE environment variable is set to "true"
	// or if this is the first run (check if any data exists)
	shouldSeed := os.Getenv("SEED_DATABASE") == "true"

	if !shouldSeed {
		// Check if database has any data
		var count int
		err := database.QueryRow("SELECT COUNT(*) FROM zip_community_map").Scan(&count)
		if err == nil && count == 0 {
			shouldSeed = true
			log.Println("Database is empty, will seed with initial data...")
		}
	}

	if shouldSeed {
		log.Println("Seeding database...")
		if err := db.SeedDatabase(database); err != nil {
			log.Printf("Warning: Database seeding encountered errors: %v", err)
		}

		// Run data preprocessing
		log.Println("Running data preprocessing...")
		if err := pipeline.StandardizeData(database); err != nil {
			log.Printf("Warning: Data standardization encountered errors: %v", err)
		}
		if err := pipeline.TransformData(database); err != nil {
			log.Printf("Warning: Data transformation encountered errors: %v", err)
		}
		if err := pipeline.EnrichData(database); err != nil {
			log.Printf("Warning: Data enrichment encountered errors: %v", err)
		}
	} else {
		log.Println("Skipping database seeding (data already exists)")
		log.Println("Set SEED_DATABASE=true to force re-seeding")
	}

	// Register routes
	handler := api.RegisterRoutes(database)

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Health check available at http://localhost:%s/health", port)
	log.Printf("API endpoints available at http://localhost:%s/api/*", port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

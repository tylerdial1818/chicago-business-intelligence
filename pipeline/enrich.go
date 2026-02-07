package pipeline

import (
	"database/sql"
	"log"
)

// EnrichData adds derived variables and calculated fields
func EnrichData(db *sql.DB) error {
	log.Println("Starting data enrichment...")

	// Most enrichment is done in the queries themselves
	// This function is a placeholder for future enrichment tasks

	log.Println("Data enrichment completed!")
	return nil
}

// CalculateAlertLevels determines COVID alert levels based on case rates and trip volumes
// Alert logic:
// - LOW: case_rate < 100 AND taxi_trips < 500/week
// - MEDIUM: case_rate 100-300 OR taxi_trips 500-1500/week
// - HIGH: case_rate > 300 OR taxi_trips > 1500/week
func CalculateAlertLevels(caseRate float64, tripCount int) string {
	if caseRate > 300 || tripCount > 1500 {
		return "HIGH"
	} else if caseRate > 100 || tripCount > 500 {
		return "MEDIUM"
	}
	return "LOW"
}

// IsInvestmentTarget identifies if a neighborhood qualifies as an investment target
// Criteria: high unemployment + high poverty
func IsInvestmentTarget(unemploymentRate float64, povertyRate float64) bool {
	return unemploymentRate > 15.0 && povertyRate > 20.0
}

// IsSmallBusinessLoanEligible determines if a zip code qualifies for small business loans
// Criteria: per capita income < $30,000 and low new construction permits
func IsSmallBusinessLoanEligible(perCapitaIncome int, newConstructionCount int) bool {
	return perCapitaIncome < 30000 && newConstructionCount < 50
}

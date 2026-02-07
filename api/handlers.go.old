package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tylerdial1818/chicago-business-intelligence/queries"
)

// CORS middleware
func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// JSON response helper
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error response helper
func respondError(w http.ResponseWriter, message string, status int) {
	respondJSON(w, map[string]string{"error": message}, status)
}

// Health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status":  "healthy",
		"service": "Chicago Business Intelligence API",
	}, http.StatusOK)
}

// Get COVID alerts by zip code
func COVIDAlertsByZipHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zip")
		if zipCode == "" {
			respondError(w, "zip parameter is required", http.StatusBadRequest)
			return
		}

		results, err := queries.GetCOVIDAlertsByZip(db, zipCode)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"zip_code": zipCode,
			"data":     results,
		}, http.StatusOK)
	})
}

// Get airport traffic patterns
func AirportTrafficHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		results, err := queries.GetAirportTrafficByZip(db)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"data": results,
		}, http.StatusOK)
	})
}

// Get high CCVI neighborhoods
func HighCCVINeighborhoodsHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		results, err := queries.GetHighCCVINeighborhoods(db)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"data": results,
		}, http.StatusOK)
	})
}

// Get investment target neighborhoods
func InvestmentTargetsHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		results, err := queries.GetInvestmentTargetNeighborhoods(db)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"data": results,
		}, http.StatusOK)
	})
}

// Get small business loan eligibility
func SmallBusinessLoansHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		results, err := queries.GetSmallBusinessLoanEligibility(db)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"data": results,
		}, http.StatusOK)
	})
}

// Get traffic patterns
func TrafficPatternsHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zip")
		if zipCode == "" {
			respondError(w, "zip parameter is required", http.StatusBadRequest)
			return
		}

		// Default to 90 days
		days := 90
		results, err := queries.GetTrafficPatternsByZip(db, zipCode, days)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"zip_code": zipCode,
			"days":     days,
			"data":     results,
		}, http.StatusOK)
	})
}

// Get active zip codes
func ActiveZipCodesHandler(db *sql.DB) http.HandlerFunc {
	return EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		results, err := queries.GetActiveZipCodes(db)
		if err != nil {
			respondError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"data": results,
		}, http.StatusOK)
	})
}

// Register all API routes
func RegisterRoutes(db *sql.DB) {
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/api/covid-alerts", COVIDAlertsByZipHandler(db))
	http.HandleFunc("/api/airport-traffic", AirportTrafficHandler(db))
	http.HandleFunc("/api/high-ccvi", HighCCVINeighborhoodsHandler(db))
	http.HandleFunc("/api/investment-targets", InvestmentTargetsHandler(db))
	http.HandleFunc("/api/small-business-loans", SmallBusinessLoansHandler(db))
	http.HandleFunc("/api/traffic-patterns", TrafficPatternsHandler(db))
	http.HandleFunc("/api/zip-codes", ActiveZipCodesHandler(db))
}

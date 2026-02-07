package api

import (
	"database/sql"
	"net/http"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", HealthHandler(db))

	// API endpoints
	mux.HandleFunc("/api/pipeline-status", PipelineStatusHandler(db))
	mux.HandleFunc("/api/zip-codes", ZipCodesHandler(db))
	mux.HandleFunc("/api/covid-alerts", CovidAlertsHandler(db))
	mux.HandleFunc("/api/airport-traffic", AirportTrafficHandler(db))
	mux.HandleFunc("/api/high-ccvi", HighCCVIHandler(db))
	mux.HandleFunc("/api/traffic-patterns", TrafficPatternsHandler(db))
	mux.HandleFunc("/api/forecast", ForecastHandler(db))
	mux.HandleFunc("/api/investment-targets", InvestmentTargetsHandler(db))
	mux.HandleFunc("/api/small-business-loans", SmallBusinessLoansHandler(db))
	mux.HandleFunc("/api/building-permits", BuildingPermitsHandler(db))

	// Apply middleware
	return Chain(mux, CORS, JSONContentType, Logger)
}

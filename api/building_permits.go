package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type BuildingPermit struct {
	ID                   string    `json:"id"`
	PermitNumber         string    `json:"permit_number"`
	PermitType           string    `json:"permit_type"`
	ReviewType           string    `json:"review_type"`
	ApplicationStartDate time.Time `json:"application_start_date,omitempty"`
	IssueDate            time.Time `json:"issue_date,omitempty"`
	StreetAddress        string    `json:"street_address"`
	WorkDescription      string    `json:"work_description"`
	SubtotalPaid         float64   `json:"subtotal_paid"`
	ZipCode              string    `json:"zip_code"`
}

// BuildingPermitsHandler returns building permit data for a specific zip code
func BuildingPermitsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zip")
		if zipCode == "" {
			http.Error(w, "zip parameter is required", http.StatusBadRequest)
			return
		}

		query := `
			SELECT
				id,
				COALESCE(permit_number, '') as permit_number,
				COALESCE(permit_type, '') as permit_type,
				COALESCE(review_type, '') as review_type,
				application_start_date,
				issue_date,
				CONCAT(
					COALESCE(street_number, ''),
					' ',
					COALESCE(street_direction, ''),
					' ',
					COALESCE(street_name, ''),
					' ',
					COALESCE(suffix, '')
				) as street_address,
				COALESCE(work_description, '') as work_description,
				COALESCE(subtotal_paid, 0) as subtotal_paid,
				COALESCE(zip_code, '') as zip_code
			FROM building_permits
			WHERE zip_code = $1
			ORDER BY issue_date DESC NULLS LAST
			LIMIT 100
		`

		rows, err := db.Query(query, zipCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var permits []BuildingPermit
		for rows.Next() {
			var p BuildingPermit
			var applicationStartDate, issueDate sql.NullTime
			var subtotalPaid sql.NullFloat64

			err := rows.Scan(
				&p.ID,
				&p.PermitNumber,
				&p.PermitType,
				&p.ReviewType,
				&applicationStartDate,
				&issueDate,
				&p.StreetAddress,
				&p.WorkDescription,
				&subtotalPaid,
				&p.ZipCode,
			)
			if err != nil {
				continue
			}

			if applicationStartDate.Valid {
				p.ApplicationStartDate = applicationStartDate.Time
			}
			if issueDate.Valid {
				p.IssueDate = issueDate.Time
			}
			if subtotalPaid.Valid {
				p.SubtotalPaid = subtotalPaid.Float64
			}

			permits = append(permits, p)
		}

		json.NewEncoder(w).Encode(permits)
	}
}

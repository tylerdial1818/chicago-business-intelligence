package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type SmallBusinessLoan struct {
	ZipCode                string  `json:"zip_code"`
	CommunityAreaName      string  `json:"community_area_name"`
	PerCapitaIncome        int     `json:"per_capita_income"`
	UnemploymentRate       float64 `json:"unemployment_rate"`
	PovertyRate            float64 `json:"poverty_rate"`
	NewConstructionPermits int     `json:"new_construction_permits"`
	MaxLoanAmount          int     `json:"max_loan_amount"`
}

// SmallBusinessLoansHandler returns eligible zip codes for small business loans
func SmallBusinessLoansHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			WITH new_construction_counts AS (
				SELECT
					COALESCE(bp.zip_code, zcm.zip_code) as zip_code,
					COUNT(*) as permit_count
				FROM building_permits bp
				LEFT JOIN zip_community_map zcm ON bp.community_area = zcm.community_area
				WHERE bp.permit_type LIKE '%NEW CONSTRUCTION%'
				GROUP BY COALESCE(bp.zip_code, zcm.zip_code)
			),
			eligible_areas AS (
				SELECT
					zcm.zip_code,
					u.community_area_name,
					u.per_capita_income,
					u.percent_aged_16_unemployed,
					u.percent_below_poverty
				FROM unemployment u
				JOIN zip_community_map zcm ON u.community_area = zcm.community_area
				WHERE u.per_capita_income < 30000
			)
			SELECT
				ea.zip_code,
				ea.community_area_name,
				ea.per_capita_income,
				ea.percent_aged_16_unemployed,
				ea.percent_below_poverty,
				COALESCE(nc.permit_count, 0) as new_construction_permits,
				250000 as max_loan_amount
			FROM eligible_areas ea
			LEFT JOIN new_construction_counts nc ON ea.zip_code = nc.zip_code
			ORDER BY COALESCE(nc.permit_count, 0) ASC, ea.per_capita_income ASC
		`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var loans []SmallBusinessLoan
		for rows.Next() {
			var loan SmallBusinessLoan
			var unemploymentRate, povertyRate sql.NullFloat64
			var perCapitaIncome sql.NullInt64

			err := rows.Scan(
				&loan.ZipCode,
				&loan.CommunityAreaName,
				&perCapitaIncome,
				&unemploymentRate,
				&povertyRate,
				&loan.NewConstructionPermits,
				&loan.MaxLoanAmount,
			)
			if err != nil {
				continue
			}

			if perCapitaIncome.Valid {
				loan.PerCapitaIncome = int(perCapitaIncome.Int64)
			}
			if unemploymentRate.Valid {
				loan.UnemploymentRate = unemploymentRate.Float64
			}
			if povertyRate.Valid {
				loan.PovertyRate = povertyRate.Float64
			}

			loans = append(loans, loan)
		}

		json.NewEncoder(w).Encode(loans)
	}
}

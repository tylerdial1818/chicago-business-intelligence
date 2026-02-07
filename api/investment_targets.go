package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type InvestmentTarget struct {
	CommunityArea          int     `json:"community_area"`
	CommunityAreaName      string  `json:"community_area_name"`
	UnemploymentRate       float64 `json:"unemployment_rate"`
	PovertyRate            float64 `json:"poverty_rate"`
	PerCapitaIncome        int     `json:"per_capita_income"`
	HardshipIndex          float64 `json:"hardship_index"`
	TotalPermits           int     `json:"total_permits"`
	NewConstructionPermits int     `json:"new_construction_permits"`
}

// InvestmentTargetsHandler returns top 5 neighborhoods for investment based on unemployment and poverty
func InvestmentTargetsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT
				u.community_area,
				u.community_area_name,
				u.percent_aged_16_unemployed as unemployment_rate,
				u.percent_below_poverty as poverty_rate,
				u.per_capita_income,
				u.hardship_index,
				COUNT(bp.id) as total_permits,
				COUNT(CASE WHEN bp.permit_type LIKE '%NEW CONSTRUCTION%' THEN 1 END) as new_construction_permits
			FROM unemployment u
			LEFT JOIN zip_community_map zcm ON u.community_area = zcm.community_area
			LEFT JOIN building_permits bp ON bp.community_area = u.community_area
			GROUP BY u.community_area, u.community_area_name, u.percent_aged_16_unemployed,
					 u.percent_below_poverty, u.per_capita_income, u.hardship_index
			ORDER BY u.percent_aged_16_unemployed DESC, u.percent_below_poverty DESC
			LIMIT 5
		`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var targets []InvestmentTarget
		for rows.Next() {
			var t InvestmentTarget
			var unemploymentRate, povertyRate, hardshipIndex sql.NullFloat64
			var perCapitaIncome sql.NullInt64

			err := rows.Scan(
				&t.CommunityArea,
				&t.CommunityAreaName,
				&unemploymentRate,
				&povertyRate,
				&perCapitaIncome,
				&hardshipIndex,
				&t.TotalPermits,
				&t.NewConstructionPermits,
			)
			if err != nil {
				continue
			}

			if unemploymentRate.Valid {
				t.UnemploymentRate = unemploymentRate.Float64
			}
			if povertyRate.Valid {
				t.PovertyRate = povertyRate.Float64
			}
			if perCapitaIncome.Valid {
				t.PerCapitaIncome = int(perCapitaIncome.Int64)
			}
			if hardshipIndex.Valid {
				t.HardshipIndex = hardshipIndex.Float64
			}

			targets = append(targets, t)
		}

		json.NewEncoder(w).Encode(targets)
	}
}

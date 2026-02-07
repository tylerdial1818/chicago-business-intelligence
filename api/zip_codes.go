package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type ZipCode struct {
	ZipCode        string `json:"zip_code"`
	Neighborhood   string `json:"neighborhood"`
	CommunityArea  string `json:"community_area"`
}

// ZipCodesHandler returns a list of all zip codes with data available
func ZipCodesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT DISTINCT zip_code, neighborhood_name, community_area_name
			FROM zip_community_map
			ORDER BY zip_code
		`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var zipCodes []ZipCode
		for rows.Next() {
			var zc ZipCode
			err := rows.Scan(&zc.ZipCode, &zc.Neighborhood, &zc.CommunityArea)
			if err != nil {
				continue
			}
			zipCodes = append(zipCodes, zc)
		}

		json.NewEncoder(w).Encode(zipCodes)
	}
}

package repostitories

import (
	"cyber-go/internal/models"
	"database/sql"
)

func GetAllParadigms(db *sql.DB) ([]models.Paradigm, error) {
	rows, err := db.Query("SELECT id, name, description FROM paradigms")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var paradigms []models.Paradigm

	for rows.Next() {
		var p models.Paradigm
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		paradigms = append(paradigms, p)
	}
	return paradigms, nil
}

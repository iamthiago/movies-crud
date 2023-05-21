package repository

import (
	"database/sql"
	"fmt"

	"github.com/iamthiago/movies-crud/internal/models"
)

var GetMovies = func(db *sql.DB) ([]models.Movie, error) {
	var movies []models.Movie

	rows, err := db.Query("select m.id, m.isbn, m.title, d.first_name, d.last_name from movies m, directors d where m.director_id = d.id")
	if err != nil {
		return nil, fmt.Errorf("getMovies %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Movie
		var d models.Director

		if err := rows.Scan(&m.ID, &m.Isbn, &m.Title, &d.Firstname, &d.LastName); err != nil {
			return nil, fmt.Errorf("getMovies %v", err)
		}
		m.Director = &d
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getMovies %v", err)
	}

	return movies, nil
}

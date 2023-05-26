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

var GetMovieById = func(db *sql.DB, id int64) (models.Movie, error) {
	var movie models.Movie
	var director models.Director

	row := db.QueryRow("select m.id, m.isbn, m.title, d.first_name, d.last_name from movies m, directors d where m.director_id = d.id and m.id = ?", id)
	if err := row.Scan(&movie.ID, &movie.Isbn, &movie.Title, &director.Firstname, &director.LastName); err != nil {
		if err == sql.ErrNoRows {
			return movie, fmt.Errorf("getMovieById %d: no such movie", id)
		}
		return movie, fmt.Errorf("getMovieById %d: %v", id, err)
	}
	movie.Director = &director
	return movie, nil
}

var CreateMovie = func(db *sql.DB, movie models.Movie) (models.Movie, error) {
	directorResult, dirErr := db.Exec("insert into directors (first_name, last_name) values (?, ?)", movie.Director.Firstname, movie.Director.LastName)

	if dirErr != nil {
		return models.Movie{}, fmt.Errorf("add directors: %v", dirErr)
	}
	directorsId, dirLastInsertErr := directorResult.LastInsertId()
	if dirLastInsertErr != nil {
		return models.Movie{}, fmt.Errorf("get directors last inserted id %v", dirLastInsertErr)
	}

	movieResult, movErr := db.Exec("insert into movies (isbn, title, director_id) values (?, ?, ?)", movie.Isbn, movie.Title, directorsId)
	if movErr != nil {
		return models.Movie{}, fmt.Errorf("add movies: %v", movErr)
	}
	movieId, movLAstInsertErr := movieResult.LastInsertId()
	if movLAstInsertErr != nil {
		return models.Movie{}, fmt.Errorf("get movie last inserted id %v", movLAstInsertErr)
	}

	movie.ID = movieId
	return movie, nil
}

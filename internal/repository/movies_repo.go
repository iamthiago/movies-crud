package repository

import (
	"database/sql"
	"fmt"

	"github.com/iamthiago/movies-crud/internal/models"
)

var GetMovies = func(db *sql.DB) ([]models.Movie, error) {
	var movies []models.Movie

	rows, err := db.Query("select id, isbn, title, director from movies")
	if err != nil {
		return nil, fmt.Errorf("getMovies %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Movie

		if err := rows.Scan(&m.ID, &m.Isbn, &m.Title, &m.Director); err != nil {
			return nil, fmt.Errorf("getMovies %v", err)
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getMovies %v", err)
	}

	return movies, nil
}

var GetMovieById = func(db *sql.DB, id int64) (models.Movie, error) {
	var movie models.Movie

	row := db.QueryRow("select id, isbn, title, director from movies where id = ?", id)
	if err := row.Scan(&movie.ID, &movie.Isbn, &movie.Title, &movie.Director); err != nil {
		if err == sql.ErrNoRows {
			return movie, fmt.Errorf("getMovieById %d: no such movie", id)
		}
		return movie, fmt.Errorf("getMovieById %d: %v", id, err)
	}
	return movie, nil
}

var CreateMovie = func(db *sql.DB, movie models.Movie) (models.Movie, error) {
	movieResult, movErr := db.Exec("insert into movies (isbn, title, director) values (?, ?, ?)", movie.Isbn, movie.Title, movie.Director)
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

var UpdateMovie = func(db *sql.DB, id int64, movie models.Movie) (models.Movie, error) {
	_, movErr := db.Exec("update movies set isbn = ?, title = ?, director = ? where id = ?", movie.Isbn, movie.Title, movie.Director, id)
	if movErr != nil {
		return models.Movie{}, fmt.Errorf("update movies: %v", movErr)
	}

	movie.ID = id

	return movie, nil
}

var DeleteMovie = func(db *sql.DB, id int64) error {
	_, err := db.Exec("delete from movies where id = ?", id)
	if err != nil {
		return fmt.Errorf("delete movies: %v", err)
	}
	return nil
}

package repository

import (
	"database/sql"
	"fmt"

	"github.com/iamthiago/movies-crud/pkg/models"
)

type MoviesRepository interface {
	GetMovies() ([]models.Movie, error)
	GetMovieById(id int64) (*models.Movie, error)
	CreateMovie(movie *models.Movie) (*models.Movie, error)
	UpdateMovie(id int64, movie *models.Movie) (*models.Movie, error)
	DeleteMovie(id int64) error
}

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetMovies() ([]models.Movie, error) {
	var movies []models.Movie

	rows, err := r.DB.Query("select id, isbn, title, director from movies")
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

func (r *Repository) GetMovieById(id int64) (*models.Movie, error) {
	var movie models.Movie

	row := r.DB.QueryRow("select id, isbn, title, director from movies where id = ?", id)
	if err := row.Scan(&movie.ID, &movie.Isbn, &movie.Title, &movie.Director); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("getMovieById %d: no such movie", id)
		}
		return nil, fmt.Errorf("getMovieById %d: %v", id, err)
	}
	return &movie, nil
}

func (r *Repository) CreateMovie(movie *models.Movie) (*models.Movie, error) {
	movieResult, movErr := r.DB.Exec("insert into movies (isbn, title, director) values (?, ?, ?)", movie.Isbn, movie.Title, movie.Director)
	if movErr != nil {
		return nil, fmt.Errorf("add movies: %v", movErr)
	}
	movieId, movLAstInsertErr := movieResult.LastInsertId()
	if movLAstInsertErr != nil {
		return nil, fmt.Errorf("get movie last inserted id %v", movLAstInsertErr)
	}

	movie.ID = movieId
	return movie, nil
}

func (r *Repository) UpdateMovie(id int64, movie *models.Movie) (*models.Movie, error) {
	_, movErr := r.DB.Exec("update movies set isbn = ?, title = ?, director = ? where id = ?", movie.Isbn, movie.Title, movie.Director, id)
	if movErr != nil {
		return nil, fmt.Errorf("update movies: %v", movErr)
	}

	movie.ID = id

	return movie, nil
}

func (r *Repository) DeleteMovie(id int64) error {
	_, err := r.DB.Exec("delete from movies where id = ?", id)
	if err != nil {
		return fmt.Errorf("delete movies: %v", err)
	}
	return nil
}

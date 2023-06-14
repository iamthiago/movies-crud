package service

import (
	"github.com/iamthiago/movies-crud/internal/movies-crud/repository"
	"github.com/iamthiago/movies-crud/pkg/models"
)

type MoviesService interface {
	GetMovies() ([]models.Movie, error)
	GetMovieById(id int64) (models.Movie, error)
	CreateMovie(movie models.Movie) (models.Movie, error)
	UpdateMovie(id int64, movie models.Movie) (models.Movie, error)
	DeleteMovie(id int64) error
}

type Service struct {
	Repository *repository.Repository
}

func (s Service) GetMovies() ([]models.Movie, error) {
	return s.Repository.GetMovies()
}

func (s Service) GetMovieById(id int64) (models.Movie, error) {
	return s.Repository.GetMovieById(id)
}

func (s Service) CreateMovie(movie models.Movie) (models.Movie, error) {
	return s.Repository.CreateMovie(movie)
}

func (s Service) UpdateMovie(id int64, movie models.Movie) (models.Movie, error) {
	return s.Repository.UpdateMovie(id, movie)
}

func (s Service) DeleteMovie(id int64) error {
	return s.Repository.DeleteMovie(id)
}

package service

import (
	"fmt"
	"log"

	"github.com/iamthiago/movies-crud/internal/movies/events"
	"github.com/iamthiago/movies-crud/internal/movies/producer"
	"github.com/iamthiago/movies-crud/internal/movies/repository"
	"github.com/iamthiago/movies-crud/pkg/models"
	"google.golang.org/protobuf/proto"
)

type MoviesService interface {
	GetMovies() ([]models.Movie, error)
	GetMovieById(id int64) (*models.Movie, error)
	CreateMovie(movie *models.Movie) (*models.Movie, error)
	UpdateMovie(id int64, movie *models.Movie) (*models.Movie, error)
	DeleteMovie(id int64) error
}

type Service struct {
	Repository    repository.MoviesRepository
	KafkaProducer producer.KafkaProducer
}

func (s *Service) GetMovies() ([]models.Movie, error) {
	return s.Repository.GetMovies()
}

func (s *Service) GetMovieById(id int64) (*models.Movie, error) {
	return s.Repository.GetMovieById(id)
}

func (s *Service) CreateMovie(movie *models.Movie) (*models.Movie, error) {
	m, err := s.Repository.CreateMovie(movie)
	if err != nil {
		return nil, fmt.Errorf("error when creating movie %v", err)
	}

	eventBytes, protoErr := toProtoEvent(m)
	if protoErr != nil {
		log.Fatalln("Failed to encode movie event", protoErr)
	}

	s.KafkaProducer.SendMovieEvent(eventBytes)

	return m, err
}

func (s *Service) UpdateMovie(id int64, movie *models.Movie) (*models.Movie, error) {
	return s.Repository.UpdateMovie(id, movie)
}

func (s *Service) DeleteMovie(id int64) error {
	return s.Repository.DeleteMovie(id)
}

func toProtoEvent(movie *models.Movie) (eventBytes []byte, err error) {
	event := events.MovieEvent{
		Id:       movie.ID,
		Isbn:     movie.Isbn,
		Title:    movie.Title,
		Director: movie.Director,
	}

	eventBytes, err = proto.Marshal(&event)

	return
}

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/iamthiago/movies-crud/internal/movies-crud/service"
	"github.com/iamthiago/movies-crud/pkg/models"
)

func GetMovies(w http.ResponseWriter, r *http.Request, service service.Service) {
	w.Header().Set("Content-Type", "application/json")

	movies, err := service.GetMovies()
	if err != nil {
		fmt.Println("Error fetching movies", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request, service service.Service) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	movie, err := service.GetMovieById(id)
	if err != nil {
		if movie.IsEmpty() {
			fmt.Println("Movie is empty", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Println("Error fetching movie by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movie)
}

func CreateMovie(w http.ResponseWriter, r *http.Request, service service.Service) {
	w.Header().Set("Content-Type", "application/json")
	var movie models.Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movieWithId, err := service.CreateMovie(movie)
	if err != nil {
		fmt.Println("Error creating movie", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movieWithId)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request, service service.Service) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	var movie models.Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	updatedMovie, dbErr := service.UpdateMovie(id, movie)
	if dbErr != nil {
		fmt.Println("Error updating movie", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedMovie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request, service service.Service) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	dbErr := service.DeleteMovie(id)
	if dbErr != nil {
		fmt.Println("Error deleting movie", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

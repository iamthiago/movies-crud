package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iamthiago/movies-crud/configs"
	"github.com/iamthiago/movies-crud/internal/models"
	"github.com/iamthiago/movies-crud/internal/repository"
)

func main() {
	db, err := configs.GetMySQLDB()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to Mysql!")

	r := mux.NewRouter()

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		getMovies(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		getMovie(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		createMovie(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateMovie(w, r, db)
	}).Methods("PUT")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteMovie(w, r, db)
	}).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getMovies(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	movies, err := repository.GetMovies(db)
	if err != nil {
		fmt.Println("Error fetching movies", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	movie, err := repository.GetMovieById(db, id)
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

func createMovie(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	var movie models.Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movieWithId, err := repository.CreateMovie(db, movie)
	if err != nil {
		fmt.Println("Error creating movie", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movieWithId)
}

func updateMovie(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	var movie models.Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	updatedMovie, dbErr := repository.UpdateMovie(db, id, movie)
	if dbErr != nil {
		fmt.Println("Error updating movie", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedMovie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	dbErr := repository.DeleteMovie(db, id)
	if dbErr != nil {
		fmt.Println("Error deleting movie", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

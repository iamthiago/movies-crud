package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/iamthiago/movies-crud/internal/models"
	"github.com/iamthiago/movies-crud/internal/repository"
)

type Movie struct {
	ID       int64     `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var moviesMap = make(map[int64]Movie)

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "movies",
		AllowNativePasswords: true,
	}

	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
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

	moviesMap[1] = Movie{ID: 1, Isbn: "123", Title: "Movie 1", Director: &Director{Firstname: "John", LastName: "Doe"}}
	moviesMap[2] = Movie{ID: 2, Isbn: "456", Title: "Movie 2", Director: &Director{Firstname: "Alex", LastName: "Smith"}}
	moviesMap[3] = Movie{ID: 3, Isbn: "789", Title: "Movie 3", Director: &Director{Firstname: "Jane", LastName: "Black"}}

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

	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

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

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	if movie, exists := moviesMap[id]; exists {
		delete(moviesMap, movie.ID)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

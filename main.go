package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getMoviesFromDB(db *sql.DB) ([]Movie, error) {
	var movies []Movie

	rows, err := db.Query("select m.id, m.isbn, m.title, d.first_name, d.last_name from movies m, directors d where m.director_id = d.id")
	if err != nil {
		return nil, fmt.Errorf("getMovies %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m Movie
		var d Director

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

func getMovies(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	movies, err := getMoviesFromDB(db)
	if err != nil {
		fmt.Println("Error fetching movies", err)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	if movie, exists := moviesMap[id]; exists {
		json.NewEncoder(w).Encode(movie)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = rand.Int63n(1000000)
	moviesMap[movie.ID] = movie
	json.NewEncoder(w).Encode(movie)
}

// it can work as an upsert
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = id
	moviesMap[id] = movie
	json.NewEncoder(w).Encode(movie)
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

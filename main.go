package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var moviesMap = make(map[int]Movie)
var movies = make([]Movie, 0, len(moviesMap))

func main() {
	r := mux.NewRouter()

	moviesMap[1] = Movie{ID: "1", Isbn: "123", Title: "Movie 1", Director: &Director{Firstname: "John", LastName: "Doe"}}
	moviesMap[2] = Movie{ID: "2", Isbn: "456", Title: "Movie 2", Director: &Director{Firstname: "Alex", LastName: "Smith"}}
	moviesMap[3] = Movie{ID: "3", Isbn: "789", Title: "Movie 3", Director: &Director{Firstname: "Jane", LastName: "Black"}}

	for _, m := range moviesMap {
		movies = append(movies, m)
	}

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	json.NewEncoder(w).Encode(moviesMap[id])
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamthiago/movies-crud/configs"
	"github.com/iamthiago/movies-crud/internal/movies-crud/rest"
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
		rest.GetMovies(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.GetMovie(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		rest.CreateMovie(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UpdateMovie(w, r, db)
	}).Methods("PUT")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.DeleteMovie(w, r, db)
	}).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

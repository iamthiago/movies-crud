package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamthiago/movies-crud/internal/movies/controller"
	"github.com/iamthiago/movies-crud/internal/movies/mysql"
	"github.com/iamthiago/movies-crud/internal/movies/producer"
	"github.com/iamthiago/movies-crud/internal/movies/repository"
	"github.com/iamthiago/movies-crud/internal/movies/service"
)

func main() {
	db, err := mysql.GetMySQLDB()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	topic := "movies"
	kafkaProducer, err := producer.GetKafkaProducer(&topic)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer db.Close()
	defer kafkaProducer.Producer.Close()

	movieRepo := repository.Repository{DB: db}
	movieService := service.Service{Repository: &movieRepo, KafkaProducer: &kafkaProducer}

	r := mux.NewRouter()

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		controller.GetMovies(w, r, &movieService)
	}).Methods("GET")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		controller.GetMovie(w, r, &movieService)
	}).Methods("GET")

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		controller.CreateMovie(w, r, &movieService)
	}).Methods("POST")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		controller.UpdateMovie(w, r, &movieService)
	}).Methods("PUT")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		controller.DeleteMovie(w, r, &movieService)
	}).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

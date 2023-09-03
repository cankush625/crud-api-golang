package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Isbn     string    `json:"isbn"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

func listMovies(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		return
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			return
		}
	}
	http.Error(w, "Movie Not Found", http.StatusNotFound)
	return
}

func retrieveMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				return
			}
			return
		}
	}
	http.Error(w, "Movie Not Found", http.StatusNotFound)
	return
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000))
	movies = append(movies, movie)
	err := json.NewEncoder(w).Encode(movie)
	if err != nil {
		return
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			_ = json.NewDecoder(r.Body).Decode(&item)
			movies[index] = item
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				http.Error(w, "Failed to update the movie", http.StatusBadRequest)
				return
			}
			return
		}
	}
	http.Error(w, "Movie Not Found", http.StatusNotFound)
	return
}

func main() {
	r := mux.NewRouter()

	movies = append(
		movies,
		Movie{ID: "1", Name: "Movie One", Isbn: "347328", Director: &Director{FirstName: "John", LastName: "Doe"}},
		Movie{ID: "2", Name: "Movie Two", Isbn: "437783", Director: &Director{FirstName: "Adam", LastName: "Smith"}},
	)

	r.HandleFunc("/movies", listMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", retrieveMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT", "PATCH")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

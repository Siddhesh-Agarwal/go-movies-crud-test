package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	ID    int
	Title string
	Year  int
	Genre string
}

func to_int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func getDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./temp.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movies []Movie
	db := getDB()
	db.Find(&movies)
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	params := mux.Vars(r)
	movie_id := to_int(params["id"])
	db := getDB()
	db.Delete(&movie, movie_id)
	json.NewEncoder(w).Encode(movie)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	params := mux.Vars(r)
	movie_id := to_int(params["id"])
	db := getDB()
	db.First(&movie, movie_id)
	json.NewEncoder(w).Encode(movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	db := getDB()
	db.Create(&movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	movie_id := to_int(params["id"])
	db := getDB()
	db.First(&movie, movie_id)
	if db.RowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&movie)
	db.Save(&movie)
	json.NewEncoder(w).Encode(movie)
}

func main() {
	// Connect to the database
	db, err := gorm.Open(sqlite.Open("./temp.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Auto migrate the database to match the struct definitions
	db.AutoMigrate(&Movie{})
	var movies []Movie
	// Add some movies to the database if they don't exist
	if db.Find(&movies).RowsAffected == 0 {
		db.Create(&Movie{ID: 1, Title: "The Shawshank Redemption", Year: 1994, Genre: "Drama"})
		db.Create(&Movie{ID: 2, Title: "The Godfather", Year: 1972, Genre: "Crime"})
		db.Create(&Movie{ID: 3, Title: "The Dark Knight", Year: 2008, Genre: "Action"})
	}

	port, port_exists := os.LookupEnv("PORT")
	if !port_exists {
		port = "8000"
	}
	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Printf("Starting server at port %s\n", port)
	fmt.Printf("Local: http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}

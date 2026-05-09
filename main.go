package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/pets") //user:password@/dbname
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Database connected")

	r := chi.NewRouter()
	r.Get("/pets/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		query := "SELECT id, name, age, kind, breed FROM pets WHERE id = ?"
		pet := &Pet{}
		err := db.QueryRow(query, id).Scan(&pet.ID, &pet.Name, &pet.Age, &pet.Kind, &pet.Breed)
		if err != nil {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pet)
	})

	r.Post("/pets", func(w http.ResponseWriter, r *http.Request) {
		var p CreateRequest
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		query := "INSERT INTO pets (name, age, kind, breed) VALUES (?, ?, ?, ?)"
		if strings.TrimSpace(p.Kind) == "" || strings.TrimSpace(p.Breed) == "" {
			http.Error(w, "Kind and Breed are required fields", http.StatusBadRequest)
			return
		}
		if p.Age.Valid && p.Age.Int64 < 0 {
			http.Error(w, "Age cannot be negative", http.StatusBadRequest)
			return
		}
		result, err := db.Exec(query, p.Name, p.Age, p.Kind, p.Breed)
		if err != nil {
			http.Error(w, "Error creating pet", http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		pet := &Pet{
			ID:    int(id),
			Name:  p.Name,
			Age:   p.Age,
			Kind:  p.Kind,
			Breed: p.Breed,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pet)
	})

	http.ListenAndServe(":8080", r)
}

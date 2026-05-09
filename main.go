package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//Connect to database
	db, err := InitDB("root:@tcp(localhost:3306)/pets")
	if err != nil {
		panic(fmt.Errorf("Failed to initialise database: %v", err))
	}
	defer db.Close()
	//

	r := chi.NewRouter()
	handler := NewHandler(NewPetRepository(db))

	r.Get("/pets/{id}", handler.GetByID)
	r.Post("/pets", handler.CreatePet)

	http.ListenAndServe(":8080", r)
}

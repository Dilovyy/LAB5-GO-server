package main

import (
	"database/sql"
	"fmt"
	"net/http"

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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var greeting string
		err := db.QueryRow("SELECT 'Hello, World!'").Scan(&greeting)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(greeting))
	})
	http.ListenAndServe(":8080", r)
}

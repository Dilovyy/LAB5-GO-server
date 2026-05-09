package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res := []byte("Hello, World!")
		w.Write(res)
	})
	http.ListenAndServe(":8080", r)
}

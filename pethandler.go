package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	pet, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pet)
}

func (h *Handler) CreatePet(w http.ResponseWriter, r *http.Request) {
	var p CreateRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(p.Kind) == "" || strings.TrimSpace(p.Breed) == "" {
		http.Error(w, "Kind and Breed are required fields", http.StatusBadRequest)
		return
	}
	if p.Age.Valid && p.Age.Int64 < 0 {
		http.Error(w, "Age cannot be negative", http.StatusBadRequest)
		return
	}
	result, err := h.repo.CreatePet(&p)
	if err != nil {
		http.Error(w, "Error creating pet", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

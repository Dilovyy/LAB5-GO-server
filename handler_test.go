package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockPetRepository struct {
	pets map[int]*Pet
}

func NewMockPetRepository() *mockPetRepository {
	return &mockPetRepository{
		pets: make(map[int]*Pet),
	}
}

func URLbyID(r *http.Request, id string) *http.Request { //creating url for tests
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func toJSON(v any) *bytes.Buffer { // text to json
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func (m *mockPetRepository) CreatePet(p *CreateRequest) (*Pet, error) {
	pet := &Pet{
		Name:  p.Name,
		Age:   p.Age,
		Kind:  p.Kind,
		Breed: p.Breed,
	}
	return pet, nil
}

func (m *mockPetRepository) GetByID(id int) (*Pet, error) {
	pet, ok := m.pets[id]
	if !ok {
		return nil, fmt.Errorf("Pet not found")
	}
	return pet, nil
}

func (m *mockPetRepository) GetList() ([]*Pet, error) {
	var list []*Pet
	for _, pet := range m.pets {
		list = append(list, pet)
	}
	return list, nil
}

func (m *mockPetRepository) UpdatePet(id int, p *UpdateRequest) (*Pet, error) {
	if _, ok := m.pets[id]; !ok {
		return nil, fmt.Errorf("pet not found")
	}
	pet := &Pet{ID: id, Name: p.Name, Age: p.Age, Kind: p.Kind, Breed: p.Breed}
	m.pets[id] = pet
	return pet, nil
}

func (m *mockPetRepository) PatchPet(id int, p *PatchRequest) (*Pet, error) {
	pet, ok := m.pets[id]
	if !ok {
		return nil, fmt.Errorf("pet not found")
	}
	if p.Name != nil {
		pet.Name = p.Name
	}
	if p.Age != nil {
		pet.Age = p.Age
	}
	if p.Kind != nil {
		pet.Kind = *p.Kind
	}
	if p.Breed != nil {
		pet.Breed = *p.Breed
	}
	return pet, nil
}

func (m *mockPetRepository) DeletePet(id int) error {
	if _, ok := m.pets[id]; !ok {
		return fmt.Errorf("pet not found")
	}
	delete(m.pets, id)
	return nil
}

// tests
// CREATE
func TestCreateHandlerInvalidJSON(t *testing.T) {

	repo := NewMockPetRepository()
	handler := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/pets", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()
	handler.CreatePet(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func intToiint(i int) *int { return &i } // helps convert int to *int for tests (because in all pet struct fields are pointers)
func TestCreate_NegativeAge(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	var a int = -1
	body := toJSON(CreateRequest{Kind: "dog", Breed: "Husky", Age: intToiint(a)})
	req := httptest.NewRequest(http.MethodPost, "/pets", body)
	rec := httptest.NewRecorder()
	h.CreatePet(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "negative") {
		t.Errorf("expected 'negative' in error, got: %s", rec.Body.String())
	}
}

func TestCreate_Valid(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	body := toJSON(CreateRequest{Kind: "dog", Breed: "Pug"})
	req := httptest.NewRequest(http.MethodPost, "/pets", body)
	rec := httptest.NewRecorder()
	h.CreatePet(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var pet Pet
	json.NewDecoder(rec.Body).Decode(&pet)
	if pet.Kind != "dog" {
		t.Errorf("expected kind=dog, got %s", pet.Kind)
	}
}

// GETLIST
func TestGetList_Empty(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	req := httptest.NewRequest(http.MethodGet, "/pets", nil)
	rec := httptest.NewRecorder()
	h.GetList(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetList_ReturnsPets(t *testing.T) {
	repo := NewMockPetRepository()
	repo.pets[1] = &Pet{ID: 1, Kind: "dog", Breed: "Husky"}
	repo.pets[2] = &Pet{ID: 2, Kind: "cat", Breed: "British"}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/pets", nil)
	rec := httptest.NewRecorder()
	h.GetList(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var pets []*Pet
	json.NewDecoder(rec.Body).Decode(&pets)
	if len(pets) != 2 {
		t.Errorf("expected 2 pets, got %d", len(pets))
	}
}

// GETBYID
func TestGetByID_Found(t *testing.T) {
	repo := NewMockPetRepository()
	repo.pets[1] = &Pet{ID: 1, Kind: "dog", Breed: "Husky"}
	h := NewHandler(repo)

	req := URLbyID(httptest.NewRequest(http.MethodGet, "/pets/1", nil), "1")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var pet Pet
	json.NewDecoder(rec.Body).Decode(&pet)
	if pet.ID != 1 {
		t.Errorf("expected ID 1, got %d", pet.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	req := URLbyID(httptest.NewRequest(http.MethodGet, "/pets/99", nil), "99")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	req := URLbyID(httptest.NewRequest(http.MethodGet, "/pets/abc", nil), "abc")
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// UPDATE
func TestUpdate_Valid(t *testing.T) {
	repo := NewMockPetRepository()
	repo.pets[1] = &Pet{ID: 1, Kind: "dog", Breed: "Husky"}
	h := NewHandler(repo)

	body := toJSON(UpdateRequest{Kind: "cat", Breed: "British"})
	req := URLbyID(httptest.NewRequest(http.MethodPut, "/pets/1", body), "1")
	rec := httptest.NewRecorder()
	h.UpdatePet(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var pet Pet
	json.NewDecoder(rec.Body).Decode(&pet)
	if pet.Kind != "cat" {
		t.Errorf("expected kind=cat, got %s", pet.Kind)
	}
}

func strToistr(i string) *string { return &i }

// PATCH
func TestPatch_BreedOnly(t *testing.T) {
	repo := NewMockPetRepository()
	repo.pets[1] = &Pet{ID: 1, Kind: "dog", Breed: "Husky"}
	h := NewHandler(repo)

	body := toJSON(PatchRequest{Breed: strToistr("Pug")})
	req := URLbyID(httptest.NewRequest(http.MethodPatch, "/pets/1", body), "1")
	rec := httptest.NewRecorder()
	h.PatchPet(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var pet Pet
	json.NewDecoder(rec.Body).Decode(&pet)
	if pet.Breed != "Pug" {
		t.Errorf("expected breed=Pug, got %s", pet.Breed)
	}
	if pet.Kind != "dog" {
		t.Errorf("kind should not change, got %s", pet.Kind)
	}
}

// DELETE
func TestDelete_Valid(t *testing.T) {
	repo := NewMockPetRepository()
	repo.pets[1] = &Pet{ID: 1, Kind: "dog", Breed: "Husky"}
	h := NewHandler(repo)

	req := URLbyID(httptest.NewRequest(http.MethodDelete, "/pets/1", nil), "1")
	rec := httptest.NewRecorder()
	h.DeletePet(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
	if _, exists := repo.pets[1]; exists {
		t.Error("pet should have been deleted")
	}
}

func TestDelete_NotFound(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	req := URLbyID(httptest.NewRequest(http.MethodDelete, "/pets/99", nil), "99")
	rec := httptest.NewRecorder()
	h.DeletePet(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

func TestDelete_InvalidID(t *testing.T) {
	h := NewHandler(NewMockPetRepository())
	req := URLbyID(httptest.NewRequest(http.MethodDelete, "/pets/abc", nil), "abc")
	rec := httptest.NewRecorder()
	h.DeletePet(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

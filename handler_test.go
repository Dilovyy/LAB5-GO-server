package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockPetRepository struct {
	pets map[int]*Pet
}

func NewMockPetRepository() *mockPetRepository {
	return &mockPetRepository{
		pets: make(map[int]*Pet),
	}
}

func (m *mockPetRepository) CreatePet(p *CreateRequest) (*Pet, error) {
	pet := &Pet{}
	m.pets[pet.ID] = pet
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
	return nil, nil
}

func (m *mockPetRepository) PatchPet(id int, p *PatchRequest) (*Pet, error) {
	return nil, nil
}

func (m *mockPetRepository) DeletePet(id int) error {
	return nil
}

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

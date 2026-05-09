package main

import "database/sql"

type Repository interface {
	GetByID(id int) (*Pet, error)
	CreatePet(p *CreateRequest) (*Pet, error)
}

type petRepository struct {
	db *sql.DB
}

func NewPetRepository(db *sql.DB) Repository {
	return &petRepository{db: db}
}

func (r *petRepository) GetByID(id int) (*Pet, error) {
	query := "SELECT id, name, age, kind, breed FROM pets WHERE id = ?"
	pet := &Pet{}
	err := r.db.QueryRow(query, id).Scan(&pet.ID, &pet.Name, &pet.Age, &pet.Kind, &pet.Breed)
	if err != nil {
		return nil, err
	}
	return pet, nil
}

func (r *petRepository) CreatePet(p *CreateRequest) (*Pet, error) {
	query := "INSERT INTO pets (name, age, kind, breed) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, p.Name, p.Age, p.Kind, p.Breed)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	pet := &Pet{
		ID:    int(id),
		Name:  p.Name,
		Age:   p.Age,
		Kind:  p.Kind,
		Breed: p.Breed,
	}
	return pet, nil
}

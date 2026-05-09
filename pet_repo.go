package main

import (
	"database/sql"
)

type Repository interface {
	GetList() ([]*Pet, error)
	GetByID(id int) (*Pet, error)
	CreatePet(p *CreateRequest) (*Pet, error)
	UpdatePet(id int, p *UpdateRequest) (*Pet, error)
	PatchPet(id int, p *PatchRequest) (*Pet, error)
	DeletePet(id int) error
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

func (r *petRepository) GetList() ([]*Pet, error) {
	query := "SELECT id, name, age, kind, breed FROM pets"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []*Pet
	for rows.Next() {
		pet := &Pet{}
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Age, &pet.Kind, &pet.Breed); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, rows.Err()
}

func (r *petRepository) UpdatePet(id int, p *UpdateRequest) (*Pet, error) {
	query := "UPDATE pets SET name = ?, age = ?, kind = ?, breed = ? WHERE id = ?"
	_, err := r.db.Exec(query, p.Name, p.Age, p.Kind, p.Breed, id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *petRepository) PatchPet(id int, p *PatchRequest) (*Pet, error) {
	pet, err := r.GetByID(id)
	if err != nil {
		return nil, err
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
	query := "UPDATE pets SET name = ?, age = ?, kind = ?, breed = ? WHERE id = ?"
	_, err = r.db.Exec(query, pet.Name, pet.Age, pet.Kind, pet.Breed, id)
	if err != nil {
		return nil, err
	}
	return pet, nil
}

func (r *petRepository) DeletePet(id int) error {
	query := "DELETE FROM pets WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

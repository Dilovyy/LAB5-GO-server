package main

import "database/sql"

type Pet struct {
	ID    int            `json:"id"`
	Name  sql.NullString `json:"name"`
	Age   sql.NullInt64  `json:"age"`
	Kind  string         `json:"kind"`
	Breed string         `json:"breed"`
}

type CreateRequest struct {
	Name  sql.NullString `json:"name"`
	Age   sql.NullInt64  `json:"age"`
	Kind  string         `json:"kind"`
	Breed string         `json:"breed"`
}

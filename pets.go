package main

type Pet struct {
	ID    int     `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Kind  string  `json:"kind"`
	Breed string  `json:"breed"`
}

type CreateRequest struct {
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Kind  string  `json:"kind"`
	Breed string  `json:"breed"`
}

type UpdateRequest struct {
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Kind  string  `json:"kind"`
	Breed string  `json:"breed"`
}

type PatchRequest struct {
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Kind  *string `json:"kind"`
	Breed *string `json:"breed"`
}

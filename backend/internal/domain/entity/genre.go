package entity

import "github.com/google/uuid"

type Genre struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

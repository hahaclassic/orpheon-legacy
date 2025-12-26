package entity

import "github.com/google/uuid"

type SearchRequest struct {
	Query   string
	Filters Filters
	Limit   int
	Offset  int
}

type Filters struct {
	GenreID uuid.UUID
	Country string
}

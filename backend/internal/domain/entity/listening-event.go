package entity

import "github.com/google/uuid"

type ListeningEvent struct {
	TrackID uuid.UUID `json:"track_id"`
	UserID  uuid.UUID `json:"user_id"`
	Ranges  []*Range  `json:"ranges"` // e.g. [[2, 39], [55, 141]] - listened from 2 to 39 seconds, then 55 to 141
}

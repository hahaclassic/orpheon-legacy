package dto

import "github.com/google/uuid"

type TrackAdditionRequest struct {
	TrackID uuid.UUID `json:"track_id"`
}

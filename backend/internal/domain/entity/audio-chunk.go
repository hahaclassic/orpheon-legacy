package entity

import "github.com/google/uuid"

type AudioChunk struct {
	Data    []byte    `json:"data"`
	TrackID uuid.UUID `json:"track_id"`
	Start   int64     `json:"start"`
	End     int64     `json:"end"`
}

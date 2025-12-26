package entity

import "github.com/google/uuid"

type PlaylistTrack struct {
	PlaylistID uuid.UUID `json:"playlist_id"`
	TrackID    uuid.UUID `json:"track_id"`
	Position   int       `json:"position"`
}

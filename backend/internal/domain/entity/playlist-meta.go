package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlaylistMeta struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Rating      int       `json:"rating"`
}

type PlaylistAccessMeta struct {
	OwnerID   uuid.UUID `json:"owner_id"`
	IsPrivate bool      `json:"is_private"`
}

type PlaylistMetaAggregated struct {
	ID          uuid.UUID    `json:"id"`
	Owner       *UserInfo    `json:"owner"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	IsPrivate   bool         `json:"is_private"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Rating      int          `json:"rating"`
	IsFavorite  bool         `json:"is_favorite"`
	TracksCount int          `json:"tracks_count"`
	Tracks      []*TrackMeta `json:"tracks"`
}

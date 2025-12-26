package entity

import (
	"time"

	"github.com/google/uuid"
)

type AlbumMeta struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Label       string    `json:"label"`
	LicenseID   uuid.UUID `json:"license_id"`
	ReleaseDate time.Time `json:"release_date"`
}

type AlbumMetaAggregated struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Label       string        `json:"label"`
	License     *License      `json:"license"`
	ReleaseDate time.Time     `json:"release_date"`
	Artists     []*ArtistMeta `json:"artists"`
	Genres      []*Genre      `json:"genres"`
}

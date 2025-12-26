package entity

import "github.com/google/uuid"

type TrackMeta struct {
	ID           uuid.UUID `json:"id"`
	GenreID      uuid.UUID `json:"genre_id"`
	Name         string    `json:"name"`
	Duration     int       `json:"duration"`
	Explicit     bool      `json:"explicit"`
	LicenseID    uuid.UUID `json:"license_id"`
	AlbumID      uuid.UUID `json:"album_id"`
	TrackNumber  int       `json:"track_number"`
	TotalStreams int       `json:"total_streams"`
}

type TrackMetaAggregated struct {
	ID           uuid.UUID     `json:"id"`
	Genre        *Genre        `json:"genre"`
	Name         string        `json:"name"`
	Duration     int           `json:"duration"`
	Explicit     bool          `json:"explicit"`
	License      *License      `json:"license"`
	TrackNumber  int           `json:"track_number"`
	TotalStreams int           `json:"total_streams"`
	Album        *AlbumMeta    `json:"album"`
	Artists      []*ArtistMeta `json:"artists"`
}

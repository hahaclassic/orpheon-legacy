package entity

import "github.com/google/uuid"

type CoverObjectType string

const (
	CoverAlbum    CoverObjectType = "album"
	CoverPlaylist CoverObjectType = "playlist"
)

type Cover struct {
	ObjectID uuid.UUID `json:"object_id"`
	Data     []byte    `json:"data"`
}

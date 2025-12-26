package entity

import (
	"github.com/google/uuid"
)

// meta information about segment
// !!! UPDATED: NO SEGMENT ID ONLY IDX
type Segment struct {
	TrackID      uuid.UUID `json:"track_id"`
	Idx          int       `json:"idx"`
	TotalStreams uint64    `json:"total_streams"`
	Range        *Range    `json:"range"`
}

type SegmentsIdxs struct {
	TrackID uuid.UUID `json:"track_id"`
	Idxs    []int     `json:"idxs"`
}

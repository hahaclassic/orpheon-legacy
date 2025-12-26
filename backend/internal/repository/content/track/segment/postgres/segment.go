package segment_postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackSegmentRepository struct {
	pool *pgxpool.Pool
}

func NewTrackSegmentRepository(pool *pgxpool.Pool) *TrackSegmentRepository {
	return &TrackSegmentRepository{pool: pool}
}

func (r *TrackSegmentRepository) IncrementTotalStreams(ctx context.Context, trackID uuid.UUID, segmentsIdxs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("rollback error", "err", err)
		}
	}()

	for _, idx := range segmentsIdxs {
		_, err := tx.Exec(ctx, `
			UPDATE track_segments
			SET total_streams = total_streams + 1
			WHERE track_id = $1 AND index = $2
		`, trackID, idx)

		if err != nil {
			return fmt.Errorf("increment stream count: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *TrackSegmentRepository) GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.index, t.total_streams, t.start_time, t.end_time
		FROM track_segments t
		WHERE t.track_id = $1
		ORDER BY t.index
	`, trackID)
	if err != nil {
		return nil, fmt.Errorf("get segments: %w", err)
	}
	defer rows.Close()

	var segments []*entity.Segment

	for rows.Next() {
		var seg entity.Segment
		var start, end int
		err := rows.Scan(&seg.Idx, &seg.TotalStreams, &start, &end)
		if err != nil {
			return nil, fmt.Errorf("scan segment: %w", err)
		}
		seg.TrackID = trackID
		seg.Range = &entity.Range{Start: start, End: end}
		segments = append(segments, &seg)
	}

	return segments, nil
}

func (r *TrackSegmentRepository) DeleteSegments(ctx context.Context, trackID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM track_segments WHERE track_id = $1
	`, trackID)
	if err != nil {
		return fmt.Errorf("delete segments: %w", err)
	}
	return nil
}

func (r *TrackSegmentRepository) CreateSegments(ctx context.Context, trackID uuid.UUID, segments []*entity.Segment) error {
	if len(segments) <= 0 {
		return fmt.Errorf("number of segments must be positive")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("rollback error", "err", err)
		}
	}()

	for i, seg := range segments {
		_, err := tx.Exec(ctx, `
			INSERT INTO track_segments (track_id, index, total_streams, start_time, end_time)
			VALUES ($1, $2, 0, $3, $4)
		`, trackID, seg.Idx, seg.Range.Start, seg.Range.End)
		if err != nil {
			return fmt.Errorf("create segment %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

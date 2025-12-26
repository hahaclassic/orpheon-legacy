package tracks

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type PlaylistTracksRepository interface {
	AddTrackToPlaylist(ctx context.Context, playlistTrack *entity.PlaylistTrack) error
	DeleteTrackFromPlaylist(ctx context.Context, playlistTrack *entity.PlaylistTrack) error
	DeleteAllTracksFromPlaylist(ctx context.Context, playlistID uuid.UUID) error
	GetAllPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]*entity.TrackMeta, error)
	ChangeTrackPosition(ctx context.Context, playlistTrack *entity.PlaylistTrack) error
}

type PlaylistTrackService struct {
	repo   PlaylistTracksRepository
	policy usecase.PlaylistPolicyService
}

func NewPlaylistTrackService(repo PlaylistTracksRepository, policy usecase.PlaylistPolicyService) *PlaylistTrackService {
	return &PlaylistTrackService{
		repo:   repo,
		policy: policy,
	}
}

func (s *PlaylistTrackService) AddTrack(ctx context.Context, claims *entity.Claims, playlistTrack *entity.PlaylistTrack) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAddTrack, err)
	}()

	if err = s.policy.CanEdit(ctx, claims, playlistTrack.PlaylistID); err != nil {
		return err
	}

	return s.repo.AddTrackToPlaylist(ctx, playlistTrack)
}

func (s *PlaylistTrackService) GetAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (tracks []*entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllTracks, err)
	}()

	if err = s.policy.CanView(ctx, claims, playlistID); err != nil {
		return nil, err
	}

	return s.repo.GetAllPlaylistTracks(ctx, playlistID)
}

func (s *PlaylistTrackService) DeleteTrack(ctx context.Context, claims *entity.Claims, playlistTrack *entity.PlaylistTrack) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteTrack, err)
	}()

	if err = s.policy.CanEdit(ctx, claims, playlistTrack.PlaylistID); err != nil {
		return err
	}

	return s.repo.DeleteTrackFromPlaylist(ctx, playlistTrack)
}

func (s *PlaylistTrackService) DeleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAllTracks, err)
	}()

	if err = s.policy.CanEdit(ctx, claims, playlistID); err != nil {
		return err
	}

	return s.repo.DeleteAllTracksFromPlaylist(ctx, playlistID)
}

func (s *PlaylistTrackService) ChangeTrackPosition(ctx context.Context, claims *entity.Claims, playlistTrack *entity.PlaylistTrack) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrChangeTrackPosition, err)
	}()

	if err = s.policy.CanEdit(ctx, claims, playlistTrack.PlaylistID); err != nil {
		return err
	}

	return s.repo.ChangeTrackPosition(ctx, playlistTrack)
}

func (s *PlaylistTrackService) RestoreAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackIDs []uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrRestoreAllTracks, err)
	}()

	if err = s.policy.CanEdit(ctx, claims, playlistID); err != nil {
		return err
	}

	for i := range trackIDs {
		trackID := trackIDs[i]
		playlistTrack := &entity.PlaylistTrack{
			PlaylistID: playlistID,
			TrackID:    trackID,
		}
		if err := s.AddTrack(ctx, claims, playlistTrack); err != nil {
			return err
		}
	}

	return nil
}

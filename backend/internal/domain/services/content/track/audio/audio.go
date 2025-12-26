package audio

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrInvalidChunkParams = errors.New("invalid chunk parameters")
	ErrInvalidTrackID     = errors.New("invalid track id")
)

type AudioFileRepository interface {
	GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error)
	UploadAudioFile(ctx context.Context, chunk *entity.AudioChunk) error
	DeleteFile(ctx context.Context, trackID uuid.UUID) error
}

type AudioConverter interface {
	ChangeBitrate(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error)
}

type AudioFileService struct {
	converter AudioConverter
	repo      AudioFileRepository
}

func New(repo AudioFileRepository, converter AudioConverter) *AudioFileService {
	return &AudioFileService{
		repo:      repo,
		converter: converter,
	}
}

func (a *AudioFileService) GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (result *entity.AudioChunk, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAudioChunk, err)
	}()

	if chunk.End <= chunk.Start {
		return nil, ErrInvalidChunkParams
	}

	return a.repo.GetAudioChunk(ctx, chunk)
}

func (a *AudioFileService) UploadAudioFile(ctx context.Context, claims *entity.Claims, chunk *entity.AudioChunk) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUploadAudioFile, err)
	}()

	switch {
	case claims == nil || claims.AccessLvl != entity.Admin:
		return commonerr.ErrForbidden
	case chunk.End <= chunk.Start || chunk.Start != 0 || chunk.End != int64(len(chunk.Data)):
		return ErrInvalidChunkParams
	case chunk.TrackID == uuid.Nil:
		return ErrInvalidTrackID
	}

	converted, err := a.converter.ChangeBitrate(ctx, chunk)
	if err != nil {
		return err
	}

	return a.repo.UploadAudioFile(ctx, converted)
}

func (a *AudioFileService) DeleteAudioFile(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAudioFile, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.DeleteFile(ctx, trackID)
}

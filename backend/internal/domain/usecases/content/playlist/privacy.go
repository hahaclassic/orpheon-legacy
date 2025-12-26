package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrChangePrivacy = errors.New("failed to change privacy")
)

type PlaylistPrivacyChanger interface {
	ChangePrivacy(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, isPrivate bool) error
}

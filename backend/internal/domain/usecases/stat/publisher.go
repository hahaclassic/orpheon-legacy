package stats

import (
	"context"
	"errors"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrPublishListeningEvent = errors.New("failed to publish listening event")
)

type ListeningEventPublisher interface {
	SendListeningEvent(ctx context.Context, event *entity.ListeningEvent) error
}

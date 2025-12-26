package stats

import (
	"context"
	"errors"
)

var (
	ErrListen                = errors.New("failed to listen events")
	ErrConsumeListeningEvent = errors.New("failed to consume listening event")
)

type ListeningEventConsumer interface {
	Start(ctx context.Context) error
	//ConsumeListeningEvent(ctx context.Context, event *entity.ListeningEvent) error
}

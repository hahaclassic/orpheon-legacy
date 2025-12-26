package consumer

import (
	"context"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/stat"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type ListeningStatService interface {
	UpdateStat(ctx context.Context, event *entity.ListeningEvent) error
}

type EventBusSub interface {
	Subscribe(ctx context.Context, handler func(ctx context.Context, event *entity.ListeningEvent) error) error
}

type ListeningEventConsumer struct {
	bus  EventBusSub
	stat ListeningStatService
}

func New(bus EventBusSub, statService ListeningStatService) *ListeningEventConsumer {
	return &ListeningEventConsumer{bus: bus, stat: statService}
}

func (c *ListeningEventConsumer) Start(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrListen, err)
		}
	}()

	return c.bus.Subscribe(ctx, c.consumeListeningEvent)
}

func (c *ListeningEventConsumer) consumeListeningEvent(ctx context.Context, event *entity.ListeningEvent) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrConsumeListeningEvent, err)
		}
	}()

	return c.stat.UpdateStat(ctx, event)
}

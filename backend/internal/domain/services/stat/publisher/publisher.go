package publisher

import (
	"context"
	"errors"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/stat"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

const (
	MinTotalDuration = 15
)

var (
	ErrShortListeningTime = errors.New("error: the listening time is too short")
)

type EventBusPub interface {
	Publish(ctx context.Context, event *entity.ListeningEvent) error
}

type ListeningEventPublisher struct {
	bus EventBusPub
}

func New(bus EventBusPub) *ListeningEventPublisher {
	return &ListeningEventPublisher{bus: bus}
}

func (p *ListeningEventPublisher) PublishListeningEvent(ctx context.Context, event *entity.ListeningEvent) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrPublishListeningEvent, err)
		}
	}()

	if totalDuration(event) < MinTotalDuration {
		return ErrShortListeningTime
	}

	if err = p.bus.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}

func totalDuration(event *entity.ListeningEvent) int {
	total := 0
	for _, r := range event.Ranges {
		total += r.End - r.Start
	}

	return total
}

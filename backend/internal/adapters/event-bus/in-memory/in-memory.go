package inmemory

import (
	"context"
	"sync"
)

type EventHandler[T any] func(ctx context.Context, event T) error

type InMemoryEventBus[T any] struct {
	mu       sync.RWMutex
	handlers []EventHandler[T]
}

func NewInMemoryEventBus[T any]() *InMemoryEventBus[T] {
	return &InMemoryEventBus[T]{
		handlers: make([]EventHandler[T], 0),
	}
}

func (e *InMemoryEventBus[T]) Subscribe(ctx context.Context, handler EventHandler[T]) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, handler)

	return nil
}

func (e *InMemoryEventBus[T]) Publish(ctx context.Context, event T) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, handler := range e.handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

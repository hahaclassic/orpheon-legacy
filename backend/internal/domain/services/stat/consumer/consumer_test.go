package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestListeningEventConsumer_Start(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(bus *mocks.EventBusSub)
		expectErr bool
	}{
		{
			name: "success",
			setupMock: func(bus *mocks.EventBusSub) {
				bus.On("Subscribe", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			expectErr: false,
		},
		{
			name: "subscribe error",
			setupMock: func(bus *mocks.EventBusSub) {
				bus.On("Subscribe", mock.Anything, mock.Anything).
					Return(errors.New("subscribe failed")).
					Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := mocks.NewEventBusSub(t)
			stat := mocks.NewListeningStatService(t)

			tt.setupMock(bus)

			c := New(bus, stat)
			err := c.Start(context.Background())
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListeningEventConsumer_consumeListeningEvent(t *testing.T) {
	tests := []struct {
		name      string
		event     *entity.ListeningEvent
		setupMock func(stat *mocks.ListeningStatService)
		expectErr bool
	}{
		{
			name: "success",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 2, End: 39},
					{Start: 55, End: 141},
				},
			},
			setupMock: func(stat *mocks.ListeningStatService) {
				stat.On("UpdateStat", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "update error",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 60},
				},
			},
			setupMock: func(stat *mocks.ListeningStatService) {
				stat.On("UpdateStat", mock.Anything, mock.Anything).
					Return(errors.New("stat error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := mocks.NewEventBusSub(t)
			stat := mocks.NewListeningStatService(t)

			tt.setupMock(stat)

			c := New(bus, stat)
			err := c.consumeListeningEvent(context.Background(), tt.event)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

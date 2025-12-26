package publisher

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/stat"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestListeningEventPublisher_PublishListeningEvent(t *testing.T) {
	tests := []struct {
		name      string
		event     *entity.ListeningEvent
		setupMock func(bus *mocks.EventBusPub)
		expectErr error
	}{
		{
			name: "successfully published",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 10},
					{Start: 15, End: 25},
				},
			},
			setupMock: func(bus *mocks.EventBusPub) {
				bus.On("Publish", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectErr: nil,
		},
		{
			name: "listening time too short",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 5},
					{Start: 6, End: 9},
				},
			},
			setupMock: func(bus *mocks.EventBusPub) {
				// without any action
			},
			expectErr: ErrShortListeningTime,
		},
		{
			name: "publish failed",
			event: &entity.ListeningEvent{
				TrackID: uuid.New(),
				UserID:  uuid.New(),
				Ranges: []*entity.Range{
					{Start: 0, End: 20},
				},
			},
			setupMock: func(bus *mocks.EventBusPub) {
				bus.On("Publish", mock.Anything, mock.Anything).
					Return(errors.New("publish failed")).Once()
			},
			expectErr: usecase.ErrPublishListeningEvent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := mocks.NewEventBusPub(t)
			if tt.setupMock != nil {
				tt.setupMock(bus)
			}

			publisher := New(bus)
			err := publisher.PublishListeningEvent(context.Background(), tt.event)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectErr),
					"expected: %v, received err: %v", tt.expectErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

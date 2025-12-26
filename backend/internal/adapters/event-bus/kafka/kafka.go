package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
	"github.com/segmentio/kafka-go"
)

const topic = "listening_events"

var (
	ErrKafkaPublish   = errors.New("kafka publish error")
	ErrKafkaSubscribe = errors.New("kafka subscribe error")
)

type KafkaEventBus struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaEventBus(brokers []string, groupID string) *KafkaEventBus {
	return &KafkaEventBus{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			GroupID:        groupID,
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
		}),
	}
}

func (k *KafkaEventBus) Publish(ctx context.Context, event *entity.ListeningEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrKafkaPublish, err)
	}

	msg := kafka.Message{
		Key:   []byte(event.TrackID.String()),
		Value: data,
	}

	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("%w: %v", ErrKafkaPublish, err)
	}

	return nil
}

func (k *KafkaEventBus) Subscribe(ctx context.Context, handler func(ctx context.Context, event *entity.ListeningEvent) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := k.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return err
				}
				return errwrap.Wrap(ErrKafkaSubscribe, err)
			}

			event := &entity.ListeningEvent{}
			if err := json.Unmarshal(m.Value, event); err != nil {
				fmt.Printf("unmarshal error: %v\n", err)
				continue
			}

			if err := handler(ctx, event); err != nil {
				fmt.Printf("handler error: %v\n", err)
				continue
			}
		}
	}
}

func (k *KafkaEventBus) Close() error {
	if err := k.reader.Close(); err != nil {
		return err
	}
	return k.writer.Close()
}

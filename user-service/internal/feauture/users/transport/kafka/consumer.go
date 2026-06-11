package transport_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	core_logger "user-service/internal/core/logger"
	core_kafka "user-service/internal/core/transport/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer *kafka.Consumer
	service  Service
	log      *core_logger.Logger
}


func NewConsumer(cfg core_kafka.ConsumerCfg, service Service, topic string) (*Consumer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":          cfg.BrokersString(),
		"group.id":                   core_kafka.RegisterUserConsumerGroup,
		"auto.offset.reset":          "earliest",
		"enable.auto.offset.store":   false,
		"queued.max.messages.kbytes": 1000,
		"enable.auto.commit":         true,
		"auto.commit.interval.ms":    5000,
		// FIX
		"session.timeout.ms": 6000,
	}

	consumer, err := kafka.NewConsumer(&conf)

	if err != nil {
		return nil, fmt.Errorf("Failed to create consumer: %w", err)
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to subscribe to the topic: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		service:  service,
	}, nil
}

type Service interface {
	SaveUser(
		ctx context.Context,
		ID uuid.UUID,
		username string,
		phoneNumber string, 
	) (error)
}

func (c *Consumer) Run(ctx context.Context) error {
	var event UserRegisterEvent

	defer c.consumer.Close()
	for {
		select{
		case <-ctx.Done():
			c.log.Info("Kafka consumer stopped...")
			return nil

		default:
		}

		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			c.log.Error("Read kafka message error: %v", zap.Error(err))
			continue
		}

		err = json.Unmarshal(
			msg.Value,
			&event,
		)
		if err != nil {
			c.log.Error("Unmarshal user register event: %w", zap.Error(err))
		}

		err = c.service.SaveUser(
			ctx,
			event.ID,
			event.Username,
			event.PhoneNumber,
		)
		if err != nil {
			c.log.Fatal("Save user error: %w", zap.Error(err))
			continue
		}

		_, err = c.consumer.CommitMessage(
			msg,
		)
		if err != nil {
			c.log.Error("commit kafka message error", zap.Error(err))
		}
		
	}
}

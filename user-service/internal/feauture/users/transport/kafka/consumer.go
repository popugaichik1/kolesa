package transport_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
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


func NewConsumer(cfg core_kafka.ConsumerCfg, service Service, topic string, log *core_logger.Logger) (*Consumer, error) {
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
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to the topic: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		service:  service,
		log:      log,
	}, nil
}

type Service interface {
	SaveUser(
		ctx context.Context,
		ID uuid.UUID,
		username string,
		phoneNumber string,
	) error
}

func (c *Consumer) Run(ctx context.Context) error {
	var event UserRegisterEvent

	defer func() {
		if err := c.consumer.Close(); err != nil {
			c.log.Error("consumer close error: ", zap.Error(err))
		}
	}()

	for {
		select{
		case <-ctx.Done():
			c.log.Info("Kafka consumer stopped...")
			return nil

		default:
		}

		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				continue
			}
			c.log.Error("read kafka message error", zap.Error(err))
			continue
		}

		if err = json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("unmarshal user register event error", zap.Error(err))
			continue
		}

		userID, err := uuid.Parse(event.ID)
		if err != nil {
			c.log.Error("parse user id from event error", zap.String("id", event.ID), zap.Error(err))
			continue
		}

		if err = c.service.SaveUser(ctx, userID, event.Username, event.PhoneNumber); err != nil {
			c.log.Error("save user error", zap.Error(err))
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

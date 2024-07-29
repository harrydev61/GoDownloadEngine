package producer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"time"
)

type producerClient struct {
	id       string
	logger   core.Logger
	producer sarama.SyncProducer
}

func NewProducerClient(id string, logger core.Logger, config common.MqConfig) (*producerClient, error) {
	producer, err := sarama.NewSyncProducer(strings.Split(config.Brokers, ","), newSaramaConfig(config))
	if err != nil {
		return &producerClient{}, errors.Errorf("Error creating sync produce: %v", err)
	}
	return &producerClient{id: id, logger: logger, producer: producer}, nil
}

func (c *producerClient) ID() string {
	return c.id
}
func (c *producerClient) Stop() error {
	if err := c.producer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %s", err)
	}
	return nil
}

func newSaramaConfig(mqConfig common.MqConfig) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = mqConfig.ClientId
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func (c producerClient) Produce(ctx context.Context, queueName string, payload []byte) error {
	if _, _, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     queueName,
		Value:     sarama.ByteEncoder(payload),
		Timestamp: time.Now(),
	}); err != nil {
		c.logger.Errorf("failed to produce message: %s", err)
		return status.Error(codes.Internal, "failed to produce message")
	}
	return nil
}

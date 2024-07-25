package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"os"
	"os/signal"
	"strings"
)

type consumer struct {
	id                        string
	logger                    core.Logger
	config                    *sarama.Config
	client                    sarama.ConsumerGroup
	consumerGroup             *sarama.ConsumerGroupHandler
	queueNameToHandlerFuncMap map[string]common.MqHandlerFunc
	clientId                  string
}

func NewConsumerC(id string, sctx core.ServiceContext, config common.MqConfig) (*consumer, error) {
	logger := sctx.Logger(id)
	client, err := sarama.NewConsumerGroup(strings.Split(config.Brokers, ","), config.Group, newSaramaConfig(config))
	if err != nil {
		return nil, errors.Errorf("Error creating consumer group client: %v", err)
	}
	return &consumer{id: id, logger: logger, config: newSaramaConfig(config), client: client, queueNameToHandlerFuncMap: make(map[string]common.MqHandlerFunc)}, nil
}

func (c *consumer) ID() string {
	return c.id
}

func newSaramaConfig(mqConfig common.MqConfig) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = mqConfig.ClientId
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func (c *consumer) Stop() error {
	c.client.Close()
	return nil
}

func (c *consumer) GetConfig() *sarama.Config {
	return c.config
}

type consumerHandler struct {
	handlerFunc       common.MqHandlerFunc
	exitSignalChannel chan os.Signal
}

func newConsumerHandler(
	handlerFunc common.MqHandlerFunc,
	exitSignalChannel chan os.Signal,
) *consumerHandler {
	return &consumerHandler{
		handlerFunc:       handlerFunc,
		exitSignalChannel: exitSignalChannel,
	}
}

func (h consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				session.Commit()
				return nil
			}

			if err := h.handlerFunc(session.Context(), message.Topic, message.Value); err != nil {
				return err
			}

		case <-h.exitSignalChannel:
			session.Commit()
			break
		}
	}
}

func (c *consumer) RegisterHandler(queueName string, handlerFunc common.MqHandlerFunc) {
	c.queueNameToHandlerFuncMap[queueName] = handlerFunc
}

func (c consumer) Start(ctx context.Context) error {
	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, os.Interrupt)
	for queueName, handlerFunc := range c.queueNameToHandlerFuncMap {
		go func(queueName string, handlerFunc common.MqHandlerFunc) {
			if err := c.client.Consume(
				context.Background(),
				[]string{queueName},
				newConsumerHandler(handlerFunc, exitSignalChannel),
			); err != nil {
				c.logger.Errorf("Error failed to consume message from queue topic :%s, error: %v", queueName, err)
			}
		}(queueName, handlerFunc)
		c.logger.Infof("Successfully consumed message from queue topic :%s", queueName)
	}
	<-exitSignalChannel
	c.logger.Warnf("Consumer [%s] has stopped", c.id)
	return nil
}

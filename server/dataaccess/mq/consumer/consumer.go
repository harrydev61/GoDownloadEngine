package consumer

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"os"
	"os/signal"
	"strings"
)

type Consumer interface {
	RegisterHandler(queueName string, handlerFunc common.MqHandlerFunc)
	Start(ctx context.Context) error
	Stop() error
}

type consumer struct {
	id                        string
	logger                    core.Logger
	config                    *sarama.Config
	client                    sarama.ConsumerGroup
	consumerGroup             *sarama.ConsumerGroupHandler
	queueNameToHandlerFuncMap map[string]common.MqHandlerFunc
}

func newSaramaConfig(mqConfig common.MqConfig) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.ClientID = mqConfig.ClientId
	saramaConfig.Metadata.Full = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	return saramaConfig
}

func NewConsumer(id string, sctx core.ServiceContext, config common.MqConfig) (Consumer, error) {
	logger := sctx.Logger(id)
	client, err := sarama.NewConsumerGroup(strings.Split(config.Brokers, ","), config.Group, newSaramaConfig(config))
	if err != nil {

		return nil, errors.Errorf("Error creating consumer group client: %v", err)
	}
	return &consumer{
		id:                        id,
		logger:                    logger,
		config:                    newSaramaConfig(config),
		client:                    client,
		queueNameToHandlerFuncMap: make(map[string]common.MqHandlerFunc),
	}, nil
}

func (c *consumer) ID() string {
	return c.id
}

func (c *consumer) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
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
	}
	<-exitSignalChannel
	c.logger.Warnf("Consumer [%s] has stopped", c.id)
	return nil
}

type consumerHandler struct {
	handlerFunc       common.MqHandlerFunc
	exitSignalChannel chan os.Signal
}

func newConsumerHandler(handlerFunc common.MqHandlerFunc, exitSignalChannel chan os.Signal) *consumerHandler {
	return &consumerHandler{
		handlerFunc:       handlerFunc,
		exitSignalChannel: exitSignalChannel,
	}
}
func (h *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("setup consumer", session.Claims())
	return nil
}
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				session.Commit()
				return nil
			}
			_ = h.handlerFunc(session.Context(), message.Topic, message.Value)
			session.MarkMessage(message, "")
		case <-h.exitSignalChannel:
			session.Commit()
			break
		}
	}
}

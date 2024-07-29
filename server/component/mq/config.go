package mq

import (
	"flag"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

type Configs struct {
	id       string
	logger   core.Logger
	brokers  string
	group    string
	assignor string
	oldest   bool
	verbose  bool
	clientId string
}

func NewConfig(id string) *Configs {
	return &Configs{id: id}
}

func (c *Configs) ID() string {
	return c.id
}

func (c *Configs) InitFlags() {

	flag.StringVar(&c.brokers, "kafka_brokers", "", "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&c.group, "kafka_group", "", "Kafka consumer group definition")
	flag.StringVar(&c.assignor, "kafka_assignor", "range", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	flag.BoolVar(&c.oldest, "kafka_oldest", true, "Kafka consumer consume initial offset from oldest")
	flag.BoolVar(&c.verbose, "kafka_verbose", false, "Sarama logging")
	flag.StringVar(&c.clientId, "kafka_client_id", "goDownload", "mq client id")
	flag.Parse()
}

func (c *Configs) Activate(serviceContext core.ServiceContext) error {
	if len(c.brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(c.group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
	return nil
}

func (c *Configs) Stop() error {
	return nil
}

func (c *Configs) GetMqConfig() common.MqConfig {
	return common.MqConfig{
		Brokers:  c.brokers,
		Group:    c.group,
		Assignor: c.assignor,
		Oldest:   c.oldest,
		Verbose:  c.verbose,
		ClientId: c.clientId,
	}
}

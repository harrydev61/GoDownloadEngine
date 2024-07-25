package producer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"time"
)

type ProducerC struct {
	id       string
	logger   core.Logger
	producer sarama.SyncProducer
}

func NewProducerC(id string, sctx core.ServiceContext, config common.MqConfig) (*ProducerC, error) {
	logger := sctx.Logger(id)
	producer, err := sarama.NewSyncProducer(strings.Split(config.Brokers, ","), newSaramaConfig(config, logger))
	if err != nil {
		return nil, errors.Errorf("Error creating sync produce: %v", err)
	}
	return &ProducerC{id: id, logger: logger, producer: producer}, nil
}

func (c *ProducerC) ID() string {
	return c.id
}
func (c *ProducerC) InitFlags() {
	return
}
func (c *ProducerC) Activate(serviceContext core.ServiceContext) error {
	return nil
}
func (c *ProducerC) Stop() error {
	if err := c.producer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %s", err)
	}
	return nil
}

func newSaramaConfig(mqConfig common.MqConfig, logger core.Logger) *sarama.Config {
	version, err := sarama.ParseKafkaVersion(mqConfig.Version)
	if err != nil {
		logger.Panicf("Error parsing Kafka version: %v", err)
		return nil
	}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = version                                // Chọn phiên bản của Kafka mà Sarama sẽ giả định đang chạy
	saramaConfig.Producer.Idempotent = true                       //Kích hoạt chế độ idempotent giúp đảm bảo rằng mỗi tin nhắn chỉ được gửi một lần duy nhất.
	saramaConfig.Producer.Return.Successes = true                 //Bật cấu hình này để nhận lại tin nhắn đã gửi thành công qua kênh Successes.
	saramaConfig.Producer.Return.Errors = true                    //Bật cấu hình này để nhận lại các lỗi phát sinh khi gửi tin nhắn qua kênh Errors.
	saramaConfig.Net.MaxOpenRequests = 1                          //Số yêu cầu tối đa được gửi đồng thời
	saramaConfig.Producer.MaxMessageBytes = 1024 * 1024           //Kích thước tối đa của mỗi tin nhắn.
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll        // Xác nhận cần thiết từ broker.
	saramaConfig.Producer.Timeout = 10 * time.Second              //Thời gian tối đa mà producer sẽ chờ xác nhận từ broker trước khi gửi tin nhắn tiếp theo.
	saramaConfig.Producer.Partitioner = sarama.NewHashPartitioner // Phân vùng cho tin nhắn.
	saramaConfig.Producer.Retry.Max = 3                           //Cấu hình thử lại khi yêu cầu Producer gặp lỗi. Cài đặt này đảm bảo rằng producer sẽ thử lại gửi tin nhắn nhiều lần nếu cần.
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond
	saramaConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault //Đặt mức độ nén dữ liệu cho tin nhắn.
	saramaConfig.Producer.Transaction.Timeout = 1 * time.Minute             //Đặt thời gian tối đa mà một giao dịch có thể giữ không được giải quyết.
	saramaConfig.Producer.Transaction.Retry.Max = 50                        //Cấu hình thử lại khi giao dịch gặp lỗi.
	saramaConfig.Producer.Transaction.Retry.Backoff = 100 * time.Millisecond
	saramaConfig.ClientID = mqConfig.ClientId
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func (c ProducerC) Produce(ctx context.Context, queueName string, payload []byte) error {
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

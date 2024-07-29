package factory

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type s3ClientReadWriteCloser struct {
	writtenData []byte
	isClosed    bool
	logger      core.Logger
}

func newS3ClientReadWriteCloser(
	ctx context.Context,
	minioClient *minio.Client,
	logger core.Logger,
	bucketName,
	objectName string,
) io.ReadWriteCloser {
	readWriteCloser := &s3ClientReadWriteCloser{
		writtenData: make([]byte, 0),
		isClosed:    false,
		logger:      logger,
	}

	go func() {
		if _, err := minioClient.PutObject(
			ctx, bucketName, objectName, readWriteCloser, -1, minio.PutObjectOptions{},
		); err != nil {
			logger.Errorf("s3 put object failed", zap.Error(err))
		}
	}()

	return readWriteCloser
}

func (s *s3ClientReadWriteCloser) Close() error {
	s.isClosed = true
	return nil
}

func (s *s3ClientReadWriteCloser) Read(p []byte) (int, error) {
	if len(s.writtenData) > 0 {
		writtenLength := copy(p, s.writtenData)
		s.writtenData = s.writtenData[writtenLength:]
		return writtenLength, nil
	}

	if s.isClosed {
		return 0, io.EOF
	}

	return 0, nil
}

func (s *s3ClientReadWriteCloser) Write(p []byte) (int, error) {
	s.writtenData = append(s.writtenData, p...)
	return len(p), nil
}

type S3Client struct {
	minioClient *minio.Client
	bucket      string
	logger      core.Logger
}

func NewS3Client(
	downloadConfig ClientConfig,
	logger core.Logger,
) (common.FileClient, error) {
	minioClient, err := minio.New(downloadConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(downloadConfig.AccessKeyID, downloadConfig.SecretAccessKey, ""),
		Secure: downloadConfig.UseSSL,
	})
	if err != nil {
		logger.Errorf("Failed to create MinIO client: %v", err)
		return nil, err
	}

	return &S3Client{
		minioClient: minioClient,
		bucket:      downloadConfig.Bucket,
		logger:      logger,
	}, nil
}

func (s S3Client) Read(ctx context.Context, filePath string) (io.ReadCloser, error) {
	logger := s.logger.With("file_path", filePath)

	object, err := s.minioClient.GetObject(ctx, s.bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("Failed to get S3 object", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get S3 object")
	}

	return object, nil
}

func (s S3Client) Write(ctx context.Context, filePath string) (io.WriteCloser, error) {
	return newS3ClientReadWriteCloser(ctx, s.minioClient, s.logger, s.bucket, filePath), nil
}

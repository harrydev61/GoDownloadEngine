package factory

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path"
)

type bufferedFileReader struct {
	file           *os.File
	bufferedReader io.Reader
}

func newBufferedFileReader(
	file *os.File,
) io.ReadCloser {
	return &bufferedFileReader{
		file:           file,
		bufferedReader: bufio.NewReader(file),
	}
}

func (b bufferedFileReader) Close() error {
	return b.file.Close()
}

func (b bufferedFileReader) Read(p []byte) (int, error) {
	return b.bufferedReader.Read(p)
}

type LocalClient struct {
	downloadDirectory string
	logger            core.Logger
}

func newLocalClient(
	downloadConfig ClientConfig,
	logger core.Logger,
) (common.FileClient, error) {
	if err := os.MkdirAll(downloadConfig.DownloadDirectory, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("failed to create download directory: %w", err)
		}
	}

	return &LocalClient{
		downloadDirectory: downloadConfig.DownloadDirectory,
		logger:            logger,
	}, nil
}

func (l LocalClient) Read(ctx context.Context, filePath string) (io.ReadCloser, error) {
	absolutePath := path.Join(l.downloadDirectory, filePath)
	file, err := os.Open(absolutePath)
	if err != nil {
		l.logger.Error("failed to open file")
		return nil, status.Error(codes.Internal, "failed to open file")
	}

	return newBufferedFileReader(file), nil
}

func (l *LocalClient) Write(ctx context.Context, filePath string) (io.WriteCloser, error) {
	absolutePath := path.Join(l.downloadDirectory, filePath)
	l.logger.Debugf("file path: %s", absolutePath)
	file, err := os.Create(absolutePath)
	if err != nil {
		l.logger.Error("failed to open file")
		return nil, status.Error(codes.Internal, "failed to open file")
	}
	return file, nil
}

package filec

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path"
)

type fileC struct {
	id                string
	logger            core.Logger
	client            common.FileClient
	mode              string
	downloadDirectory string
	bucket            string
	address           string
	username          string
	password          string
}

func NewFileC(id string) *fileC {
	return &fileC{id: id, downloadDirectory: common.DownDirectory}
}

func (f *fileC) ID() string {
	return f.id
}

func (f *fileC) InitFlags() {
	flag.StringVar(
		&f.mode,
		"download_mode",
		"local",
		"Mode download",
	)
	flag.StringVar(
		&f.bucket,
		"download_bucket",
		"",
		"bucket name",
	)
	flag.StringVar(
		&f.address,
		"download_address",
		"",
		"download address",
	)
	flag.StringVar(
		&f.username,
		"download_username",
		"",
		"download username",
	)
	flag.StringVar(
		&f.password,
		"download_password",
		"",
		"download password",
	)

}

func (f *fileC) Activate(sctx core.ServiceContext) error {
	f.logger = sctx.Logger(f.id)
	client, err := newClient(*f, f.logger)
	if err != nil {
		f.logger.Errorf("New file client error: %v", err)
	}
	f.client = client
	return nil
}

func (r *fileC) Stop() error {
	return nil
}

func (r *fileC) GetClient() common.FileClient {
	return r.client
}

func newClient(downloadConfig fileC, logger core.Logger) (common.FileClient, error) {
	switch downloadConfig.mode {
	case common.DownloadModeLocal:
		return newLocalClient(downloadConfig, logger)
	default:
		return nil, fmt.Errorf("unsupported download mode: %s", downloadConfig.mode)
	}
}

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
	downloadConfig fileC,
	logger core.Logger,
) (common.FileClient, error) {
	if err := os.MkdirAll(downloadConfig.downloadDirectory, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("failed to create download directory: %w", err)
		}
	}

	return &LocalClient{
		downloadDirectory: downloadConfig.downloadDirectory,
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

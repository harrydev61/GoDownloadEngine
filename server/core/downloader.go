package core

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	HTTPResponseHeaderContentType = "Content-Type"
	HTTPMetadataKeyContentType    = "content-type"
	HttpMetadataKeyContentLength  = "content-length"
)

type Downloader interface {
	Download(ctx context.Context, writer io.Writer) (map[string]any, error)
}

type HTTPDownloader struct {
	url    string
	logger Logger
}

func NewHTTPDownloader(
	url string,
	logger Logger,
) Downloader {
	return &HTTPDownloader{
		url:    url,
		logger: logger,
	}
}

func (h HTTPDownloader) Download(ctx context.Context, writer io.Writer) (map[string]any, error) {

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, http.NoBody)
	if err != nil {
		h.logger.Error("failed to create http get request")
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		h.logger.Error("failed to make http get request")
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		h.logger.Error("failed to download file", zap.Int("status_code", response.StatusCode))
		return nil, fmt.Errorf("failed to download file, status_code: %d", response.StatusCode)
	}
	// Get summary information from headers
	contentType := response.Header.Get(HTTPResponseHeaderContentType)
	contentLength := response.ContentLength

	h.logger.Info("Downloading file",
		zap.String("content_type", contentType),
		zap.Int64("content_length", contentLength),
	)

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		h.logger.Error("failed to read response and write to writer")
		return nil, err
	}

	metadata := map[string]any{
		HTTPMetadataKeyContentType:   contentType,
		HttpMetadataKeyContentLength: contentLength,
	}

	return metadata, nil
}

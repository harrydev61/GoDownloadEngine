package factory

import (
	"errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

type ClientConfig struct {
	Mode              string
	DownloadDirectory string
	Bucket            string
	Endpoint          string
	AccessKeyID       string
	SecretAccessKey   string
	UseSSL            bool
}

func FileFactory(config ClientConfig, logger core.Logger) (common.FileClient, error) {
	switch config.Mode {
	case common.DownloadModeLocal:
		return newLocalClient(config, logger)
	case common.DownloadModeS3:
		return NewS3Client(config, logger)
	default:
		return nil, errors.New("not supported download mode")
	}
}

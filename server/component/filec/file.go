package filec

import (
	"flag"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/component/filec/factory"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

type fileC struct {
	id                string
	logger            core.Logger
	client            common.FileClient
	mode              string
	downloadDirectory string
	bucket            string
	endpoint          string
	accessKeyID       string
	secretAccessKey   string
	useSSL            bool
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
		&f.endpoint,
		"download_endpoint",
		"",
		"download endpoint",
	)
	flag.StringVar(
		&f.accessKeyID,
		"download_access_Key_id",
		"",
		"download accessKeyID",
	)
	flag.StringVar(
		&f.secretAccessKey,
		"download_secret_access_key",
		"",
		"download secretAccessKey",
	)
	flag.BoolVar(
		&f.useSSL,
		"download_use_ssl",
		false,
		"download use ssl",
	)

}

func (f *fileC) Activate(sctx core.ServiceContext) error {
	f.logger = sctx.Logger(f.id)
	config := factory.ClientConfig{
		Mode:              f.mode,
		DownloadDirectory: f.downloadDirectory,
		Bucket:            f.bucket,
		Endpoint:          f.endpoint,
		AccessKeyID:       f.accessKeyID,
		SecretAccessKey:   f.secretAccessKey,
		UseSSL:            f.useSSL,
	}
	client, err := factory.FileFactory(config, f.logger)
	if err != nil {
		f.logger.Errorf("New file client error: %v", err)
		return err
	}
	f.logger.Infof("New file client with mod: %v", f.mode)
	f.client = client
	return nil
}

func (r *fileC) Stop() error {
	return nil
}

func (r *fileC) GetClient() common.FileClient {
	return r.client
}

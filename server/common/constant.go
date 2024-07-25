package common

const (
	RecordIsDeleted  = 1
	RecordNotDeleted = 0

	ResponseStatusSuccess = 1
	ResponseStatusError   = 0

	UserStatusProcess  = -1
	UserStatusUnActive = 0
	UserStatusActive   = 1

	KeyCompMySQL      = "mysql"
	KeyCompGIN        = "gin"
	KeyCompJWT        = "jwt"
	KeyCompRedis      = "redis"
	KeyCompGRPC       = "grpc"
	KeyCompMqConfig   = "mqConfig"
	KeyCompProducer   = "producer"
	KeyCompConsumer   = "consumer"
	KeyCompFileClient = "fileClient"

	MaskTypeUser = 1
	MaskTypeTask = 2

	AuthTypeEmailPassword = 1

	//kafka
	DownloadTaskTopic = "downloadTask"

	//download status
	DownloadTaskPending = 0
	DownloadTaskLoading = 1
	DownloadTaskSuccess = 2
	DownloadTaskFailed  = 3
	DownloadTaskExpired = 4

	//download task type
	DownloadTaskTypeHTTP = 1

	DownDirectory = "download_sources"

	//download mode
	DownloadModeLocal = "local"
	DownloadModeS3    = "S3"

	//file
	DownloadFileNamePrefix = "download_file"
)

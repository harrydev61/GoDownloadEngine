package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"time"
)

var tableName = "download_tasks"

type DownloadTask struct {
	DownloadID     string    `json:"download_id" form:"download_id" gorm:"type:char(36);primary_key;column:download_id"`
	Name           string    `json:"name" form:"download_id" gorm:"type:char(500);column:name;not null"`
	Description    string    `json:"description" form:"download_id" gorm:"type:char(1000);column:description"`
	UserID         string    `json:"user_id" form:"user_id" gorm:"type:char(36);not null;column:user_id"`
	DownloadType   int16     `json:"download_type" form:"download_type" gorm:"type:smallint;not null;column:download_type"`
	URL            string    `json:"url" form:"url" gorm:"type:text;not null;column:url"`
	DownloadStatus int16     `json:"download_status" form:"download_status" gorm:"type:smallint;not null;column:download_status"`
	IsDeleted      int       `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;not null;column:is_deleted"`
	Metadata       JSON      `json:"metadata" form:"metadata" gorm:"type:text;not null:column:metadata"`
	CreatedAt      time.Time `json:"created_at" form:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" form:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at"`
}

func (DownloadTask) GetTableName() string {
	return tableName
}

func NewDownloadTask(userId, name string, downloadType int, url string, status int, des *string) *DownloadTask {
	now := time.Now().UTC()
	return &DownloadTask{
		DownloadID:     common.GentNewUuid().String(),
		Name:           name,
		Description:    *des,
		UserID:         userId,
		DownloadType:   int16(downloadType),
		URL:            url,
		DownloadStatus: int16(status),
		IsDeleted:      0,
		Metadata:       JSON{Data: make(map[string]any)},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

type JSON struct {
	Data any
}

func (j *JSON) Scan(src any) error {
	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, &j.Data)

	case string:
		return json.Unmarshal([]byte(src), &j.Data)

	default:
		return fmt.Errorf("unsupported type for json scan: %T", src)
	}
}

func (j JSON) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}

	return json.Marshal(j.Data)
}

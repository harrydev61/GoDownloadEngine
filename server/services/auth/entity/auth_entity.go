package entity

import (
	"time"
)

const tableName = "auths"

type AuthEntity struct {
	UserId     string     `json:"userId" gorm:"column:user_id;" db:"user_id"`
	Email      string     `json:"email" gorm:"column:email;" db:"email"`
	Salt       string     `json:"salt" gorm:"column:salt;" db:"salt" `
	Password   string     `json:"password" gorm:"column:password;" db:"password" `
	PublicKey  string     `json:"publicKey" gorm:"column:public_key;" db:"public_key" `
	PrivateKey string     `json:"privateKey" gorm:"column:private_key;" db:"private_key" `
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"  db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"  db:"updated_at"`
}

func (AuthEntity) GetTableName() string {
	return tableName
}

func NewAuthEntity(userid, email string, salt, password string) *AuthEntity {
	now := time.Now().UTC()
	return &AuthEntity{
		UserId:     userid,
		Email:      email,
		Salt:       salt,
		Password:   password,
		PublicKey:  "",
		PrivateKey: "",
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
}

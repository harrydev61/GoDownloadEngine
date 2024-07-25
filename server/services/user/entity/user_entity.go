package entity

import (
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"time"
)

const tableName = "users"

type UserEntity struct {
	UserId    string     `json:"userId" gorm:"column:user_id" db:"column:user_id"`
	FirstName string     `json:"firstName" gorm:"column:first_name" db:"column:first_name"`
	LastName  string     `json:"lastName" gorm:"column:last_name" db:"column:last_name"`
	Email     string     `json:"email" gorm:"column:email" db:"column:email"`
	Phone     string     `json:"phone" gorm:"column:phone" db:"column:phone"`
	Avatar    string     `json:"avatar" gorm:"column:avatar" db:"column:avatar"`
	Role      int        `json:"role" gorm:"column:role" db:"column:role"`
	Gender    string     `json:"gender" gorm:"column:gender" db:"column:gender"`
	Dob       *time.Time `json:"dob" gorm:"column:dob" db:"column:dob"`
	Status    int        `json:"status" gorm:"column:status" db:"column:status"`
	Ip        string     `json:"ip" gorm:"column:ip" db:"column:ip"`
	IsDeleted int        `json:"isDeleted" gorm:"column:is_deleted" db:"column:is_deleted"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at" db:"column:created_at"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at" db:"column:updated_at"`
}

func NewUserEntity(firstName, lastName, email string, phone string, avatar string, role int, gender string, dob *time.Time, status int, ip string) *UserEntity {
	return &UserEntity{
		UserId:    common.GentNewUuid().String(),
		FirstName: firstName, LastName: lastName, Email: email, Phone: phone, Avatar: avatar, Role: role, Gender: gender, Dob: dob, Status: status, Ip: ip,
		IsDeleted: 0,
		CreatedAt: dateNow(),
		UpdatedAt: dateNow(),
	}
}

func dateNow() *time.Time {
	now := time.Now().UTC()
	return &now
}

func (entity *UserEntity) GetTableName() string {
	return tableName
}

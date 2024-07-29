package entity

import "time"

type UserDataResponse struct {
	UserId    string     `json:"userId"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Avatar    string     `json:"avatar"`
	Role      int        `json:"role"`
	Gender    string     `json:"gender"`
	Dob       *time.Time `json:"dob"`
	Status    int        `json:"status"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at" db:"column:created_at"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at" db:"column:updated_at"`
}

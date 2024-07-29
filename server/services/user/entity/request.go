package entity

import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
	Email     string `validate:"required,email"`
	Ip        string `validate:"required,ip"`
	Firstname string `validate:"required"`
	Lastname  string `validate:"required"`
}

type GetDetail struct {
	UserID string `uri:"userID" validate:"required,uuid"`
}

func (d *GetDetail) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		// custom validate
		return err
	}
	return nil
}

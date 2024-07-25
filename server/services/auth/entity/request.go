package entity

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"regexp"
)

type AuthLogging struct {
	AuthType int    `json:"authType" form:"authType" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password"  validate:"required,min=8,max=100"`
}

func (ar *AuthLogging) Validate() error {
	regex := `^[a-zA-Z0-9]{6,32}$`
	re := regexp.MustCompile(regex)
	if !re.MatchString(ar.Password) {
		return errors.New("email or password is valid")
	}
	validate := validator.New()
	if err := validate.Struct(ar); err != nil {
		// custom validate
		return err
	}
	return nil
}

type AuthRegister struct {
	AuthType        int    `json:"authType" form:"authType" validate:"required"`
	FirstName       string `json:"firstName" form:"firstName" validate:"required"`
	LastName        string `json:"lastName" form:"lastName" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
}

func (ar *AuthRegister) Validate() error {
	regex := `^[a-zA-Z0-9]{6,32}$`
	re := regexp.MustCompile(regex)
	if !re.MatchString(ar.FirstName) || !re.MatchString(ar.LastName) || !re.MatchString(ar.Password) {
		return errors.New("firstName, lastName or password is valid")
	}
	switch ar.AuthType {
	case common.AuthTypeEmailPassword:
		validate := validator.New()
		if err := validate.Struct(ar); err != nil {
			// custom validate
			return err
		}
	default:
		return ErrAuthTypeValidation

	}

	return nil
}

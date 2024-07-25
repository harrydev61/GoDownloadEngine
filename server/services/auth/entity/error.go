package entity

import "errors"

var (
	ErrPasswordIsNotValid     = errors.New("password must have from 8 to 30 characters")
	ErrEmailIsNotValid        = errors.New("email is not valid")
	ErrEmailHasExisted        = errors.New("email has existed")
	ErrOtpHasExisted          = errors.New("otp has existed")
	ErrLoginFailed            = errors.New("email and password are not valid")
	ErrFirstNameIsEmpty       = errors.New("first name can not be blank")
	ErrFirstNameTooLong       = errors.New("first name too long, max character is 30")
	ErrLastNameIsEmpty        = errors.New("last name can not be blank")
	ErrLastNameTooLong        = errors.New("last name too long, max character is 30")
	ErrCannotRegister         = errors.New("cannot register")
	ErrAuthTypeValidation     = errors.New("auth type validation failed")
	ErrAuthTypeIsNotValid     = errors.New("auth type is not valid")
	ErrRefreshTokenIsNotValid = errors.New("refresh token is not valid")
	ErrAccessTokenIsNotValid  = errors.New("access token is not valid")
)

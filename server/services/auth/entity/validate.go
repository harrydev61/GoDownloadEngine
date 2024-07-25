package entity

import (
	"net/mail"
)

func EmailIsValid(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func CheckPassword(s string) error {
	if len(s) < 8 || len(s) > 30 {
		return ErrPasswordIsNotValid
	}

	return nil
}

func CheckFirstName(s string) error {
	if s == "" {
		return ErrFirstNameIsEmpty
	}

	if len(s) > 30 {
		return ErrFirstNameTooLong
	}

	return nil
}

func CheckLastName(s string) error {
	if s == "" {
		return ErrLastNameIsEmpty
	}

	if len(s) > 30 {
		return ErrLastNameTooLong
	}

	return nil
}

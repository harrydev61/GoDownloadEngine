package entity

import (
	"github.com/go-playground/validator/v10"
)

type Create struct {
	Name         string `json:"name" form:"name"  validate:"required,min=3,max=500"`
	Description  string `json:"description" form:"description"  validate:"required,min=8,max=1000"`
	URL          string `json:"url" form:"url"  validate:"required,url"`
	DownloadType int    `json:"downloadType" form:"downloadType"  validate:"required,oneof=1"`
}

func (d *Create) Validate() error {
	//regex := `^[a-zA-Z0-9 ]*$`
	//re := regexp.MustCompile(regex)
	//if !re.MatchString(d.Name) {
	//	return errors.New("name is valid")
	//}
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		// custom validate
		return err
	}
	return nil
}

type GetDetail struct {
	DownloadID string `uri:"downloadID" validate:"required,uuid"`
}

func (d *GetDetail) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		// custom validate
		return err
	}
	return nil
}

type GetLists struct {
	Page  int `form:"page" validate:"required,min=1,numeric"`
	Limit int `form:"limit" validate:"required,min=1,numeric"`
}

func (d *GetLists) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		// custom validate
		return err
	}
	return nil
}

type Delete struct {
	DownloadID string `uri:"downloadID" validate:"required,uuid"`
}

func (d *Delete) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		// custom validate
		return err
	}
	return nil
}

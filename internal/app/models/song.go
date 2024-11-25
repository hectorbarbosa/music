package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Song struct {
	ID int32 `example:"1"`
	// Group name
	Group string `validate:"required" example:"Muse"`
	// Song name
	Name string `validate:"required" example:"Supermassive Black Hole"`
	// Release date in 02.01.2006 format
	ReleaseDate time.Time `validate:"required" example:"16.07.2006"`
	// Song text
	Text string `validate:"required" example:"Some text\n"`
	// URL link
	Link string `validate:"required" example:"http://example.org"`
}

func (s *Song) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

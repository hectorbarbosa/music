package models

import (
	"time"

	"github.com/go-playground/validator"
)

type SongDetails struct {
	// Group name
	Group string `json:"group" validate:"required" example:"Muse"`
	// Song name
	Name string `json:"song" validate:"required" example:"Supermassive Black Hole"`
}

func (s *SongDetails) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

type CreateParams struct {
	Group       string `validate:"required"`
	Name        string `validate:"required"`
	ReleaseDate string `json:"releaseDate" validate:"required"`
	Text        string `json:"text" validate:"required"`
	Link        string `json:"link" validate:"required"`
}

func (s *CreateParams) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	_, err := time.Parse("02.01.2006", s.ReleaseDate)
	if err != nil {
		return err
	}

	return nil
}

type UpdateParams struct {
	// Group name
	Group string `json:"group_name" validate:"required" example:"Muse"`
	// Song name
	Name string `json:"song_name" validate:"required" example:"Supermassive Black Hole"`
	// Release date in 02.01.2006 format
	ReleaseDate string `json:"release_date" validate:"required" example:"16.07.2006"`
	// Song text
	Text string `json:"song_text" validate:"required" example:"Some text\n\n Some text2\n"`
	// URL link
	Link string `json:"link" validate:"required" example:"http://example.org"`
}

func (s *UpdateParams) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	_, err := time.Parse("02.01.2006", s.ReleaseDate)
	if err != nil {
		return err
	}

	return nil
}

type Verse struct {
	// verse number
	Num string `json:"num" example:"1"`
	// verse text
	Text string `json:"text" example:"Some text\n"`
}

package service

import (
	"fmt"
	"log/slog"
	"music/internal"
	"music/internal/app/models"
	"music/internal/config"
	m "music/internal/rest/models"
	"net/url"
	"strings"
)

type SongRepository interface {
	Create(params m.CreateParams) (models.Song, error)
	Delete(id int32) error
	Update(id int32, s m.UpdateParams) (models.Song, error)
	SelectText(id int32) (string, error)
	Search(vals url.Values, pageNum, perPage int) ([]models.Song, error)
}

type SongService struct {
	cfg    config.Config
	logger *slog.Logger
	repo   SongRepository
}

func NewSongService(cfg config.Config, logger *slog.Logger, repo SongRepository) *SongService {
	return &SongService{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
	}
}

func (s *SongService) Create(params m.CreateParams) (models.Song, error) {
	if err := params.Validate(); err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "service create")
	}

	song, err := s.repo.Create(params)
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (s *SongService) Update(id int32, p m.UpdateParams) (models.Song, error) {
	if err := p.Validate(); err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "service update")
	}

	song, err := s.repo.Update(id, p)
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (s *SongService) Delete(id int32) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *SongService) SelectVerse(id int32, v int) (string, error) {
	text, err := s.repo.SelectText(id)
	if err != nil {
		return "", err
	}

	verses := strings.Split(text, "\n\n")
	if len(verses) < v {
		return "", internal.NewErrorf(internal.ErrorCodeNotFound, "the song has only %d verses", len(verses))
	}

	// numeration from 1!
	if v < 1 {
		return "", internal.NewErrorf(internal.ErrorCodeInvalidArgument, "verses must be numerated from 1")
	}

	return verses[v-1], nil
}

func (s *SongService) Search(vals url.Values, pageNum, perPage int) ([]models.Song, error) {
	if err := validateURLParams(vals); err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "service search")
	}

	songs, err := s.repo.Search(vals, pageNum, perPage)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func validateURLParams(vals url.Values) error {
	// Only one key-val pair required
	for _, val := range vals {
		if len(val) > 1 {
			return fmt.Errorf("invalid url params, %v", val)
		}
	}
	return nil
}

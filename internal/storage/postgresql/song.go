package postgresql

import (
	"database/sql"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"music/internal"
	"music/internal/app/models"
	m "music/internal/rest/models"
)

type SongRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewSongRepo(db *sql.DB, logger *slog.Logger) *SongRepository {
	return &SongRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SongRepository) Create(p m.CreateParams) (models.Song, error) {
	var id int32
	release, err := time.Parse("02.01.2006", p.ReleaseDate)
	if err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid date")
	}

	if err := r.db.QueryRow(
		`INSERT INTO public.songs 
		    (group_name, song_name, release_date, song_text, link) 
		VALUES 
		    ($1, $2, $3, $4, $5) 
		RETURNING id;`,
		p.Group, p.Name, release, p.Text, p.Link,
	).Scan(&id); err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo create")
	}

	r.logger.Debug("record created", "id", id)

	return models.Song{
		ID:          id,
		Group:       p.Group,
		Name:        p.Name,
		ReleaseDate: release,
		Text:        p.Text,
		Link:        p.Link,
	}, nil
}

func (r *SongRepository) Delete(id int32) error {
	result, err := r.db.Exec("DELETE FROM public.songs WHERE id = $1", id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo delete")
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo delete")
	}
	if deleted != 1 {
		return internal.NewErrorf(internal.ErrorCodeNotFound, "resourse with id %d not found", id)
	}

	r.logger.Debug("record deleted", "id", id)
	return nil
}

func (r *SongRepository) Update(id int32, p m.UpdateParams) (models.Song, error) {
	release, err := time.Parse("02.01.2006", p.ReleaseDate)
	if err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid date")
	}
	result, err := r.db.Exec(
		`UPDATE
		    public.songs 
		SET 
		    group_name = $1, song_name = $2, release_date = $3, 
			song_text = $4, link = $5
		WHERE
		   id = $6;`,
		p.Group, p.Name, release, p.Text, p.Link, id,
	)
	if err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo update")
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return models.Song{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo update")
	}
	if updated != 1 {
		return models.Song{}, internal.NewErrorf(internal.ErrorCodeNotFound, "resourse with id %d not found", id)
	}

	r.logger.Debug("record updated", "id", id)

	return models.Song{
		ID:          id,
		Group:       p.Group,
		Name:        p.Name,
		ReleaseDate: release,
		Text:        p.Text,
		Link:        p.Link,
	}, nil
}

func (r *SongRepository) SelectText(id int32) (string, error) {
	var text string
	if err := r.db.QueryRow(
		`SELECT 
		    song_text from public.songs 
		WHERE 
		    id = $1;`,
		id,
	).Scan(&text); err != nil {
		return "", internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo select")
	}

	r.logger.Debug("text selected", "id", id)

	return text, nil
}

func (r *SongRepository) Search(vals url.Values, pageNum, perPage int) ([]models.Song, error) {
	songs := make([]models.Song, 0)

	offset := strconv.Itoa(pageNum * perPage)
	limit := strconv.Itoa(perPage)
	fields := []string{"group_name", "song_name", "release_date", "song_text", "link"}
	query, err := NewQuery(fields, "SELECT * FROM public.songs", limit, offset, vals)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo search")
	}
	q := query.GetQuery()
	// fmt.Println(q)
	r.logger.Debug("Search", "query", q)

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo search")
	}
	defer rows.Close()

	var s models.Song
	for rows.Next() {
		if err := rows.Scan(
			&s.ID,
			&s.Group,
			&s.Name,
			&s.ReleaseDate,
			&s.Text,
			&s.Link,
		); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	if len(songs) == 0 {
		return nil, internal.NewErrorf(internal.ErrorCodeNotFound, "no records found")
	}

	r.logger.Debug("records selected", "count", len(songs))

	return songs, nil
}

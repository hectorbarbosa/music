package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"music/internal"
	"music/internal/app/models"
	"music/internal/config"
	m "music/internal/rest/models"
)

type SongService interface {
	Create(params m.CreateParams) (models.Song, error)
	Delete(id int32) error
	Update(id int32, f m.UpdateParams) (models.Song, error)
	SelectVerse(id int32, v int) (string, error)
	Search(params url.Values, pageNum, perPage int) ([]models.Song, error)
}

type SongHandler struct {
	cfg    config.Config
	logger *slog.Logger
	svc    SongService
}

func NewSongHandler(cfg config.Config, logger *slog.Logger, svc SongService) *SongHandler {
	return &SongHandler{
		cfg:    cfg,
		logger: logger,
		svc:    svc,
	}
}

func (h *SongHandler) Register(r *mux.Router) {
	r.HandleFunc("/songs", h.create).Methods(http.MethodPost)
	r.HandleFunc("/songs/{id}", h.update).Methods(http.MethodPut)
	r.HandleFunc("/songs/{id}", h.delete).Methods(http.MethodDelete)
	r.HandleFunc("/songs/{id}/verse/{vid}", h.getVerse).Methods(http.MethodGet)
	r.HandleFunc("/songs/page/{page_num}/records/{per_page}", h.search).Methods(http.MethodGet)
}

//	@Tags Фонотека
//
// @Description Create new record
// @Accept		json
// @Produce		json
// @Param		json	body		m.SongDetails	true	    "input data"
// @Success		201		{object}	models.Song			        "Created"
// @Failure		400		{object}	rest.ErrorResponse	        "Bad request"
// @Failure		500		{object}	rest.ErrorResponse	        "Internal error"
// @Failure		502		{object}	rest.ErrorResponse	        "Bad Gateway"
// @Router		/songs [post]
func (h *SongHandler) create(w http.ResponseWriter, r *http.Request) {
	var sd m.SongDetails
	if err := json.NewDecoder(r.Body).Decode(&sd); err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid song details params")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	defer r.Body.Close()

	if err := sd.Validate(); err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid song details params")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	createParams, err := fetchDetails(h.cfg, h.logger, sd)
	// fmt.Println("createSong: %v", createParams)
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeBadGateWay, "create failed")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	song, err := h.svc.Create(createParams)
	if err != nil {
		msg := fmt.Errorf("create failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	h.logger.Info("POST request success, record created", "id", song.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

//	@Tags Фонотека
//
// @Description Delete record
// @Accept		json
// @Produce		json
// @Param		id		path		int		            true	    "Song ID"
// @Success		200		{object}	nil      			"ok"
// @Failure		400		{object}	rest.ErrorResponse	"Bad request"
// @Failure		404		{object}	rest.ErrorResponse	"Not found"
// @Failure		500		{object}	rest.ErrorResponse	"Internal error"
// @Router		/songs/{id} [delete]
func (h *SongHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid id")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	if err := h.svc.Delete(int32(id)); err != nil {
		msg := fmt.Errorf("delete failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	h.logger.Info("DELETE request success, record deleted", "id", id)
	w.WriteHeader(http.StatusOK)
}

//	@Tags Фонотека
//
// @Description Update record
// @Accept		json
// @Produce		json
// @Param		id		path		int		            true    	    "Song ID"
// @Param		json	body		m.UpdateParams  	true	        "input data"
// @Success		200		{object}	models.Song			"ok"
// @Failure		400		{object}	rest.ErrorResponse	"Bad request"
// @Failure		404		{object}	rest.ErrorResponse	"Not found"
// @Failure		500		{object}	rest.ErrorResponse	"Internal error"
// @Router		/songs/{id} [put]
func (h *SongHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid id")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	var updateParams m.UpdateParams
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid update params")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	defer r.Body.Close()

	song, err := h.svc.Update(int32(id), updateParams)
	if err != nil {
		msg := fmt.Errorf("update failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	h.logger.Info("PUT request success, record updated", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}

//	@Tags Фонотека
//
// @Description Получить куплет песни
// @Param		id		path		int		true	    "Song ID"
// @Param		vid		path		int		true	    "Номер куплета, начиная с 1"
// @Accept		json
// @Produce		json
// @Success		200		{object}	m.Verse	            "ok"
// @Failure		400		{object}	rest.ErrorResponse	"Bad request"
// @Failure		404		{object}	rest.ErrorResponse	"Not found"
// @Failure		500		{object}	rest.ErrorResponse	"Internal error"
// @Router		/songs/{id}/verse/{vid}  [get]
func (h *SongHandler) getVerse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid id")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	vid, err := strconv.Atoi(mux.Vars(r)["vid"])
	if err != nil {
		// msg := fmt.Errorf("getVerse failed: %w", err)
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid vid")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	v, err := h.svc.SelectVerse(int32(id), vid)
	if err != nil {
		msg := fmt.Errorf("getVerse failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	h.logger.Info("GET request success, text selected", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m.Verse{Num: strconv.Itoa(vid), Text: v})
}

func fetchDetails(cfg config.Config, logger *slog.Logger, sd m.SongDetails) (m.CreateParams, error) {
	q := make(url.Values)
	q.Add("group", sd.Group)
	q.Add("song", sd.Name)
	url := url.URL{
		Scheme:   "http",
		Host:     cfg.ApiAddr,
		Path:     cfg.ApiPath,
		RawQuery: q.Encode(),
	}
	logger.Debug("remote api request", "url", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return m.CreateParams{}, fmt.Errorf("NewRequest error: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return m.CreateParams{}, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	logger.Debug("remote api request response", "status", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return m.CreateParams{}, fmt.Errorf("remote response status is: %s", resp.Status)
	}

	var p m.CreateParams
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return m.CreateParams{}, fmt.Errorf("json decoder error: %w", err)
	}

	p.Group = sd.Group
	p.Name = sd.Name
	return p, nil
}

//	@Tags Фонотека
//
// @Description Поиск по фонотеке
// @Param		page_num		path		int		true	    "Page number from 0"
// @Param		per_page		path		int		true	    "Records per page"
// @Param		group_name		query		string	false	    "Searching group name"
// @Param		song_name		query		string	false	    "Song name"
// @Param		release_date	query		string	false	    "Release date (example 17.06.2006)"
// @Param		song_text		query		string	false	    "Song text"
// @Param		link    		query		string	false	    "Link"
// @Accept		json
// @Produce		json
// @Success		200		{object}	[]models.Song	            "ok"
// @Failure		400		{object}	rest.ErrorResponse      	"Bad request"
// @Failure		404		{object}	rest.ErrorResponse  	    "Not found"
// @Failure		500		{object}	rest.ErrorResponse      	"Internal error"
// @Router		/songs/page/{page_num}/records/{per_page}  [get]
func (h *SongHandler) search(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(mux.Vars(r)["page_num"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid page_num")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	perPage, err := strconv.Atoi(mux.Vars(r)["per_page"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid per_page")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid url query")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	songs, err := h.svc.Search(m, pageNum, perPage)
	if err != nil {
		msg := fmt.Errorf("search failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	h.logger.Info("GET request success, records found", "number", len(songs))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(songs)
}

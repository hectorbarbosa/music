package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"music/api"
	"music/internal/config"
)

type songsService struct {
	songs map[api.InfoGetParams]api.SongDetail
	mux   sync.Mutex
}

func (s *songsService) InfoGet(ctx context.Context, req api.InfoGetParams) (api.InfoGetRes, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// fmt.Printf("request: %v", req)

	if res, ok := s.songs[req]; ok {
		return &res, nil
	}

	return &api.InfoGetBadRequest{}, nil
}

func main() {
	// Find path for env file
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// config file must be in project root dir, compiled bin must be in /bin dir!!!
	configPath, err := filepath.Abs(dir + "/../.env")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Config path:", configPath)

	// Create config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create songs map.
	songs := make(map[api.InfoGetParams]api.SongDetail)
	song1 := api.InfoGetParams{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	details1 := api.SongDetail{
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	songs[song1] = details1

	song2 := api.InfoGetParams{
		Group: "Group2",
		Song:  "Song2",
	}
	details2 := api.SongDetail{
		ReleaseDate: "27.02.2015",
		Text:        "Text2\n\nText22",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	songs[song2] = details2
	song3 := api.InfoGetParams{
		Group: "Group3",
		Song:  "Song3",
	}
	details3 := api.SongDetail{
		ReleaseDate: "20.01.2010",
		Text:        "Text3\n\nText33",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	songs[song3] = details3

	service := &songsService{
		songs: songs,
	}
	// Create generated server.
	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(cfg.ApiAddr, srv); err != nil {
		log.Fatal(err)
	}
}

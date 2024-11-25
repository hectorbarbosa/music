package main

import (
	"database/sql"
	"errors"
	"fmt"
	"music/internal/app/service"
	"music/internal/config"
	"music/internal/logging"
	"music/internal/rest"
	"music/internal/storage/postgresql"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "music/docs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			Swagger Songs
//	@version		1.0
//	@description	This is a sample music service.
//	@schemes        http
//	@host           localhost:8080

func main() {
	// Find path for env file
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, "abs path", err)
		os.Exit(1)
	}

	// config file must be in project root dir, compiled bin must be in /bin dir!!!
	configPath, err := filepath.Abs(dir + "/../.env")
	if err != nil {
		fmt.Fprintln(os.Stderr, "composing path", err)
		os.Exit(1)
	}
	// fmt.Println("Config path:", configPath)

	// Create config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading config", err)
		os.Exit(1)
	}

	// Create db connection
	db, err := newDB(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error connecting db", err)
		os.Exit(1)
	}

	// Make Up migration
	if err := makeMigration(db); err != nil {
		fmt.Fprintln(os.Stderr, "migration error", err)
		os.Exit(1)
	}

	// Create logger
	logger, err := logging.GetLogger(cfg.LogLevel)
	if err != nil {
		fmt.Printf("getting logger: %s\n", err)
		os.Exit(1)
	}
	// fmt.Println("Log Level:", cfg.LogLevel)

	r := mux.NewRouter()
	repo := postgresql.NewSongRepo(db, logger)
	svc := service.NewSongService(*cfg, logger, repo)
	rest.NewSongHandler(*cfg, logger, svc).Register(r)

	swagUrl := "./docs/doc.json"

	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swagUrl), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	server := http.Server{
		Handler:           r,
		Addr:              cfg.ServerAddr,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
	logger.Info("Server start", "listening on address:", cfg.ServerAddr)
	server.ListenAndServe()
}

func newDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func makeMigration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations/",
		"music",
		driver,
	)
	if err != nil {
		return err
	}

	if _, _, err := m.Version(); err == nil {
		// Migration has been made already
		return nil
	} else if errors.Is(err, migrate.ErrNilVersion) {
		// No migrations has been applied, make migration
		if err := m.Up(); err != nil {
			return err
		}
		return nil
	} else {
		// Some error
		return err
	}
}

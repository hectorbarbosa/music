package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DbUrl      string `env:"DB_URL"`
	ServerAddr string `env:"SERVER_ADDR"`
	ApiAddr    string `env:"API_ADDR"`
	LogLevel   int    `env:"LOG_LEVEL"`
	ApiPath    string `env:"API_PATH"`
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

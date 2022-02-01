package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres   PostgresConfig
	HttpServer HttpServerConfig
	Jwt        JwtConfig
}

type PostgresConfig struct {
	User     string `env:"DATABASE_USER" env-default:"postgres"`
	Password string `env:"DATABASE_PASSWORD" env-default:"good-password"`
	Host     string `env:"DATABASE_HOST" env-default:"localhost"`
	Port     string `env:"DATABASE_PORT" env-default:"5433"`
	Instance string `env:"DATABASE_NAME" env-default:"suricate"`
}

type HttpServerConfig struct {
	Host string `env:"SERVER_HOST" env-default:"localhost"`
	Port string `env:"SERVER_PORT" env-default:"3000"`
}

type JwtConfig struct {
	Secret string `env:"JWT_SECRET" env-default:"BAD_SECRET"`
}

func ReadConfigFromEnv() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Panicf("failed to load config: %s", err)
	}

	return &cfg
}

func ReadConfigFromFile(filename string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		log.Panicf("failed to load config: %s", err)
	}

	return &cfg
}

func ReadConfig(filename string) *Config {
	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		log.Printf("%s file not found, attempting to load from env", filename)
		return ReadConfigFromEnv()
	}

	return ReadConfigFromFile(filename)
}

func (cfg PostgresConfig) Url() string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Instance,
	)
	return url
}

func (cfg HttpServerConfig) ServerPort() string {
	return fmt.Sprintf(":%s", cfg.Port)
}

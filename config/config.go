package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Log        ConfLog
	Postgres   PostgresConfig
	HTTPServer HTTPServerConfig
	Jwt        JwtConfig
	Dockertest DockertestConfig
	Redis      RedisConfig
}

type ConfLog struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

type PostgresConfig struct {
	User     string `env:"DATABASE_USER" env-default:"postgres"`
	Password string `env:"DATABASE_PASSWORD" env-default:"good-password"`
	Host     string `env:"DATABASE_HOST" env-default:"localhost"`
	Port     string `env:"DATABASE_PORT" env-default:"5432"`
	Instance string `env:"DATABASE_NAME" env-default:"suricate"`
}

type HTTPServerConfig struct {
	Host string `env:"SERVER_HOST" env-default:"localhost"`
	Port string `env:"SERVER_PORT" env-default:"8080"`
}

type JwtConfig struct {
	Secret  string `env:"JWT_SECRET" env-default:"BAD_SECRET"`
	Timeout string `env:"JWT_TIMEOUT" env-default:"30"`
}

type DockertestConfig struct {
	Timeout string `env:"DOCKERTEST_TIMEOUT" env-default:"30"`
}

type RedisConfig struct {
	Host              string `env:"REDIS_HOST" env-default:"localhost"`
	Port              string `env:"REDIS_PORT" env-default:"6379"`
	IdempotencyKeyTTL string `env:"IDEMPOTENCY_TTL" env-default:"86400"`
}

func ReadConfigFromEnv() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}

	cfg.setEnvs()
	log.Info().Msg("successfully loaded from env")

	return &cfg
}

func ReadConfigFromFile(filename string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}

	cfg.setEnvs()
	log.Info().Msg("successfully loaded env variables from .env file")

	return &cfg
}

func ReadConfig(filename string) *Config {
	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		log.Info().Msg(" file not found, attempting to load from env")

		return ReadConfigFromEnv()
	}

	return ReadConfigFromFile(filename)
}

func (cfg PostgresConfig) URL() string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Instance,
	)

	return url
}

func (cfg HTTPServerConfig) ServerPort() string {
	return fmt.Sprintf(":%s", cfg.Port)
}

func (cfg HTTPServerConfig) HostAndPort() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func (cfg RedisConfig) URL() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func (cfg Config) setEnvs() {
	os.Setenv("JWT_SECRET", cfg.Jwt.Secret)
	os.Setenv("JWT_TIMEOUT", cfg.Jwt.Timeout)
	os.Setenv("DOCKERTEST_TIMEOUT", cfg.Dockertest.Timeout)
	os.Setenv("IDEMPOTENCY_TTL", cfg.Redis.IdempotencyKeyTTL)
}

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ConfLog
	PostgresConfig
	HTTPServerConfig
	JwtConfig
	DockertestConfig
	RedisConfig
}

type ConfLog struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

type PostgresConfig struct {
	PgUser     string `envconfig:"DATABASE_USER" default:"postgres"`
	PgPassword string `envconfig:"DATABASE_PASSWORD" default:"good-password"`
	PgHost     string `envconfig:"DATABASE_HOST" default:"localhost"`
	PgPort     string `envconfig:"DATABASE_PORT" default:"5432"`
	PgDatabase string `envconfig:"DATABASE_NAME" default:"suricate"`
	PgURL      string `envconfig:"DATABASE_URL"`
}

type HTTPServerConfig struct {
	HTTPHost string `envconfig:"HOST" default:"localhost"`
	HTTPPort string `envconfig:"PORT" default:"8080"`
}

type JwtConfig struct {
	JWTSecret  string `envconfig:"JWT_SECRET" default:"BAD_SECRET"`
	JWTTimeout int    `envconfig:"JWT_TIMEOUT" default:"30"`
}

type DockertestConfig struct {
	DockertestTimeout int `envconfig:"DOCKERTEST_TIMEOUT" default:"30"`
}

type RedisConfig struct {
	RedisHost         string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort         string `envconfig:"REDIS_PORT" default:"6379"`
	IdempotencyKeyTTL int    `envconfig:"IDEMPOTENCY_TTL" default:"86400"`
	RedisURL          string `envconfig:"REDIS_URL"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	var config Config
	noPrefix := ""
	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed loading config: %w", err)
	}
	log.Info().Msgf("HTTP server address: %s%s", config.GetHTTPHost(), config.GetHTTPPort())

	return config, nil
}

func (cfg PostgresConfig) GetPgURL() string {
	if cfg.PgURL != "" {
		return cfg.PgURL
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDatabase,
	)
}

func (cfg HTTPServerConfig) GetHTTPPort() string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return fmt.Sprintf(":%s", envPort)
	}

	return cfg.HTTPPort
}

func (cfg HTTPServerConfig) GetHTTPHost() string {
	if envHost := os.Getenv("HOST"); envHost != "" {
		return envHost
	}

	return cfg.HTTPHost
}

func (cfg RedisConfig) GetRedisURL() string {
	if cfg.RedisURL != "" {
		return cfg.RedisURL
	}

	return fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
}

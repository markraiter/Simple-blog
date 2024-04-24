package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"development"`
	Server
	Postgres
	Auth
}

type Postgres struct {
	Driver   string `env:"POSTGRES_DRIVER" env-default:"postgres"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string `env:"POSTGRES_DB" env-default:"blogdb"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`
}

type Server struct {
	Port         string        `env:"PORT" env-default:"9000"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" env-default:"120s"`
}

type Auth struct {
	SigningKey string        `env:"SIGNING_KEY" env-required:"true"`
	AccessTTL  time.Duration `env:"ACCESS_TTL" env-default:"1h"`
	RefreshTTL time.Duration `env:"REFRESH_TTL" env-default:"720h"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("error loading env variables: " + err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	cfg.Server.Port = ":" + cfg.Server.Port

	return &cfg
}

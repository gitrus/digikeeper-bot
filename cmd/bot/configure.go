package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `env:"ENVIRONMENT_NAME" env-default:"dev"`
	CommonTimeout  time.Duration `env:"COMMON_TIMEOUT" env-default:"5s"`
	LocalPort      string        `env:"LOCAL_PORT" env-default:"9000"`
	LocalHost      string        `env:"LOCAL_HOST" env-default:"localhost"`
	BotKey         string        `env:"TELEGRAM_BOT_TOKEN" env-default:""`
	BotPublicURL   string        `env:"TELEGRAM_BOT_PUBLIC_URL" env-default:"localhost"`
	AllowedUpdates []string      `env:"TELEGRAM_ALLOWED_UPDATES" env-default:"message"`
}

func (c *Config) IsDevEnv() bool {
	return strings.HasPrefix(strings.ToLower(c.Env), "dev")
}

func configure() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to fetch env vars: %v", err)
	}

	err = configureLogger(cfg.Env)
	if err != nil {
		log.Fatalf("Failed to configure logger: %v", err)
	}

	if cfg.IsDevEnv() {
		absPath, err := filepath.Abs("../../.env")
		if err != nil {
			log.Fatalf("Failed to get .env: %v", err)
		}
		err = cleanenv.ReadConfig(absPath, &cfg)
		if err != nil {
			log.Fatalf("Failed to get .env: %v", err)
		}
	}

	return cfg
}

func configureLogger(environ string) error {
	err := loggingctx.InitLogger(environ)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	return nil
}

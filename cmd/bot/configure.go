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
	Env           string        `env:"ENVIRONMENT_NAME" env-default:"dev"`
	BotKey        string        `env:"TELEGRAM_BOT_TOKEN" env-default:""`
	CommonTimeout time.Duration `env:"COMMON_TIMEOUT" env-default:"5s"`
}

func Configure() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to fetch env vars: %v", err)
	}

	err = configureLogger(cfg.Env)
	if err != nil {
		log.Fatalf("Failed to configure logger: %v", err)
	}

	if strings.HasPrefix(strings.ToLower(cfg.Env), "dev") {
		absPath, err := filepath.Abs("../../.env")
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

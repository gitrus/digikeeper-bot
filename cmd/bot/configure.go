package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/joho/godotenv"
)

type Config struct {
	BotKey        string
	CommonTimeout time.Duration
}

func Configure() Config {
	environ := "dev"
	err := configureLogger(environ)
	if err != nil {
		log.Fatalf("Failed to configure logger: %v", err)
	}

	if environ == "dev" {
		absPath, err := filepath.Abs("../../.env")
		err = godotenv.Load(absPath)
		if err != nil {
			log.Fatalf("Failed to get .env: %v", err)
		}
	}

	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
	}
	commonTimeout := 5 * time.Second

	return Config{
		BotKey:        tgToken,
		CommonTimeout: commonTimeout,
	}
}

func configureLogger(environ string) error {
	err := loggingctx.InitLogger(environ)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	return nil
}

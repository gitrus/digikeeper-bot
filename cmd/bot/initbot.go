package main

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
)

func initBot(ctx context.Context, cfg Config) (*telego.Bot, <-chan telego.Update, error) {
	bot, err := telego.NewBot(cfg.BotKey, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	if !cfg.IsDevEnv() {
		upd, err := initWebHookBot(ctx, cfg, bot)
		return bot, upd, err
	} else {
		upd, err := initPollingBot(ctx, cfg, bot)
		return bot, upd, err
	}
}

func initWebHookBot(ctx context.Context, cfg Config, bot *telego.Bot) (<-chan telego.Update, error) {
	err := bot.SetWebhook(ctx, &telego.SetWebhookParams{
		URL:            cfg.BotPublicURL,
		SecretToken:    bot.Token(),
		AllowedUpdates: cfg.AllowedUpdates,
	})
	if err != nil {
		return nil, err
	}

	srv := &fasthttp.Server{}

	updates, err := bot.UpdatesViaWebhook(
		ctx,
		telego.WebhookFastHTTP(
			srv, "/bot", "secret_digi",
		),
	)

	return updates, err

}

func initPollingBot(ctx context.Context, cfg Config, bot *telego.Bot) (<-chan telego.Update, error) {
	updates, err := bot.UpdatesViaLongPolling(
		ctx,
		&telego.GetUpdatesParams{AllowedUpdates: cfg.AllowedUpdates},
	)
	return updates, err
}

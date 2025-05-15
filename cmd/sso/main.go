package main

import (
	"log/slog"
	"os"
	"sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()
	// TODO: инициализировать логгер
	log := setupLoger(cfg.Env)
	log.Info(
		"starting application",
		slog.Any("cfg", cfg),
	)
	// TODO: инициализировать приложение (app)
	app := app.New(log, cfg.GRPCConfig.Port, cfg.StoragePath, cfg.GRPCConfig.Timeout)
	err := app.GRPCServer.Run()
	if err != nil {
		log.Error(err.Error())
		panic("init failed")
	}
	// TODO: инициализировать сервер gRPC

}

func setupLoger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

package main

import (
	"go-api/internal/auth"
	"go-api/internal/config"
	handlerAuth "go-api/internal/handlers/auth"
	"go-api/internal/storage"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad() //init config cleanenv

	logger := setupLogger(cfg.Env) //init logger slog
	logger.Info("starting", slog.String("env", cfg.Env))

	store, err := storage.NewDB(cfg.DB, logger) // init storage postgresql
	auth.Init(cfg.JWT.SecretKey)                //init secret key

	if err != nil {
		logger.Error("DB init failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer store.Close()

	http.HandleFunc("/api/login", handlerAuth.LoginHandler(store.DB(), logger)) //временно, потом заменить
	http.ListenAndServe(":8080", nil)

	// TODO: init router chi

	// TODO: run server
}

// константы логгера
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// создание логгера
func setupLogger(env string) *slog.Logger {
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
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

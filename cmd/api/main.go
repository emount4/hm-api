package main

import (
	"go-api/internal/auth"
	"go-api/internal/config"
	handlerAuth "go-api/internal/handlers/auth"
	"go-api/internal/storage"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// http.HandleFunc("/api/login", handlerAuth.LoginHandler(store.DB(), logger)) //временно, потом заменить
	// http.HandleFunc("/api/register", handlerAuth.RegisterHandler(store.DB(), logger))
	// http.ListenAndServe(":8080", nil)

	router := chi.NewRouter() // init router chi

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Post("/api/login", handlerAuth.LoginHandler(store.DB(), logger))
	router.Post("/api/register", handlerAuth.RegisterHandler(store.DB(), logger))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	logger.Info("server started", slog.String("port", ":8080"))
	http.ListenAndServe(":8080", router)

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

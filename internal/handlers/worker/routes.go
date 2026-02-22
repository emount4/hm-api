package worker

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	r.Route("/handyman", func(r chi.Router) {
		r.Get("/", AllWorkersHandler(db, logger))
	})
}

package ads

import (
	"go-api/internal/middleware"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	public := chi.NewRouter()
	protected := chi.NewRouter()

	//  ПУБЛИЧНЫЕ (мастера смотрят без токена)
	public.Get("/", AdsHandler(db, logger))       // GET /api/v1/ads - список
	public.Get("/{adID}", AdsHandler(db, logger)) // GET /api/v1/ads/123

	//  ЗАЩИЩЁННЫЕ (клиент управляет)
	protected.Use(middleware.AuthMiddleware(logger))
	protected.Post("/", AdsHandler(db, logger))
	protected.Patch("/{adID}", AdsHandler(db, logger))
	protected.Delete("/{adID}", AdsHandler(db, logger))

	r.Mount("/ads", public)       // /api/v1/ads → публичные
	r.Mount("/my-ads", protected) // /api/v1/my-ads → личный кабинет
}

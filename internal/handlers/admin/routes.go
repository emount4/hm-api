package admin

import (
	"go-api/internal/middleware"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	admin := chi.NewRouter()

	// Защита: требуется аутентификация + роль администратора
	admin.Use(middleware.AuthMiddleware(logger))
	admin.Use(middleware.AdminMiddleware(db, logger))

	// Управление пользователями
	admin.Get("/users", GetUsersHandler(db, logger))                       // GET /admin/users - список пользователей
	admin.Get("/users/{userID}", GetUserHandler(db, logger))               // GET /admin/users/123 - пользователь по ID
	admin.Delete("/users/{userID}", DeleteUserHandler(db, logger))         // DELETE /admin/users/123 - удалить пользователя
	admin.Patch("/users/{userID}/role", UpdateUserRoleHandler(db, logger)) // PATCH /admin/users/123/role - изменить роль

	// Модерация объявлений
	admin.Get("/ads", GetAllAdsHandler(db, logger))          // GET /admin/ads - все объявления
	admin.Delete("/ads/{adID}", DeleteAdHandler(db, logger)) // DELETE /admin/ads/123 - удалить объявление

	// Модерация откликов
	admin.Get("/responses", GetAllResponsesHandler(db, logger))                // GET /admin/responses - все отклики
	admin.Delete("/responses/{responseID}", DeleteResponseHandler(db, logger)) // DELETE /admin/responses/123 - удалить отклик

	// Статистика
	admin.Get("/stats", GetStatsHandler(db, logger)) // GET /admin/stats - общая статистика

	// Черный список
	admin.Get("/blacklist", GetBlacklistHandler(db, logger))                   // GET /admin/blacklist
	admin.Post("/blacklist", AddToBlacklistHandler(db, logger))                // POST /admin/blacklist
	admin.Delete("/blacklist/{email}", RemoveFromBlacklistHandler(db, logger)) // DELETE /admin/blacklist/email@example.com

	r.Mount("/admin", admin)
}

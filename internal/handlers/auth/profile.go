package auth

import (
	"encoding/json"
	"errors"
	"go-api/internal/models"
	"go-api/internal/storage"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func ProfileHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProfile(db, logger)(w, r)
		case http.MethodPatch:
			editProfile(db, logger)(w, r)
		default:
			http.Error(w, "Method not allows", http.StatusMethodNotAllowed)
			logger.Error("Ошибка метода в хендлера роутера", "Метод", r.Method)
			return
		}
	}

}

func getProfile(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		user, err := storage.UserById(db, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Role.RoleName == "worker" {
			// Загружаем Worker данные
			var worker models.WorkerProfile
			err := db.Where("user_id = ?", userID).First(&worker).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role.RoleName,
				"name":  user.Name,
				"worker": map[string]interface{}{
					"specialization": worker.Categories,
					"experience":     worker.ExpYears,
					// "rating":         worker.Rating,
					"phone":       worker.Phone,
					"description": worker.Description,
					"is_busy":     worker.IsBusy,
					"reviews":     worker.ReviewsReceived,
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role.RoleName,
			"name":  user.Name,
		})
	}
}

func editProfile(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

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
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		type ProfileInput struct {
			// User поля
			Name  *string `json:"name,omitempty"`
			Phone *string `json:"phone,omitempty"`

			// WorkerProfile поля (только если роль worker)
			ExpYears    *int    `json:"exp_years,omitempty"`
			Description *string `json:"description,omitempty"`
			IsBusy      *bool   `json:"is_busy,omitempty"`
		}

		var input ProfileInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		userUpdates := map[string]interface{}{}
		if input.Name != nil {
			userUpdates["name"] = *input.Name
		}

		if len(userUpdates) > 0 {
			result := tx.Model(&models.User{}).Where("id = ?", userID).Updates(userUpdates)
			if result.Error != nil {
				tx.Rollback()
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if result.RowsAffected == 0 {
				tx.Rollback()
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
		}

		var user models.User
		if err := tx.Preload("Role").First(&user, userID).Error; err != nil {
			tx.Rollback()
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Role.RoleName == "worker" {
			workerUpdates := map[string]interface{}{}
			if input.ExpYears != nil {
				workerUpdates["exp_years"] = *input.ExpYears
			}
			if input.Description != nil {
				workerUpdates["description"] = *input.Description
			}
			if input.IsBusy != nil {
				workerUpdates["is_busy"] = *input.IsBusy
			}
			if input.Phone != nil {
				workerUpdates["phone"] = *input.Phone
			}
			if len(workerUpdates) > 0 {
				result := tx.Model(&models.WorkerProfile{}).Where("user_id = ?", userID).Updates(workerUpdates)
				if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					tx.Rollback()
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}
			}
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		var updatedUser models.User
		if err := db.Preload("Role").Preload("WorkerProfile").First(&updatedUser, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     updatedUser.ID,
			"email":  updatedUser.Email,
			"name":   updatedUser.Name,
			"role":   updatedUser.Role.RoleName,
			"worker": updatedUser.WorkerProfile != nil, // Есть ли профиль worker
		})
	}
}

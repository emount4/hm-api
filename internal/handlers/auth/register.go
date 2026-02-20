package auth

import (
	"encoding/json"
	"go-api/internal/auth"
	"go-api/internal/models"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Error("Неправильный метод",
				"method", r.Method,
				"path", r.URL.Path)
			http.Error(w, "Method not allowes", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
			Role     uint   `json:"role"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Ошибка парсинга JSON",
				"req body", r.Body)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if len(input.Email) == 0 || len(input.Name) == 0 {
			logger.Error("Пустые поля ввода",
				"email", input.Email,
				"name", input.Name,
			)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		}

		if userExists(db, input.Email) {
			logger.Info("Пользователь с таким email уже существует", "email", input.Email)
			http.Error(w, "User already exist", http.StatusConflict)
			return
		}

		passwordHash, err := auth.HashPassword(input.Password)

		if err != nil {
			logger.Error("Ошибка генерации хэша", "err:", err)
			http.Error(w, "Ошибка сервера при создании пароля", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Name:         input.Name,
			Email:        input.Email,
			PasswordHash: passwordHash,
			RoleID:       input.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			logger.Error("Ошибка вставки пользователя", "err:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, logger)
		if err != nil {
			http.Error(w, "Ошибка токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"message": "Пользователь зарегистрирован",
			"token":   token,
		})
	}
}

func userExists(db *gorm.DB, email string) bool {
	var user models.User
	result := db.Where("email = ?", email).First(&user)

	return result.Error == nil
}

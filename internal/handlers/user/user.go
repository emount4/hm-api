package user

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func UserHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet{

		}
	}

}

//Смена роли, бан пользователя
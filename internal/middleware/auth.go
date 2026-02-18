package middleware

import (
	"context"
	"go-api/internal/auth"
	"log/slog"
	"net/http"
)

func AuthMiddlware(next http.HandlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Невалидный заголовок, пустой",
				slog.String("Header", authHeader),
			)
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			logger.Error("Невалидный заголовок, начало не Bearer",
				slog.String("Header", authHeader),
			)
			http.Error(w, "Invalid header format", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]

		claims, err := auth.ValidateToken(tokenString, logger)

		if err != nil {
			logger.Error("Невалидный токен", slog.String("Токен:", tokenString))
			http.Error(w, "Invalid token"+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(r.Context(), "email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

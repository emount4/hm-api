package storage

import (
	"go-api/internal/models"

	"gorm.io/gorm"
)

func UserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).Preload("Role").First(&user)
	return &user, result.Error
}

func UserById(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	result := db.Where("id = ?", id).Preload("Role").First(&user)
	return &user, result.Error
}

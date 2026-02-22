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

type WorkerResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       *string `json:"phone,omitempty"`
	ExpYears    *int    `json:"exp_years,omitempty"`
	Description *string `json:"description,omitempty"`
	IsBusy      bool    `json:"is_busy"`
}

func ListApprovedWorkers(db *gorm.DB, limit, offset int) ([]WorkerResponse, int64, error) {
	var total int64

	db.Model(&models.User{}).
		Joins("JOIN worker_profiles ON worker_profiles.user_id = users.id").
		Where("users.role_id = ?", 2).
		Count(&total)

	var workers []WorkerResponse
	db.Table("users u").
		Select("u.id, u.name, u.email, wp.phone, wp.exp_years, wp.description, wp.is_busy").
		Joins("JOIN worker_profiles wp ON u.id = wp.user_id").
		Where("u.role_id = ?", 2).
		Order("u.id ASC").
		Offset(offset).
		Limit(limit).
		Scan(&workers)

	return workers, total, nil
}

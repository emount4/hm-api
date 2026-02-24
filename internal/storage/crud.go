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

type CategoryJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type WorkerResponse struct {
	ID          uint           `json:"id"`                                // user_id == worker_id
	WorkerID    uint           `json:"worker_id" gorm:"column:worker_id"` // дублируем для явного worker_id
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone,omitempty"`
	ExpYears    *int           `json:"exp_years,omitempty"`
	Description *string        `json:"description,omitempty"`
	IsBusy      bool           `json:"is_busy"`
	Categories  []CategoryJSON `json:"categories,omitempty" gorm:"-"`
}

func ListApprovedWorkers(db *gorm.DB, limit, offset int) ([]WorkerResponse, int64, error) {
	var total int64

	db.Model(&models.User{}).
		Joins("JOIN worker_profiles ON worker_profiles.user_id = users.id").
		Where("users.role_id = ?", 2).
		Count(&total)

	var workers []WorkerResponse
	err := db.Table("users u").
		Select("u.id, u.id as worker_id, u.name, u.email, u.phone, wp.exp_years, wp.description, wp.is_busy").
		Joins("JOIN worker_profiles wp ON u.id = wp.user_id").
		Where("u.role_id = ?", 2).
		Order("u.id ASC").
		Offset(offset).
		Limit(limit).
		Scan(&workers).Error

	if err != nil {
		return nil, 0, err
	}

	for i := range workers {
		var categories []models.Category
		db.Table("categories c").
			Joins("JOIN worker_categories wc ON c.id = wc.category_id").
			Where("wc.worker_id = ?", workers[i].ID).
			Find(&categories)

		catJSON := make([]CategoryJSON, len(categories))
		for j, cat := range categories {
			catJSON[j] = CategoryJSON{
				ID:   cat.ID,
				Name: cat.Name,
			}
		}
		workers[i].Categories = catJSON
	}

	return workers, total, nil
}

func WorkerByID(db *gorm.DB, id uint) (*WorkerResponse, error) {
	var worker WorkerResponse

	err := db.Table("users u").
		Select("u.id, u.id as worker_id, u.name, u.email, u.phone, wp.exp_years, wp.description, wp.is_busy").
		Joins("JOIN worker_profiles wp ON u.id = wp.user_id").
		Where("u.id = ? AND u.role_id = ?", id, 2).
		Scan(&worker).Error

	if err != nil || worker.ID == 0 {
		return nil, err
	}

	var categories []models.Category
	db.Table("categories c").
		Joins("JOIN worker_categories wc ON c.id = wc.category_id").
		Where("wc.worker_id = ?", worker.ID).
		Find(&categories)

	catJSON := make([]CategoryJSON, len(categories))
	for j, cat := range categories {
		catJSON[j] = CategoryJSON{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}
	worker.Categories = catJSON

	return &worker, nil
}

func WorkerByUserID(db *gorm.DB, id uint) (*WorkerResponse, error) {
	var worker WorkerResponse

	err := db.Table("users u").
		Select("u.id, u.id as worker_id, u.name, u.email, u.phone, wp.exp_years, wp.description, wp.is_busy").
		Joins("JOIN worker_profiles wp ON u.id = wp.user_id").
		Where("u.id = ? AND u.role_id = ?", id, 2).
		Scan(&worker).Error

	if err != nil || worker.ID == 0 {
		return nil, err
	}

	var categories []models.Category
	db.Table("categories c").
		Joins("JOIN worker_categories wc ON c.id = wc.category_id").
		Where("wc.worker_id = ?", worker.ID).
		Find(&categories)

	catJSON := make([]CategoryJSON, len(categories))
	for j, cat := range categories {
		catJSON[j] = CategoryJSON{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}
	worker.Categories = catJSON

	return &worker, nil
}

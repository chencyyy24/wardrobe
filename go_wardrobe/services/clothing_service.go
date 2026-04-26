package services

import (
	"go_wardrobe/database"
	"go_wardrobe/models"
	"gorm.io/gorm"
)

type ClothingService struct{}

func NewClothingService() *ClothingService {
	return &ClothingService{}
}

func (s *ClothingService) Create(item *models.ClothingItem) error {
	return database.DB.Create(item).Error
}

func (s *ClothingService) GetByID(id string) (*models.ClothingItem, error) {
	var item models.ClothingItem
	err := database.DB.Preload("Category").Preload("Subcategory").
		Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ClothingService) List(categoryID string, userID uint) ([]models.ClothingItem, error) {
	query := database.DB.Where("user_id = ?", userID)
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	var items []models.ClothingItem
	err := query.Preload("Category").Preload("Subcategory").
		Order("created_at DESC").Find(&items).Error
	return items, err
}

func (s *ClothingService) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.ClothingItem{}).Error
}

func (s *ClothingService) GetByIDs(ids []string) ([]models.ClothingItem, error) {
	var items []models.ClothingItem
	err := database.DB.Preload("Category").Preload("Subcategory").
		Where("id IN ?", ids).Find(&items).Error
	return items, err
}

func (s *ClothingService) ExistsByCategory(categoryID uint, userID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.ClothingItem{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Count(&count).Error
	return count > 0, err
}

func (s *ClothingService) WithTrashed(id string) (*models.ClothingItem, error) {
	var item models.ClothingItem
	err := database.DB.Unscoped().Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ClothingService) Update(id string, updates map[string]interface{}) error {
	return database.DB.Model(&models.ClothingItem{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ClothingService) DB() *gorm.DB {
	return database.DB
}

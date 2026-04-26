package services

import (
	"go_wardrobe/database"
	"go_wardrobe/models"

	"gorm.io/gorm"
)

type OutfitService struct{}

func NewOutfitService() *OutfitService {
	return &OutfitService{}
}

// Create 创建搭配（含关联衣物）
func (s *OutfitService) Create(outfit *models.Outfit, items []models.OutfitItem) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(outfit).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OutfitID = outfit.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByID 获取搭配详情
func (s *OutfitService) GetByID(id string) (*models.Outfit, error) {
	var outfit models.Outfit
	err := database.DB.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Clothing.Category").Preload("Clothing.Subcategory")
	}).Where("id = ?", id).First(&outfit).Error
	if err != nil {
		return nil, err
	}
	return &outfit, nil
}

// List 获取所有搭配
func (s *OutfitService) List(userID uint) ([]models.Outfit, error) {
	var outfits []models.Outfit
	err := database.DB.Where("user_id = ?", userID).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Clothing.Category").Preload("Clothing.Subcategory")
		}).
		Order("created_at DESC").Find(&outfits).Error
	return outfits, err
}

// Delete 软删除搭配
func (s *OutfitService) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Outfit{}).Error
}

// DeletePermanent 永久删除搭配（含关联项）
func (s *OutfitService) DeletePermanent(id string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("outfit_id = ?", id).Delete(&models.OutfitItem{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Where("id = ?", id).Delete(&models.Outfit{}).Error
	})
}

// CloneOutfit 复制搭配（用于"从历史复制"功能）
func (s *OutfitService) CloneOutfit(id string, newOutfit *models.Outfit) error {
	original, err := s.GetByID(id)
	if err != nil {
		return err
	}

	newItems := make([]models.OutfitItem, len(original.Items))
	for i, item := range original.Items {
		newItems[i] = models.OutfitItem{
			ClothingID: item.ClothingID,
			Slot:       item.Slot,
		}
	}

	return s.Create(newOutfit, newItems)
}

// ReplaceItems 全量替换搭配衣物（删除旧关联，创建新关联）
func (s *OutfitService) ReplaceItems(outfitID string, newItems []models.OutfitItem) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("outfit_id = ?", outfitID).Delete(&models.OutfitItem{}).Error; err != nil {
			return err
		}
		for i := range newItems {
			newItems[i].OutfitID = outfitID
		}
		if len(newItems) > 0 {
			if err := tx.Create(&newItems).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateCardImage 更新搭配卡片图
func (s *OutfitService) UpdateCardImage(id, cardImage string) error {
	return database.DB.Model(&models.Outfit{}).
		Where("id = ?", id).Update("card_image", cardImage).Error
}

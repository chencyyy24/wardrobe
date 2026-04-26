package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `gorm:"primarykey;autoIncrement;type:int" json:"id"`
	Username  string         `gorm:"size:100;not null;uniqueIndex" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:100" json:"nickname,omitempty"`
	Avatar    string         `gorm:"size:500" json:"avatar,omitempty"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
}

// Category 大类表（外套/上衣/裤子/裙子/鞋子/配饰）
type Category struct {
	ID        uint      `gorm:"primarykey;autoIncrement;type:int" json:"id"`
	Name      string    `gorm:"size:50;not null;uniqueIndex" json:"name"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

func (Category) TableName() string {
	return "categories"
}

// Subcategory 小类表
type Subcategory struct {
	ID         uint      `gorm:"primarykey;autoIncrement;type:int" json:"id"`
	CategoryID uint      `gorm:"not null;index;type:int" json:"category_id"`
	Name       string    `gorm:"size:50;not null" json:"name"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
}

// ClothingItem 衣物表
type ClothingItem struct {
	ID            string         `gorm:"type:char(36);primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index:idx_user_category;type:int" json:"user_id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	CategoryID    uint           `gorm:"not null;index:idx_user_category;type:int" json:"category_id"`
	SubcategoryID *uint          `gorm:"index;type:int" json:"subcategory_id,omitempty"`
	OriginalImage string         `gorm:"size:500;not null" json:"original_image"`
	MaskedImage   string         `gorm:"size:500" json:"masked_image,omitempty"`
	Status        string         `gorm:"size:20;not null;default:pending;index" json:"status"`
	Price         *float64       `gorm:"type:decimal(10,2)" json:"price,omitempty"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Category    Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Subcategory Subcategory `gorm:"foreignKey:SubcategoryID" json:"subcategory,omitempty"`
}

func (ClothingItem) TableName() string {
	return "clothing_items"
}

// Outfit 搭配表
type Outfit struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_user_created;type:int" json:"user_id"`
	Name      string         `gorm:"size:255" json:"name,omitempty"`
	CardImage string         `gorm:"size:500" json:"card_image,omitempty"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Items []OutfitItem `gorm:"foreignKey:OutfitID" json:"items,omitempty"`
}

// OutfitItem 搭配-衣物关联表
type OutfitItem struct {
	ID         uint   `gorm:"primarykey;autoIncrement" json:"id"`
	OutfitID   string `gorm:"type:char(36);not null;uniqueIndex:uk_outfit_slot" json:"outfit_id"`
	ClothingID string `gorm:"type:char(36);not null;index" json:"clothing_id"`
	Slot       string `gorm:"size:20;not null;uniqueIndex:uk_outfit_slot" json:"slot"`

	// 关联
	Clothing ClothingItem `gorm:"foreignKey:ClothingID" json:"clothing,omitempty"`
}

func (OutfitItem) TableName() string {
	return "outfit_items"
}

// WearRecord 穿着记录表
type WearRecord struct {
	ID         uint      `gorm:"primarykey;autoIncrement;type:int" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_user_date;type:int" json:"user_id"`
	ClothingID *string   `gorm:"type:char(36)" json:"clothing_id,omitempty"`
	OutfitID   *string   `gorm:"type:char(36)" json:"outfit_id,omitempty"`
	WearDate   string    `gorm:"type:date;not null;index:idx_wear_date" json:"wear_date"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
}

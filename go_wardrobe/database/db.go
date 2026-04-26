package database

import (
	"log"

	"go_wardrobe/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dsn string) {
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DisableWithReturning:     true,
	}), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 尝试自动迁移，忽略外键相关的 ALTER 错误（表结构已通过 SQL 文件创建）
	migrateTables()

	// 初始化默认数据
	seedData()

	log.Println("database connected and migrated successfully")
}

func migrateTables() {
	// 临时关闭外键检查
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	tables := []interface{}{
		&models.User{},
		&models.Category{},
		&models.Subcategory{},
		&models.ClothingItem{},
		&models.Outfit{},
		&models.OutfitItem{},
		&models.WearRecord{},
	}

	for _, table := range tables {
		if err := DB.AutoMigrate(table); err != nil {
			log.Printf("migrate %T (non-fatal): %v", table, err)
		}
	}

	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func seedData() {
	// 检查是否已有数据
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("seeding default categories and subcategories...")

	categories := []models.Category{
		{Name: "外套", SortOrder: 1},
		{Name: "上衣", SortOrder: 2},
		{Name: "裤子", SortOrder: 3},
		{Name: "裙子", SortOrder: 4},
		{Name: "鞋子", SortOrder: 5},
		{Name: "配饰", SortOrder: 6},
	}
	for i := range categories {
		DB.Create(&categories[i])
	}

	subcategories := []models.Subcategory{
		{CategoryID: 1, Name: "风衣", SortOrder: 1},
		{CategoryID: 1, Name: "牛仔夹克", SortOrder: 2},
		{CategoryID: 1, Name: "西装", SortOrder: 3},
		{CategoryID: 1, Name: "针织开衫", SortOrder: 4},
		{CategoryID: 2, Name: "T恤", SortOrder: 1},
		{CategoryID: 2, Name: "衬衫", SortOrder: 2},
		{CategoryID: 2, Name: "卫衣", SortOrder: 3},
		{CategoryID: 2, Name: "毛衣", SortOrder: 4},
		{CategoryID: 3, Name: "牛仔裤", SortOrder: 1},
		{CategoryID: 3, Name: "休闲裤", SortOrder: 2},
		{CategoryID: 3, Name: "短裤", SortOrder: 3},
		{CategoryID: 4, Name: "JK裙", SortOrder: 1},
		{CategoryID: 4, Name: "百褶裙", SortOrder: 2},
		{CategoryID: 4, Name: "连衣裙", SortOrder: 3},
		{CategoryID: 5, Name: "运动鞋", SortOrder: 1},
		{CategoryID: 5, Name: "帆布鞋", SortOrder: 2},
		{CategoryID: 5, Name: "靴子", SortOrder: 3},
		{CategoryID: 6, Name: "帽子", SortOrder: 1},
		{CategoryID: 6, Name: "项链", SortOrder: 2},
		{CategoryID: 6, Name: "耳环", SortOrder: 3},
		{CategoryID: 6, Name: "包包", SortOrder: 4},
	}
	for i := range subcategories {
		DB.Create(&subcategories[i])
	}

	// 创建默认用户
	DB.Create(&models.User{
		Username: "default",
		Password: "default",
		Nickname: "默认用户",
	})

	log.Println("seed data created")
}

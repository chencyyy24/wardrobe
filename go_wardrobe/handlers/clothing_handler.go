package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go_wardrobe/config"
	"go_wardrobe/models"
	"go_wardrobe/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClothingHandler struct {
	service  *services.ClothingService
	cfg      *config.Config
}

func NewClothingHandler(cfg *config.Config) *ClothingHandler {
	return &ClothingHandler{
		service: services.NewClothingService(),
		cfg:     cfg,
	}
}

// Upload 上传衣物
// POST /api/clothing
func (h *ClothingHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传图片文件"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写衣物名称"})
		return
	}

	categoryID := c.PostForm("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择大类"})
		return
	}

	subcategoryID := c.PostForm("subcategory_id")

	clothingID := uuid.New().String()

	// 保存原图
	ext := filepath.Ext(file.Filename)
	originName := clothingID + ext
	originPath := filepath.Join(h.cfg.UploadDir, "origin", originName)
	if err := c.SaveUploadedFile(file, originPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		return
	}

	// 创建衣物记录（直接完成，无需抠图）
	userID := c.GetUint("user_id")
	item := &models.ClothingItem{
		ID:            clothingID,
		UserID:        userID,
		Name:          name,
		CategoryID:    parseUintDefault(categoryID),
		SubcategoryID: parseUintPtr(subcategoryID),
		OriginalImage: "/uploads/origin/" + originName,
		Status:        "done",
	}

	if err := h.service.Create(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建衣物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "上传成功",
		"data":    item,
	})
}

// List 获取衣物列表
// GET /api/clothing?category_id=1
func (h *ClothingHandler) List(c *gin.Context) {
	categoryID := c.Query("category_id")
	userID := c.GetUint("user_id")

	items, err := h.service.List(categoryID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取衣物列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

// GetByID 获取单个衣物详情
// GET /api/clothing/:id
func (h *ClothingHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少衣物 ID"})
		return
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "衣物不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Delete 删除衣物
// DELETE /api/clothing/:id
func (h *ClothingHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少衣物 ID"})
		return
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "衣物不存在"})
		return
	}

	// 删除关联文件
	deleteFileFromPath(item.OriginalImage, h.cfg.UploadDir)
	if item.MaskedImage != "" {
		deleteFileFromPath(item.MaskedImage, h.cfg.UploadDir)
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetCategories 获取所有分类及其小类
// GET /api/categories
func (h *ClothingHandler) GetCategories(c *gin.Context) {
	svc := services.NewClothingService()
	var categories []models.Category
	if err := svc.DB().Order("sort_order ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		return
	}

	type subcategoryResp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	type categoryResp struct {
		ID            uint              `json:"id"`
		Name          string            `json:"name"`
		Subcategories []subcategoryResp `json:"subcategories"`
	}

	var result []categoryResp
	for _, cat := range categories {
		var subs []models.Subcategory
		svc.DB().Where("category_id = ?", cat.ID).Order("sort_order ASC").Find(&subs)
		subList := make([]subcategoryResp, len(subs))
		for i, s := range subs {
			subList[i] = subcategoryResp{ID: s.ID, Name: s.Name}
		}
		result = append(result, categoryResp{
			ID:            cat.ID,
			Name:          cat.Name,
			Subcategories: subList,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// CreateTestData 创建测试数据
// POST /api/clothing/test-data
func (h *ClothingHandler) CreateTestData(c *gin.Context) {
	testItems := []struct {
		Name          string
		CategoryID    uint
		SubcategoryID *uint
	}{
		{Name: "米色风衣", CategoryID: 1, SubcategoryID: uintPtr(1)},
		{Name: "蓝色牛仔夹克", CategoryID: 1, SubcategoryID: uintPtr(2)},
		{Name: "白色T恤", CategoryID: 2, SubcategoryID: uintPtr(5)},
		{Name: "蓝色衬衫", CategoryID: 2, SubcategoryID: uintPtr(6)},
		{Name: "浅蓝牛仔裤", CategoryID: 3, SubcategoryID: uintPtr(9)},
		{Name: "黑色休闲裤", CategoryID: 3, SubcategoryID: uintPtr(10)},
		{Name: "灰色百褶裙", CategoryID: 4, SubcategoryID: uintPtr(13)},
		{Name: "白色运动鞋", CategoryID: 5, SubcategoryID: uintPtr(15)},
		{Name: "黑色帆布鞋", CategoryID: 5, SubcategoryID: uintPtr(16)},
		{Name: "棒球帽", CategoryID: 6, SubcategoryID: uintPtr(18)},
	}

	userID := c.GetUint("user_id")
	for i, ti := range testItems {
		id := uuid.New().String()
		item := &models.ClothingItem{
			ID:            id,
			UserID:        userID,
			Name:          ti.Name,
			CategoryID:    ti.CategoryID,
			SubcategoryID: ti.SubcategoryID,
			OriginalImage: fmt.Sprintf("/uploads/origin/test_%d.jpg", i+1),
			MaskedImage:   fmt.Sprintf("/uploads/masked/test_%d.png", i+1),
			Status:        "done",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := h.service.Create(item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测试数据失败: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "测试数据创建成功", "count": len(testItems)})
}

func uintPtr(v uint) *uint {
	return &v
}

// parseUintDefault 解析字符串为 uint，失败返回 0
func parseUintDefault(s string) uint {
	var v uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		v = v*10 + uint(c-'0')
	}
	return v
}

// parseUintPtr 解析字符串为 *uint
func parseUintPtr(s string) *uint {
	if s == "" {
		return nil
	}
	var v uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return nil
		}
		v = v*10 + uint(c-'0')
	}
	return &v
}

// deleteFileFromPath 根据 URL 路径删除文件
func deleteFileFromPath(urlPath, uploadDir string) {
	filename := filepath.Base(urlPath)
	for _, dir := range []string{"origin", "masked", "cards"} {
		fullPath := filepath.Join(uploadDir, dir, filename)
		if err := removeFile(fullPath); err == nil {
			return
		}
	}
}

func removeFile(path string) error {
	return os.Remove(path)
}

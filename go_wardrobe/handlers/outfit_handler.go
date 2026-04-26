package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"go_wardrobe/config"
	"go_wardrobe/models"
	"go_wardrobe/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OutfitHandler struct {
	service *services.OutfitService
	cfg     *config.Config
}

func NewOutfitHandler(cfg *config.Config) *OutfitHandler {
	return &OutfitHandler{
		service: services.NewOutfitService(),
		cfg:     cfg,
	}
}

// Create 创建搭配
// POST /api/outfit
func (h *OutfitHandler) Create(c *gin.Context) {
	var req struct {
		Name      string            `json:"name"`
		Items     map[string]string `json:"items"`
		CardImage string            `json:"card_image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误: " + err.Error()})
		return
	}

	if req.Items == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供搭配衣物"})
		return
	}

	topID := req.Items["top"]
	bottomID := req.Items["bottom"]
	if topID == "" || bottomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上衣和裤子为必选"})
		return
	}

	skirtID := req.Items["skirt"]
	if bottomID != "" && skirtID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "裤子和裙子不能同时选择"})
		return
	}

	outfitID := uuid.New().String()
	userID := c.GetUint("user_id")
	outfit := &models.Outfit{
		ID:        outfitID,
		UserID:    userID,
		Name:      req.Name,
		CardImage: req.CardImage,
	}

	slotMap := map[string]string{
		"outer":     req.Items["outer"],
		"top":       topID,
		"bottom":    bottomID,
		"skirt":     skirtID,
		"shoes":     req.Items["shoes"],
		"accessory": req.Items["accessory"],
	}

	var outfitItems []models.OutfitItem
	for slot, clothingID := range slotMap {
		if clothingID == "" {
			continue
		}
		outfitItems = append(outfitItems, models.OutfitItem{
			OutfitID:   outfitID,
			ClothingID: clothingID,
			Slot:       slot,
		})
	}

	if err := h.service.Create(outfit, outfitItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建搭配失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "搭配创建成功",
		"data":    outfit,
	})
}

// List 获取所有搭配
// GET /api/outfit
func (h *OutfitHandler) List(c *gin.Context) {
	outfits, err := h.service.List(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取搭配列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": outfits})
}

// GetByID 获取搭配详情
// GET /api/outfit/:id
func (h *OutfitHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少搭配 ID"})
		return
	}

	outfit, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "搭配不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": outfit})
}

// Update 更新搭配字段
// PATCH /api/outfit/:id
func (h *OutfitHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少搭配 ID"})
		return
	}

	var req struct {
		CardImage string `json:"card_image"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.CardImage == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供卡片图片 URL"})
		return
	}

	if err := h.service.UpdateCardImage(id, req.CardImage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新搭配封面失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// UpdateItems 更新搭配的衣物（全量替换）
// PUT /api/outfit/:id/items
func (h *OutfitHandler) UpdateItems(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少搭配 ID"})
		return
	}

	var req struct {
		Items map[string]string `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.Items == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供搭配衣物"})
		return
	}

	// 验证必填
	if req.Items["top"] == "" || req.Items["bottom"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上衣和裤子为必选"})
		return
	}

	if req.Items["bottom"] != "" && req.Items["skirt"] != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "裤子和裙子不能同时选择"})
		return
	}

	// 检查搭配存在
	_, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "搭配不存在"})
		return
	}

	// 构建新的 items
	var newItems []models.OutfitItem
	for _, slot := range []string{"outer", "top", "bottom", "skirt", "shoes", "accessory"} {
		clothingID := req.Items[slot]
		if clothingID == "" {
			continue
		}
		newItems = append(newItems, models.OutfitItem{
			OutfitID:   id,
			ClothingID: clothingID,
			Slot:       slot,
		})
	}

	if err := h.service.ReplaceItems(id, newItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新搭配衣物失败: " + err.Error()})
		return
	}

	// 返回更新后的搭配
	updated, _ := h.service.GetByID(id)
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// Delete 删除搭配
// DELETE /api/outfit/:id
func (h *OutfitHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少搭配 ID"})
		return
	}

	outfit, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "搭配不存在"})
		return
	}

	if outfit.CardImage != "" {
		cardPath := filepath.Join(h.cfg.UploadDir, "cards", filepath.Base(outfit.CardImage))
		_ = removeFile(cardPath)
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除搭配失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// Clone 复制搭配
// POST /api/outfit/:id/clone
func (h *OutfitHandler) Clone(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少搭配 ID"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Name = ""
	}

	newOutfit := &models.Outfit{
		ID:     uuid.New().String(),
		UserID: c.GetUint("user_id"),
		Name:   req.Name,
	}

	if err := h.service.CloneOutfit(id, newOutfit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "复制搭配失败: " + err.Error()})
		return
	}

	// 重新查询（含 Items），前端需要读取 items.length
	result, err := h.service.GetByID(newOutfit.ID)
	if err != nil {
		result = newOutfit
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "复制成功",
		"data":    result,
	})
}

// GetFlatList 获取扁平的搭配列表
// GET /api/outfit/summary
func (h *OutfitHandler) GetFlatList(c *gin.Context) {
	outfits, err := h.service.List(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取搭配列表失败"})
		return
	}

	type outfitSummary struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CardImage string `json:"card_image"`
		CreatedAt string `json:"created_at"`
		ItemCount int    `json:"item_count"`
	}

	summaries := make([]outfitSummary, 0)
	for _, o := range outfits {
		summaries = append(summaries, outfitSummary{
			ID:        o.ID,
			Name:      o.Name,
			CardImage: o.CardImage,
			CreatedAt: o.CreatedAt.Format(time.RFC3339),
			ItemCount: len(o.Items),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": summaries})
}

// CreateTestData 创建测试搭配数据
// POST /api/outfit/test-data
func (h *OutfitHandler) CreateTestData(c *gin.Context) {
	clothingSvc := services.NewClothingService()
	items, err := clothingSvc.List("", 1)
	if err != nil || len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先创建衣物测试数据"})
		return
	}

	byCategory := make(map[uint][]models.ClothingItem)
	for _, item := range items {
		byCategory[item.CategoryID] = append(byCategory[item.CategoryID], item)
	}

	userID := c.GetUint("user_id")
	for i, outfitName := range []string{"春日出游风", "通勤简约风", "周末休闲风"} {
		outfitID := uuid.New().String()
		outfit := &models.Outfit{
			ID:     outfitID,
			UserID: userID,
			Name:   outfitName,
		}

		var outfitItems []models.OutfitItem
		if tops, ok := byCategory[2]; ok && len(tops) > 0 {
			outfitItems = append(outfitItems, models.OutfitItem{
				OutfitID: outfitID, ClothingID: tops[i%len(tops)].ID, Slot: "top",
			})
		}
		if bottoms, ok := byCategory[3]; ok && len(bottoms) > 0 {
			outfitItems = append(outfitItems, models.OutfitItem{
				OutfitID: outfitID, ClothingID: bottoms[i%len(bottoms)].ID, Slot: "bottom",
			})
		}
		if outers, ok := byCategory[1]; ok && len(outers) > 0 {
			outfitItems = append(outfitItems, models.OutfitItem{
				OutfitID: outfitID, ClothingID: outers[i%len(outers)].ID, Slot: "outer",
			})
		}
		if shoes, ok := byCategory[5]; ok && len(shoes) > 0 {
			outfitItems = append(outfitItems, models.OutfitItem{
				OutfitID: outfitID, ClothingID: shoes[i%len(shoes)].ID, Slot: "shoes",
			})
		}

		if err := h.service.Create(outfit, outfitItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建搭配 %s 失败: %s", outfitName, err.Error())})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "测试搭配数据创建成功"})
}

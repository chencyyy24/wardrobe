package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go_wardrobe/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	cfg *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

// UploadCard 上传搭配卡片图片
// POST /api/upload/card
// Content-Type: multipart/form-data
// Fields: image (file)
func (h *UploadHandler) UploadCard(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传卡片图片"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".png"
	}
	fileName := fmt.Sprintf("card_%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	savePath := filepath.Join(h.cfg.UploadDir, "cards", fileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存卡片图片失败"})
		return
	}

	imageURL := "/uploads/cards/" + fileName
	c.JSON(http.StatusCreated, gin.H{
		"message": "上传成功",
		"url":     imageURL,
	})
}

// UploadImage 通用图片上传接口
// POST /api/upload/image
// Fields: image (file), type (可选: origin|masked|custom)
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传图片"})
		return
	}

	imageType := c.DefaultPostForm("type", "origin")
	saveDir := filepath.Join(h.cfg.UploadDir, imageType)

	// 确保目录存在
	_ = ensureDir(saveDir)

	fileName := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), filepath.Ext(file.Filename))
	savePath := filepath.Join(saveDir, fileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		return
	}

	imageURL := fmt.Sprintf("/uploads/%s/%s", imageType, fileName)
	c.JSON(http.StatusCreated, gin.H{
		"message": "上传成功",
		"url":     imageURL,
	})
}

// ensureDir 确保目录存在
func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

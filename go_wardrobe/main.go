package main

import (
	"fmt"
	"log"
	"os"

	"go_wardrobe/config"
	"go_wardrobe/database"
	"go_wardrobe/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 确保上传目录存在
	dirs := []string{
		cfg.UploadDir,
		cfg.UploadDir + "/origin",
		cfg.UploadDir + "/masked",
		cfg.UploadDir + "/cards",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	// 初始化数据库
	database.Init(cfg.DatabaseDSN)

	// 初始化 Gin
	r := gin.Default()

	// 设置路由
	routes.Setup(r, cfg)

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Wardrobe server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

package routes

import (
	"go_wardrobe/config"
	"go_wardrobe/database"
	"go_wardrobe/handlers"
	"go_wardrobe/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, cfg *config.Config) {
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/uploads", cfg.UploadDir)

	// 公开路由
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "wardrobe-server"})
	})
	r.GET("/api/categories", handlers.NewClothingHandler(cfg).GetCategories)

	// 认证路由
	authHandler := handlers.NewAuthHandler(cfg, database.DB)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.Auth(cfg), authHandler.Me)
		auth.PUT("/profile", middleware.Auth(cfg), authHandler.UpdateProfile)
		auth.PUT("/password", middleware.Auth(cfg), authHandler.ChangePassword)
		auth.POST("/avatar", middleware.Auth(cfg), authHandler.UploadAvatar)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.Auth(cfg))
	{
		clothingHandler := handlers.NewClothingHandler(cfg)
		outfitHandler := handlers.NewOutfitHandler(cfg)
		uploadHandler := handlers.NewUploadHandler(cfg)

		clothing := protected.Group("/clothing")
		{
			clothing.POST("", clothingHandler.Upload)
			clothing.GET("", clothingHandler.List)
			clothing.POST("/test-data", clothingHandler.CreateTestData)
			clothing.GET("/:id", clothingHandler.GetByID)
			clothing.DELETE("/:id", clothingHandler.Delete)
		}

		outfit := protected.Group("/outfit")
		{
			outfit.POST("", outfitHandler.Create)
			outfit.GET("", outfitHandler.List)
			outfit.GET("/summary", outfitHandler.GetFlatList)
			outfit.POST("/test-data", outfitHandler.CreateTestData)
			outfit.GET("/:id", outfitHandler.GetByID)
			outfit.PATCH("/:id", outfitHandler.Update)
			outfit.PUT("/:id/items", outfitHandler.UpdateItems)
			outfit.DELETE("/:id", outfitHandler.Delete)
			outfit.POST("/:id/clone", outfitHandler.Clone)
		}

		upload := protected.Group("/upload")
		{
			upload.POST("/card", uploadHandler.UploadCard)
			upload.POST("/image", uploadHandler.UploadImage)
		}
	}
}

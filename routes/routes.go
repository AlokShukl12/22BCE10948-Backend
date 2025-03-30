package routes

import (
	"filesharing/handlers"
	"filesharing/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, fileHandler *handlers.FileHandler) {
	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.Use(middleware.RateLimitMiddleware())
	{
		// File routes
		authorized.POST("/upload", fileHandler.UploadFile)
		authorized.GET("/files", fileHandler.ListFiles)
		authorized.GET("/files/search", fileHandler.SearchFiles)
		authorized.GET("/files/:id", fileHandler.GetFile)
		authorized.POST("/files/:id/share", fileHandler.ShareFile)
	}

	// Background job for cleaning up expired files
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			if err := fileHandler.DeleteExpiredFiles(); err != nil {
				// Log error but continue
				log.Printf("Error cleaning up expired files: %v", err)
			}
		}
	}()
} 
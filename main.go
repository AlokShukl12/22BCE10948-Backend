package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize router
	router := gin.Default()

	// Initialize services
	initServices(router)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initServices(router *gin.Engine) {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis connection
	redisClient, err := initRedis()
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize S3 client
	s3Client, err := initS3()
	if err != nil {
		log.Fatal("Failed to initialize S3 client:", err)
	}

	// Initialize repositories
	userRepo := NewUserRepository(db)
	fileRepo := NewFileRepository(db)

	// Initialize services
	authService := NewAuthService(userRepo, redisClient)
	fileService := NewFileService(fileRepo, s3Client, redisClient)

	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	fileHandler := NewFileHandler(fileService)

	// Setup routes
	setupRoutes(router, authHandler, fileHandler)
} 
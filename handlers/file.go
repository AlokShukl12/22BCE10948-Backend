package handlers

import (
	"filesharing/models"
	"filesharing/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	fileService *services.FileService
}

func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID := userID.(uuid.UUID)

	var upload models.FileUpload
	if err := c.ShouldBind(&upload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.fileService.UploadFile(userUUID, &upload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, file)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID := userID.(uuid.UUID)
	fileID := c.Param("id")

	file, err := h.fileService.GetFile(userUUID, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID := userID.(uuid.UUID)

	files, err := h.fileService.ListFiles(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) SearchFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID := userID.(uuid.UUID)

	var search models.FileSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files, err := h.fileService.SearchFiles(userUUID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) ShareFile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID := userID.(uuid.UUID)

	var share models.FileShare
	if err := c.ShouldBindJSON(&share); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.fileService.ShareFile(userUUID, &share)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file)
} 
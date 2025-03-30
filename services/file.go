package services

import (
	"context"
	"errors"
	"filesharing/models"
	"filesharing/repositories"
	"filesharing/storage"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type FileService struct {
	fileRepo    *repositories.FileRepository
	s3Client    *s3.Client
	redisClient *redis.Client
}

func NewFileService(fileRepo *repositories.FileRepository, s3Client *s3.Client, redisClient *redis.Client) *FileService {
	return &FileService{
		fileRepo:    fileRepo,
		s3Client:    s3Client,
		redisClient: redisClient,
	}
}

func (s *FileService) UploadFile(userID uuid.UUID, upload *models.FileUpload) (*models.File, error) {
	// Generate unique file ID
	fileID := uuid.New()
	s3Key := fileID.String()

	// Upload to S3
	_, err := s.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(storage.GetS3BucketName()),
		Key:         aws.String(s3Key),
		Body:        io.NopCloser(io.Reader(io.MultiReader(io.NopCloser(io.Reader(upload.File))))),
		ContentType: aws.String(upload.ContentType),
	})
	if err != nil {
		return nil, err
	}

	// Create file record
	file := &models.File{
		ID:          fileID,
		UserID:      userID,
		Name:        upload.Name,
		Size:        int64(len(upload.File)),
		ContentType: upload.ContentType,
		S3Key:       s3Key,
		IsPublic:    upload.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if upload.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, upload.ExpiresAt)
		if err != nil {
			return nil, err
		}
		file.ExpiresAt = expiresAt
	}

	if err := s.fileRepo.Create(file); err != nil {
		// Cleanup S3 if database insert fails
		s.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(storage.GetS3BucketName()),
			Key:    aws.String(s3Key),
		})
		return nil, err
	}

	return file, nil
}

func (s *FileService) GetFile(userID uuid.UUID, fileID string) (*models.File, error) {
	// Try to get from cache first
	cacheKey := "file:" + fileID
	cachedFile, err := s.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// Cache hit
		var file models.File
		if err := json.Unmarshal([]byte(cachedFile), &file); err == nil {
			return &file, nil
		}
	}

	// Get from database
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if file.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Cache the result
	if fileJSON, err := json.Marshal(file); err == nil {
		s.redisClient.Set(context.Background(), cacheKey, fileJSON, 5*time.Minute)
	}

	return file, nil
}

func (s *FileService) ListFiles(userID uuid.UUID) ([]models.File, error) {
	return s.fileRepo.FindByUserID(userID.String())
}

func (s *FileService) SearchFiles(userID uuid.UUID, search models.FileSearch) ([]models.File, error) {
	return s.fileRepo.Search(userID.String(), search)
}

func (s *FileService) ShareFile(userID uuid.UUID, share *models.FileShare) (*models.File, error) {
	file, err := s.GetFile(userID, share.FileID.String())
	if err != nil {
		return nil, err
	}

	file.IsPublic = true
	file.ExpiresAt = share.ExpiresAt

	if err := s.fileRepo.Update(file); err != nil {
		return nil, err
	}

	// Invalidate cache
	s.redisClient.Del(context.Background(), "file:"+file.ID.String())

	return file, nil
}

func (s *FileService) DeleteExpiredFiles() error {
	// Get expired files
	files, err := s.fileRepo.FindExpiredFiles()
	if err != nil {
		return err
	}

	// Delete from S3 and database
	for _, file := range files {
		// Delete from S3
		_, err := s.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(storage.GetS3BucketName()),
			Key:    aws.String(file.S3Key),
		})
		if err != nil {
			// Log error but continue with other files
			log.Printf("Error deleting file from S3: %v", err)
		}

		// Delete from database
		if err := s.fileRepo.Delete(&file); err != nil {
			log.Printf("Error deleting file from database: %v", err)
		}

		// Invalidate cache
		s.redisClient.Del(context.Background(), "file:"+file.ID.String())
	}

	return nil
} 
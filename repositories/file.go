package repositories

import (
	"filesharing/models"
	"time"

	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) FindByID(id string) (*models.File, error) {
	var file models.File
	if err := r.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) FindByUserID(userID string) ([]models.File, error) {
	var files []models.File
	if err := r.db.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) Search(userID string, search models.FileSearch) ([]models.File, error) {
	var files []models.File
	query := r.db.Where("user_id = ?", userID)

	if search.Name != "" {
		query = query.Where("name ILIKE ?", "%"+search.Name+"%")
	}
	if !search.StartDate.IsZero() {
		query = query.Where("created_at >= ?", search.StartDate)
	}
	if !search.EndDate.IsZero() {
		query = query.Where("created_at <= ?", search.EndDate)
	}
	if search.FileType != "" {
		query = query.Where("content_type = ?", search.FileType)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) DeleteExpiredFiles() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&models.File{}).Error
}

func (r *FileRepository) Update(file *models.File) error {
	return r.db.Save(file).Error
} 
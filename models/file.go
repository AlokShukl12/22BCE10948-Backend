package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Size        int64     `json:"size" gorm:"not null"`
	ContentType string    `json:"content_type" gorm:"not null"`
	S3Key       string    `json:"-" gorm:"not null"`
	PublicURL   string    `json:"public_url"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FileUpload struct {
	File        []byte `form:"file" binding:"required"`
	Name        string `form:"name"`
	ContentType string `form:"content_type"`
	IsPublic    bool   `form:"is_public"`
	ExpiresAt   string `form:"expires_at"`
}

type FileSearch struct {
	Name      string    `form:"name"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	FileType  string    `form:"file_type"`
}

type FileShare struct {
	FileID    uuid.UUID `json:"file_id" binding:"required"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
} 
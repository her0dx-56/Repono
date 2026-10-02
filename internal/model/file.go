package model

import (
	"time"

	"github.com/google/uuid"
)

type File struct{
	ID uuid.UUID
	UserID uuid.UUID
	Name string
	StorageKey string
	Size int64
	ContentType string
	CreatedAt time.Time
	UpdatedAt time.Time
}
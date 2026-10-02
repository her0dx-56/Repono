package model

import (
	"time"
	"github.com/google/uuid"
)

type FileResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}
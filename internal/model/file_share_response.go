package model

import(
	"time"
	"github.com/google/uuid"
)

type SharedFileResponse struct {
	ShareID  uuid.UUID `json:"share_id"`
	FileID  uuid.UUID `json:"file_id"`
	OwnerID  uuid.UUID `json:"owner_id"`
	Permission string    `json:"permission"`
	Name string    `json:"name"`
	Size int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	CreatedAt    time.Time `json:"created_at"`
	SharedAt     time.Time `json:"shared_at"`
}
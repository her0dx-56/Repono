package model

import(
	"time"
	"github.com/google/uuid"
)

type FileShare struct{
	ID uuid.UUID
	FileID uuid.UUID
	OwnerID uuid.UUID
	SharedWithUserID uuid.UUID
	Permission string
	CreatedAt time.Time
}
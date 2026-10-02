package model

import(
	"time"
	 "github.com/google/uuid"
)

type PublicShare struct{
	ID uuid.UUID
	FileID uuid.UUID
	OwnerID uuid.UUID
	Token string
	ExpiresAt *time.Time
	CreatedAt  time.Time
	RevokedAt *time.Time
}
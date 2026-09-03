package entity

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	ChatType  string
	Title     *string
	CreatedAt time.Time
}

package entity

import "github.com/google/uuid"

type UserSession struct {
	UserID  uuid.UUID
	IsValid bool
}

package entity

import "github.com/google/uuid"

type LoginResult struct {
	Required2FA  bool
	SessionToken string
}

type User struct {
	ID           uuid.UUID
	Name         string
	Username     string
	Email        string
	PasswordHash string
}

package dto

import (
	"time"
	"user-service/internal/entity"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string
}
type GetUserByIDsRequest struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}
type SearchUsersResponse struct {
	Name     string    `json:"name"`
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
}
type MeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func ToDTOMeResponse(u *entity.User) *MeResponse {
	return &MeResponse{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
func ToDTOSearchResult(u *entity.User) *SearchUsersResponse {
	return &SearchUsersResponse{
		Name:     u.Name,
		Username: u.Username,
		ID:       u.ID,
	}
}

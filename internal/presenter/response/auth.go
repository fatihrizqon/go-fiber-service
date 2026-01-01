package response

import (
	"time"

	"github.com/fatihrizqon/go-fiber-service/internal/entity"
	"github.com/google/uuid"
)

type RegisterResponse struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  entity.User `json:"user"`
}

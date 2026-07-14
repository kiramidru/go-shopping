package users

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=8"`
	FirstName string  `json:"firstName" binding:"required,min=1,max=100"`
	LastName  string  `json:"lastName" binding:"required,min=1,max=100"`
	Phone     *string `json:"phone" binding:"omitempty,max=30"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateRequest struct {
	FirstName *string `json:"firstName" binding:"omitempty,min=1,max=100"`
	LastName  *string `json:"lastName" binding:"omitempty,min=1,max=100"`
	Phone     *string `json:"phone" binding:"omitempty,max=30"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Phone      *string   `json:"phone,omitempty"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
}

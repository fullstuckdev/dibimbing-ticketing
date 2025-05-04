package dto

import (
	"time"
)

// UserRegisterRequest represents the request for user registration
type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLoginRequest represents the request for user login
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the response format for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProfileResponse represents user profile data with additional information
type UserProfileResponse struct {
	UserResponse
	TicketCount int `json:"ticket_count"`
}

// TokenResponse represents JWT token response
type TokenResponse struct {
	Token string `json:"token"`
} 
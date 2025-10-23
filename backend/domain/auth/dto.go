package auth

import "github.com/guregu/null/v6/zero"

type RegisterRequest struct {
	Username zero.String `json:"username"`
	Password zero.String `json:"password"`
}

type LoginRequest struct {
	Username zero.String `json:"username"`
	Password zero.String `json:"password"`
}

type UserResponse struct {
	ID        int         `json:"id"`
	Username  zero.String `json:"username"`
	CreatedAt string      `json:"created_at"`
}

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

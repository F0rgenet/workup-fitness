package user

import "github.com/guregu/null/v6/zero"

type GetPublicProfileResponse struct {
	ID        int         `json:"id"`
	Username  zero.String `json:"username"`
	CreatedAt string      `json:"created_at"`
}

type GetPrivateProfileResponse struct {
	ID        int         `json:"id"`
	Username  zero.String `json:"username"`
	CreatedAt string      `json:"created_at"`
	// TODO: Add private info fields
}

type UpdateRequest struct {
	Username zero.String `json:"username"`
	Password zero.String `json:"password"`
}

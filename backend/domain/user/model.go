package user

import (
	"time"

	"github.com/guregu/null/v6/zero"
)

type User struct {
	ID           int         `json:"id"`
	Username     zero.String `json:"username"`
	PasswordHash zero.String `json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
}

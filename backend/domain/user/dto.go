package user

type GetPublicProfileResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type GetPrivateProfileResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	// TODO: Add private info fields
}

type UpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

package repository

type User struct {
	UserID string `json:"user_id"`
	TgID   int64  `json:"tg_id"`
}

type SaveUserRequest struct {
	TgID     int64  `json:"tg_id"`
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
}

type GetUserByTgIDRequest struct {
	TgID int64 `json:"tg_id"`
}

type GetUserByTgIDResponse struct {
	user *User
}

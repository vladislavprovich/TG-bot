package repository

type User struct {
	//todo
}

type GetUserByTgIDRequest struct {
	TgID int64 `json:"tg_id"`
}

type GetUserByTgIDResponse struct {
	user *User
}

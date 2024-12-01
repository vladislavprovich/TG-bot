package repository

type UserRepository interface {
	SaveUser()
	GetUserByTgID()
}

type UserRepo struct {
}

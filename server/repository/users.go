package repository

type UsersRepo interface {
	Login(login, password string) (string, error)
	Register(userID, login, hashPassword string) error
}

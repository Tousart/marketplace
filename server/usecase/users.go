package usecase

type AuthService interface {
	Login(login, password string) (string, error)
	Register(login, password string) (string, string, string, error)
}

package model

type User struct {
	UUID     string
	Login    string
	Password string
}

func NewUser(UUID, login, password string) *User {
	return &User{
		UUID:     UUID,
		Login:    login,
		Password: password,
	}
}

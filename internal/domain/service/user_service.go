package service

import (
	"tictactoe/internal/domain/model"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

type Encrypter interface {
	Encrypt(data string) string
}

func (us *UserService) SignUp(UUID, login, password string) (*model.User, error) {
	return model.NewUser(UUID, login, password), nil
}

func (us *UserService) SignIn(user *model.User, password string) bool {
	return user.Password == password
}

package usecase

import (
	"context"
	"fmt"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type SignInUseCase interface {
	Execute(ctx context.Context, cmd *SignInCommand) (*model.User, error)
}

type SignInCommand struct {
	Login    string
	Password string
}

type SignIn struct {
	userRepo    port.UserRepo
	userService port.UserService
}

func NewSignIn(userRepo port.UserRepo, userService port.UserService) *SignIn {
	return &SignIn{
		userRepo:    userRepo,
		userService: userService,
	}
}

func (uc *SignIn) Execute(ctx context.Context, cmd *SignInCommand) (*model.User, error) {
	wrap := func(err error) error {
		return fmt.Errorf("sign in: %w", err)
	}

	users, err := uc.userRepo.Get(ctx, port.WithLogin(cmd.Login))
	if err != nil {
		return nil, wrap(err)
	}
	if len(users) != 1 {
		return nil, wrap(port.ErrReturnedNotOne)
	}
	user := users[0]

	ok := uc.userService.SignIn(user, cmd.Password)
	if !ok {
		return nil, wrap(port.ErrInvalidCredentials)
	}
	return user, nil
}

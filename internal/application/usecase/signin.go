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

	user, err := uc.userRepo.Get(ctx, cmd.Login)
	if err != nil {
		return nil, wrap(err)
	}

	ok := uc.userService.SignIn(user, cmd.Password)
	if !ok {
		return nil, wrap(port.ErrInvalidCredentials)
	}
	return user, nil
}

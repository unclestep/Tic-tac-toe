package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type SignUpCommand struct {
	Login    string
	Password string
}

type SignUpUseCase interface {
	Execute(ctx context.Context, cmd *SignUpCommand) (*model.User, error)
}

type SignUp struct {
	userRepo    port.UserRepo
	userService port.UserService
}

func NewSignUp(userRepo port.UserRepo, userService port.UserService) *SignUp {
	return &SignUp{
		userRepo:    userRepo,
		userService: userService,
	}
}

func (uc *SignUp) Execute(ctx context.Context, cmd *SignUpCommand) (*model.User, error) {
	wrap := func(err error) error {
		return fmt.Errorf("sign up: %w", err)
	}

	uuid := uuid.NewString()
	user, err := uc.userService.SignUp(uuid, cmd.Login, cmd.Password)
	if err != nil {
		return nil, wrap(err)
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, wrap(err)
	}

	return user, nil
}

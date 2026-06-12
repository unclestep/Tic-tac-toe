package port

import (
	"context"
	"errors"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

var (
	ErrPlayerNotBelongToSession = errors.New("player is not belong to session")
	ErrGameAlreadyStarted       = errors.New("game already started")
	ErrPlayerCantMakeMove       = errors.New("player cant make move")
	ErrGameNotStarted           = errors.New("game not started")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrInvalidBoardSize         = errors.New("invalid board size")
)

type BotMover interface {
	MakeMove(session *model.Session, botMark model.Mark) error
}

type HumanMover interface {
	MakeMove(session *model.Session, player *model.Player, p geometry.Point) error
}

type UserService interface {
	SignUp(UUID, login, password string) (*model.User, error)
	SignIn(user *model.User, password string) bool
}

type CreateUseCase interface {
	Execute(ctx context.Context, cmd *CreateCommand) (*model.Session, error)
}

type ConnectUseCase interface {
	Execute(ctx context.Context, cmd *ConnectCommand) (*model.Session, error)
}

type StartUseCase interface {
	Execute(ctx context.Context, cmd *StartCommand) (*model.Session, error)
}

type MakeMoveUseCase interface {
	Execute(ctx context.Context, cmd *MakeMoveCommand) (*model.Session, error)
}

type DisconnectUseCase interface {
	Execute(ctx context.Context, cmd *DisconnectCommand) (*model.Session, error)
}

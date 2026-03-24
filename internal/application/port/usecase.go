package port

import (
	"context"
	"errors"
	"tictactoe/internal/domain/model"
)

var (
	ErrPlayerNotBelongToSession = errors.New("player is not belong to session")
	ErrGameAlreadyStarted       = errors.New("game already started")
	ErrPlayerCantMakeMove       = errors.New("player cant make move")
	ErrGameNotStarted           = errors.New("game not started")
)

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

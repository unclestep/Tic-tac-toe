package usecase

import (
	"context"
	"fmt"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"

	"github.com/google/uuid"
)

type Create struct {
	sessionRepo port.SessionRepo
}

func NewCreate(sessionRepo port.SessionRepo) *Create {
	return &Create{
		sessionRepo: sessionRepo,
	}
}

func (uc *Create) Execute(ctx context.Context, cmd *port.CreateCommand) (*model.Session, error) {
	wrap := func(err error) error {
		return fmt.Errorf("create: %w", err)
	}

	if cmd.Rules.BoardHeight <= 0 || cmd.Rules.BoardWidth <= 0 {
		return nil, wrap(port.ErrInvalidBoardSize)
	}

	board := model.NewBoard(cmd.Rules.BoardWidth, cmd.Rules.BoardHeight)
	session := model.NewSession(uuid.NewString(), cmd.SessionParams, cmd.Rules, board)

	err := uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

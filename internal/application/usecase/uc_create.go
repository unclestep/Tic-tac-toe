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
	rulesRepo   port.RulesRepo
}

func NewCreate(sessionRepo port.SessionRepo, rulesRepo port.RulesRepo) *Create {
	return &Create{
		sessionRepo: sessionRepo,
		rulesRepo:   rulesRepo,
	}
}

func (uc *Create) Execute(ctx context.Context, cmd *port.CreateCommand) (*model.Session, error) {
	wrap := func(err error) error {
		return fmt.Errorf("create: %w", err)
	}

	cmd.Rules.UUID = uuid.NewString()
	err := uc.rulesRepo.Save(ctx, cmd.Rules)
	if err != nil {
		return nil, wrap(err)
	}

	board := model.NewBoard(cmd.Rules.BoardWidth, cmd.Rules.BoardHeight)
	session := model.NewSession(uuid.NewString(), cmd.SessionParams, cmd.Rules, board)
	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

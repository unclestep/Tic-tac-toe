package usecase

import (
	"context"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
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
	// The function will assign a UUID if it is empty.
	err := uc.rulesRepo.Save(ctx, cmd.Rules)
	if err != nil {
		return nil, err
	}

	session, err := uc.sessionRepo.Create(ctx, cmd.SessionParams, cmd.Rules)
	if err != nil {
		return nil, fmt.Errorf("create: cant create session : %w", err)
	}

	err = session.AddPlayer(cmd.PlayerID, cmd.PlayerName)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("connect: save session %s: %w", session.UUID, err)
	}

	return session, nil
}

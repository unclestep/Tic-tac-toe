package usecase

import (
	"context"
	"errors"
	"fmt"
	appPort "tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	servPort "tictactoe/internal/domain/service/port"
)

type Disconnect struct {
	sessionRepo   appPort.SessionRepo
	gameMechanics servPort.GameMechanics
}

func NewDisconnect(sessionRepo appPort.SessionRepo, gameMechanics servPort.GameMechanics) *Disconnect {
	return &Disconnect{
		sessionRepo:   sessionRepo,
		gameMechanics: gameMechanics,
	}
}

func (uc *Disconnect) Execute(ctx context.Context, cmd *appPort.DisconnectCommand) (*model.Session, error) {
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionID)
	if err != nil {
		if errors.Is(err, appPort.ErrSessionNotFound) {
			return nil, fmt.Errorf("session %s not found", cmd.SessionID)
		}
		return nil, fmt.Errorf("disconnect: get session %s: %w", cmd.SessionID, err)
	}

	err = session.HidePlayer(cmd.PlayerID)
	if err != nil {
		return nil, err
	}

	err = uc.gameMechanics.Advance(session)
	if err != nil {
		return nil, err
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

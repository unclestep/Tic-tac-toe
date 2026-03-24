package usecase

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type Disconnect struct {
	sessionRepo port.SessionRepo
}

func (uc *Disconnect) Execute(ctx context.Context, cmd *port.DisconnectCommand) (*model.Session, error) {
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionID)
	if err != nil {
		if errors.Is(err, port.ErrSessionNotFound) {
			return nil, fmt.Errorf("session %s not found", cmd.SessionID)
		}
		return nil, fmt.Errorf("connect: get session %s: %w", cmd.SessionID, err)
	}

	err = session.RemovePlayer(cmd.PlayerID)
	if err != nil {
		return nil, err
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

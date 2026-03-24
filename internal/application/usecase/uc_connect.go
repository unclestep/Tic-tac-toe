package usecase

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type Connect struct {
	sessionRepo port.SessionRepo
}

func NewConnect(sessionRepo port.SessionRepo) *Connect {
	return &Connect{
		sessionRepo: sessionRepo,
	}
}

func (uc *Connect) Execute(ctx context.Context, cmd *port.ConnectCommand) (*model.Session, error) {
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionID)
	if err != nil {
		if errors.Is(err, port.ErrSessionNotFound) {
			return nil, fmt.Errorf("session %s not found", cmd.SessionID)
		}
		return nil, fmt.Errorf("connect: get session %s: %w", cmd.SessionID, err)
	}

	// if session.State != model.StateLobby {
	// 	return nil, fmt.Errorf("%w: session %s", model.ErrGameAlreadyStarted, session.UUID)
	// }

	err = session.AddPlayer(cmd.PlayerID, cmd.PlayerName)
	if err != nil {
		return nil, err
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("connect: save session %s: %w", session.UUID, err)
	}

	return session, nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	appPort "tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	servicePort "tictactoe/internal/domain/service/port"
)

type MakeMove struct {
	sessionRepo   appPort.SessionRepo
	gameMechanics servicePort.GameMechanics
}

func NewMakeMove(sessionRepo appPort.SessionRepo, gameMechanics servicePort.GameMechanics) *MakeMove {
	return &MakeMove{
		sessionRepo:   sessionRepo,
		gameMechanics: gameMechanics,
	}
}

func (uc *MakeMove) Execute(ctx context.Context, cmd *appPort.MakeMoveCommand) (*model.Session, error) {
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionID)
	if err != nil {
		if errors.Is(err, appPort.ErrSessionNotFound) {
			return nil, fmt.Errorf("session %s not found", cmd.SessionID)
		}
		return nil, fmt.Errorf("move: get session %s: %w", cmd.SessionID, err)
	}

	if session.State != model.StatePlaying {
		return nil, fmt.Errorf("%w: session %s", appPort.ErrGameNotStarted, session.UUID)
	}

	if !session.IsPlayerExist(cmd.PlayerID) {
		return nil, fmt.Errorf("%w: session %s, player %s", appPort.ErrPlayerNotBelongToSession, session.UUID, cmd.PlayerID)
	}

	player := session.GetTurnPlayer()
	if player == nil {
		return nil, fmt.Errorf("%w: session %s, player %s", model.ErrPlayerNotFound, session.UUID, cmd.PlayerID)
	}

	if player.UUID != cmd.PlayerID {
		return nil, fmt.Errorf("%w: not player %s turn, session %s", appPort.ErrPlayerCantMakeMove, cmd.PlayerID, cmd.SessionID)
	}

	err = uc.gameMechanics.MakeMove(session, player, cmd.MarkPos)
	if err != nil {
		return nil, err
	}

	if session.State == model.StatePlaying {
		err = uc.gameMechanics.Advance(session)
		if err != nil {
			return nil, err
		}
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

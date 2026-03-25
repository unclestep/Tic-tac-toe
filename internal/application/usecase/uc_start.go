package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	appPort "tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	servicePort "tictactoe/internal/domain/service/port"
)

type Start struct {
	sessionRepo   appPort.SessionRepo
	gameMechanics servicePort.GameMechanics
}

func NewStart(sessionRepo appPort.SessionRepo, gameService servicePort.GameMechanics) *Start {
	return &Start{
		sessionRepo:   sessionRepo,
		gameMechanics: gameService,
	}
}

func (uc *Start) Execute(ctx context.Context, cmd *appPort.StartCommand) (*model.Session, error) {
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionID)
	if err != nil {
		if errors.Is(err, appPort.ErrSessionNotFound) {
			return nil, fmt.Errorf("session %s not found", cmd.SessionID)
		}
		return nil, fmt.Errorf("start: get session %s: %w", cmd.SessionID, err)
	}

	if session.State == model.StatePlaying {
		return nil, fmt.Errorf("%w: session %s", appPort.ErrGameAlreadyStarted, session.UUID)
	}

	if !session.IsPlayerExist(cmd.PlayerID) {
		return nil, fmt.Errorf("%w: session %s, player %s", appPort.ErrPlayerNotBelongToSession, session.UUID, cmd.PlayerID)
	}

	rng := rand.New(rand.NewSource(session.Params.Seed))
	session.State = model.StatePlaying

	if len(session.Players) == 1 {
		session.AddBot()
	}

	uc.gameMechanics.Prepare(session, rng)
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

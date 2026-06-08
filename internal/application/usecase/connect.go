package usecase

import (
	"context"
	"fmt"
	"math/rand"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"

	"github.com/google/uuid"
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
	wrap := func(err error) error {
		return fmt.Errorf("connect (session %s): %w", cmd.SessionUUID, err)
	}
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionUUID)
	if err != nil {
		return nil, wrap(err)
	}

	am := session.GetAvailableMarks()
	if len(am) == 0 {
		return nil, wrap(model.ErrSessionFull)
	}

	player := model.NewPlayer(uuid.NewString(), cmd.PlayerName, am[rand.Intn(len(am))])
	err = session.AddPlayer(player)
	if err != nil {
		return nil, wrap(err)
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

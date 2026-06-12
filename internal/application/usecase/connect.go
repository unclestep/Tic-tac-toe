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
	sessions, err := uc.sessionRepo.Get(ctx, port.WithUUID(cmd.SessionUUID))
	if err != nil {
		return nil, wrap(err)
	}
	if len(sessions) != 1 {
		return nil, wrap(port.ErrReturnedNotOne)
	}
	session := sessions[0]

	am := session.GetAvailableMarks()
	if len(am) == 0 {
		return nil, wrap(model.ErrSessionFull)
	}

	rng := rand.New(rand.NewSource(session.Params.Seed))

	player := model.NewPlayer(uuid.NewString(), ctx.Value(port.UserUUIDKey).(string), cmd.PlayerName, am[rng.Intn(len(am))])
	err = session.AddPlayer(player)
	if err != nil {
		return nil, wrap(err)
	}

	session.Params.Seed = rng.Int63()

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

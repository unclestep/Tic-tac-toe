package usecase

import (
	"context"
	"fmt"
	"math/rand"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type Start struct {
	sessionRepo port.SessionRepo
	botMover    port.BotMover
}

func NewStart(sessionRepo port.SessionRepo, botMover port.BotMover) *Start {
	return &Start{
		sessionRepo: sessionRepo,
		botMover:    botMover,
	}
}

func (uc *Start) Execute(ctx context.Context, cmd *port.StartCommand) (*model.Session, error) {
	wrap := func(err error) error {
		return fmt.Errorf("start (session %s, player %s): %w", cmd.SessionUUID, cmd.PlayerUUID, err)
	}

	session, err := uc.sessionRepo.Get(ctx, cmd.SessionUUID)
	if err != nil {
		return nil, wrap(err)
	}

	if session.State != model.StateLobby {
		return nil, wrap(port.ErrGameAlreadyStarted)
	}

	if !session.IsPlayerExist(cmd.PlayerUUID) {
		return nil, wrap(port.ErrPlayerNotBelongToSession)
	}

	rng := rand.New(rand.NewSource(session.Params.Seed))
	session.Start()

	if session.IsFull() {
		rng.Shuffle(len(session.Players), func(i, j int) {
			session.Players[i], session.Players[j] = session.Players[j], session.Players[i]
		})
		session.Players[0].Mark = model.MarkX
		session.Players[1].Mark = model.MarkO
	} else {
		if rng.Intn(2) == 1 {
			session.Players[0].Mark = model.MarkO
			if err := uc.botMover.MakeMove(session, session.Players[0].Mark.Opposite()); err != nil {
				return nil, wrap(err)
			}
		} else {
			session.Players[0].Mark = model.MarkX
		}
	}

	session.Params.Seed = rng.Int63()

	if err = uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

package usecase

import (
	"context"
	"fmt"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type Disconnect struct {
	sessionRepo port.SessionRepo
	botMover    port.BotMover
}

func NewDisconnect(sessionRepo port.SessionRepo, botMover port.BotMover) *Disconnect {
	return &Disconnect{
		sessionRepo: sessionRepo,
		botMover:    botMover,
	}
}

func (uc *Disconnect) Execute(ctx context.Context, cmd *port.DisconnectCommand) (*model.Session, error) {
	wrap := func(err error) error {
		return fmt.Errorf("disconnect (session %s, player %s): %w", cmd.SessionUUID, cmd.PlayerUUID, err)
	}

	sessions, err := uc.sessionRepo.Get(ctx, port.WithUUID(cmd.SessionUUID))
	if err != nil {
		return nil, wrap(err)
	}
	if len(sessions) != 1 {
		return nil, wrap(port.ErrReturnedNotOne)
	}
	session := sessions[0]

	turnPlayer, err := session.GetTurnPlayer()
	if err != nil {
		return nil, wrap(err)
	}

	if err := session.RemovePlayer(cmd.PlayerUUID); err != nil {
		return nil, wrap(err)
	}

	if session.State == model.StatePlaying && turnPlayer.UUID == cmd.PlayerUUID {
		if err := uc.botMover.MakeMove(session, turnPlayer.Mark); err != nil {
			return nil, wrap(err)
		}
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

package usecase

import (
	"context"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
)

type BotMover interface {
	MakeMove(session *model.Session, botMark model.Mark) error
}

type Disconnect struct {
	sessionRepo port.SessionRepo
	botMover    BotMover
}

func NewDisconnect(sessionRepo port.SessionRepo, botMover BotMover) *Disconnect {
	return &Disconnect{
		sessionRepo: sessionRepo,
		botMover:    botMover,
	}
}

func (uc *Disconnect) Execute(ctx context.Context, cmd *port.DisconnectCommand) (*model.Session, error) {
	wrap := func(err error) error {
		return fmt.Errorf("disconnect (session %s, player %s): %w", cmd.SessionUUID, cmd.PlayerUUID, err)
	}

	session, err := uc.sessionRepo.Get(ctx, cmd.SessionUUID)
	if err != nil {
		return nil, wrap(err)
	}

	turnPlayer, err := session.GetTurnPlayer()
	if err != nil {
		return nil, wrap(err)
	}

	if err := session.RemovePlayer(cmd.PlayerUUID); err != nil {
		return nil, wrap(err)
	}

	if turnPlayer.UUID == cmd.PlayerUUID {
		if err := uc.botMover.MakeMove(session, turnPlayer.Mark); err != nil {
			return nil, wrap(err)
		}
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

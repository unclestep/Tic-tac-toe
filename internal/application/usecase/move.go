package usecase

import (
	"context"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type HumanMover interface {
	MakeMove(session *model.Session, player *model.Player, p geometry.Point) error
}

type MakeMove struct {
	sessionRepo port.SessionRepo
	humanMover  HumanMover
	botMover    BotMover
}

func NewMakeMove(sessionRepo port.SessionRepo, humanMover HumanMover, botMover BotMover) *MakeMove {
	return &MakeMove{
		sessionRepo: sessionRepo,
		humanMover:  humanMover,
		botMover:    botMover,
	}
}

func (uc *MakeMove) Execute(ctx context.Context, cmd *port.MakeMoveCommand) (*model.Session, error) {
	wrap := func(err error) error {
		return fmt.Errorf("move (session %s, player %s): %w", cmd.SessionUUID, cmd.PlayerUUID, err)
	}

	session, err := uc.sessionRepo.Get(ctx, cmd.SessionUUID)
	if err != nil {
		return nil, wrap(err)
	}

	switch session.State {
	case model.StateLobby:
		return nil, wrap(port.ErrGameNotStarted)
	case model.StateGameOver:
		return nil, wrap(model.ErrGameAlreadyOver)
	}

	if !session.IsPlayerExist(cmd.PlayerUUID) {
		return nil, wrap(port.ErrPlayerNotBelongToSession)
	}

	turnPlayer, err := session.GetTurnPlayer()
	if err != nil {
		return nil, wrap(model.ErrPlayerNotFound)
	}

	if turnPlayer.UUID != cmd.PlayerUUID {
		return nil, wrap(port.ErrPlayerCantMakeMove)
	}

	if err := uc.humanMover.MakeMove(session, turnPlayer, cmd.MarkPos); err != nil {
		return nil, wrap(err)
	}

	if !session.IsFull() && session.State == model.StatePlaying {
		if err := uc.botMover.MakeMove(session, turnPlayer.Mark.Opposite()); err != nil {
			return nil, wrap(err)
		}
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, wrap(err)
	}

	return session, nil
}

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
	sessionRepo appPort.SessionRepo
	gameService servicePort.GameService
}

func NewMakeMove(sessionRepo appPort.SessionRepo, gameService servicePort.GameService) *MakeMove {
	return &MakeMove{
		sessionRepo: sessionRepo,
		gameService: gameService,
	}
}

func (uc *MakeMove) Execute(ctx context.Context, cmd *appPort.MakeMoveCommand) (*model.Session, error) {
	session, err := uc.sessionRepo.Get(ctx, cmd.SessionID)
	if err != nil {
		if errors.Is(err, appPort.ErrSessionNotFound) {
			return nil, fmt.Errorf("session %s not found", cmd.SessionID)
		}
		return nil, fmt.Errorf("connect: get session %s: %w", cmd.SessionID, err)
	}

	if session.State != model.StatePlaying {
		return nil, fmt.Errorf("%w: session %s", appPort.ErrGameNotStarted, session.UUID)
	}

	if !session.IsPlayerExist(cmd.PlayerID) {
		return nil, fmt.Errorf("%w: session %s, player %s", appPort.ErrPlayerNotBelongToSession, session.UUID, cmd.PlayerID)
	}

	player := session.GetTurnPlayer()
	if player.UUID != cmd.PlayerID {
		return nil, fmt.Errorf("%w: not player %s turn, session %s", appPort.ErrPlayerCantMakeMove, cmd.PlayerID, cmd.SessionID)
	}

	err = uc.gameService.MakeMovePlayer(session.Board, player, cmd.MarkPos)
	if err != nil {
		return nil, err
	}

	if len(session.Players) == 1 {
		err = uc.gameService.MakeMoveBot(session.Board, -player.Mark)
		session.Turn++
	}

	mark, state := uc.gameService.CheckWin(session.Board)
	if state == model.StatePlaying && len(session.Players) == 1 {
		session.Turn++
	} else if state == model.StateGameOver {
		session.Winner = session.DetermineWinner(mark)
		session.State = model.StateGameOver
	}

	err = uc.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

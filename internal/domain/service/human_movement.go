package service

import (
	"fmt"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type HumanMovement struct {
	winChecker *WinChecker
}

func NewHumanMovement(winChecker *WinChecker) *HumanMovement {
	return &HumanMovement{
		winChecker: winChecker,
	}
}

func (m *HumanMovement) MakeMove(session *model.Session, player *model.Player, p geometry.Point) error {
	if session.State == model.StateGameOver {
		return fmt.Errorf("make a move (session UUID %s): %w", session.UUID, model.ErrGameAlreadyOver)
	}
	if err := session.Board.SetMark(player.Mark, p); err != nil {
		return fmt.Errorf("make a move (player UUID %s, point %v): %w", player.UUID, p, err)
	}

	mark, state := m.winChecker.CheckWin(session.Board, session.Rules.WinLength)
	if state == model.StateGameOver {
		session.Winner = mark.String()
		session.State = model.StateGameOver
	}

	session.Turn++
	return nil
}

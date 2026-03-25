package service

import (
	"errors"
	"fmt"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type Movement struct {
	winChecker *WinChecker
}

func NewMovementService(winChecker *WinChecker) *Movement {
	return &Movement{
		winChecker: winChecker,
	}
}

var (
	ErrPlayerWOMark = errors.New("player does not have mark")
)

func (m *Movement) Make(session *model.Session, player *model.Player, p geometry.Point) error {
	board := session.Board

	if session.State == model.StateGameOver || len(board.GetEmptyCells()) == 0 {
		return fmt.Errorf("%w: %v", model.ErrGameAlreadyOver, session.UUID)
	}

	// Check if player has mark
	playerMark := player.Mark
	if playerMark != model.O && playerMark != model.X {
		return fmt.Errorf("%w: playerUUID - %s", ErrPlayerWOMark, player.UUID)
	}

	// Check if player set his mark on empty cell
	if mark := board.GetMark(p); mark != model.Empty {
		return fmt.Errorf("cell %v is already marked", p)
	}

	err := board.SetMark(playerMark, p)

	// Setpoint is out of bounds or player has invalid mark
	if err != nil {
		return fmt.Errorf("%w: playerUUID - %v, player's mark - %v, mark's place - %v", err, player.UUID, playerMark, p)
	}

	mark, state := m.winChecker.CheckWin(board)
	if state == model.StateGameOver {
		session.Winner = session.DetermineWinner(mark)
		session.State = state
	}
	session.Turn++

	return nil
}

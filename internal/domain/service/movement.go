package service

import (
	"errors"
	"fmt"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type Movement struct{}

func NewMovementService() *Movement {
	return &Movement{}
}

var (
	ErrPlayerWOMark = errors.New("player does not have mark")
)

func (m *Movement) Make(board *model.Board, player *model.Player, p geometry.Point) error {
	// Check if player has mark
	playerMark := player.Mark
	if playerMark != model.O && playerMark != model.X {
		return fmt.Errorf("%w: playerUUID - %s", ErrPlayerWOMark, player.UUID)
	}

	// Check if player set his mark on empty cell
	if mark := board.GetMark(p); mark != model.Empty {
		return fmt.Errorf("%w: cell %v is already marked", p)
	}

	err := board.Set(playerMark, p)

	// Setpoint is out of bounds or player has invalid mark
	if err != nil {
		return fmt.Errorf("%w: playerUUID - %s, player's mark - %s, mark's place - %s", err, player.UUID, playerMark, p)
	}

	return nil
}

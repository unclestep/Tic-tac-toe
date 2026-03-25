package port

import (
	"math/rand"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type GameMechanics interface {
	MakeMove(session *model.Session, player *model.Player, p geometry.Point) error
	Advance(session *model.Session) error
	Prepare(session *model.Session, rng *rand.Rand)
}

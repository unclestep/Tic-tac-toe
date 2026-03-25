package service

import (
	"math/rand"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type GameMechanics struct {
	advancer     *Advancer
	moveMaker    *Movement
	gamePreparer *GamePreparer
}

func NewGameMechanics(advancer *Advancer, moveMaker *Movement, gamePreparer *GamePreparer) *GameMechanics {
	return &GameMechanics{
		advancer:     advancer,
		moveMaker:    moveMaker,
		gamePreparer: gamePreparer,
	}
}

func (g *GameMechanics) MakeMove(session *model.Session, player *model.Player, p geometry.Point) error {
	return g.moveMaker.Make(session, player, p)
}

func (g *GameMechanics) Advance(session *model.Session) error {
	return g.advancer.Advance(session)
}

func (g *GameMechanics) Prepare(session *model.Session, rng *rand.Rand) {
	g.gamePreparer.Prepare(session, rng)

}

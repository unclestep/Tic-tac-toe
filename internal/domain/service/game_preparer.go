package service

import (
	"math/rand"
	"tictactoe/internal/domain/model"
)

type GamePreparer struct{}

func NewGamePreparer() *GamePreparer {
	return &GamePreparer{}
}

func (g *GamePreparer) Prepare(session *model.Session, rng *rand.Rand) {
	if len(session.Players) == 0 {
		return
	}

	session.Board.Clear()

	marks := [2]model.Mark{model.O, model.X}
	rng.Shuffle(2, func(i, j int) {
		marks[i], marks[j] = marks[j], marks[i]
	})

	if len(session.Players) > 1 {
		rng.Shuffle(2, func(i, j int) {
			session.Players[i], session.Players[j] = session.Players[j], session.Players[i]
		})
	}

	for i, player := range session.Players {
		player.Mark = marks[i]
	}

	session.Params.Seed = int64(rng.Int())
}

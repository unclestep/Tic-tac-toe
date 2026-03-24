package port

import (
	"math/rand"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type GameService interface {
	MakeMoveBot(board *model.Board, botMark model.Mark) error
	MakeMovePlayer(board *model.Board, player *model.Player, p geometry.Point) error
	CheckWin(board *model.Board) (model.Mark, model.State)
	PrepareGame(session *model.Session, rng *rand.Rand)
}

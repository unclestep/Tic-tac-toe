package service

import (
	"math/rand"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type GameService struct {
	botMove      *Minimax
	playerMove   *Movement
	winChecker   *WinChecker
	gamePreparer *GamePreparer
}

func NewGameService(botMove *Minimax, playerMove *Movement, winChecker *WinChecker, gamePreparer *GamePreparer) *GameService {
	return &GameService{
		botMove:      botMove,
		playerMove:   playerMove,
		winChecker:   winChecker,
		gamePreparer: gamePreparer,
	}
}

func (g *GameService) MakeMoveBot(board *model.Board, botMark model.Mark) error {
	return g.botMove.BestMove(board, botMark)
}

func (g *GameService) MakeMovePlayer(board *model.Board, player *model.Player, p geometry.Point) error {
	return g.playerMove.Make(board, player, p)
}

func (g *GameService) CheckWin(board *model.Board) (model.Mark, model.State) {
	return g.winChecker.CheckWin(board)
}

func (g *GameService) PrepareGame(session *model.Session, rng *rand.Rand) {
	g.gamePreparer.Prepare(session, rng)

}

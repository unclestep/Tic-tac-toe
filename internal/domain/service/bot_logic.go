package service

import (
	"fmt"
	"math"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type BotLogic struct {
	winChecker *WinChecker
}

func NewBotLogic(winChecker *WinChecker) *BotLogic {
	return &BotLogic{
		winChecker: winChecker,
	}
}

func (m *BotLogic) BestMove(session *model.Session, botMark model.Mark) error {
	board := session.Board

	emptyCells := board.GetEmptyCells()
	if session.State == model.StateGameOver || len(emptyCells) == 0 {
		return fmt.Errorf("%w: %s", model.ErrGameAlreadyOver, session.UUID)
	}

	bestScore := model.Mark(math.MinInt8)
	bestPoint := geometry.Point{}
	maximizing := true
	// X - maximizing, O - minimizing
	if botMark == model.O {
		bestScore = model.Mark(math.MaxInt8)
		maximizing = false
	}

	for _, cell := range emptyCells {
		_ = board.SetMark(botMark, cell)
		score := m.perform(board, !maximizing, model.Mark(math.MinInt8), model.Mark(math.MaxInt8))
		_ = board.ClearMark(cell)

		if maximizing {
			if score > bestScore {
				bestScore = score
				bestPoint = cell
			}
		} else {
			if score < bestScore {
				bestScore = score
				bestPoint = cell
			}
		}
	}

	_ = board.SetMark(botMark, bestPoint)

	mark, state := m.winChecker.CheckWin(board)
	if state == model.StateGameOver {
		session.Winner = session.DetermineWinner(mark)
		session.State = state
	}
	session.Turn++

	return nil
}

func (m *BotLogic) perform(board *model.Board, maximizingPlayer bool, alpha, beta model.Mark) model.Mark {
	if winner, state := m.winChecker.CheckWin(board); state == model.StateGameOver {
		return winner
	}

	emptyCells := board.GetEmptyCells()
	if len(emptyCells) == 0 {
		return 0
	}

	if maximizingPlayer {
		maxEval := model.Mark(math.MinInt8)
		for _, cell := range emptyCells {
			_ = board.SetMark(model.X, cell)
			eval := m.perform(board, false, alpha, beta)
			_ = board.ClearMark(cell)

			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval
	}

	minEval := model.Mark(math.MaxInt8)
	for _, cell := range emptyCells {
		_ = board.SetMark(model.O, cell)
		eval := m.perform(board, true, alpha, beta)
		_ = board.ClearMark(cell)

		minEval = min(minEval, eval)
		beta = min(beta, eval)
		if beta <= alpha {
			break
		}
	}
	return minEval
}

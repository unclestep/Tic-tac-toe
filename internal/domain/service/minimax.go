package service

import (
	"fmt"
	"math"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type Minimax struct {
	winChecker *WinChecker
}

func NewMiniMax(winChecker *WinChecker) *Minimax {
	return &Minimax{
		winChecker: winChecker,
	}
}

func (m *Minimax) BestMove(board *model.Board, botMark model.Mark) error {
	emptyCells := board.GetEmptyCells()
	if len(emptyCells) == 0 {
		return fmt.Errorf("Game over, no available cell")
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
		board.Set(botMark, cell)
		score := m.perform(board, !maximizing, model.Mark(math.MinInt8), model.Mark(math.MaxInt8))
		board.Set(model.Empty, cell)

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

	board.Set(botMark, bestPoint)
	return nil
}

func (m *Minimax) perform(board *model.Board, maximizingPlayer bool, alpha, beta model.Mark) model.Mark {
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
			testBoard := board.Clone()
			testBoard.Set(model.X, cell)

			eval := m.perform(testBoard, false, alpha, beta)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := model.Mark(math.MaxInt8)
		for _, cell := range emptyCells {
			testBoard := board.Clone()
			testBoard.Set(model.O, cell)

			eval := m.perform(testBoard, true, alpha, beta)
			minEval = min(minEval, eval)
			beta = min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

package service

import (
	"fmt"
	"math"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type BotMovement struct {
	winChecker *WinChecker
}

func NewBotMovement(winChecker *WinChecker) *BotMovement {
	return &BotMovement{
		winChecker: winChecker,
	}
}

func (m *BotMovement) MakeMove(session *model.Session, botMark model.Mark) error {
	board := session.Board

	emptyCells := board.GetEmptyCells()
	if session.State == model.StateGameOver || len(emptyCells) == 0 {
		return fmt.Errorf("best move (session %s): %w", session.UUID, model.ErrGameAlreadyOver)
	}

	bestScore := model.Mark(math.MinInt8)
	bestPoint := geometry.Point{}
	maximizing := true
	// X - maximizing, O - minimizing
	if botMark == model.MarkO {
		bestScore = model.Mark(math.MaxInt8)
		maximizing = false
	}

	for _, cell := range emptyCells {
		err := board.SetMark(botMark, cell)
		if err != nil {
			panic("best move: set mark precondition violation")
		}
		score := m.perform(session, !maximizing, model.Mark(math.MinInt8), model.Mark(math.MaxInt8))
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

	mark, state := m.winChecker.CheckWin(board, session.Rules.WinLength)
	if state == model.StateGameOver {
		session.Winner = mark.String()
		session.State = model.StateGameOver
	}
	session.Turn++

	return nil
}

func (m *BotMovement) perform(session *model.Session, maximizingPlayer bool, alpha, beta model.Mark) model.Mark {
	board := session.Board

	if winner, state := m.winChecker.CheckWin(board, session.Rules.WinLength); state == model.StateGameOver {
		return winner
	}

	emptyCells := board.GetEmptyCells()
	if len(emptyCells) == 0 {
		return 0
	}

	if maximizingPlayer {
		maxEval := model.Mark(math.MinInt8)
		for _, cell := range emptyCells {
			_ = board.SetMark(model.MarkX, cell)
			eval := m.perform(session, false, alpha, beta)
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
		_ = board.SetMark(model.MarkO, cell)
		eval := m.perform(session, true, alpha, beta)
		_ = board.ClearMark(cell)

		minEval = min(minEval, eval)
		beta = min(beta, eval)
		if beta <= alpha {
			break
		}
	}
	return minEval
}

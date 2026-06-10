package service

import (
	"fmt"
	"math"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type BotMovement struct {
	winChecker *WinChecker
	heuristic  *Heuristic
}

func NewBotMovement(winChecker *WinChecker, heuristic *Heuristic) *BotMovement {
	return &BotMovement{
		winChecker: winChecker,
		heuristic:  heuristic,
	}
}

const (
	searchDepth int     = 4
	winScore    float64 = 1_000_000
	drawScore   float64 = 0
)

func (m *BotMovement) MakeMove(session *model.Session, botMark model.Mark) error {
	board := session.Board
	emptyCells := board.GetEmptyCells()

	if session.State == model.StateGameOver || len(emptyCells) == 0 {
		return fmt.Errorf("bot make move (session %s): %w", session.UUID, model.ErrGameAlreadyOver)
	}

	_, bestCell := m.negamax(session, botMark, 0, searchDepth, math.Inf(-1), math.Inf(1))
	if bestCell == nil {
		return fmt.Errorf("bot make move: no available cells")
	}

	err := board.SetMark(botMark, *bestCell)
	if err != nil {
		panic(fmt.Sprintf("bot make move: %s", err))
	}

	mark, state := m.winChecker.CheckWin(board, session.Rules.WinLength)
	if state == model.StateGameOver {
		if mark != model.MarkEmpty {
			session.Winner = model.NewBot(mark)
		}
		session.State = model.StateGameOver
	}
	session.Turn++

	return nil
}

func (m *BotMovement) negamax(session *model.Session, mark model.Mark, depth, remainingDepth int, alpha, beta float64) (float64, *geometry.Point) {
	board := session.Board

	if winner, state := m.winChecker.CheckWin(board, session.Rules.WinLength); state == model.StateGameOver {
		return m.terminalScore(winner, mark, depth), nil
	}

	emptyCells := board.GetEmptyCells()
	if len(emptyCells) == 0 {
		return drawScore, nil
	}

	if remainingDepth == 0 {
		return m.heuristic.Evaluate(board, mark, session.Rules.WinLength), nil
	}

	bestScore := math.Inf(-1)
	var bestCell geometry.Point

	for _, cell := range emptyCells {
		if err := board.SetMark(mark, cell); err != nil {
			panic(fmt.Sprintf("negamax: %s", err))
		}

		// Recurse from opponent's perspective: invert score and swap-negate window.
		score, _ := m.negamax(session, mark.Opposite(), depth+1, remainingDepth-1, -beta, -alpha)
		score = -score

		if err := board.ClearMark(cell); err != nil {
			panic(fmt.Sprintf("negamax: %s", err))
		}

		if score > bestScore {
			bestScore = score
			bestCell = cell
		}

		alpha = max(alpha, score)
		if alpha >= beta {
			break
		}
	}

	return bestScore, &bestCell
}

func (m *BotMovement) terminalScore(winner, mark model.Mark, depth int) float64 {
	switch winner {
	case mark:
		return winScore - float64(depth)
	case mark.Opposite():
		return -winScore + float64(depth)
	}
	return drawScore
}

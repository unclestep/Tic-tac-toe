package service

import (
	"fmt"
	"math"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type Heuristic struct{}

func NewHeuristic() *Heuristic {
	return &Heuristic{}
}

func (h *Heuristic) Evaluate(board *model.Board, mark model.Mark, win int) float64 {
	var score float64

	for r := 0; r < board.Height; r++ {
		for c := 0; c < board.Width; c++ {
			start := geometry.Point{X: c, Y: r}
			for _, dir := range geometry.GetAllDirs() {
				if h.lineFits(board, start, dir, win) {
					score += h.evaluateLine(board, start, dir, win, mark, mark.Opposite())
				}
			}
		}
	}

	return score
}

func (h *Heuristic) lineFits(board *model.Board, start, dir geometry.Point, length int) bool {
	endX := start.X + dir.X*(length-1)
	endY := start.Y + dir.Y*(length-1)
	return endX >= 0 && endX < board.Width && endY >= 0 && endY < board.Height
}

func (h *Heuristic) evaluateLine(board *model.Board, start, dir geometry.Point, length int, curMark, oppMark model.Mark) float64 {
	var selfCount, oppCount float64

	for i := range length {
		p := geometry.Point{X: start.X + dir.X*i, Y: start.Y + dir.Y*i}
		m, err := board.GetMark(p)
		if err != nil {
			panic(fmt.Sprintf("bot evaluate line: %s", err))
		}

		switch m {
		case curMark:
			selfCount++
		case oppMark:
			oppCount++
		}
	}

	if selfCount > 0 && oppCount > 0 {
		return 0
	}
	if selfCount > 0 {
		return math.Pow(10, selfCount)
	}
	if oppCount > 0 {
		return -math.Pow(10, oppCount)
	}

	return 0
}

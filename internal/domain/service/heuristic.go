package service

import (
	"fmt"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type Heuristic struct {
	weights []int
}

func NewHeuristic() *Heuristic {
	return &Heuristic{
		weights: []int{0, 1, 10, 100, 1_000, 10_000, 100_000},
	}
}

func (h *Heuristic) Evaluate(board *model.Board, mark model.Mark, win int) int {
	score := 0

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

func (h *Heuristic) evaluateLine(board *model.Board, start, dir geometry.Point, length int, curMark, oppMark model.Mark) int {
	selfCount, oppCount := 0, 0
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
		return h.weights[selfCount]
	}
	if oppCount > 0 {
		return -h.weights[oppCount]
	}

	return 0
}

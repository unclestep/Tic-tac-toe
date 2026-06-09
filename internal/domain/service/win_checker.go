package service

import (
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type WinChecker struct{}

func NewWinChecker() *WinChecker {
	return &WinChecker{}
}

func (w *WinChecker) CheckWin(board *model.Board, win int) (model.Mark, model.State) {
	width := board.Width
	height := board.Height

	// Horizontal
	if width >= win {
		for r := range height {
			mark := w.checkLine(board, win, geometry.NewPoint(0, r), geometry.NewPoint(1, 0))
			if mark != model.MarkEmpty {
				return mark, model.StateGameOver
			}
		}
	}

	// Vertical
	if height >= win {
		for c := range width {
			mark := w.checkLine(board, win, geometry.NewPoint(c, 0), geometry.NewPoint(0, 1))
			if mark != model.MarkEmpty {
				return mark, model.StateGameOver
			}
		}
	}

	// Main diagonal
	for r := range height {
		length := min(height-r, width)
		if length >= win {
			mark := w.checkLine(board, win, geometry.NewPoint(0, r), geometry.NewPoint(1, 1))
			if mark != model.MarkEmpty {
				return mark, model.StateGameOver
			}
		}
	}

	for c := 1; c < width; c++ {
		length := min(height, width-c)
		if length >= win {
			mark := w.checkLine(board, win, geometry.NewPoint(c, 0), geometry.NewPoint(1, 1))
			if mark != model.MarkEmpty {
				return mark, model.StateGameOver
			}
		}
	}

	// Secondary diagonal
	for r := range height {
		length := min(height-r, width)
		if length >= win {
			mark := w.checkLine(board, win, geometry.NewPoint(width-1, r), geometry.NewPoint(-1, 1))
			if mark != model.MarkEmpty {
				return mark, model.StateGameOver
			}
		}
	}

	for c := 0; c < width-1; c++ {
		length := min(height, c+1)
		if length >= win {
			mark := w.checkLine(board, win, geometry.NewPoint(c, 0), geometry.NewPoint(-1, 1))
			if mark != model.MarkEmpty {
				return mark, model.StateGameOver
			}
		}
	}

	if len(board.GetEmptyCells()) == 0 {
		return model.MarkEmpty, model.StateGameOver
	}

	return model.MarkEmpty, model.StatePlaying
}

func (w *WinChecker) checkLine(board *model.Board, win int, start geometry.Point, vector geometry.Point) model.Mark {
	width := board.Width
	height := board.Height

	cur := model.MarkEmpty
	count := 0
	x, y := start.X, start.Y
	for y >= 0 && y < height && x >= 0 && x < width {
		mark, _ := board.GetMark(geometry.NewPoint(x, y))

		if mark == model.MarkEmpty {
			cur = model.MarkEmpty
			count = 0
		} else if mark != cur {
			cur = mark
			count = 1
		} else {
			count++
		}

		if count == win {
			return cur
		}

		x += vector.X
		y += vector.Y
	}

	return model.MarkEmpty
}

package model

import (
	"errors"
	"fmt"
	"slices"
	"tictactoe/pkg/geometry"
)

var (
	ErrSetPointOutOfBounds = errors.New("set point is out of bounds")
	ErrInvalidMark         = errors.New("invalid mark")
)

type Board struct {
	Width  int
	Height int
	cells  []Mark
}

type Mark int8

const (
	Empty Mark = 0
	X     Mark = 1
	O     Mark = -1
)

func NewBoard(width, height int) *Board {
	return &Board{
		Width:  width,
		Height: height,
		cells:  make([]Mark, width*height),
	}
}

func NewBoardFromCells(width, height int, cells []Mark) *Board {
	if len(cells) != width*height {
		return nil
	}

	return &Board{
		Width:  width,
		Height: height,
		cells:  cells,
	}
}

func (b *Board) Clone() *Board {
	return &Board{
		Width:  b.Width,
		Height: b.Height,
		cells:  slices.Clone(b.cells),
	}
}

func (b *Board) CloneCells() []Mark {
	return slices.Clone(b.cells)
}

func (b *Board) Clear() {
	clear(b.cells)
}

func (b *Board) Set(m Mark, p geometry.Point) error {
	if !b.InBounds(p) {
		return fmt.Errorf("%w: point %v, board width %v, board height %v", ErrSetPointOutOfBounds, p, b.Width, b.Height)
	}
	if !b.IsValidMark(m) {
		return fmt.Errorf("%w: given mark %v, valid: %v, %v", ErrInvalidMark, m, X, O)
	}
	b.cells[b.Width*p.Y+p.X] = m
	return nil
}

func (b *Board) GetEmptyCells() []geometry.Point {
	emptyCells := make([]geometry.Point, 0, b.Height*b.Width)
	for i, cell := range b.cells {
		if cell == Empty {
			p := geometry.NewPoint(i%b.Width, i/b.Width)
			emptyCells = append(emptyCells, p)
		}
	}
	return emptyCells
}

func (b *Board) GetMark(p geometry.Point) Mark {
	if !b.InBounds(p) {
		return Empty
	}
	return b.cells[b.Width*p.Y+p.X]
}

func (b *Board) InBounds(p geometry.Point) bool {
	return p.X >= 0 && p.X < b.Width && p.Y >= 0 && p.Y < b.Height
}

func (b *Board) IsValidMark(m Mark) bool {
	return m == Empty || m == X || m == O
}

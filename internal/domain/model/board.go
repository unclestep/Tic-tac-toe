package model

import (
	"errors"
	"fmt"
	"slices"
	"tictactoe/pkg/geometry"
)

var (
	ErrInvalidConstructorArgs = errors.New("invalid board constructor args")
	ErrSetPointOutOfBounds    = errors.New("set point is out of bounds")
	ErrInvalidMark            = errors.New("invalid mark")
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

func NewBoard(width, height int) (*Board, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("%w: width %v, height %v", ErrInvalidConstructorArgs, width, height)
	}
	return &Board{
		Width:  width,
		Height: height,
		cells:  make([]Mark, width*height),
	}, nil
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

func (b *Board) ClearMark(p geometry.Point) error {
	if !b.InBounds(p) {
		return fmt.Errorf("%w: point %v, board width %v, board height %v", ErrSetPointOutOfBounds, p, b.Width, b.Height)
	}
	b.cells[b.Width*p.Y+p.X] = Empty
	return nil
}

func (b *Board) SetMark(m Mark, p geometry.Point) error {
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
	return m == X || m == O
}

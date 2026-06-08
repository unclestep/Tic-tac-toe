package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/domain/model"
	"tictactoe/internal/domain/service"
	"tictactoe/pkg/geometry"
)

func boardFromStrings(t *testing.T, rows []string) *model.Board {
	t.Helper()
	height := len(rows)
	width := len(rows[0])
	b := model.NewBoard(width, height)
	for y, row := range rows {
		for x, ch := range row {
			switch ch {
			case 'X':
				require.NoError(t, b.SetMark(model.MarkX, geometry.NewPoint(x, y)))
			case 'O':
				require.NoError(t, b.SetMark(model.MarkO, geometry.NewPoint(x, y)))
			}
		}
	}
	return b
}

func TestCheckWinEmptyBoard(t *testing.T) {
	b := model.NewBoard(3, 3)

	mark, state := service.NewWinChecker().CheckWin(b, 3)

	assert.Equal(t, model.MarkEmpty, mark)
	assert.Equal(t, model.StatePlaying, state)
}

func TestCheckWinHorizontal(t *testing.T) {
	cases := []struct {
		name  string
		board []string
		mark  model.Mark
	}{
		{
			"X wins row 0",
			[]string{
				"XXX",
				"O..",
				"O..",
			},
			model.MarkX,
		},
		{
			"X wins row 1",
			[]string{
				"O..",
				"XXX",
				"O..",
			},
			model.MarkX,
		},
		{
			"X wins row 2",
			[]string{
				"O..",
				"O..",
				"XXX",
			},
			model.MarkX,
		},
		{
			"O wins row 1",
			[]string{
				"X..",
				"OOO",
				"X..",
			},
			model.MarkO,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, tc.mark, mark)
			assert.Equal(t, model.StateGameOver, state)
		})
	}
}

func TestCheckWinVertical(t *testing.T) {
	cases := []struct {
		name  string
		board []string
		mark  model.Mark
	}{
		{
			"X wins col 0",
			[]string{
				"XO.",
				"XO.",
				"X..",
			},
			model.MarkX,
		},
		{
			"X wins col 1",
			[]string{
				"OX.",
				"OX.",
				".X.",
			},
			model.MarkX,
		},
		{
			"X wins col 2",
			[]string{
				"OOX",
				"..X",
				"..X",
			},
			model.MarkX,
		},
		{
			"O wins col 0",
			[]string{
				"OX.",
				"OX.",
				"O..",
			},
			model.MarkO,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, tc.mark, mark)
			assert.Equal(t, model.StateGameOver, state)
		})
	}
}

func TestCheckWinMainDiagonal(t *testing.T) {
	cases := []struct {
		name   string
		board  []string
		winLen int
		mark   model.Mark
	}{
		{
			"X wins corner diagonal on 3x3",
			[]string{
				"XO.",
				"OX.",
				"..X",
			},
			3,
			model.MarkX,
		},
		{
			"O wins corner diagonal on 3x3",
			[]string{
				"OX.",
				"XO.",
				"..O",
			},
			3,
			model.MarkO,
		},
		{
			"X wins off-corner diagonal from (0,1) on 4x4",
			[]string{
				"....",
				"X...",
				".X..",
				"..X.",
			},
			3,
			model.MarkX,
		},
		{
			"X wins off-corner diagonal from (1,0) on 4x4",
			[]string{
				".X..",
				"..X.",
				"...X",
				"....",
			},
			3,
			model.MarkX,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, tc.mark, mark)
			assert.Equal(t, model.StateGameOver, state)
		})
	}
}

func TestCheckWinSecondaryDiagonal(t *testing.T) {
	cases := []struct {
		name   string
		board  []string
		winLen int
		mark   model.Mark
	}{
		{
			"X wins corner diagonal on 3x3",
			[]string{
				"OOX",
				".X.",
				"X..",
			},
			3,
			model.MarkX,
		},
		{
			"O wins corner diagonal on 3x3",
			[]string{
				"XXO",
				".O.",
				"O..",
			},
			3,
			model.MarkO,
		},
		{
			"X wins off-corner diagonal from (3,1) on 4x4",
			[]string{
				"....",
				"...X",
				"..X.",
				".X..",
			},
			3,
			model.MarkX,
		},
		{
			"X wins off-corner diagonal from (2,0) on 4x4",
			[]string{
				"..X.",
				".X..",
				"X...",
				"....",
			},
			3,
			model.MarkX,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, tc.mark, mark)
			assert.Equal(t, model.StateGameOver, state)
		})
	}
}

func TestCheckWinDraw(t *testing.T) {
	b := boardFromStrings(t, []string{
		"XOX",
		"XOO",
		"OXX",
	})

	mark, state := service.NewWinChecker().CheckWin(b, 3)

	assert.Equal(t, model.MarkEmpty, mark)
	assert.Equal(t, model.StateGameOver, state)
}

func TestCheckWinInProgress(t *testing.T) {
	cases := []struct {
		name  string
		board []string
	}{
		{
			"single mark",
			[]string{
				"X..",
				"...",
				"...",
			},
		},
		{
			"partial board no winner",
			[]string{
				"XO.",
				"OX.",
				"...",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, model.MarkEmpty, mark)
			assert.Equal(t, model.StatePlaying, state)
		})
	}
}

func TestCheckWinBrokenSequence(t *testing.T) {
	cases := []struct {
		name  string
		board []string
	}{
		{
			"horizontal broken by opponent",
			[]string{
				"XOX",
				"...",
				"...",
			},
		},
		{
			"horizontal broken by empty",
			[]string{
				"X.X",
				"...",
				"...",
			},
		},
		{
			"vertical broken by opponent",
			[]string{
				"X..",
				"O..",
				"X..",
			},
		},
		{
			"main diagonal broken by opponent",
			[]string{
				"X..",
				".O.",
				"..X",
			},
		},
		{
			"secondary diagonal broken by opponent",
			[]string{
				"..X",
				".O.",
				"X..",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, model.MarkEmpty, mark)
			assert.Equal(t, model.StatePlaying, state)
		})
	}
}

func TestCheckWinNonSquareBoards(t *testing.T) {
	cases := []struct {
		name   string
		board  []string
		winLen int
		mark   model.Mark
		state  model.State
	}{
		{
			"narrow board (2x3): X wins vertically",
			[]string{
				"XO",
				"X.",
				"X.",
			},
			3,
			model.MarkX,
			model.StateGameOver,
		},
		{
			"short board (3x2): O wins horizontally",
			[]string{
				"OOO",
				"XX.",
			},
			3,
			model.MarkO,
			model.StateGameOver,
		},
		{
			"5x5 X wins in middle of row",
			[]string{
				"O....",
				".XXX.",
				".OO..",
				".....",
				".....",
			},
			3,
			model.MarkX,
			model.StateGameOver,
		},
		{
			"5x5 X wins in middle of col",
			[]string{
				"O....",
				".X...",
				".X...",
				".X...",
				".....",
			},
			3,
			model.MarkX,
			model.StateGameOver,
		},
		{
			"5x5 WinLength=4 no winner",
			[]string{
				"XO...",
				".X...",
				"..O..",
				".....",
				".....",
			},
			4,
			model.MarkEmpty,
			model.StatePlaying,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := boardFromStrings(t, tc.board)

			mark, state := service.NewWinChecker().CheckWin(b, 3)

			assert.Equal(t, tc.mark, mark)
			assert.Equal(t, tc.state, state)
		})
	}
}

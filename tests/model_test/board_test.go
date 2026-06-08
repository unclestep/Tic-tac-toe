package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

func newBoard(t *testing.T, width, height int) *model.Board {
	t.Helper()
	return model.NewBoard(width, height)
}

func TestNewBoardInvalidDimensions(t *testing.T) {
	cases := []struct {
		name          string
		width, height int
	}{
		{"zero width", 0, 3},
		{"zero height", 3, 0},
		{"both zero", 0, 0},
		{"negative width", -1, 3},
		{"negative height", 3, -1},
		{"both negative", -1, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Panics(t, func() {
				model.NewBoard(tc.width, tc.height)
			})
		})
	}
}

func TestNewBoardValid(t *testing.T) {
	b := model.NewBoard(3, 3)
	require.NotNil(t, b)

	assert.Equal(t, 3, b.Width)
	assert.Equal(t, 3, b.Height)

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			mark, err := b.GetMark(geometry.NewPoint(x, y))
			require.NoError(t, err)
			assert.Equal(t, model.MarkEmpty, mark)
		}
	}
}

func TestNewBoardMinimalValid(t *testing.T) {
	b := model.NewBoard(1, 1)
	assert.NotNil(t, b)
}

func TestNewBoardFromCellsWrongSize(t *testing.T) {
	cases := []struct {
		name          string
		width, height int
		cells         []model.Mark
	}{
		{"too few cells", 3, 3, make([]model.Mark, 8)},
		{"too many cells", 3, 3, make([]model.Mark, 10)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Panics(t, func() {
				model.NewBoardFromCells(tc.width, tc.height, tc.cells)
			})
		})
	}
}

func TestNewBoardFromCellsValid(t *testing.T) {
	cells := []model.Mark{model.MarkX, model.MarkO, model.MarkO, model.MarkX}

	b := model.NewBoardFromCells(2, 2, cells)
	require.NotNil(t, b)

	assert.Equal(t, 2, b.Width)
	assert.Equal(t, 2, b.Height)

	expected := []struct {
		p    geometry.Point
		mark model.Mark
	}{
		{geometry.NewPoint(0, 0), model.MarkX},
		{geometry.NewPoint(1, 0), model.MarkO},
		{geometry.NewPoint(0, 1), model.MarkO},
		{geometry.NewPoint(1, 1), model.MarkX},
	}

	for _, tc := range expected {
		got, err := b.GetMark(tc.p)
		require.NoError(t, err)
		assert.Equal(t, tc.mark, got)
	}
}

func TestCloneAllSame(t *testing.T) {
	original := newBoard(t, 3, 3)
	p := geometry.NewPoint(1, 1)
	require.NoError(t, original.SetMark(model.MarkX, p))

	clone := original.Clone()

	assert.Equal(t, *original, *clone)
}

func TestCloneMutationDoesNotAffectOriginal(t *testing.T) {
	original := newBoard(t, 3, 3)
	p := geometry.NewPoint(0, 0)
	require.NoError(t, original.SetMark(model.MarkX, p))

	clone := original.Clone()
	clone.Cells[0] = model.MarkO

	mark, err := original.GetMark(p)
	require.NoError(t, err)
	assert.Equal(t, model.MarkX, mark)
}

func TestClearAllCellsBecomeEmpty(t *testing.T) {
	b := newBoard(t, 3, 3)
	require.NoError(t, b.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, b.SetMark(model.MarkO, geometry.NewPoint(2, 2)))

	b.Clear()

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			mark, err := b.GetMark(geometry.NewPoint(x, y))
			require.NoError(t, err)
			assert.Equal(t, model.MarkEmpty, mark, "cell (%d,%d) must be empty after Clear", x, y)
		}
	}
}

func TestClearMarkValidPointWithMark(t *testing.T) {
	b := newBoard(t, 3, 3)
	p := geometry.NewPoint(1, 1)
	require.NoError(t, b.SetMark(model.MarkX, p))

	err := b.ClearMark(p)
	require.NoError(t, err)

	mark, err := b.GetMark(p)
	require.NoError(t, err)
	assert.Equal(t, model.MarkEmpty, mark)
}

func TestClearMarkValidPointAlreadyEmpty(t *testing.T) {
	b := newBoard(t, 3, 3)
	p := geometry.NewPoint(1, 1)

	err := b.ClearMark(p)
	require.NoError(t, err)

	mark, err := b.GetMark(p)
	require.NoError(t, err)
	assert.Equal(t, model.MarkEmpty, mark)
}

func TestClearMarkOutOfBounds(t *testing.T) {
	b := newBoard(t, 3, 3)
	err := b.ClearMark(geometry.NewPoint(5, 5))
	assert.ErrorIs(t, err, model.ErrOutOfBounds)
}

func TestSetMarkValid(t *testing.T) {
	cases := []struct {
		name string
		mark model.Mark
	}{
		{"X", model.MarkX},
		{"O", model.MarkO},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := newBoard(t, 3, 3)
			p := geometry.NewPoint(1, 1)

			err := b.SetMark(tc.mark, p)
			require.NoError(t, err)

			got, err := b.GetMark(p)
			require.NoError(t, err)
			assert.Equal(t, tc.mark, got)
		})
	}
}

func TestSetMarkInvalidMark(t *testing.T) {
	cases := []struct {
		name string
		mark model.Mark
	}{
		{"empty mark", model.MarkEmpty},
		{"arbitrary positive", model.Mark(2)},
		{"arbitrary negative", model.Mark(-2)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := newBoard(t, 3, 3)
			assert.Panics(t, func() {
				err := b.SetMark(tc.mark, geometry.NewPoint(0, 0))
				assert.Error(t, err)
			})
		})
	}
}

func TestSetMarkOutOfBounds(t *testing.T) {
	b := newBoard(t, 3, 3)
	err := b.SetMark(model.MarkX, geometry.NewPoint(5, 5))
	assert.ErrorIs(t, err, model.ErrOutOfBounds)
}

func TestGetEmptyCellsEmptyBoard(t *testing.T) {
	const width, height = 3, 3
	b := newBoard(t, width, height)

	empty := b.GetEmptyCells()
	assert.Len(t, empty, width*height)
}

func TestGetEmptyCellsFullBoard(t *testing.T) {
	b := newBoard(t, 2, 2)
	require.NoError(t, b.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, b.SetMark(model.MarkO, geometry.NewPoint(1, 0)))
	require.NoError(t, b.SetMark(model.MarkX, geometry.NewPoint(0, 1)))
	require.NoError(t, b.SetMark(model.MarkO, geometry.NewPoint(1, 1)))

	empty := b.GetEmptyCells()
	assert.Empty(t, empty)
}

func TestGetEmptyCellsPartialBoard(t *testing.T) {
	b := newBoard(t, 3, 3)
	require.NoError(t, b.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, b.SetMark(model.MarkO, geometry.NewPoint(2, 2)))

	empty := b.GetEmptyCells()
	assert.Len(t, empty, 7)
	assert.NotContains(t, empty, geometry.NewPoint(0, 0))
	assert.NotContains(t, empty, geometry.NewPoint(2, 2))
}

func TestGetEmptyCellsCorrectCoordinates(t *testing.T) {
	b := newBoard(t, 2, 3)
	require.NoError(t, b.SetMark(model.MarkX, geometry.NewPoint(0, 0)))

	empty := b.GetEmptyCells()

	expected := []geometry.Point{
		geometry.NewPoint(1, 0),
		geometry.NewPoint(0, 1),
		geometry.NewPoint(1, 1),
		geometry.NewPoint(0, 2),
		geometry.NewPoint(1, 2),
	}
	assert.ElementsMatch(t, expected, empty)
}

func TestGetMarkOutOfBounds(t *testing.T) {
	b := newBoard(t, 3, 3)

	_, err := b.GetMark(geometry.NewPoint(5, 5))
	assert.ErrorIs(t, err, model.ErrOutOfBounds)
}

func TestGetMarkCornerCells(t *testing.T) {
	const width, height = 3, 4
	b := newBoard(t, width, height)

	corners := []struct {
		name string
		p    geometry.Point
		mark model.Mark
	}{
		{"top-left", geometry.NewPoint(0, 0), model.MarkX},
		{"top-right", geometry.NewPoint(width-1, 0), model.MarkO},
		{"bottom-left", geometry.NewPoint(0, height-1), model.MarkX},
		{"bottom-right", geometry.NewPoint(width-1, height-1), model.MarkO},
	}

	for _, tc := range corners {
		require.NoError(t, b.SetMark(tc.mark, tc.p))
	}

	for _, tc := range corners {
		t.Run(tc.name, func(t *testing.T) {
			got, err := b.GetMark(tc.p)
			require.NoError(t, err)
			assert.Equal(t, tc.mark, got)
		})
	}
}

func TestInBounds(t *testing.T) {
	const width, height = 3, 4
	b := newBoard(t, width, height)

	cases := []struct {
		name     string
		p        geometry.Point
		expected bool
	}{
		{"top-left corner", geometry.NewPoint(0, 0), true},
		{"top-right corner", geometry.NewPoint(width-1, 0), true},
		{"bottom-left corner", geometry.NewPoint(0, height-1), true},
		{"bottom-right corner", geometry.NewPoint(width-1, height-1), true},
		{"center", geometry.NewPoint(1, 2), true},
		{"negative X", geometry.NewPoint(-1, 0), false},
		{"negative Y", geometry.NewPoint(0, -1), false},
		{"both negative", geometry.NewPoint(-1, -1), false},
		{"X equals width", geometry.NewPoint(width, 0), false},
		{"Y equals height", geometry.NewPoint(0, height), false},
		{"both out of range", geometry.NewPoint(width, height), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, b.InBounds(tc.p))
		})
	}
}

func TestIsValidMark(t *testing.T) {
	b := newBoard(t, 3, 3)

	cases := []struct {
		name     string
		mark     model.Mark
		expected bool
	}{
		{"X", model.MarkX, true},
		{"O", model.MarkO, true},
		{"Empty", model.MarkEmpty, false},
		{"arbitrary positive", model.Mark(2), false},
		{"arbitrary negative", model.Mark(-2), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, b.IsValidMark(tc.mark))
		})
	}
}

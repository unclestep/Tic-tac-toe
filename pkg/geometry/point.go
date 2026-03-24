package geometry

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

//
// -- CARDINAL DIRECTIONS --
//

func GetDirUp() Point {
	return Point{X: 0, Y: -1}
}

func GetDirDown() Point {
	return Point{X: 0, Y: 1}
}

func GetDirRight() Point {
	return Point{X: 1, Y: 0}
}

func GetDirLeft() Point {
	return Point{X: -1, Y: 0}
}

func GetCardinalDirs() []Point {
	return []Point{GetDirUp(), GetDirRight(), GetDirDown(), GetDirLeft()}
}

//
// -- DIAGONAL DIRECTIONS --
//

func GetDirUpRight() Point {
	return Point{X: 1, Y: -1}
}

func GetDirUpLeft() Point {
	return Point{X: -1, Y: -1}
}

func GetDirDownRight() Point {
	return Point{X: 1, Y: 1}
}

func GetDirDownLeft() Point {
	return Point{X: -1, Y: 1}
}

func GetDiagonalDirs() []Point {
	return []Point{GetDirUpRight(), GetDirDownRight(), GetDirDownLeft(), GetDirUpLeft()}
}

//
// -- ALL DIRECTIONS --
//

func GetAllDirs() []Point {
	return append(GetCardinalDirs(), GetDiagonalDirs()...)
}

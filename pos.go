package main

import "fmt"

// Pos is a convenience struct for coordinate pairs
type Pos struct {
	X, Y int
}

// NewPos is the constructor. It range-checks the coordinates to guarantee
// they fall within the board
func NewPos(x, y int) (*Pos, error) {
	if x >= 0 && x <= 7 && y >= 0 && y <= 7 {
		return &Pos{x, y}, nil
	}
	return nil, fmt.Errorf("Range error")
}

// NewPosFromSquareID constructs from the standard notation square number
func NewPosFromSquareID(squareNumber int) Pos {
	pair := [][]int{
		{1, 0}, {3, 0}, {5, 0}, {7, 0},
		{0, 1}, {2, 1}, {4, 1}, {6, 1},
		{1, 2}, {3, 2}, {5, 2}, {7, 2},
		{0, 3}, {2, 3}, {4, 3}, {6, 3},
		{1, 4}, {3, 4}, {5, 4}, {7, 4},
		{0, 5}, {2, 5}, {4, 5}, {6, 5},
		{1, 6}, {3, 6}, {5, 6}, {7, 6},
		{0, 7}, {2, 7}, {4, 7}, {6, 7},
	}[squareNumber-1]

	return Pos{pair[0], pair[1]}
}

// AsString returns a string representation of a coordinate for
// display and error reporting purposes
func (p Pos) AsString() string {
	return fmt.Sprintf("{%d, %d}", p.X, p.Y)
}

// IsJump compares a pair and returns true if the transition has skipped
// one square, and an indicator of the direction
func IsJump(p1, p2 Pos) (bool, int) {
	dY := p2.Y - p1.Y
	if dY == 2 {
		return true, 1
	}

	if dY == -2 {
		return true, -1
	}

	return false, 0
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// ValidDiagonal check that we moved diagonally at least one and max 2
// squares
func ValidDiagonal(p1, p2 Pos) (int, int, error) {
	dX := p2.X - p1.X
	dY := p2.Y - p1.Y

	adX := abs(dX)
	adY := abs(dY)

	// Check that we moved either 1 or 2 squares in X
	if adX != 1 && adX != 2 {
		return 0, 0, fmt.Errorf("Incorrect X for valid diagonal move")
	}

	// Check that we moved either 1 or 2 squares in Y
	if adY != 1 && adY != 2 {
		return 0, 0, fmt.Errorf("Incorrect Y for valid diagonal move")
	}

	return dX, dY, nil
}

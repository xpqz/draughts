package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Move is holding the list of squares a piece will visit when making a move
type Move struct {
	Squares []Pos // List of squares visited in this move
	Player  int
}

// NewMove is the constructor
func NewMove(player int, pos ...Pos) *Move {
	m := &Move{Player: player}
	for _, p := range pos {
		m.addSquare(p)
	}

	return m
}

// Captures returns true if the move captures any of the opposition's pieces.
// It is assumed to be valid
func (m Move) Captures() bool {
	if len(m.Squares) > 2 {
		return true
	}
	if len(m.Squares) == 2 {
		dY := m.Squares[1].Y - m.Squares[0].Y
		if dY > 1 || dY < -1 {
			return true
		}
	}

	return false
}

// ValidJumpSequence tests if all transitions in a move skips a square
// diagonally and in the same Y-direction
func (m Move) ValidJumpSequence() bool {
	direction := 0
	for i := 0; i < len(m.Squares)-1; i++ {
		first := m.Squares[i]
		second := m.Squares[i+1]
		jumped, dir := IsJump(first, second)
		if jumped == false {
			return false
		}

		if direction == 0 {
			direction = dir
		} else if direction != dir {
			return false
		}
	}
	return true
}

// AsString returns a string representation of a move chain for
// display and error reporting purposes. It distinguishes between
// simple moves (arrow separator) and capture moves (cross separator)
func (m Move) AsString() string {
	separator := ` ➤ `
	if m.Captures() {
		separator = ` ☓ `
	}
	data := []string{}
	for _, square := range m.Squares {
		data = append(data, square.AsString())
	}

	return fmt.Sprintf("%d: %s", m.Player, strings.Join(data, separator))
}

// JumpedSquares returns a list of all squares that were jumped over
// This works for either up or down the board
func (m Move) JumpedSquares() []Pos {
	list := []Pos{}

	for index := 0; index < len(m.Squares)-1; index++ {
		start := m.Squares[index]
		end := m.Squares[index+1]

		skipped := Pos{
			start.X + (end.X-start.X)/2,
			start.Y + (end.Y-start.Y)/2,
		}

		if skipped != start {
			list = append(list, skipped)
		}
	}

	return list
}

func (m *Move) addSquare(pos Pos) {
	m.Squares = append(m.Squares, pos)
}

// ParseMove reads a list of squares separated by space. A square is
// represented by its numbering in the standard notation.
func ParseMove(moveStr string, player int) (*Move, error) {
	squares := strings.Split(moveStr, " ")
	move := &Move{Player: player}
	for _, squareStr := range squares {
		squareNumber, err := strconv.Atoi(squareStr)
		if err != nil || squareNumber < 1 || squareNumber > 32 {
			return nil, fmt.Errorf("Bad square '%s'", squareStr)
		}
		move.addSquare(NewPosFromSquareID(squareNumber))
	}

	return move, nil
}

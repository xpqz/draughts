package main

import (
	"bufio"
	"fmt"
	"os"
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

// ReadMove creates a move from stdin input
func ReadMove(player int) (*Move, error) {
	reader := bufio.NewReader(os.Stdin)
	moveStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	moveStr = strings.Trim(moveStr, " \n")
	return ParseMove(moveStr, player)
}

// ValidJumpSequence tests if all transitions in a move skips a square
// diagonally and in the same Y-direction
func (m Move) ValidJumpSequence() bool {
	direction := 0
	for i := 0; i < m.Length()-1; i++ {
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
	if m.IsJump() {
		separator = ` ☓ `
	}
	data := []string{}
	for _, square := range m.Squares {
		data = append(data, square.AsString())
	}

	return fmt.Sprintf("%d: %s", m.Player, strings.Join(data, separator))
}

// AsPDNString returns a string representation of a move chain for
// display and error reporting purposes, using the PDN square. It
// distinguishes between simple moves (arrow separator) and capture moves
// (cross separator)
func (m Move) AsPDNString() string {
	separator := ` ➤ `
	if m.IsJump() {
		separator = ` ☓ `
	}
	data := []string{}
	for _, square := range m.Squares {
		data = append(data, fmt.Sprintf("%d", square.AsPDNSquare()))
	}

	return fmt.Sprintf("%s", strings.Join(data, separator))
}

// JumpedSquares returns a list of all squares that were jumped over
// This works for either up or down the board
func (m Move) JumpedSquares() []Pos {
	list := []Pos{}

	for index := 0; index < m.Length()-1; index++ {
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

// IsJump returns true if the move is a capture sequence
func (m Move) IsJump() bool {
	jumps := false
	if len(m.Squares) > 1 {
		jumps, _ = IsJump(m.Squares[0], m.Squares[1])
	}

	return jumps
}

// Length returns the number of squares in the move
func (m Move) Length() int {
	return len(m.Squares)
}

// Equals compares two moves for equality
func (m Move) Equals(move *Move) bool {
	// Only check positions
	if m.Length() != move.Length() {
		return false
	}

	for i, pos := range m.Squares {
		if pos != move.Squares[i] {
			return false
		}
	}
	return true
}

// CapturesAvailable returns true if any move in a list is a capture sequence
func CapturesAvailable(movesList []*Move) bool {
	for _, move := range movesList {
		if move.IsJump() {
			return true
		}
	}

	return false
}

// ContainsMove checks a move against a list of moves to determine if it is
// present
func ContainsMove(movesList []*Move, move *Move) bool {
	for _, candidate := range movesList {
		if candidate.Equals(move) {
			return true
		}
	}

	return false
}

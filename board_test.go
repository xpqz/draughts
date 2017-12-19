package main

import (
	"fmt"
	"testing"
)

func TestApply(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0, 0},
		{2, 0, 2, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	expected := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{2, 0, 2, 0, 1, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	move := NewMove(player, Pos{2, 3}, Pos{4, 5})

	newBoard := board.Apply(move)

	if expected.state != newBoard.state {
		t.Errorf("Apply didn't produce expected result after %s",
			move.AsString())
	}
}

func TestValidate(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0, 0},
		{2, 0, 2, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	// non-diagonal move
	move := NewMove(player, Pos{5, 2}, Pos{5, 3})
	err := board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 1 as legal")
	}

	// jump to occupied square
	move = NewMove(Opposition(player), Pos{3, 3}, Pos{2, 3})
	err = board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 2 as legal")
	}

	board = &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 2, 0, 0, 0, 0, 0},
		{2, 0, 2, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	// non-jump following jump
	move = NewMove(player, Pos{3, 3}, Pos{1, 5}, Pos{0, 6})
	err = board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 3 as legal")
	}

	board = &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{2, 2, 2, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	// jump following non-jump
	move = NewMove(player, Pos{3, 3}, Pos{2, 4}, Pos{0, 6})
	err = board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 4 as legal")
	}

	// non-king backwards diagonal
	move = NewMove(player, Pos{3, 3}, Pos{4, 2})
	err = board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 5 as legal")
	}

	// non-jump move not to adjacent square
	move = NewMove(player, Pos{5, 2}, Pos{3, 4})
	err = board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 6 as legal")
	}
}

func TestValidateKing(t *testing.T) {
	player := 2

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 0, 0, 1},
		{0, 0, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, -2, 0, 0, 0, 0},
		{2, 0, 2, 0, 1, 0, 2, 0},
		{0, 2, 0, 2, 0, 0, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	move := NewMove(player, Pos{3, 4}, Pos{5, 6})
	err := board.Validate(move)
	if err != nil {
		t.Errorf("Validate reported legal move 1 as illegal, %s", err)
	}

	move = NewMove(player, Pos{3, 4}, Pos{5, 2})
	err = board.Validate(move)
	if err != nil {
		t.Errorf("Validate reported legal move 2 as illegal, %s", err)
	}

	move = NewMove(player, Pos{3, 4}, Pos{6, 1})
	err = board.Validate(move)
	if err == nil {
		t.Error("Validate reported illegal move 2 as legal")
	}
}

func TestAvailableJumps(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 2, 0, 2, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	jumps := board.singleJumps(player, Pos{2, 3})

	if len(jumps) != 2 {
		t.Errorf("Expected 2 available jumps, found %d", len(jumps))
	}
}

func TestAvailableJumpsKing(t *testing.T) {
	player := 2

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 1, 0, 1, 0, 0, 0},
		{0, 2, 0, -2, 0, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 2, 0},
		{0, 0, 0, 2, 0, 0, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	jumps := board.singleJumps(player, Pos{3, 4})

	if len(jumps) != 4 {
		t.Errorf("Expected 4 available jumps, found %d", len(jumps))
	}
}

func TestJumpMoves(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 2, 0, 2, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 0, 0, 2, 0, 2, 0},
	}}

	pos := Pos{2, 3}
	movesList := board.JumpMoves(player, pos)

	if len(movesList) != 2 {
		t.Errorf("Expected 2 available jump moves, found %d", len(movesList))
	}

	for _, move := range movesList {
		if len(move.Squares) != 3 {
			t.Errorf("Expected each jump move to have 3 squares, found %d",
				len(move.Squares))
			fmt.Printf("%s\n", move.AsString())
		}
	}
}

func TestNonCaptureMoves(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0, 0},
		{2, 0, 2, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
	}}

	movesList1 := board.nonCaptureMoves(player, Pos{5, 2})
	if len(movesList1) != 2 {
		t.Errorf("Expected to find 2 non-cap moves, found %d", len(movesList1))
	}
	movesList2 := board.nonCaptureMoves(player, Pos{0, 5})
	if len(movesList2) != 0 {
		t.Errorf("Expected to find 0 non-cap moves, found %d", len(movesList2))
	}
	movesList3 := board.nonCaptureMoves(player, Pos{2, 3})
	if len(movesList3) != 1 {
		t.Errorf("Expected to find 1 non-cap moves, found %d", len(movesList3))
	}
	movesList4 := board.nonCaptureMoves(Opposition(player), Pos{0, 2})
	if len(movesList4) != 1 {
		t.Errorf("Expected to find 1 non-cap moves, found %d", len(movesList4))
	}
}

func TestAllMoves(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 2, 0, 2, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 0, 0, 2, 0, 2, 0},
	}}

	// If jump moves are present only they count, so only 2
	movesList := board.AllMoves(player)

	if len(movesList) != 2 {
		t.Errorf("Expected to find 2 valid moves, found %d", len(movesList))
	}
}

func TestAllMovesGeneratedAreValid(t *testing.T) {
	player := 1

	board := &Board{[8][8]int{
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, -1, 0, 0, 0, 0, 0},
		{0, 2, 0, 2, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 0, 0, 2, 0, 2, 0},
	}}

	for _, move := range board.AllMoves(player) {
		err := board.Validate(move)
		if err != nil {
			t.Errorf("Generated move %s reported invalid", move.AsString())
		}
	}
}

func TestContainsMove(t *testing.T) {
	player := 2

	board := &Board{[8][8]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, -2, 0, 1, 0},
		{0, -1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, -1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, -1, 0, 0, 0},
	}}

	move := NewMove(player, Pos{4, 1}, Pos{2, 3})
	err := board.Validate(move)
	if err != nil {
		t.Errorf("Generated move %s reported invalid", move.AsString())
	}

	allMoves := board.AllMoves(player)
	if !ContainsMove(allMoves, move) {
		t.Errorf("Complete capture move %s reported incomplete", move.AsString())

		fmt.Printf("All moves:\n")
		for _, m := range allMoves {
			fmt.Printf("%s\n", m.AsString())
		}
	}
}

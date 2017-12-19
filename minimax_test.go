package main

import (
	"testing"
)

func TestEval(t *testing.T) {
	player := 1

	board := NewBoard()
	be := eval(board, player)

	if be.PieceDiff != 0 {
		t.Errorf("Expected a piece difference of 0, found %d", be.PieceDiff)
	}
	if be.KingDiff != 0 {
		t.Errorf("Expected a king difference of 0, found %d", be.KingDiff)
	}

	board = &Board{[8][8]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 1, 0},
		{0, -1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, -1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0, 0},
	}}

	be = eval(board, player)

	if be.PieceDiff != 4 {
		t.Errorf("Expected a piece difference of 4, found %d", be.PieceDiff)
	}
	if be.KingDiff != 2 {
		t.Errorf("Expected a king difference of 2, found %d", be.KingDiff)
	}

	board = &Board{[8][8]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, -1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{2, 0, 0, 0, 0, 0, 0, 0},
	}}

	be = eval(board, player)

	if be.MovesCountDiff != 3 {
		t.Errorf("Expected a movescount difference of 3, found %d in %+v",
			be.MovesCountDiff, be)
	}
}

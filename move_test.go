package main

import "testing"

func TestJumpedSquares(t *testing.T) {
	player := 1
	square, _ := NewPos(1, 2)

	move := NewMove(player, *square)
	square2, _ := NewPos(3, 4)

	move.addSquare(*square2)
	square3, _ := NewPos(5, 6)

	move.addSquare(*square3)

	jumped := move.JumpedSquares()

	if len(jumped) != 2 {
		t.Errorf("Expected 2 jumped squares, found %d", len(jumped))
	}
}

func TestParseMove(t *testing.T) {
	player := 1

	type parseTest struct {
		data     string
		expected *Move
	}

	var parseTests = []parseTest{
		{"14 23 30", &Move{
			[]Pos{{2, 3}, {4, 5}, {2, 7}},
			player,
		}},

		{"15 19", &Move{
			[]Pos{{4, 3}, {5, 4}},
			player,
		}},
	}

	for _, tt := range parseTests {
		actual, err := ParseMove(tt.data, player)
		if err != nil {
			t.Errorf("ParseMove(\"%s\", %d) returned error %s",
				tt.data, player, err)
		} else if actual.AsString() != tt.expected.AsString() {
			t.Errorf("ParseMove(\"%s\", %d): expected %s, actual %s",
				tt.data, player, tt.expected.AsString(), actual.AsString())
		}
	}
}

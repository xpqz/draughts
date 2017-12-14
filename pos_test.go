package main

import "testing"

func TestNewPos(t *testing.T) {
	type newPosTest struct {
		x        int
		y        int
		expected *Pos
	}

	var newPosTests = []newPosTest{
		{3, 0, &Pos{3, 0}},
		{2, 3, &Pos{2, 3}},
		{8, 0, nil},
		{7, -1, nil},
	}

	for _, tt := range newPosTests {
		actual, err := NewPos(tt.x, tt.y)
		if err != nil {
			if tt.expected != nil {
				t.Errorf("NewPos(%d, %d) returned error %s", tt.x, tt.y, err)
			}
		} else if actual.AsString() != tt.expected.AsString() {
			t.Errorf("NewPos(%d, %d): expected %s, actual %s",
				tt.x, tt.y, tt.expected.AsString(), actual.AsString())
		}
	}
}

package main

import "fmt"

// Opposition is the other player
func Opposition(player int) int {
	if player == 1 {
		return 2
	}
	return 1
}

// Direction says if "forward" is up or down the board. This allows us to
// use the same logic for player 1 and 2.
func direction(player int) int {
	if player == 1 {
		return 1
	}
	return -1
}

// Board represents the state of the game as an 8x8 array of integers
// where player 1 is the value 1, player 2 the value 2 and an empty
// space is a 0
type Board struct {
	state [8][8]int
}

// NewBoard sets up a fresh board in the initial configuration
func NewBoard() *Board {
	b := &Board{}
	b.state[0] = [8]int{0, 1, 0, 1, 0, 1, 0, 1}
	b.state[1] = [8]int{1, 0, 1, 0, 1, 0, 1, 0}
	b.state[2] = [8]int{0, 1, 0, 1, 0, 1, 0, 1}

	b.state[5] = [8]int{2, 0, 2, 0, 2, 0, 2, 0}
	b.state[6] = [8]int{0, 2, 0, 2, 0, 2, 0, 2}
	b.state[7] = [8]int{2, 0, 2, 0, 2, 0, 2, 0}

	return b
}

// Get returns the state of the board Get position `p`
func (b Board) Get(p Pos) int {
	return b.state[p.Y][p.X]
}

// Set sets the the state of the board at position `p`
func (b *Board) Set(p Pos, state int) {
	b.state[p.Y][p.X] = state
}

// Apply moves a piece according to `move` and returns a new
// state. The move must be legal.
func (b Board) Apply(move *Move) *Board {
	startPos := move.Squares[0]
	startPiece := b.Get(startPos)
	newBoard := &Board{b.state}

	// Blank out the starting square of the move
	newBoard.Set(startPos, 0)

	// Remove any captured pieces we jumped from the board
	for _, square := range move.JumpedSquares() {
		newBoard.Set(square, 0)
	}

	// Land on the final square
	final := move.Squares[move.Length()-1]

	// If not already a king, and we've landed on either baseline
	// we get a coronation.
	if startPiece > 0 && final.Y == 0 || final.Y == 7 {
		newBoard.Set(final, -startPiece) // Coronation
	} else {
		newBoard.Set(final, startPiece)
	}

	return newBoard
}

// Validate traverses a chain of squares to dermine if the move is legal
func (b Board) Validate(move *Move) error {
	// Start square has to be held by `player`
	startPos := move.Squares[0]
	startPiece := b.Get(startPos)
	if abs(startPiece) != move.Player {
		return fmt.Errorf("Start position %s isn't valid", startPos.AsString())
	}

	// If we have more that one transition (more than two squares), then all
	// transitions must be jumps. Note, this doesn't verify the actual state
	// of the board, just the coordinates.
	if move.Length() > 2 && !move.ValidJumpSequence() {
		return fmt.Errorf("Move is not a valid jump sequence")
	}

	for index := 1; index < move.Length(); index++ {
		endPos := move.Squares[index]

		// endPos has to be available
		if b.Get(endPos) != 0 {
			return fmt.Errorf("Position %s isn't available", endPos.AsString())
		}

		// Check that we moved diagonally, either 1 or 2 squares
		dX, dY, err := ValidDiagonal(startPos, endPos)
		if err != nil {
			return err
		}

		// Check we're moving forwards, if not king. King-case is covered
		// by the above diagonal check.
		if dY > 0 && startPiece == 2 || dY < 0 && startPiece == 1 {
			return fmt.Errorf("Non-king move must move forward")
		}

		// If I jumped a square, it has to be occupied by the opponent.
		betweenX := dX / 2
		betweenY := dY / 2

		middle := Pos{
			startPos.X + betweenX,
			startPos.Y + betweenY,
		}

		if betweenX != 0 && betweenY != 0 {
			if abs(b.Get(middle)) != Opposition(move.Player) {
				return fmt.Errorf("Jumped square not held by opponent")
			}
		}

		startPos = endPos
	}

	return nil // we're good
}

// nonCaptureMoves returns a list of non capture moves from `square` --
// 0, 1, 2 (4, 5 for kings) possible squares.
func (b Board) nonCaptureMoves(player int, square Pos) []*Move {
	cols := []int{-1, 1}
	rows := []int{direction(player)}
	if b.Get(square) < 0 {
		// Kings can go both backwards and forwards
		rows = append(rows, -direction(player))
	}

	moves := []*Move{}
	for _, dY := range rows {
		for _, dX := range cols {
			if candidate, err := NewPos(square.X+dX, square.Y+dY); err == nil {
				if b.Get(*candidate) == 0 {
					moves = append(moves, NewMove(player, square, *candidate))
				}
			}
		}
	}

	return moves
}

// singleJumps returns a list of single jumpable positions from `square` --
// 0, 1, 2 (4, 5 for kings) possible squares
func (b Board) singleJumps(player int, square Pos) []Pos {
	cols := []int{-1, 1}
	rows := []int{direction(player)}
	if b.Get(square) < 0 {
		// Kings can go both backwards and forwards
		rows = append(rows, -direction(player))
	}

	jumps := []Pos{}

	for _, dY := range rows {
		for _, dX := range cols {
			if first, err := NewPos(square.X+dX, square.Y+dY); err == nil {
				if abs(b.Get(*first)) == Opposition(player) {
					if second, err := NewPos(first.X+dX, first.Y+dY); err == nil {
						if b.Get(*second) == 0 {
							jumps = append(jumps, *second)
						}
					}
				}
			}
		}
	}

	return jumps
}

// JumpMoves finds all available capture moves starting at `square`, including
// multi-hops and king moves
func (b Board) JumpMoves(player int, square Pos) []*Move {
	movesList := []*Move{}
	b.jumpMoves(&Board{b.state}, square, NewMove(player, square), &movesList)

	return movesList
}

// validDirection discards potential jump squares based on the prevaling
// direction already set in the move. This is to ensure that kings, whilst
// allowed to go both forward and back still stick to a single direction
// within each single move.
func validDirection(jumps []Pos, move *Move) []Pos {
	if move.Length() == 1 {
		// We have no prevaling direction, so any jumps may be considered
		return jumps
	}

	dY := move.Squares[1].Y - move.Squares[0].Y

	sameDirection := []Pos{}
	for _, jmp := range jumps {
		if jmp.Y-move.Squares[1].Y == dY {
			sameDirection = append(sameDirection, jmp)
		}
	}

	return sameDirection
}

func (b Board) jumpMoves(state *Board, square Pos, move *Move, moves *[]*Move) {
	player := move.Player
	jumps := validDirection(b.singleJumps(player, square), move)

	if len(jumps) == 0 && move.Length() > 1 {
		*moves = append(*moves, move)
		return
	}

	for _, jmp := range jumps {
		newState := state.Apply(NewMove(player, square, jmp))
		move.addSquare(jmp)                     // Record this in current move..
		b.jumpMoves(newState, jmp, move, moves) // Walk the tree depth-first

		// Now make a fresh move for when we follow the next branch
		move = NewMove(player, square)
	}
}

// AllMoves returns a list of all valid moves for `player`
func (b *Board) AllMoves(player int) []*Move {
	movesList := []*Move{}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			pos := Pos{x, y}
			if abs(b.Get(pos)) == player {
				movesList = append(movesList, b.nonCaptureMoves(player, pos)...)
				movesList = append(movesList, b.JumpMoves(player, pos)...)
			}
		}
	}

	return movesList
}

package main

// This is the main interactive game loop for two-people game play.

import (
	"fmt"
)

func main() {
	board := NewBoard()
	player := 1
	for {
		// Check that player has any move options available
		allMoves := board.AllMoves(player)
		if len(allMoves) == 0 { // Opponent wins
			break
		}

		Display(board)

		// Read-Validate-Apply move for `player`
		for {
			if player == 1 {
				fmt.Printf("\033[31mRed's move: \033[39m")
			} else {
				fmt.Printf("\033[32mGreen's move: \033[39m")
			}

			move, err := ReadMove(player)
			if err != nil {
				fmt.Printf("\n--> incorrectly entered move; try again (%s)\n",
					err)
				continue
			}

			// Check if the move is legal under the rules of the game
			err = board.Validate(move)
			if err != nil {
				fmt.Printf("\n--> illegal move; try again (%s)\n", err)
				continue
			}

			// If a capture move is available, it must be taken. If several
			// capture moves are available, it's sufficient to pick one; it does
			// not have to be the longest.
			if CapturesAvailable(allMoves) && !move.IsJump() {
				fmt.Printf("\n--> capture moves available and not taken; try again (%s)\n",
					err)
				continue
			}

			// A capture move must be taken in its entirety. In this case it's
			// sufficient to check that the move is present in `allMoves` as this
			// generates the longest possible jumps.
			if !ContainsMove(allMoves, move) {
				fmt.Printf("\n--> capture moves must be taken entirely; try again\n")
				continue
			}

			// Finally we can make this move, and swap turns.
			board = board.Apply(move)
			player = Opposition(player)
			break
		}
	}

	if player == 2 { // Note: opponent is the winner
		fmt.Printf("Winner: \033[31mRed\033[39m\n")
	} else {
		fmt.Printf("Winner: \033[32mGreen\033[39m\n")
	}
}

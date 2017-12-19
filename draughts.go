package main

// This is the main interactive game loop for two-people game play.

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CheckUserMove validates a move entered by a human player
func CheckUserMove(board *Board, allMoves []*Move, move *Move) error {
	// Check if the move is legal under the rules of the game
	err := board.Validate(move)
	if err != nil {
		return err
	}

	// If a capture move is available, it must be taken. If several
	// capture moves are available, it's sufficient to pick one; it does
	// not have to be the longest.
	if CapturesAvailable(allMoves) && !move.IsJump() {
		return fmt.Errorf("capture moves available and not taken")
	}

	// A capture move must be taken in its entirety. In this case it's
	// sufficient to check that the move is present in `allMoves` as this
	// generates the longest possible jumps.
	if !ContainsMove(allMoves, move) {
		return fmt.Errorf("capture moves must be taken entirely")
	}

	return nil
}

// TwoPersonGame pits Human against Human
func TwoPersonGame() {
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

			err = CheckUserMove(board, allMoves, move)
			if err != nil {
				fmt.Printf("\n--> incorrect move; try again (%s)\n", err)
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

// OnePersonGame pits Human vs Machine. The minimaxer plays a sort of ok
// opening and middle game, but pretty poor endgame
func OnePersonGame() {
	board := NewBoard()
	player := 1
	for {
		// Check that player has any move options available
		allMoves := board.AllMoves(player)
		if len(allMoves) == 0 { // Opponent wins
			break
		}

		score := minimax(board, player, 0, 6)

		Display(board)
		var (
			move *Move
			err  error
		)
		// Read-Validate-Apply move for `player`
		for {
			if player == 1 {
				fmt.Printf("\033[31mRed's move: %s\n\033[39m", score.move.AsPDNString())
				move = score.move
			} else {
				fmt.Printf("\033[32mGreen's move: \033[39m")

				move, err = ReadMove(player)
				if err != nil {
					fmt.Printf("\n--> incorrectly entered move; try again (%s)\n",
						err)
					continue
				}

				err = CheckUserMove(board, allMoves, move)
				if err != nil {
					fmt.Printf("\n--> incorrect move; try again (%s)\n", err)
					continue
				}
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

func main() {
	fmt.Printf("Press c for computer opponent or anything else for human: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.Trim(choice, " \n")

	if choice == "c" {
		OnePersonGame()
	} else {
		TwoPersonGame()
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	board := NewBoard()
	player := 1
	for !board.Winner(player) {
		displayStandardNotation(board)
		for {
			if player == 1 {
				fmt.Printf("\033[31mRed's move: \033[39m")
			} else {
				fmt.Printf("\033[32mGreen's move: \033[39m")
			}
			reader := bufio.NewReader(os.Stdin)
			moveStr, err := reader.ReadString('\n')
			moveStr = strings.Trim(moveStr, " \n")
			move, err := ParseMove(moveStr, player)
			if err != nil {
				fmt.Printf("\n--> incorrectly entered move; try again (%s)\n",
					err)
				continue
			}

			newBoard, err := board.MakeMove(move)
			if err != nil {
				fmt.Printf("\n--> illegal move; try again (%s)\n", err)
				continue
			}

			player = Opposition(player)
			board = newBoard
			break
		}
	}

	fmt.Printf("Winner: player %d", player)
}

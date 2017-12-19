package main

import "fmt"

func coloured(squareState, key int) {
	padding := "  "
	if key >= 10 {
		padding = " "
	}
	king := " "
	if squareState < 0 {
		king = "*"
	}
	if abs(squareState) == 1 {
		fmt.Printf("\033[31m%s%d%s \033[39m", padding, key, king) // RED
	} else if abs(squareState) == 2 {
		fmt.Printf("\033[32m%s%d%s \033[39m", padding, key, king) // GREEN
	} else {
		fmt.Printf("%s%d  ", padding, key)
	}

	fmt.Printf("|")
}

// Display prints out square numbers, coloured by player
func Display(b *Board) {
	rCount, gCount := 0, 0
	fmt.Print("  +-----+-----+-----+-----+-----+-----+-----+-----+\n")
	key := 1
	for y := 0; y < 8; y++ {
		fmt.Print("  |")
		for x := 0; x < 8; x++ {
			if y%2 == 0 && x%2 != 0 || // even y, odd x
				y%2 != 0 && x%2 == 0 { // odd y, even x
				squareState := b.Get(Pos{x, y})
				coloured(squareState, key)
				if abs(squareState) == 1 {
					rCount++
				} else if abs(squareState) == 2 {
					gCount++
				}
				key++
			} else {
				fmt.Print("     |")
			}
		}
		fmt.Print("\n  +-----+-----+-----+-----+-----+-----+-----+-----+\n")
	}
	fmt.Printf("  \033[31m[Red: %d]\033[39m \033[32m[Green: %d]\033[39m\n",
		rCount, gCount)
	// fmt.Printf("  \033[31m[Score: %d]\033[39m \033[32m[Score: %d]\033[39m\n",
	// 	eval(b, 1), eval(b, 2))

}

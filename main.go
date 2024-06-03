package main

import (
	"fmt"
	"os"

	"github.com/pavelanni/hexapawn-go/hexapawn"
)

// Define constants for board dimensions
const (
	boardRows = 3
	boardCols = 3
)

func main() {
	board := hexapawn.NewBoard(boardCols, boardRows)
	if len(os.Args) > 1 {
		if os.Args[1] == "generate" {
			board.Generate()
			os.Exit(0)
		} else {
			fmt.Println("Unknown command")
			os.Exit(1)
		}
	}
	board.PlayGame()
}

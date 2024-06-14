package main

import (
	"log"
	"os"

	"github.com/pavelanni/hexapawn-go/hexapawn"
	flag "github.com/spf13/pflag"
)

// Define constants for board dimensions
var boardRows, boardCols, numPlayers int
var generate bool

func main() {
	flag.IntVarP(&boardRows, "rows", "r", 3, "Number of rows in the board, must be at least 3, at most 9. Default is 3.")
	flag.IntVarP(&numPlayers, "players", "p", 2, "Number of human players: 0, 1, or 2. Default is 2.")
	flag.BoolVarP(&generate, "generate", "g", false, "Generate possible moves and exit.")
	flag.Parse()
	game, err := hexapawn.NewGame(boardRows, numPlayers)
	if err != nil {
		log.Fatal(err)
	}
	if generate {
		game.Board.Generate()
		os.Exit(0)
	}
	game.Play()
}

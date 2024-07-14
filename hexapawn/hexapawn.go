package hexapawn

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyz"
)

// Initialize a new board
func NewBoard(boardRows int) *Board {
	b := &Board{
		Rows: boardRows,
		Cols: boardRows,
		Grid: make([][]string, boardRows),
	}
	for row := 0; row < b.Rows; row++ {
		b.Grid[row] = make([]string, b.Cols)
		for col := 0; col < b.Cols; col++ {
			b.Grid[row][col] = "." // fill empty cells with dots
		}
	}
	// Place white pawns (W) at the bottom row
	for col := 0; col < b.Cols; col++ {
		b.Grid[0][col] = "W"
	}
	// Place black pawns (B) at the top row
	for col := 0; col < b.Cols; col++ {
		b.Grid[b.Rows-1][col] = "B"
	}
	return b
}

// Initialize a new game
func NewGame(boardRows, numPlayers int) (*Game, error) {
	// check if the board dimensions are valid
	if boardRows < 3 || boardRows > 9 {
		return nil, fmt.Errorf("Invalid board dimensions. Rows and Cols must be equal and at least 3, at most 9.")
	}
	if numPlayers < 0 || numPlayers > 2 {
		return nil, fmt.Errorf("Invalid number of players. Must be 0, 1, or 2.")
	}
	g := &Game{
		NumPlayers: numPlayers,
		Board:      NewBoard(boardRows),
	}
	return g, nil
}

// Check if a move is valid
func (b *Board) IsValidMove(ms string, player string) bool {
	m, err := b.MoveFromString(ms)
	if err != nil {
		return false
	}
	// Check if move's from position is in bounds
	if m.FromRow < 0 || m.FromRow >= b.Rows || m.FromCol < 0 || m.FromCol >= b.Cols {
		return false
	}
	// Check if move's to position is in bounds
	if m.ToRow < 0 || m.ToRow >= b.Rows || m.ToCol < 0 || m.ToCol >= b.Cols {
		return false
	}
	// Check if we move our piece
	if b.Grid[m.FromRow][m.FromCol] != player {
		return false
	}
	// Check if we move for one row only
	if player == "W" && m.ToRow != m.FromRow+1 {
		return false
	}
	if player == "B" && m.ToRow != m.FromRow-1 {
		return false
	}
	// We can't move non-vertically on an empty space
	if m.ToCol != m.FromCol && b.Grid[m.ToRow][m.ToCol] == "." {
		return false
	}
	// Check if we move to an empty position or it's a capture move
	if b.Grid[m.ToRow][m.ToCol] != "." {
		// check if the move is diagonal
		if m.ToCol != m.FromCol+1 && m.ToCol != m.FromCol-1 {
			return false
		}
		if b.Grid[m.ToRow][m.ToCol] == player { // we can't capture our own piece
			return false
		}
	}
	return true
}

// Apply a move to the board
func (b *Board) ApplyMove(ms string) {
	m, err := b.MoveFromString(ms)
	if err != nil {
		log.Fatal(err)
	}
	b.Grid[m.ToRow][m.ToCol] = b.Grid[m.FromRow][m.FromCol]
	b.Grid[m.FromRow][m.FromCol] = "."
}

// Check for a win condition
func (b *Board) CheckWin() string {
	for col := 0; col < b.Cols; col++ {
		if b.Grid[0][col] == "B" {
			return "B"
		}
		if b.Grid[b.Rows-1][col] == "W" {
			return "W"
		}
	}
	return ""
}

// Play the game
func (g *Game) Play() {
	currentPlayer := "W"
	winner := ""
	var move string
	var stepNumber int

	for {
		fmt.Printf("Step %d\n", stepNumber+1)
		moves := g.Steps[stepNumber].Moves[g.Board.String()]
		if len(moves) == 0 {
			if currentPlayer == "W" {
				winner = "B"
			} else {
				winner = "W"
			}
			fmt.Printf("Player %s has no moves; player %s wins!\n", currentPlayer, winner)
			break
		}
		switch g.NumPlayers {
		case 2:
			fmt.Printf("Player %s, enter your move: ", currentPlayer)
			fmt.Scan(&move)
		case 1:
			if currentPlayer == "W" {
				fmt.Print("Player W, enter your move: ")
				fmt.Scan(&move)
			} else {
				move = moves[rand.Intn(len(moves))]
			}
		case 0:
			move = moves[rand.Intn(len(moves))]
		default:
			fmt.Println("Invalid number of players")
			continue
		}

		if g.Board.IsValidMove(move, currentPlayer) {
			fmt.Printf("Player %s moves %s\n", currentPlayer, move)
			g.MovesPlayed = append(g.MovesPlayed, BoardMove{
				BoardStr: g.Board.String(),
				MoveStr:  move})
			g.Board.ApplyMove(move)
			g.Board.Print()
			winner = g.Board.CheckWin()
			if winner != "" {
				fmt.Printf("Player %s wins!\n", winner)
				break
			}
			if currentPlayer == "W" {
				currentPlayer = "B"
			} else {
				currentPlayer = "W"
			}
		} else {
			fmt.Println("Invalid move, try again.")
		}
		stepNumber++
		if stepNumber >= len(g.Steps) {
			break
		}
	}
	g.Winner = winner
}

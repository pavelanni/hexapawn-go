package hexapawn

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyz"
)

// Define the Board struct
type Board struct {
	Rows int
	Cols int
	grid [][]string
}

type Piece struct {
	Player string
	Row    int
	Col    int
}

// Move defines a move
type Move struct {
	FromRow int
	FromCol int
	ToRow   int
	ToCol   int
}

// Initialize a new board with pawns
func NewBoard(boardCols, boardRows int) *Board {
	// check if the board dimensions are valid
	if boardRows < 3 || boardCols < 3 || boardRows > 9 || boardCols > 9 || boardRows != boardCols {
		log.Fatalf("Invalid board dimensions. Rows and Cols must be equal and at least 3, at most 9.")
	}
	b := &Board{
		Rows: boardRows,
		Cols: boardCols,
		grid: make([][]string, boardRows)}
	for row := 0; row < b.Rows; row++ {
		b.grid[row] = make([]string, b.Cols)
		for col := 0; col < b.Cols; col++ {
			b.grid[row][col] = "." // fill empty cells with dots
		}
	}
	// Place white pawns (W) at the bottom row
	for col := 0; col < b.Cols; col++ {
		b.grid[0][col] = "W"
	}
	// Place black pawns (B) at the top row
	for col := 0; col < b.Cols; col++ {
		b.grid[b.Rows-1][col] = "B"
	}
	return b
}

// Print the board state
func (b *Board) Print() {
	for i := b.Rows - 1; i >= 0; i-- {
		fmt.Printf("%2d  ", i+1)
		for _, cell := range b.grid[i] {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
	fmt.Printf("   ")
	for i := 0; i < b.Cols; i++ {
		fmt.Printf("%2c", letters[i])
	}
	fmt.Println()
}

// NewMove creates a new Move from a string of format "a1-a2", similar to chess
func (b *Board) NewMove(moveStr string) (Move, error) {
	if len(moveStr) != 5 {
		return Move{}, fmt.Errorf("Invalid length of move string: %s", moveStr)
	}
	fromCol := strings.Index(letters, moveStr[0:1])
	if fromCol == -1 {
		return Move{}, fmt.Errorf("Invalid from column in move string: %s", moveStr)
	}
	if fromCol < 0 || fromCol >= b.Cols {
		return Move{}, fmt.Errorf("From column out of bounds in move string: %s", moveStr)
	}
	fromRow, err := strconv.Atoi(moveStr[1:2])
	if err != nil {
		return Move{}, fmt.Errorf("Invalid from row in move string: %s", moveStr)
	}
	fromRow-- // Convert from 1-indexed to 0-indexed
	if fromRow < 0 || fromRow >= b.Rows {
		return Move{}, fmt.Errorf("From row out of bounds in move string: %s", moveStr)
	}
	toCol := strings.Index(letters, moveStr[3:4])
	if toCol == -1 {
		return Move{}, fmt.Errorf("Invalid to column in move string: %s", moveStr)
	}
	if toCol < 0 || toCol >= b.Cols {
		return Move{}, fmt.Errorf("To column out of bounds in move string: %s", moveStr)
	}
	toRow, err := strconv.Atoi(moveStr[4:5])
	if err != nil {
		return Move{}, fmt.Errorf("Invalid to row in move string: %s", moveStr)
	}
	toRow-- // Convert from 1-indexed to 0-indexed
	if toRow < 0 || toRow >= b.Rows {
		return Move{}, fmt.Errorf("To row out of bounds in move string: %s", moveStr)
	}
	return Move{fromRow, fromCol, toRow, toCol}, nil
}

// Check if a move is valid
func (b *Board) IsValidMove(m Move, player string) bool {
	// Check if move's from position is in bounds
	if m.FromRow < 0 || m.FromRow >= b.Rows || m.FromCol < 0 || m.FromCol >= b.Cols {
		return false
	}
	// Check if move's to position is in bounds
	if m.ToRow < 0 || m.ToRow >= b.Rows || m.ToCol < 0 || m.ToCol >= b.Cols {
		return false
	}
	// Check if we move our piece
	if b.grid[m.FromRow][m.FromCol] != player {
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
	if m.ToCol != m.FromCol && b.grid[m.ToRow][m.ToCol] == "." {
		return false
	}
	// Check if we move to an empty position or it's a capture move
	if b.grid[m.ToRow][m.ToCol] != "." {
		// check if the move is diagonal
		if m.ToCol != m.FromCol+1 && m.ToCol != m.FromCol-1 {
			return false
		}
	}
	return true
}

// Apply a move to the board
func (b *Board) ApplyMove(m Move) {
	b.grid[m.ToRow][m.ToCol] = b.grid[m.FromRow][m.FromCol]
	b.grid[m.FromRow][m.FromCol] = "."
}

// Check for a win condition
func (b *Board) CheckWin() string {
	for col := 0; col < b.Cols; col++ {
		if b.grid[0][col] == "B" {
			return "B"
		}
		if b.grid[2][col] == "W" {
			return "W"
		}
	}
	return ""
}

// Play the game
func (b *Board) PlayGame() {
	currentPlayer := "W"
	var moveStr string

	for {
		b.Print()
		moves := b.GetMoves(currentPlayer)
		if len(moves) == 0 {
			fmt.Printf("Player %s has no moves, it loses!\n", currentPlayer)
			break
		}
		fmt.Printf("Player %s, enter your move: ", currentPlayer)
		fmt.Scan(&moveStr)

		move, err := b.NewMove(moveStr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if b.IsValidMove(move, currentPlayer) {
			b.ApplyMove(move)
			winner := b.CheckWin()
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
	}
}

// String is a stringer for Move
func (m Move) String() string {
	return fmt.Sprintf("%s%d-%s%d", letters[m.FromCol:m.FromCol+1], m.FromRow+1, letters[m.ToCol:m.ToCol+1], m.ToRow+1)
}

// String is a stringer for Board
func (b Board) String() string {
	return fmt.Sprintf("%s%s%s", strings.Join(b.grid[0], ""),
		strings.Join(b.grid[1], ""), strings.Join(b.grid[2], ""))
}

// GetMoves returns possible moves for the given board and player
func (b *Board) GetMoves(player string) []Move {
	pieces := make([]Piece, 0)
	moves := make([]Move, 0)

	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			if b.grid[row][col] == player {
				pieces = append(pieces, Piece{Player: player, Row: row, Col: col})
			}
		}
	}

	for _, piece := range pieces {
		for row := 0; row < b.Rows; row++ {
			for col := 0; col < b.Cols; col++ {
				if b.IsValidMove(Move{FromRow: piece.Row, FromCol: piece.Col, ToRow: row, ToCol: col}, player) {
					moves = append(moves, Move{FromRow: piece.Row, FromCol: piece.Col, ToRow: row, ToCol: col})
				}
			}
		}
	}
	fmt.Printf("Possible moves for player %s: ", player)
	for _, move := range moves {
		fmt.Printf("%s ", move)
	}
	fmt.Println()
	return moves
}

// Generate generates possible positions and moves
func (b *Board) Generate() {
	fmt.Println(b)
	b.GetMoves("B")
	b.GetMoves("W")

}

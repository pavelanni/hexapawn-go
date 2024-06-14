package hexapawn

import (
	"fmt"
	"math"
	"math/rand"
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

type Game struct {
	NumPlayers int
	Board      *Board
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

// Position defines a position with a board, a player, and a list of possible moves
type Position struct {
	Board  *Board
	Player string
	Moves  []Move
}

type Step struct {
	Player    string              // current player
	Positions map[string]Position // each position is a string like "BBB...WWW"; for each position, there is a list of possible moves
}

// Initialize a new board
func NewBoard(boardRows int) *Board {
	b := &Board{
		Rows: boardRows,
		Cols: boardRows,
		grid: make([][]string, boardRows),
	}
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
		if b.grid[b.Rows-1][col] == "W" {
			return "W"
		}
	}
	return ""
}

// Play the game
func (g *Game) Play() {
	currentPlayer := "W"
	winner := ""
	var moveStr string
	var move Move
	var moveNumber int
	var err error

	for {
		moveNumber++
		fmt.Printf("Move %d\n", moveNumber)
		moves := g.Board.GetMoves(currentPlayer)
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
			fmt.Scan(&moveStr)
			move, err = g.Board.NewMove(moveStr)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case 1:
			if currentPlayer == "W" {
				fmt.Print("Player W, enter your move: ")
				fmt.Scan(&moveStr)
				move, err = g.Board.NewMove(moveStr)
				if err != nil {
					fmt.Println(err)
					continue
				}
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
	}
}

// String is a stringer for Move
func (m Move) String() string {
	return fmt.Sprintf("%s%d-%s%d", letters[m.FromCol:m.FromCol+1], m.FromRow+1, letters[m.ToCol:m.ToCol+1], m.ToRow+1)
}

// String is a stringer for Board
func (b Board) String() string {
	var output string
	for i := 0; i < b.Rows; i++ {
		output += strings.Join(b.grid[i], "")
	}
	return output
}

func BoardFromString(boardStr string) *Board {
	// Board string must be a square of integer
	rows := int(math.Sqrt(float64(len(boardStr))))
	if rows*rows != len(boardStr) {
		panic("Board string's length must be a square of integer: 9, 16, 25, etc.")
	}
	b := &Board{
		grid: make([][]string, rows),
		Rows: rows,
		Cols: rows,
	}
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			b.grid[i][j] = boardStr[i*b.Cols+j : i*b.Cols+j+1]
		}
	}
	return b
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
	ps := NewPositionStore()
	GeneratePositions(ps, b, "W")
}

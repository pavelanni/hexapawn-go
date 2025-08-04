package hexapawn

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	minBoardRows  = 3
	maxBoardRows  = 9
	minNumPlayers = 0
	maxNumPlayers = 2
	EmptyCell     = '.'
	WhitePlayer   = 'W'
	BlackPlayer   = 'B'
)

// Initialize a new board
func NewBoard(boardRows int) *Board {
	b := &Board{
		Rows:  boardRows,
		Cols:  boardRows,
		Grid:  make([][]string, boardRows),
		State: GameInProgress,
	}
	for row := 0; row < b.Rows; row++ {
		b.Grid[row] = make([]string, b.Cols)
		for col := 0; col < b.Cols; col++ {
			b.Grid[row][col] = string(EmptyCell)
		}
	}
	// Place white pawns at the bottom row
	for col := 0; col < b.Cols; col++ {
		b.Grid[0][col] = string(WhitePlayer)
	}
	// Place black pawns at the top row
	for col := 0; col < b.Cols; col++ {
		b.Grid[b.Rows-1][col] = string(BlackPlayer)
	}
	return b
}

// Initialize a new game
func NewGame(boardRows, numPlayers int) (*Game, error) {
	if boardRows < minBoardRows || boardRows > maxBoardRows {
		return nil, NewGameError(ErrInvalidBoard, "board dimensions must be between 3 and 9")
	}
	if numPlayers < minNumPlayers || numPlayers > maxNumPlayers {
		return nil, NewGameError(ErrInvalidPlayer, "number of players must be 0, 1, or 2")
	}

	g := &Game{
		NumPlayers:    numPlayers,
		Board:         NewBoard(boardRows),
		CurrentPlayer: string(WhitePlayer),
	}
	return g, nil
}

// Check if a move is valid
func (b *Board) IsValidMove(ms string, player string) error {
	if b.State != GameInProgress {
		return NewGameError(ErrGameOver, "game is already over")
	}

	m, err := b.MoveFromString(ms)
	if err != nil {
		return NewGameError(ErrInvalidMove, err.Error())
	}

	// Check if move's positions are in bounds
	if !b.isInBounds(m.FromRow, m.FromCol) || !b.isInBounds(m.ToRow, m.ToCol) {
		return NewGameError(ErrInvalidMove, "move is out of bounds")
	}

	// Check if we're moving our own piece
	if b.Grid[m.FromRow][m.FromCol] != player {
		return NewGameError(ErrInvalidMove, "cannot move opponent's piece")
	}

	// Check movement direction based on player
	if player == string(WhitePlayer) && m.ToRow != m.FromRow+1 {
		return NewGameError(ErrInvalidMove, "white can only move forward")
	}
	if player == string(BlackPlayer) && m.ToRow != m.FromRow-1 {
		return NewGameError(ErrInvalidMove, "black can only move forward")
	}

	// Check diagonal capture and forward movement rules
	targetCell := b.Grid[m.ToRow][m.ToCol]
	if m.ToCol != m.FromCol {
		// Diagonal move must be a capture
		if targetCell == string(EmptyCell) {
			return NewGameError(ErrInvalidMove, "diagonal moves must capture")
		}
		if abs(m.ToCol-m.FromCol) != 1 {
			return NewGameError(ErrInvalidMove, "can only move diagonally one space")
		}
		if targetCell == player {
			return NewGameError(ErrInvalidMove, "cannot capture own piece")
		}
	} else {
		// Forward move must be to empty cell
		if targetCell != string(EmptyCell) {
			return NewGameError(ErrInvalidMove, "forward moves must be to empty space")
		}
	}

	return nil
}

// isInBounds checks if the given coordinates are within the board boundaries
func (b *Board) isInBounds(row, col int) bool {
	return row >= 0 && row < b.Rows && col >= 0 && col < b.Cols
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Apply a move to the board
func (b *Board) ApplyMove(ms string) {
	m, err := b.MoveFromString(ms)
	if err != nil {
		log.Fatal(err)
	}
	b.Grid[m.ToRow][m.ToCol] = b.Grid[m.FromRow][m.FromCol]
	b.Grid[m.FromRow][m.FromCol] = string(EmptyCell)
}

// Check for a win condition
func (b *Board) CheckWin() string {
	for col := 0; col < b.Cols; col++ {
		if b.Grid[0][col] == string(BlackPlayer) {
			return string(BlackPlayer)
		}
		if b.Grid[b.Rows-1][col] == string(WhitePlayer) {
			return string(WhitePlayer)
		}
	}
	return ""
}

// GenerateAvailableMoves generates all valid moves for the current position
func (p *Position) GenerateAvailableMoves() ([]Move, error) {
	if p.Board == nil {
		return nil, NewGameError(ErrInvalidBoard, "board is nil")
	}

	var moves []Move
	isWhite := p.Player == string(WhitePlayer)

	// For each square on the board
	for row := 0; row < p.Board.Rows; row++ {
		for col := 0; col < p.Board.Cols; col++ {
			// If we find our piece
			if p.Board.Grid[row][col] == p.Player {
				// Forward move
				newRow := row + 1
				if !isWhite {
					newRow = row - 1
				}

				// Check forward move
				if p.Board.isInBounds(newRow, col) &&
					p.Board.Grid[newRow][col] == string(EmptyCell) {
					moves = append(moves, Move{
						FromRow: row,
						FromCol: col,
						ToRow:   newRow,
						ToCol:   col,
					})
				}

				// Check diagonal captures
				for _, colOffset := range []int{-1, 1} {
					newCol := col + colOffset
					if p.Board.isInBounds(newRow, newCol) {
						targetCell := p.Board.Grid[newRow][newCol]
						if targetCell != string(EmptyCell) && targetCell != p.Player {
							moves = append(moves, Move{
								FromRow: row,
								FromCol: col,
								ToRow:   newRow,
								ToCol:   newCol,
							})
						}
					}
				}
			}
		}
	}

	return moves, nil
}

// Play the game
func (g *Game) Play() {
	currentPlayer := string(WhitePlayer)
	winner := ""
	var move string
	var stepNumber int

	for {
		fmt.Printf("Step %d\n", stepNumber+1)
		moves := g.Steps[stepNumber].Moves[g.Board.String()]
		if len(moves) == 0 {
			if currentPlayer == string(WhitePlayer) {
				winner = string(BlackPlayer)
			} else {
				winner = string(WhitePlayer)
			}
			fmt.Printf("Player %s has no moves; player %s wins!\n", currentPlayer, winner)
			break
		}
		switch g.NumPlayers {
		case 2:
			fmt.Printf("Player %s, enter your move: ", currentPlayer)
			fmt.Scan(&move)
		case 1:
			if currentPlayer == string(WhitePlayer) {
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

		err := g.Board.IsValidMove(move, currentPlayer)
		if err != nil {
			fmt.Println(err)
			continue
		}

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
		if currentPlayer == string(WhitePlayer) {
			currentPlayer = string(BlackPlayer)
		} else {
			currentPlayer = string(WhitePlayer)
		}
		stepNumber++
		if stepNumber >= len(g.Steps) {
			break
		}
	}
	g.Winner = winner
}

package hexapawn

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

// String is a stringer for Move
func (m Move) String() string {
	return fmt.Sprintf("%s%d-%s%d", letters[m.FromCol:m.FromCol+1], m.FromRow+1, letters[m.ToCol:m.ToCol+1], m.ToRow+1)
}

// String is a stringer for Board
func (b Board) String() string {
	var output string
	for i := 0; i < b.Rows; i++ {
		output += strings.Join(b.Grid[i], "")
	}
	return output
}

// Print the board state
func (b *Board) Print() {
	for i := b.Rows - 1; i >= 0; i-- {
		fmt.Printf("%2d  ", i+1)
		for _, cell := range b.Grid[i] {
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

func BoardFromString(boardStr string) *Board {
	if len(boardStr) == 0 {
		log.Fatal("Board string cannot be empty")
	}
	// Board string must be a square of integer
	rows := int(math.Sqrt(float64(len(boardStr))))
	if rows*rows != len(boardStr) {
		log.Fatal("Board string's length must be a square of integer: 9, 16, 25, etc.")
	}
	b := NewBoard(rows)
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			b.Grid[i][j] = boardStr[i*b.Cols+j : i*b.Cols+j+1]
		}
	}
	return b
}

// MoveFromString creates a new Move from a string of format "a1-a2", similar to chess
func (b *Board) MoveFromString(moveStr string) (Move, error) {
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

package hexapawn

import (
	"testing"
)

const (
	boardRows = 3
	boardCols = 3
)

func TestNewBoard(t *testing.T) {
	board := NewBoard(boardCols, boardRows, 0)
	if board.Rows != boardRows {
		t.Errorf("Expected Rows to be %d, got %d", boardRows, board.Rows)
	}
	if board.Cols != boardCols {
		t.Errorf("Expected Cols to be %d, got %d", boardCols, board.Cols)
	}
}

func TestNewMove(t *testing.T) {
	b := &Board{Cols: 8, Rows: 8}

	// Test valid move
	move, err := b.NewMove("a1-b2")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedMove := Move{0, 0, 1, 1}
	if move != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, move)
	}

	// Test invalid move string length
	_, err = b.NewMove("a1-b2-c3")
	if err == nil {
		t.Errorf("Expected error for invalid move string length, got nil")
	}

	// Test invalid move string format
	_, err = b.NewMove("a1-b2-c3")
	if err == nil {
		t.Errorf("Expected error for invalid move string format, got nil")
	}

	// Test invalid fromCol
	_, err = b.NewMove("i1-b2")
	if err == nil {
		t.Errorf("Expected error for invalid fromCol, got nil")
	}

	// Test invalid fromRow
	_, err = b.NewMove("a9-b2")
	if err == nil {
		t.Errorf("Expected error for invalid fromRow, got nil")
	}

	// Test invalid toCol
	_, err = b.NewMove("a1-i2")
	if err == nil {
		t.Errorf("Expected error for invalid toCol, got nil")
	}

	// Test invalid toRow
	_, err = b.NewMove("a1-b9")
	if err == nil {
		t.Errorf("Expected error for invalid toRow, got nil")
	}
}

package hexapawn

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	var boardRows = 3
	board := NewBoard(boardRows)
	if board.Rows != boardRows {
		t.Errorf("Expected Rows to be %d, got %d", boardRows, board.Rows)
	}
	if board.Cols != boardRows {
		t.Errorf("Expected Cols to be %d, got %d", boardRows, board.Cols)
	}
}

func TestBoardString(t *testing.T) {
	var boardRows = 3
	board := NewBoard(boardRows)
	boardString := board.String()
	if boardString != "WWW...BBB" {
		t.Errorf("Expected board string to be WWW...BBB, got %s", boardString)
	}
}

func TestNewGame(t *testing.T) {
	tests := []struct {
		boardRows   int
		numPlayers  int
		expectError bool
	}{
		{2, 2, true},  // Invalid board dimensions (too small)
		{10, 2, true}, // Invalid board dimensions (too large)
		{3, 3, true},  // Invalid player count
		{3, 2, false}, // Valid board dimensions
		{5, 1, false}, // Valid board dimensions and player count
	}

	for _, test := range tests {
		game, err := NewGame(test.boardRows, test.numPlayers)
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for boardRows %d, numPlayers %d, but got none", test.boardRows, test.numPlayers)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for boardRows %d, numPlayers %d, but got %v", test.boardRows, test.numPlayers, err)
			}
			if game == nil {
				t.Errorf("Expected valid game object for boardRows %d, numPlayers %d, but got nil", test.boardRows, test.numPlayers)
			}
		}
	}
}

func TestNewMove(t *testing.T) {
	b := &Board{Cols: 8, Rows: 8}

	// Test valid move
	move, err := b.MoveFromString("a1-b2")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedMove := Move{0, 0, 1, 1}
	if move != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, move)
	}

	// Test invalid move string length
	_, err = b.MoveFromString("a1-b2-c3")
	if err == nil {
		t.Errorf("Expected error for invalid move string length, got nil")
	}

	// Test invalid move string format
	_, err = b.MoveFromString("a1-b2-c3")
	if err == nil {
		t.Errorf("Expected error for invalid move string format, got nil")
	}

	// Test invalid fromCol
	_, err = b.MoveFromString("i1-b2")
	if err == nil {
		t.Errorf("Expected error for invalid fromCol, got nil")
	}

	// Test invalid fromRow
	_, err = b.MoveFromString("a9-b2")
	if err == nil {
		t.Errorf("Expected error for invalid fromRow, got nil")
	}

	// Test invalid toCol
	_, err = b.MoveFromString("a1-i2")
	if err == nil {
		t.Errorf("Expected error for invalid toCol, got nil")
	}

	// Test invalid toRow
	_, err = b.MoveFromString("a1-b9")
	if err == nil {
		t.Errorf("Expected error for invalid toRow, got nil")
	}
}

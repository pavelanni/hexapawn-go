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
	tests := []struct {
		name        string
		moveStr     string
		board       *Board
		expected    Move
		expectError bool
	}{
		{
			name:        "Valid move",
			moveStr:     "a1-b2",
			board:       &Board{Cols: 8, Rows: 8},
			expected:    Move{FromRow: 0, FromCol: 0, ToRow: 1, ToCol: 1},
			expectError: false,
		},
		{
			name:        "Invalid move string length",
			moveStr:     "a1-b2-c3",
			board:       &Board{Cols: 8, Rows: 8},
			expectError: true,
		},
		{
			name:        "Invalid from column",
			moveStr:     "i1-b2",
			board:       &Board{Cols: 8, Rows: 8},
			expectError: true,
		},
		{
			name:        "Invalid from row",
			moveStr:     "a9-b2",
			board:       &Board{Cols: 8, Rows: 8},
			expectError: true,
		},
		{
			name:        "Invalid to column",
			moveStr:     "a1-i2",
			board:       &Board{Cols: 8, Rows: 8},
			expectError: true,
		},
		{
			name:        "Invalid to row",
			moveStr:     "a1-b9",
			board:       &Board{Cols: 8, Rows: 8},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			move, err := tt.board.MoveFromString(tt.moveStr)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
				if move != tt.expected {
					t.Errorf("Expected move %v, got %v", tt.expected, move)
				}
			}
		})
	}
}

func TestPosition(t *testing.T) {
	board := BoardFromString("WWW...BBB")

	tests := []struct {
		name          string
		board         *Board
		player        string
		expectedMoves []Move
		expectError   bool
	}{
		{
			name:   "Initial position white",
			board:  board,
			player: "W",
			expectedMoves: []Move{
				{FromRow: 0, FromCol: 0, ToRow: 1, ToCol: 0},
				{FromRow: 0, FromCol: 1, ToRow: 1, ToCol: 1},
				{FromRow: 0, FromCol: 2, ToRow: 1, ToCol: 2},
			},
			expectError: false,
		},
		{
			name:   "Initial position black",
			board:  board,
			player: "B",
			expectedMoves: []Move{
				{FromRow: 2, FromCol: 0, ToRow: 1, ToCol: 0},
				{FromRow: 2, FromCol: 1, ToRow: 1, ToCol: 1},
				{FromRow: 2, FromCol: 2, ToRow: 1, ToCol: 2},
			},
			expectError: false,
		},
		{
			name:        "Nil board",
			board:       nil,
			player:      "W",
			expectError: true,
		},
		{
			name:   "Position with capture moves",
			board:  BoardFromString("W.W.B.B.B"),
			player: "W",
			expectedMoves: []Move{
				{FromRow: 0, FromCol: 0, ToRow: 1, ToCol: 1}, // Diagonal capture
				{FromRow: 0, FromCol: 2, ToRow: 1, ToCol: 1}, // Diagonal capture
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &Position{
				Board:  tt.board,
				Player: tt.player,
			}
			moves, err := pos.GenerateAvailableMoves()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
				if len(moves) != len(tt.expectedMoves) {
					t.Errorf("Expected %d moves, got %d", len(tt.expectedMoves), len(moves))
				}
				// Compare moves
				for i, move := range moves {
					if move != tt.expectedMoves[i] {
						t.Errorf("Move %d: expected %v, got %v", i, tt.expectedMoves[i], move)
					}
				}
			}
		})
	}
}

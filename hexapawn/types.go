// Package hexapawn implements the core game logic for Hexapawn, a simplified chess variant
// played with only pawns. The package provides:
//   - Board representation and game state management
//   - Move validation and execution
//   - Machine learning capabilities for AI players
//   - Game state persistence
package hexapawn

import (
	"log/slog"
)

// Player represents a game player
type Player string

// GameState represents the current state of the game
type GameState string

const (
	// GameInProgress indicates the game is still being played
	GameInProgress GameState = "IN_PROGRESS"
	// GameWonByWhite indicates white player has won
	GameWonByWhite GameState = "WHITE_WON"
	// GameWonByBlack indicates black player has won
	GameWonByBlack GameState = "BLACK_WON"
	// GameDrawn indicates the game ended in a draw
	GameDrawn GameState = "DRAWN"
)

// Define the Board struct
type Board struct {
	Rows  int        `json:"rows"`
	Cols  int        `json:"cols"`
	Grid  [][]string `json:"grid"`
	State GameState  `json:"stat"`
}

type BoardMove struct {
	BoardStr string `json:"board_str"`
	MoveStr  string `json:"move_str"`
}

type Game struct {
	NumPlayers    int         `json:"num_players"`
	Board         *Board      `json:"board"`
	Steps         []Step      `json:"steps"`
	MovesPlayed   []BoardMove `json:"moves_played"`
	CurrentPlayer string      `json:"current_player"`
	Winner        string      `json:"winner"`
}

type GamePlayed struct {
	MovesPlayed []BoardMove `json:"moves_played"`
	Winner      string      `json:"winner"`
}

type Machine struct {
	MachineFile string       `json:"machine_file"`
	NumRows     int          `json:"rows"`
	Steps       []Step       `json:"steps"`
	GamesPlayed []GamePlayed `json:"games_played"`
	Logger      *slog.Logger
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

// Position represents a game position with a board state, current player, and available moves
type Position struct {
	Board          *Board `json:"board"`
	Player         string `json:"player"`
	AvailableMoves []Move `json:"available_moves"`
}

type Step struct {
	Player string              `json:"player"` // current player
	Moves  map[string][]string `json:"moves"`  // each position is a string like "BBB...WWW"; for each position, there is a list of possible moves
}

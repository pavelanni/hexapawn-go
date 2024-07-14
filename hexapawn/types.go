package hexapawn

import (
	"log/slog"
)

// Define the Board struct
type Board struct {
	Rows int        `json:"rows"`
	Cols int        `json:"cols"`
	Grid [][]string `json:"grid"`
}

type BoardMove struct {
	BoardStr string `json:"board_str"`
	MoveStr  string `json:"move_str"`
}

type Game struct {
	NumPlayers  int         `json:"num_players"`
	Board       *Board      `json:"board"`
	Steps       []Step      `json:"steps"`
	MovesPlayed []BoardMove `json:"moves_played"`
	Winner      string      `json:"winner"`
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

type Step struct {
	Player string              `json:"player"` // current player
	Moves  map[string][]string `json:"moves"`  // each position is a string like "BBB...WWW"; for each position, there is a list of possible moves
}

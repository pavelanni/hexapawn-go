package hexapawn

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
)

// NewMachine creates a new machine learning instance
func NewMachine() *Machine {
	return &Machine{
		Steps:       make([]Step, 0),
		GamesPlayed: make([]GamePlayed, 0),
	}
}

// Play executes the specified number of games
func (m *Machine) Play(numGames int) error {
	if m.Logger == nil {
		return NewGameError(ErrMachineLearning, "logger is not initialized")
	}

	for i := 0; i < numGames; i++ {
		m.Logger.Info("starting game", slog.Int("game", i+1), slog.Int("total", numGames))

		g, err := NewGame(m.NumRows, 0)
		if err != nil {
			return fmt.Errorf("failed to create game: %w", err)
		}

		g.Steps = m.Steps
		g.Play()

		m.GamesPlayed = append(m.GamesPlayed, GamePlayed{
			MovesPlayed: g.MovesPlayed,
			Winner:      g.Winner,
		})

		m.Logger.Info("game completed",
			slog.Int("game", i+1),
			slog.String("winner", g.Winner),
			slog.Int("moves", len(g.MovesPlayed)))

		if err := m.Train(string(BlackPlayer)); err != nil {
			return fmt.Errorf("training failed: %w", err)
		}
	}
	return nil
}

// Train updates the machine learning model based on game results
func (m *Machine) Train(player string) error {
	lastGame := m.GamesPlayed[len(m.GamesPlayed)-1]
	if lastGame.Winner == player {
		return nil // no need to train on winning games
	}

	if len(lastGame.MovesPlayed) < 2 {
		return NewGameError(ErrMachineLearning, "insufficient moves for training")
	}

	// Get the losing move (second to last move)
	lastMove := lastGame.MovesPlayed[len(lastGame.MovesPlayed)-2]
	lastStep := len(lastGame.MovesPlayed) - 2

	m.Logger.Debug("training on losing move",
		slog.String("board", lastMove.BoardStr),
		slog.String("move", lastMove.MoveStr))

	// Remove the losing move from available moves
	moves, exists := m.Steps[lastStep].Moves[lastMove.BoardStr]
	if !exists {
		return NewGameError(ErrMachineLearning, "board state not found in training data")
	}

	lastMoveIndex := slices.Index(moves, lastMove.MoveStr)
	if lastMoveIndex == -1 {
		return NewGameError(ErrMachineLearning, "move not found in training data")
	}

	m.Steps[lastStep].Moves[lastMove.BoardStr] = slices.Delete(moves, lastMoveIndex, lastMoveIndex+1)
	return m.Save()
}

func (m *Machine) Load() error {
	fmt.Println("Loading machine...")

	data, err := os.ReadFile(m.MachineFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	fmt.Printf("Machine loaded from %s.\n", m.MachineFile)

	return nil
}

func (m *Machine) Save() error {
	fmt.Println("Saving machine...")
	f, err := os.Create(m.MachineFile)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	fmt.Printf("Machine saved to %s.\n", m.MachineFile)

	return nil
}

// Generate generates possible positions and moves
func (m *Machine) Init(numRows, numPlayers int) error {
	m.NumRows = numRows
	// Generate the first game
	game, err := NewGame(numRows, numPlayers)
	if err != nil {
		return err
	}
	// Generate the first step
	var newPlayer string
	steps := make([]Step, 1)
	moves := game.Board.ValidMoves("W") // Generate the first set of moves for the white player
	steps[0] = Step{
		Player: "W",
		Moves:  map[string][]string{game.Board.String(): moves},
	}
	// Generate the rest of the steps
	for step := 1; step < 10; step++ {
		if steps[step-1].Player == "W" {
			newPlayer = "B"
		} else {
			newPlayer = "W"
		}

		newStep := Step{
			Player: newPlayer,
			Moves:  make(map[string][]string),
		}
		fmt.Printf("Generating moves for step %d, player %s\n", step+1, newPlayer)
		for board, prevMoves := range steps[step-1].Moves { // for each board in the previous step
			for _, prevMove := range prevMoves { // for each move in the previous step
				newBoard := BoardFromString(board)
				newBoard.ApplyMove(prevMove)
				newMoves := []string{}
				if newBoard.CheckWin() != newPlayer {
					newMoves = newBoard.ValidMoves(newPlayer)
				}
				if len(newMoves) > 0 {
					newStep.Moves[newBoard.String()] = newMoves
					fmt.Printf("Added Position: %s, Player: %s, Moves: %v\n", newBoard.String(), newPlayer, newMoves)
				}
			}
		}
		steps = append(steps, newStep)
	}

	m.Steps = steps
	return nil
}

// ValidMoves returns possible moves for the given board and player
func (b *Board) ValidMoves(player string) []string {
	pieces := make([]Piece, 0)
	moves := make([]string, 0)

	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			if b.Grid[row][col] == player {
				pieces = append(pieces, Piece{Player: player, Row: row, Col: col})
			}
		}
	}

	for _, piece := range pieces {
		for row := 0; row < b.Rows; row++ {
			for col := 0; col < b.Cols; col++ {
				m := Move{FromRow: piece.Row, FromCol: piece.Col, ToRow: row, ToCol: col}
				if err := b.IsValidMove(m.String(), player); err == nil {
					moves = append(moves, m.String())
				} else {
					continue
				}
			}
		}
	}
	/*
		fmt.Printf("Possible moves for player %s: ", player)
		for _, move := range moves {
			fmt.Printf("%s ", move)
		}
		fmt.Println()
	*/
	return moves
}

// Package hexapawn provides the core game logic and machine learning capabilities
// for the Hexapawn game.
package hexapawn

import "fmt"

// GameError represents a game-specific error
type GameError struct {
	Code    ErrorCode
	Message string
}

func (e *GameError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ErrorCode represents different types of game errors
type ErrorCode string

const (
	// ErrInvalidMove indicates an illegal move attempt
	ErrInvalidMove ErrorCode = "INVALID_MOVE"
	// ErrInvalidBoard indicates an invalid board configuration
	ErrInvalidBoard ErrorCode = "INVALID_BOARD"
	// ErrInvalidPlayer indicates an invalid player number or type
	ErrInvalidPlayer ErrorCode = "INVALID_PLAYER"
	// ErrGameOver indicates an attempt to make a move in a finished game
	ErrGameOver ErrorCode = "GAME_OVER"
	// ErrMachineLearning indicates an error in the machine learning process
	ErrMachineLearning ErrorCode = "MACHINE_LEARNING_ERROR"
)

// NewGameError creates a new GameError with the given code and message
func NewGameError(code ErrorCode, message string) *GameError {
	return &GameError{
		Code:    code,
		Message: message,
	}
}

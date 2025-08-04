# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development commands

### Running the application
```bash
go run main.go                                    # Run with default settings (3x3 board, 2 players, 20 games)
go run main.go -r 5 -p 1 -g 10                   # Run with 5x5 board, 1 human player, 10 games
go run main.go --filename custom_machine.json    # Use custom machine learning file
go run main.go --logfile custom_log.json         # Use custom log file
```

### Visualization modes
```bash
go run main.go -g 10 -p 0 -v                     # Text summary with statistics (works everywhere)
go run main.go -g 5 -p 0 -i                      # Interactive TUI with auto-play (requires compatible terminal)
go run main.go -g 1 -p 0 -R                      # Step-by-step game replay
go run main.go -g 3 -p 0 -v -R                   # Combined text summary and replay
```

**Important for TUI mode (`-i`)**: Must be run as compiled binary directly in terminal (Terminal.app, iTerm2, Ghostty), NOT through IDEs or Claude Code:
```bash
go build
./hexapawn-go -g 5 -p 0 -i
```

### Testing
```bash
go test ./...                      # Run all tests
go test ./hexapawn                 # Run tests for hexapawn package only
go test -v ./hexapawn              # Run tests with verbose output
go test -run TestNewBoard ./hexapawn  # Run specific test
```

### Building
```bash
go build                           # Build executable in current directory
go build -o hexapawn-game         # Build with custom executable name
```

### Code quality
```bash
go fmt ./...                       # Format all Go code
go vet ./...                       # Run Go static analysis
go mod tidy                        # Clean up module dependencies
```

## Architecture overview

### Core components

**Game engine (`hexapawn/hexapawn.go`)**
- `Board`: Represents the game board with NxN grid (default 3x3)
- `Game`: Manages game state, player turns, and win conditions
- `Position`: Represents board state with available moves for machine learning
- Move validation and execution system

**Machine learning (`hexapawn/machine.go`)**
- `Machine`: Implements reinforcement learning that improves by removing losing moves
- `Step`: Represents decision points in the game tree with available moves per board state
- Training system that learns from losses by eliminating bad moves from future games
- Persistence to JSON files for learning continuity

**Game types (`hexapawn/types.go`)**
- Core data structures for board representation, moves, and game state
- JSON serialization support for machine learning persistence
- Error handling types in `hexapawn/errors.go`

**Visualization system (`hexapawn/ui.go`)**
- `GameViewer`: Interactive TUI using Bubbletea framework
- `ShowGamesText()`: Text-based summary with statistics
- `ShowGameReplay()`: Step-by-step game replay functionality
- TTY compatibility checking for robust cross-environment support

### Game mechanics

**Board representation**: String format like "WWW...BBB" where W=white pawns, B=black pawns, .=empty
**Move notation**: Algebraic notation like "a1-b2" (from square to square)
**Win conditions**: 
1. Advance pawn to opposite end
2. Capture all enemy pieces  
3. Block opponent from moving

### Machine learning approach

The AI uses a simplified reinforcement learning approach:
1. Generates all possible game positions and moves during initialization
2. Plays games by randomly selecting from available moves
3. After losing, removes the losing move from that board position
4. Over time, eliminates bad moves to improve play quality

### File structure

- `main.go`: CLI interface and application entry point
- `hexapawn/`: Core game logic package
  - `hexapawn.go`: Game engine and board logic
  - `machine.go`: Machine learning implementation
  - `types.go`: Core data structures
  - `ui.go`: Visualization and TUI components
  - `conversions.go`: String/board conversion utilities
  - `errors.go`: Custom error types
- `machine.json`: Default machine learning state persistence
- `hexapawn_log.json`: Game logging in JSON format
- `notes.md`: Development notes and game analysis
- `SESSION_SUMMARY.md`: Implementation session documentation

## Visualization features

The project includes multiple visualization modes to demonstrate machine learning progress:

### Interactive TUI (`-i` flag)
- Full-screen terminal interface with Bubbletea framework
- Real-time game replay with auto-play functionality
- Controls: ←→ (moves), ↑↓ (games), Space (play/pause), r (reset), e (end), q (quit)
- Requires direct terminal execution (not through IDEs)

### Text summary (`-v` flag)
- Statistics showing win percentages and learning progress
- Recent games with move sequences and final board states  
- Works in all environments including IDEs and scripts

### Step-by-step replay (`-R` flag)
- Manual advancement through moves of the most recent game
- Interactive demonstration of auto-play functionality
- Fallback when TUI is not available

## Dependencies

- `github.com/spf13/pflag`: Command-line flag parsing
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Terminal styling
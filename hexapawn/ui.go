package hexapawn

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI colors and styles
var (
	boardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)

	whiteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240")).
			Bold(true)

	blackStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("252")).
			Bold(true)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Align(lipgloss.Center)

	moveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Italic(true)
)

// GameViewer represents the TUI model for viewing games
type GameViewer struct {
	games       []GamePlayed
	currentGame int
	currentMove int
	board       *Board
	width       int
	height      int
	playing     bool
	speed       time.Duration
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*1200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// NewGameViewer creates a new game viewer with the provided games
func NewGameViewer(games []GamePlayed, boardSize int) *GameViewer {
	return &GameViewer{
		games:       games,
		currentGame: 0,
		currentMove: -1, // Start before first move to show initial board
		board:       NewBoard(boardSize),
		speed:       time.Millisecond * 800,
	}
}

func (gv GameViewer) Init() tea.Cmd {
	return nil
}

func (gv GameViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		gv.width = msg.Width
		gv.height = msg.Height
		return gv, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return gv, tea.Quit

		case "left", "h":
			if gv.currentMove > -1 {
				gv.currentMove--
				gv.updateBoard()
			}
			return gv, nil

		case "right", "l":
			if gv.currentGame < len(gv.games) && gv.currentMove < len(gv.games[gv.currentGame].MovesPlayed)-1 {
				gv.currentMove++
				gv.updateBoard()
			}
			return gv, nil

		case "up", "k":
			if gv.currentGame > 0 {
				gv.currentGame--
				gv.currentMove = -1
				gv.updateBoard()
			}
			return gv, nil

		case "down", "j":
			if gv.currentGame < len(gv.games)-1 {
				gv.currentGame++
				gv.currentMove = -1
				gv.updateBoard()
			}
			return gv, nil

		case " ":
			// Space key - toggle play/pause
			gv.playing = !gv.playing
			if gv.playing {
				return gv, tick()
			}
			return gv, nil

		case "r":
			gv.currentMove = -1
			gv.updateBoard()
			return gv, nil

		case "e":
			if gv.currentGame < len(gv.games) {
				gv.currentMove = len(gv.games[gv.currentGame].MovesPlayed) - 1
				gv.updateBoard()
			}
			return gv, nil
		}

	case tickMsg:
		if gv.playing {
			if gv.currentGame < len(gv.games) && gv.currentMove < len(gv.games[gv.currentGame].MovesPlayed)-1 {
				gv.currentMove++
				gv.updateBoard()
				return gv, tick()
			} else {
				gv.playing = false
			}
		}
		return gv, nil
	}

	return gv, nil
}

func (gv *GameViewer) updateBoard() {
	// Reset board to initial state
	gv.board = NewBoard(gv.board.Rows)

	// Apply moves up to current position
	if gv.currentGame < len(gv.games) && gv.currentMove >= 0 {
		game := gv.games[gv.currentGame]
		for i := 0; i <= gv.currentMove && i < len(game.MovesPlayed); i++ {
			move := game.MovesPlayed[i]
			// Apply the move to our board
			gv.board.ApplyMove(move.MoveStr)
		}
	}
}

func (gv GameViewer) View() string {
	if len(gv.games) == 0 {
		return "No games to display. Press 'q' to quit."
	}

	var sections []string

	// Header
	header := headerStyle.Render(fmt.Sprintf("Hexapawn Game Viewer - Game %d/%d", gv.currentGame+1, len(gv.games)))
	sections = append(sections, header)

	// Game info
	if gv.currentGame < len(gv.games) {
		game := gv.games[gv.currentGame]
		moveDisplay := "Start"
		if gv.currentMove >= 0 {
			moveDisplay = fmt.Sprintf("%d/%d", gv.currentMove+1, len(game.MovesPlayed))
		} else {
			moveDisplay = fmt.Sprintf("0/%d", len(game.MovesPlayed))
		}
		gameInfo := fmt.Sprintf("Winner: %s | Move: %s", game.Winner, moveDisplay)
		if gv.playing {
			gameInfo += " | ▶ PLAYING (auto-advance every 1.2s)"
		} else {
			gameInfo += " | ⏸ PAUSED (press SPACE to play)"
		}
		sections = append(sections, statusStyle.Render(gameInfo))
	}

	// Board
	boardDisplay := gv.renderBoard()
	sections = append(sections, boardStyle.Render(boardDisplay))

	// Current move info
	if gv.currentGame < len(gv.games) && gv.currentMove >= 0 && gv.currentMove < len(gv.games[gv.currentGame].MovesPlayed) {
		move := gv.games[gv.currentGame].MovesPlayed[gv.currentMove]
		moveInfo := fmt.Sprintf("Current move: %s", moveStyle.Render(move.MoveStr))
		sections = append(sections, moveInfo)
	}

	// Controls
	controls := `Controls:
  ← → (h/l): Previous/Next move    ↑ ↓ (k/j): Previous/Next game
  Space: Play/Pause                r: Reset to start    e: End of game
  q: Quit`
	sections = append(sections, controls)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (gv GameViewer) renderBoard() string {
	var lines []string
	
	// Render from top to bottom (reverse row order for display)
	for i := gv.board.Rows - 1; i >= 0; i-- {
		line := fmt.Sprintf("%2d │ ", i+1)
		for j := 0; j < gv.board.Cols; j++ {
			cell := gv.board.Grid[i][j]
			var styled string
			switch cell {
			case "W":
				styled = whiteStyle.Render(" W ")
			case "B":
				styled = blackStyle.Render(" B ")
			default:
				styled = emptyStyle.Render(" · ")
			}
			line += styled + " "
		}
		lines = append(lines, line)
	}
	
	// Add column labels
	colLine := "   └─"
	for i := 0; i < gv.board.Cols; i++ {
		colLine += "────"
	}
	lines = append(lines, colLine)
	
	labelLine := "     "
	for i := 0; i < gv.board.Cols; i++ {
		labelLine += fmt.Sprintf(" %c  ", 'a'+i)
	}
	lines = append(lines, labelLine)
	
	return strings.Join(lines, "\n")
}

// RunGameViewerTUI starts the interactive TUI game viewer
func RunGameViewerTUI(games []GamePlayed, boardSize int) error {
	if len(games) == 0 {
		return fmt.Errorf("no games to display")
	}
	
	// Check TTY availability first
	fmt.Printf("Checking terminal compatibility...\n")
	if _, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err != nil {
		fmt.Printf("❌ TTY not available: %v\n", err)
		fmt.Printf("\n🔧 To use interactive mode:\n")
		fmt.Printf("  1. Compile: go build\n")
		fmt.Printf("  2. Run directly in terminal: ./hexapawn-go -g 5 -p 0 -i\n")
		fmt.Printf("  3. NOT through IDE, Claude Code, or scripts\n\n")
		return ShowGamesText(games, boardSize)
	}
	
	fmt.Printf("✅ TTY available, starting TUI...\n")
	
	viewer := NewGameViewer(games, boardSize)
	
	fmt.Printf("Controls: ←→ (moves), ↑↓ (games), SPACE (play/pause), r (reset), e (end), q (quit)\n\n")
	
	// Try without alt screen first for better macOS compatibility
	p := tea.NewProgram(viewer)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("TUI failed: %v\n", err)
		return ShowGamesText(games, boardSize)
	}
	return err
}

// ShowGamesText provides a simple text-based game viewer fallback
func ShowGamesText(games []GamePlayed, boardSize int) error {
	fmt.Printf("\n=== HEXAPAWN GAMES SUMMARY ===\n")
	fmt.Printf("Total games played: %d\n\n", len(games))
	
	// Show last 5 games in detail first
	startIdx := len(games) - 5
	if startIdx < 0 {
		startIdx = 0
	}
	
	fmt.Printf("Recent games (last %d):\n", len(games)-startIdx)
	for i := startIdx; i < len(games); i++ {
		game := games[i]
		fmt.Printf("\nGame %d: Winner = %s\n", i+1, game.Winner)
		fmt.Printf("Moves: ")
		for j, move := range game.MovesPlayed {
			if j > 0 {
				fmt.Print(" -> ")
			}
			fmt.Print(move.MoveStr)
		}
		fmt.Printf("\n")
		
		// Show final board state
		if len(game.MovesPlayed) > 0 {
			lastMove := game.MovesPlayed[len(game.MovesPlayed)-1]
			finalBoard := BoardFromString(lastMove.BoardStr)
			finalBoard.ApplyMove(lastMove.MoveStr)
			fmt.Printf("Final board:\n")
			finalBoard.Print()
		}
	}
	
	// Show statistics at the end so they're visible
	wWins := 0
	bWins := 0
	for _, game := range games {
		if game.Winner == "W" {
			wWins++
		} else {
			bWins++
		}
	}
	fmt.Printf("\n=== OVERALL STATISTICS ===\n")
	fmt.Printf("Total games: %d\n", len(games))
	fmt.Printf("White wins: %d (%.1f%%)\n", wWins, float64(wWins)/float64(len(games))*100)
	fmt.Printf("Black wins: %d (%.1f%%)\n", bWins, float64(bWins)/float64(len(games))*100)
	fmt.Printf("\nMachine learning progress: Black is winning %.1f%% of games!\n", float64(bWins)/float64(len(games))*100)
	
	return nil
}

// ShowGameReplay provides an animated text-based replay of the most recent game
func ShowGameReplay(games []GamePlayed, boardSize int) error {
	if len(games) == 0 {
		return fmt.Errorf("no games to display")
	}
	
	// Get the most recent game
	game := games[len(games)-1]
	
	fmt.Printf("\n=== GAME REPLAY ===\n")
	fmt.Printf("Replaying Game %d (Winner: %s)\n\n", len(games), game.Winner)
	
	// Show initial board
	board := NewBoard(boardSize)
	fmt.Printf("Initial position:\n")
	board.Print()
	fmt.Printf("\nPress Enter to advance moves, or Ctrl+C to stop...\n\n")
	
	// Wait for user input to start
	fmt.Scanln()
	
	// Show each move
	for i, move := range game.MovesPlayed {
		player := "W"
		if i%2 == 1 {
			player = "B"
		}
		
		fmt.Printf("Move %d: Player %s plays %s\n", i+1, player, move.MoveStr)
		board.ApplyMove(move.MoveStr)
		board.Print()
		
		if i < len(game.MovesPlayed)-1 {
			fmt.Printf("\nPress Enter for next move...\n")
			fmt.Scanln()
		}
		fmt.Printf("\n")
	}
	
	fmt.Printf("Game complete! Winner: %s\n", game.Winner)
	return nil
}
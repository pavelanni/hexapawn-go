# Task: Implement Comprehensive Visualization System for Hexapawn Machine Learning Game

## Objective
Add user-facing visualization features to the hexapawn-go machine learning game to help users understand the learning progress and replay games. The goal is to create multiple visualization modes that work across different environments while preserving the original Martin Gardner machine learning algorithm.

## Background
The current hexapawn-go project implements Martin Gardner's matchbox learning machine algorithm where an AI learns by eliminating losing moves over time. The project needs visualization features to make the machine learning progress visible and engaging for educational purposes.

## Requirements

### 1. Text-Based Visualization Mode (`-v` flag)
**Purpose**: Provide comprehensive game statistics and recent game summaries that work in all environments.

**Implementation Requirements**:
- Add `--visualize` / `-v` command-line flag
- Display overall statistics at the END of output (not beginning) to ensure visibility
- Show total games played, win percentages for both players
- Display recent games (last 5) with move sequences and final board states
- Include a "machine learning progress" message highlighting the learning effect
- Ensure output is properly formatted and easy to read

**Expected Output Format**:
```
=== HEXAPAWN GAMES SUMMARY ===
Total games played: 250

Recent games (last 5):
[detailed game information with moves and final boards]

=== OVERALL STATISTICS ===
Total games: 250
White wins: 20 (8.0%)
Black wins: 230 (92.0%)

Machine learning progress: Black is winning 92.0% of games!
```

### 2. Interactive TUI Mode (`-i` flag)
**Purpose**: Create an interactive terminal user interface with auto-play functionality for real-time game viewing.

**Technical Requirements**:
- Use Bubbletea framework (github.com/charmbracelet/bubbletea) for TUI
- Use Lipgloss (github.com/charmbracelet/lipgloss) for styling
- Implement TTY compatibility checking with clear error messages
- Provide graceful fallback to text mode when TUI unavailable

**Functionality Requirements**:
- Navigate between games and moves with arrow keys
- **Space key auto-play**: Toggle between PAUSED and PLAYING modes
- When PLAYING: automatically advance through moves every 1.2 seconds
- Display current game number, move number, winner, and play status
- Show styled game board with colored pieces (White/Black pawns)
- Controls: ←→ (moves), ↑↓ (games), Space (play/pause), r (reset), e (end), q (quit)

**Environment Compatibility**:
- Must work in direct terminal execution (Terminal.app, iTerm2, Ghostty)
- Should detect TTY availability and provide clear instructions when unavailable
- Should NOT work through IDEs or script execution contexts (by design)

### 3. Step-by-Step Replay Mode (`-R` flag)
**Purpose**: Provide manual game replay functionality that works everywhere as a fallback demonstration.

**Implementation Requirements**:
- Show initial board position
- Allow manual advancement through moves with Enter key
- Display each move with player indication and resulting board state
- Show game completion message with winner
- Work in all environments including those without TTY access

### 4. Command-Line Integration
**Flag Specifications**:
- `-v, --visualize`: Text summary mode
- `-i, --interactive`: TUI mode with auto-play
- `-R, --replay`: Step-by-step replay mode
- Flags should be combinable (e.g., `-v -R` for both modes)

**Help Documentation**:
- Update `--help` output with clear descriptions
- Include compatibility notes for TUI mode

## Technical Implementation Details

### Dependencies to Add
```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
```

### Key Components to Implement

#### 1. TUI Game Viewer (`hexapawn/ui.go`)
```go
type GameViewer struct {
    games       []GamePlayed
    currentGame int
    currentMove int
    playing     bool
    // ... other fields
}
```

**Critical Implementation Notes**:
- Space key detection: Use `case " ":` in key switch statement
- Auto-play timing: Use `tea.Tick(time.Millisecond*1200, ...)` for consistent timing
- Board update logic: Start from initial board and apply moves sequentially (avoid double-applying moves)
- TTY checking: Use `os.OpenFile("/dev/tty", os.O_RDWR, 0)` to verify availability

#### 2. Statistics Display Function
- Must place statistics at END of output for visibility
- Calculate win percentages with proper formatting
- Show learning progress indication

#### 3. Replay Functionality
- Manual step-through with user input
- Clear move annotations with player indicators
- Board state display after each move

### Error Handling and Fallbacks

#### TTY Compatibility Issues
**Problem**: TUI frameworks require direct TTY access
**Solution**: Implement checking with clear user guidance:
```
❌ TTY not available: open /dev/tty: device not configured

🔧 To use interactive mode:
  1. Compile: go build
  2. Run directly in terminal: ./hexapawn-go -g 5 -p 0 -i
  3. NOT through IDE, Claude Code, or scripts
```

#### Environment Detection
- Check TTY availability before attempting TUI initialization
- Provide specific instructions for different environments
- Graceful fallback with explanation when TUI fails

## Debugging Considerations

### Common Issues to Anticipate
1. **Space key not working**: Usually due to TTY access problems, not key handling
2. **Statistics not visible**: Ensure they appear at END of output
3. **TUI not starting**: Environment doesn't provide TTY access
4. **Auto-play not advancing**: Check tick message handling and move update logic

### Testing Strategy
- Test in multiple terminal environments (Terminal.app, iTerm2, third-party terminals)
- Verify graceful fallback behavior in non-compatible environments
- Test all key combinations and navigation
- Validate statistics accuracy and visibility

## Success Criteria

### Functional Requirements ✅
- [ ] Text mode (`-v`) shows statistics at bottom of output
- [ ] TUI mode (`-i`) launches with proper TTY checking
- [ ] Space key toggles auto-play between PAUSED/PLAYING states
- [ ] Auto-play advances moves automatically every 1.2 seconds
- [ ] All navigation controls work (arrows, r, e, q)
- [ ] Replay mode (`-R`) works in all environments
- [ ] Graceful fallback when TUI unavailable

### Educational Value ✅
- [ ] Machine learning progress clearly visible (Black win percentage increasing)
- [ ] Game replay demonstrates strategy evolution
- [ ] Statistics show learning effect over time
- [ ] User can observe how AI eliminates losing moves

### Cross-Platform Compatibility ✅
- [ ] Works on macOS (Terminal.app, iTerm2, Ghostty)
- [ ] Provides clear error messages and instructions
- [ ] Fallback modes ensure functionality everywhere
- [ ] No breaking changes to existing functionality

## Deliverables

### Code Files
- `hexapawn/ui.go`: Complete TUI implementation
- `main.go`: Updated with new flags and integration
- `go.mod`: Updated dependencies

### Documentation
- Updated `CLAUDE.md` with visualization commands
- Session summary documenting implementation details
- Clear usage examples and troubleshooting guide

## Additional Context

### Project Philosophy
- Preserve Martin Gardner's original learning algorithm unchanged
- Focus on educational value and learning demonstration
- Ensure accessibility across different environments
- Maintain clean, educational code that others can learn from

### User Experience Goals
- Make machine learning progress immediately visible
- Provide engaging way to watch AI learn over time
- Ensure tool works for students, educators, and enthusiasts
- Create "wow factor" when users see Black's win rate climb to 90%+

This task represents a complete enhancement of the hexapawn project from a basic CLI tool to a comprehensive educational visualization system that effectively demonstrates machine learning principles through Martin Gardner's elegant matchbox algorithm.
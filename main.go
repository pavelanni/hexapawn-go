// Package main implements a Hexapawn game with machine learning capabilities.
// Hexapawn is a simplified chess variant played on a 3x3 board (or larger) with only pawns.
// The game includes options for human players and AI opponents that learn from their mistakes.
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/pavelanni/hexapawn-go/hexapawn"
	flag "github.com/spf13/pflag"
)

const (
	// Default configuration values
	defaultBoardRows   = 3
	defaultNumPlayers  = 2
	defaultNumGames    = 20
	defaultMachineFile = "machine.json"
	defaultLogFile     = "hexapawn_log.json"

	// Validation constraints
	minBoardRows  = 3
	maxBoardRows  = 9
	minNumPlayers = 0
	maxNumPlayers = 2
	minNumGames   = 1
)

// Configuration holds all game settings
type Config struct {
	boardRows   int
	numPlayers  int
	numGames    int
	machineFile string
	logFile     string
	visualize   bool
	interactive bool
	replay      bool
}

func main() {
	var config Config

	flag.IntVarP(&config.boardRows, "rows", "r", defaultBoardRows, "Number of rows in the board, must be at least 3, at most 9. Default is 3.")
	flag.IntVarP(&config.numPlayers, "players", "p", defaultNumPlayers, "Number of human players: 0, 1, or 2. Default is 2.")
	flag.IntVarP(&config.numGames, "games", "g", defaultNumGames, "Number of games to play: at least 1. Default is 20.")
	flag.StringVarP(&config.machineFile, "filename", "f", defaultMachineFile, "Load the machine from this file. If it doesn't exist, a new machine will be created and saved into this file. Default is 'machine.json'.")
	flag.StringVarP(&config.logFile, "logfile", "l", defaultLogFile, "Log the game into this file. Default is 'hexapawn_log.json'.")
	flag.BoolVarP(&config.visualize, "visualize", "v", false, "Show text summary of played games after completion.")
	flag.BoolVarP(&config.interactive, "interactive", "i", false, "Show interactive TUI game viewer (requires compatible terminal).")
	flag.BoolVarP(&config.replay, "replay", "R", false, "Show step-by-step replay of the most recent game.")

	flag.Parse()

	if config.boardRows < minBoardRows || config.boardRows > maxBoardRows {
		log.Fatalf("Invalid board dimensions. Rows and Cols must be equal and at least %d, at most %d.", minBoardRows, maxBoardRows)
	}
	if config.numPlayers < minNumPlayers || config.numPlayers > maxNumPlayers {
		log.Fatalf("Invalid number of players. Must be %d, %d, or %d.", minNumPlayers, 1, maxNumPlayers)
	}
	if config.numGames < minNumGames {
		log.Fatalf("Invalid number of games. Must be at least %d.", minNumGames)
	}
	if config.machineFile == "" {
		log.Printf("Machine file is not specified. Using '%s'.", defaultMachineFile)
		config.machineFile = defaultMachineFile
	}
	if config.logFile == "" {
		log.Printf("Log file is not specified. Using '%s'.", defaultLogFile)
		config.logFile = defaultLogFile
	}

	machine := hexapawn.NewMachine()
	// if machineFile exists, load it
	config.machineFile = os.ExpandEnv(config.machineFile)
	machine.MachineFile = config.machineFile
	_, err := os.Stat(machine.MachineFile)
	if !os.IsNotExist(err) && err != nil {
		log.Fatal(err)
	}

	if err != nil {
		// if it doesn't exist, create a new machine and save it
		err := machine.Init(config.boardRows, config.numPlayers)
		if err != nil {
			log.Fatal(err)
		}
		err = machine.Save()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Machine saved to %s", machine.MachineFile)
	} else {
		// if it exists, load it
		err := machine.Load()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Machine loaded from %s", machine.MachineFile)
	}

	logFile := os.ExpandEnv(config.logFile)
	l, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewJSONHandler(l, nil))
	if logger == nil {
		log.Fatal("slog.NewJSONHandler failed")
	}
	logger.Info("Starting hexapawn", slog.String("filename", machine.MachineFile))
	machine.Logger = logger

	err = machine.Play(config.numGames)
	if err != nil {
		log.Fatal(err)
	}
	err = machine.Save()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Machine saved to %s", machine.MachineFile)

	// Show visualization if requested
	if config.visualize {
		log.Printf("Starting game text summary...")
		err = hexapawn.ShowGamesText(machine.GamesPlayed, config.boardRows)
		if err != nil {
			log.Printf("Visualization error: %v", err)
		}
	}
	
	// Show interactive TUI if requested
	if config.interactive {
		log.Printf("Starting interactive game viewer...")
		err = hexapawn.RunGameViewerTUI(machine.GamesPlayed, config.boardRows)
		if err != nil {
			log.Printf("Interactive viewer error: %v", err)
		}
	}
	
	// Show step-by-step replay if requested
	if config.replay {
		log.Printf("Starting game replay...")
		err = hexapawn.ShowGameReplay(machine.GamesPlayed, config.boardRows)
		if err != nil {
			log.Printf("Replay error: %v", err)
		}
	}
}

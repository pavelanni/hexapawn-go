package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/pavelanni/hexapawn-go/hexapawn"
	flag "github.com/spf13/pflag"
)

// Define constants for board dimensions
var boardRows, numPlayers, numGames int
var machineFile, logFile string

func main() {
	flag.IntVarP(&boardRows, "rows", "r", 3, "Number of rows in the board, must be at least 3, at most 9. Default is 3.")
	flag.IntVarP(&numPlayers, "players", "p", 2, "Number of human players: 0, 1, or 2. Default is 2.")
	flag.IntVarP(&numGames, "games", "g", 20, "Number of games to play: at least 1. Default is 20.")
	flag.StringVarP(&machineFile, "filename", "f", "machine.json", "Load the machine from this file. If it doesn't exist, a new machine will be created and saved into this file. Default is 'machine.json'.")
	flag.StringVarP(&logFile, "logfile", "l", "hexapawn_log.json", "Log the game into this file. Default is 'hexapawn_log.json'.")

	flag.Parse()

	if boardRows < 3 || boardRows > 9 {
		log.Fatalf("Invalid board dimensions. Rows and Cols must be equal and at least 3, at most 9.")
	}
	if numPlayers < 0 || numPlayers > 2 {
		log.Fatalf("Invalid number of players. Must be 0, 1, or 2.")
	}
	if numGames < 1 {
		log.Fatalf("Invalid number of games. Must be at least 1.")
	}
	if machineFile == "" {
		log.Printf("Machine file is not specified. Using 'machine.json'.")
	}
	if logFile == "" {
		log.Printf("Log file is not specified. Using 'hexapawn_log.json'.")
	}

	machine := hexapawn.NewMachine()
	// if machineFile exists, load it
	machineFile = os.ExpandEnv(machineFile)
	machine.MachineFile = machineFile
	_, err := os.Stat(machineFile)
	if !os.IsNotExist(err) && err != nil {
		log.Fatal(err)
	}

	if err != nil {
		// if it doesn't exist, create a new machine and save it
		err := machine.Init(boardRows, numPlayers)
		if err != nil {
			log.Fatal(err)
		}
		err = machine.Save()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Machine saved to %s", machineFile)
	} else {
		// if it exists, load it
		err := machine.Load()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Machine loaded from %s", machineFile)
	}

	logFile = os.ExpandEnv(logFile)
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

	err = machine.Play(numGames)
	if err != nil {
		log.Fatal(err)
	}
	err = machine.Save()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Machine saved to %s", machineFile)
}

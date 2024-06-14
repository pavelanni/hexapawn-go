package hexapawn

import "fmt"

var positionNumber int

// Define the PositionStore struct
type PositionStore struct {
	positions map[string][]Move
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		positions: make(map[string][]Move),
	}
}

func (ps *PositionStore) AddPosition(board *Board, moves []Move) {
	position := board.String()
	ps.positions[position] = moves
}

func GeneratePositions(store *PositionStore, board *Board, player string) {
	fmt.Printf("Generating positions for position number %d for %s\n", positionNumber, player)
	positionNumber++
	moves := board.GetMoves(player)
	if len(moves) == 0 {
		return
	}
	store.AddPosition(board, moves)

	for _, move := range moves {
		newBoard := BoardFromString(board.String())
		newBoard.ApplyMove(move)

		fmt.Printf("Move: %s\n", move)
		fmt.Println(newBoard)
		if player == "W" {
			player = "B"
		} else {
			player = "W"
		}
		GeneratePositions(store, newBoard, player)
	}
}

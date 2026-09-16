package connectfour

// GameState represents the current state of the game.
type GameState int

const (
	InProgress GameState = iota // game is still being played
	Won                         // a player has won
	Draw                        // board is full with no winner
)

func (g GameState) String() string {
	switch g {
	case InProgress:
		return "In Progress"
	case Won:
		return "Won"
	case Draw:
		return "Draw"
	default:
		return "Unknown"
	}
}

// IsOver returns true if the game has ended (won or draw).
func (g GameState) IsOver() bool {
	return g == Won || g == Draw
}

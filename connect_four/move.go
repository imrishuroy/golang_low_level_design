package connectfour

// Move records a single move in the game history.
// Used for implementing undo functionality.
type Move struct {
	Player *Player
	Row    int
	Col    int
}

// NewMove creates a new move record.
func NewMove(player *Player, row, col int) Move {
	return Move{Player: player, Row: row, Col: col}
}

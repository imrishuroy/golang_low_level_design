package connectfour

// DiscColor represents the color of a disc on the board.
type DiscColor int

const (
	Empty  DiscColor = iota // zero value marks an unoccupied cell
	Red                     // first player's disc
	Yellow                  // second player's disc
)

func (d DiscColor) String() string {
	switch d {
	case Red:
		return "Red"
	case Yellow:
		return "Yellow"
	default:
		return "Empty"
	}
}

// Symbol returns a single character representation for display.
func (d DiscColor) Symbol() string {
	switch d {
	case Red:
		return "R"
	case Yellow:
		return "Y"
	default:
		return "."
	}
}

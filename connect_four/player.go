package connectfour

// Player represents a participant in the game.
// It holds identifying information (name) and the disc color they play with.
type Player struct {
	name  string
	color DiscColor
}

// NewPlayer creates a new player with the given name and disc color.
func NewPlayer(name string, color DiscColor) *Player {
	return &Player{name: name, color: color}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Color() DiscColor {
	return p.color
}

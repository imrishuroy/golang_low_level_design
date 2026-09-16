package connectfour

import (
	"math/rand"
)

// BotEngine defines the interface for AI opponents.
// Different implementations can provide varying levels of difficulty.
type BotEngine interface {
	ChooseMove(game *Game) int
}

// SimpleBotEngine picks the first valid column.
// This is the minimal implementation from the article.
type SimpleBotEngine struct{}

func NewSimpleBotEngine() *SimpleBotEngine {
	return &SimpleBotEngine{}
}

func (b *SimpleBotEngine) ChooseMove(game *Game) int {
	moves := game.ValidMoves()
	if len(moves) == 0 {
		return -1
	}
	return moves[0]
}

// SmartBotEngine uses a simple strategy: win if possible, block opponent,
// prefer center columns, otherwise pick randomly.
type SmartBotEngine struct{}

func NewSmartBotEngine() *SmartBotEngine {
	return &SmartBotEngine{}
}

func (b *SmartBotEngine) ChooseMove(game *Game) int {
	moves := game.ValidMoves()
	if len(moves) == 0 {
		return -1
	}

	currentPlayer := game.CurrentPlayer()
	opponent := game.Opponent(currentPlayer)

	// 1. check for winning move
	if col := b.findWinningMove(game, currentPlayer); col != -1 {
		return col
	}

	// 2. block opponent's winning move
	if col := b.findWinningMove(game, opponent); col != -1 {
		return col
	}

	// 3. prefer center columns (better strategic position)
	centerCol := game.Board().Cols() / 2
	for _, offset := range []int{0, -1, 1, -2, 2, -3, 3} {
		col := centerCol + offset
		if col >= 0 && col < game.Board().Cols() && game.Board().CanPlace(col) {
			return col
		}
	}

	// 4. fallback: random valid move
	return moves[rand.Intn(len(moves))]
}

// findWinningMove checks if the player can win in one move.
func (b *SmartBotEngine) findWinningMove(game *Game, player *Player) int {
	board := game.Board()

	for _, col := range game.ValidMoves() {
		// simulate placing a disc
		boardCopy := board.Copy()
		row := boardCopy.PlaceDisc(col, player.Color())
		if row != -1 && boardCopy.CheckWin(row, col, player.Color()) {
			return col
		}
	}
	return -1
}

// RandomBotEngine picks a random valid column.
// Useful for testing and as a baseline.
type RandomBotEngine struct{}

func NewRandomBotEngine() *RandomBotEngine {
	return &RandomBotEngine{}
}

func (b *RandomBotEngine) ChooseMove(game *Game) int {
	moves := game.ValidMoves()
	if len(moves) == 0 {
		return -1
	}
	return moves[rand.Intn(len(moves))]
}

package connectfour

import (
	"testing"
)

func TestSimpleBotEngine(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)
	game := NewGame(p1, p2)

	bot := NewSimpleBotEngine()

	// simple bot should always pick the first valid column
	col := bot.ChooseMove(game)
	if col != 0 {
		t.Errorf("simple bot should pick column 0, got %d", col)
	}

	// fill column 0
	for i := 0; i < 6; i++ {
		game.MakeMoveByColumn(0)
	}

	// now should pick column 1
	col = bot.ChooseMove(game)
	if col != 1 {
		t.Errorf("simple bot should pick column 1, got %d", col)
	}
}

func TestSimpleBotEngineNoMoves(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)

	board := NewBoard(1, 1)
	game := NewGameWithBoard(p1, p2, board)

	game.MakeMoveByColumn(0)

	bot := NewSimpleBotEngine()
	col := bot.ChooseMove(game)
	if col != -1 {
		t.Errorf("should return -1 when no moves, got %d", col)
	}
}

func TestSmartBotWinsWhenPossible(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)
	game := NewGame(p1, p2)

	// set up bot with 3 in a row horizontally at bottom
	// p1 plays column 6, p2 plays 0,1,2
	game.MakeMove(p1, 6)
	game.MakeMove(p2, 0)
	game.MakeMove(p1, 6)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 6)
	game.MakeMove(p2, 2)
	game.MakeMove(p1, 5) // p1's turn, plays elsewhere

	// now it's bot's turn (p2) and can win with column 3
	bot := NewSmartBotEngine()
	col := bot.ChooseMove(game)
	if col != 3 {
		t.Errorf("smart bot should win at column 3, got %d", col)
	}
}

func TestSmartBotBlocksOpponent(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)
	game := NewGame(p1, p2)

	// p1 has 3 in a row, p2 should block
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 1)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 2)
	// p2's turn, p1 can win at column 3

	bot := NewSmartBotEngine()
	col := bot.ChooseMove(game)
	if col != 3 {
		t.Errorf("smart bot should block at column 3, got %d", col)
	}
}

func TestSmartBotPrefersCenter(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)
	game := NewGame(p1, p2)

	// make first move elsewhere so bot goes second
	game.MakeMove(p1, 0)

	bot := NewSmartBotEngine()
	col := bot.ChooseMove(game)

	// should prefer center column (3) when no immediate win/block
	if col != 3 {
		t.Errorf("smart bot should prefer center column 3, got %d", col)
	}
}

func TestSmartBotNoMoves(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)

	board := NewBoard(1, 1)
	game := NewGameWithBoard(p1, p2, board)

	game.MakeMoveByColumn(0)

	bot := NewSmartBotEngine()
	col := bot.ChooseMove(game)
	if col != -1 {
		t.Errorf("should return -1 when no moves, got %d", col)
	}
}

func TestRandomBotEngine(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)
	game := NewGame(p1, p2)

	bot := NewRandomBotEngine()

	// just verify it returns a valid column
	col := bot.ChooseMove(game)
	if col < 0 || col > 6 {
		t.Errorf("random bot should return valid column, got %d", col)
	}
}

func TestRandomBotEngineNoMoves(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)

	board := NewBoard(1, 1)
	game := NewGameWithBoard(p1, p2, board)

	game.MakeMoveByColumn(0)

	bot := NewRandomBotEngine()
	col := bot.ChooseMove(game)
	if col != -1 {
		t.Errorf("should return -1 when no moves, got %d", col)
	}
}

func TestBotEngineInterface(t *testing.T) {
	// verify all bots implement the interface
	var _ BotEngine = NewSimpleBotEngine()
	var _ BotEngine = NewSmartBotEngine()
	var _ BotEngine = NewRandomBotEngine()
}

func TestSmartBotPrioritizesWinOverBlock(t *testing.T) {
	p1 := NewPlayer("Human", Red)
	p2 := NewPlayer("Bot", Yellow)
	game := NewGame(p1, p2)

	// setup: both players have 3 in a row, bot can win at 3 or block at 4
	// p2 (bot): 0,1,2 in bottom row
	// p1: stacked in column 5,6

	game.MakeMove(p1, 5)
	game.MakeMove(p2, 0)
	game.MakeMove(p1, 5)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 5)
	game.MakeMove(p2, 2)
	// p1 has vertical threat at 5, p2 can win at 3

	game.MakeMove(p1, 6) // p1 plays elsewhere

	bot := NewSmartBotEngine()
	col := bot.ChooseMove(game)

	// bot should win rather than block
	if col != 3 {
		t.Errorf("bot should prioritize winning at 3, got %d", col)
	}
}

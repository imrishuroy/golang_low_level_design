package connectfour

import (
	"errors"
	"testing"
)

func TestNewGame(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	if game.State() != InProgress {
		t.Error("new game should be in progress")
	}
	if game.Winner() != nil {
		t.Error("new game should have no winner")
	}
	if game.CurrentPlayer() != p1 {
		t.Error("player1 should go first")
	}
	if game.MoveCount() != 0 {
		t.Error("new game should have no moves")
	}
}

func TestMakeMove(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	err := game.MakeMove(p1, 0)
	if err != nil {
		t.Errorf("valid move should succeed: %v", err)
	}

	if game.CurrentPlayer() != p2 {
		t.Error("turn should switch to player2")
	}
	if game.MoveCount() != 1 {
		t.Error("move count should be 1")
	}
}

func TestMakeMoveByColumn(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	err := game.MakeMoveByColumn(3)
	if err != nil {
		t.Errorf("valid move should succeed: %v", err)
	}

	if game.Board().Cell(5, 3) != Red {
		t.Error("disc should be at (5,3)")
	}
}

func TestMakeMoveWrongTurn(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	err := game.MakeMove(p2, 0)
	if !errors.Is(err, ErrNotYourTurn) {
		t.Errorf("expected ErrNotYourTurn, got %v", err)
	}
}

func TestMakeMoveInvalidColumn(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	err := game.MakeMove(p1, -1)
	if !errors.Is(err, ErrInvalidMove) {
		t.Errorf("expected ErrInvalidMove, got %v", err)
	}

	err = game.MakeMove(p1, 7)
	if !errors.Is(err, ErrInvalidMove) {
		t.Errorf("expected ErrInvalidMove, got %v", err)
	}
}

func TestMakeMoveFullColumn(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	// fill column 0 by alternating players
	for i := 0; i < 6; i++ {
		player := p1
		if i%2 == 1 {
			player = p2
		}
		err := game.MakeMove(player, 0)
		if err != nil {
			t.Fatalf("move %d failed: %v", i, err)
		}
	}

	// next move to full column should fail
	err := game.MakeMove(p1, 0)
	if !errors.Is(err, ErrInvalidMove) {
		t.Errorf("expected ErrInvalidMove for full column, got %v", err)
	}
}

func TestMakeMoveAfterGameOver(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	// p1 wins with horizontal
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 0)
	game.MakeMove(p1, 1)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 2)
	game.MakeMove(p2, 2)
	game.MakeMove(p1, 3) // p1 wins

	if game.State() != Won {
		t.Error("game should be won")
	}

	// try to make another move
	err := game.MakeMove(p2, 4)
	if !errors.Is(err, ErrGameOver) {
		t.Errorf("expected ErrGameOver, got %v", err)
	}
}

func TestWinHorizontal(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	// p1: 0,1,2,3 (bottom row)
	// p2: 0,1,2 (stacking on top)
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 0)
	game.MakeMove(p1, 1)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 2)
	game.MakeMove(p2, 2)
	game.MakeMove(p1, 3)

	if game.State() != Won {
		t.Error("game should be won")
	}
	if game.Winner() != p1 {
		t.Error("player1 should be the winner")
	}
}

func TestWinVertical(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	// p1 stacks 4 in column 0
	// p2 plays in column 1
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 1)
	game.MakeMove(p1, 0)

	if game.State() != Won {
		t.Error("game should be won")
	}
	if game.Winner() != p1 {
		t.Error("player1 should be the winner")
	}
}

func TestDraw(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)

	// use small board to make draw easier to test
	board := NewBoard(2, 2)
	game := NewGameWithBoard(p1, p2, board)

	// fill 2x2 board without winning (impossible to win on 2x2 anyway)
	game.MakeMove(p1, 0) // (1,0)
	game.MakeMove(p2, 1) // (1,1)
	game.MakeMove(p1, 1) // (0,1)
	game.MakeMove(p2, 0) // (0,0)

	if game.State() != Draw {
		t.Errorf("game should be draw, got %v", game.State())
	}
	if game.Winner() != nil {
		t.Error("draw should have no winner")
	}
}

func TestUndoLastMove(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	game.MakeMove(p1, 0)
	game.MakeMove(p2, 1)

	if game.CurrentPlayer() != p1 {
		t.Error("should be p1's turn")
	}

	err := game.UndoLastMove()
	if err != nil {
		t.Errorf("undo should succeed: %v", err)
	}

	if game.CurrentPlayer() != p2 {
		t.Error("after undo, should be p2's turn")
	}
	if game.Board().Cell(5, 1) != Empty {
		t.Error("cell should be empty after undo")
	}
	if game.MoveCount() != 1 {
		t.Error("move count should be 1 after undo")
	}
}

func TestUndoNoMoves(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	err := game.UndoLastMove()
	if !errors.Is(err, ErrNoMovesToUndo) {
		t.Errorf("expected ErrNoMovesToUndo, got %v", err)
	}
}

func TestUndoAfterWin(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	// p1 wins
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 1)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 2)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 3) // win

	if game.State() != Won {
		t.Error("game should be won")
	}

	// undo the winning move
	err := game.UndoLastMove()
	if err != nil {
		t.Errorf("undo should succeed: %v", err)
	}

	if game.State() != InProgress {
		t.Error("game should be in progress after undo")
	}
	if game.Winner() != nil {
		t.Error("winner should be nil after undo")
	}
}

func TestLastMove(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	_, ok := game.LastMove()
	if ok {
		t.Error("should have no last move initially")
	}

	game.MakeMove(p1, 3)
	move, ok := game.LastMove()
	if !ok {
		t.Error("should have a last move")
	}
	if move.Col != 3 || move.Player != p1 {
		t.Error("last move should be p1 at column 3")
	}
}

func TestGameValidMoves(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	moves := game.ValidMoves()
	if len(moves) != 7 {
		t.Errorf("expected 7 valid moves, got %d", len(moves))
	}
}

func TestValidMovesAfterGameOver(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	// p1 wins
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 1)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 2)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 3)

	moves := game.ValidMoves()
	if moves != nil {
		t.Error("should have no valid moves after game over")
	}
}

func TestIsOver(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	if game.IsOver() {
		t.Error("new game should not be over")
	}

	// p1 wins
	game.MakeMove(p1, 0)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 1)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 2)
	game.MakeMove(p2, 6)
	game.MakeMove(p1, 3)

	if !game.IsOver() {
		t.Error("game should be over after win")
	}
}

func TestOpponent(t *testing.T) {
	p1 := NewPlayer("Alice", Red)
	p2 := NewPlayer("Bob", Yellow)
	game := NewGame(p1, p2)

	if game.Opponent(p1) != p2 {
		t.Error("opponent of p1 should be p2")
	}
	if game.Opponent(p2) != p1 {
		t.Error("opponent of p2 should be p1")
	}
}
